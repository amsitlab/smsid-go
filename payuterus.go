package smsid

import (
	"errors"
	"fmt"
	//	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Payuterus struct {
	Adapter // Embed interface

	verbose     Verbose
	message     string // Text message
	phone       string // Phone number
	captcha     string // Captcha value
	initial     bool   // Initialize state
	url         string // Provider url
	client      *http.Client
	err         error // Error
	successText string
}

//
//
//
func (this *Payuterus) IsInitialized() bool {
	return this.initial
}

//
//
//
func (this *Payuterus) Initialize() {
	if nil == this.verbose {
		this.SetVerbose(new(NilVerbose))
	}

	this.verbose.Start()
	this.verbose.Success("[Initialize]")
	this.verbose.Info("Setting url")
	this.url = "http://sms.payuterus.biz/alpha/"

	this.verbose.Info("Initialize client")
	this.client = new(http.Client)

	this.verbose.Info("Finish..")
	this.successText = "Agar kami dapat mempertahankan kualitas layanan ini, mohon sudi kiranya untuk memberi rating bintang 5 pada play store, terima kasih"
	this.resolveCaptcha()
	this.initial = true
}

//
//
//
func (this *Payuterus) Terminate() {
	if nil != this.err {
		panic(this.err)
	}
	this.message = ""
	this.phone = ""
	this.captcha = ""
	this.url = ""
	this.verbose = nil
	this.initial = false

}

//
//
//
func (this *Payuterus) SetVerbose(verb Verbose) {
	this.verbose = verb
}

//
//
//
func (this *Payuterus) Send(phone, text string) Status {
	if !this.initial {
		this.err = errors.New("smsid.Payuterus: Missing initialize")
		return Failed
	}

	this.verbose.Success("[Send]")

	if this.captcha == "" {
		this.err = errors.New("smsid.Payuterus: Missing captcha, call func resolveCaptcha() on func Initialize()")
		return Failed
	}
	if this.url == "" {
		this.err = errors.New("smsid.Payuterus: Empty url")
		return Failed
	}

	this.phone = phone
	this.message = text

	this.verbose.Info("Make data..")
	this.verbose.Info("To: %s", phone)
	this.verbose.Info("Message: %s", text)

	data := url.Values{
		"nohp":    {this.phone},
		"pesan":   {this.message},
		"captcha": {this.captcha},
	}

	this.verbose.Info("Make request...")

	req, err := http.NewRequest(http.MethodPost, this.url+"send.php", strings.NewReader(data.Encode()))
	this.errMsg("smsid.Payuterus: Request error: %s", err)

	// use for Header Content-Length
	length := len(data.Encode()) + 21

	this.verbose.Info("Create request header")

	// Set(s) Header
	req.Header.Set("Host", "sms.payuterus.biz")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Length", strconv.Itoa(length))
	req.Header.Set("Save-Data", "on")
	req.Header.Set("Origin", "http://sms.payuterus.biz")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 5.0; ASUS_T00G Build/LRX21V) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.98 Mobile Safari/537.36")

	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Referer", "http://sms.payuterus.biz/alpha/")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "id-ID,id;q=0.8,en-US;q=0.6,en;q=0.4")
	req.Header.Set("Cookie", "PHPSESSID=f7h2k9jj1ugqf3cou3mem5bq35; _ga=GA1.2.1592713489.1519458607; _gid=GA1.2.1299671487.1519458611")
	this.verbose.Info("Requesting...")
	this.verbose.Info("Getting response...")

	response, err := this.client.Do(req)
	this.errMsg("smsid.Payuterus: Response error: %s", err)

	defer response.Body.Close()
	//	body, _ := ioutil.ReadAll(response.Body)
	//	body = []byte("")
	//	fmt.Println(string(body))
	status := this.isSuccess(response)

	this.verbose.Info("Is Sending?")
	if status == Success {
		this.verbose.Success("[?] Yes")
	} else {
		this.verbose.Warn("[?] No")
	}

	this.verbose.Info("Finished")
	return status
}

// ||||||||||||||||||||||||||||||||| [ PRIVATE ] ||||||||||||||||||||||||||||||||||

//
//
//
func (this *Payuterus) isSuccess(r *http.Response) Status {
	doc, err := html.Parse(r.Body)
	this.errMsg("smsid.Payuterus: Parsing error %s", err)

	var matcher func(*html.Node) (keep, exit bool)
	matcher = func(n *html.Node) (keep, exit bool) {
		if n.Type == html.TextNode && strings.TrimSpace(n.Data) != "" {
			keep = true
		}

		if n.DataAtom == atom.P {
			exit = false
		}

		return
	}

	nodes := this.traverse(doc, matcher)

	for _, node := range nodes {
		if strings.TrimSpace(node.Data) == this.successText {
			return Success
		}
	}

	return Failed
}

//
//
//
func (this *Payuterus) resolveCaptcha() {
	var result *html.Node

	this.verbose.Success("[Resolve captcha]")
	if this.url == "" {
		this.err = errors.New("smsid.Payuterus: Empty url")
		return
	}

	this.verbose.Info("Make request..")

	req, err := http.NewRequest(http.MethodGet, this.url, nil)
	this.errMsg("smsid.Payuterus: Request error: %s", err)

	this.verbose.Info("Create request header")

	req.Header.Set("Host", "sms.payuterus.biz")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Save-Data", "on")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 5.0; ASUS_T00G Build/LRX21V) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.98 Mobile Safari/537.36")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "id-ID,id;q=0.8,en-US;q=0.6,en;q=0.4")
	req.Header.Set("Cookie", "PHPSESSID=f7h2k9jj1ugqf3cou3mem5bq35; _ga=GA1.2.1592713489.1519458607; _gid=GA1.2.1299671487.1519458611")

	this.verbose.Info("Requesting..")

	response, err := this.client.Do(req)
	this.verbose.Info("Get response..")

	this.errMsg("smsid.Payuterus: Response error; %s", err)
	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		this.errMsg("smsid.Payuterus: Error response %s", err)
		return
	}

	this.verbose.Info("Parsing contents...")

	doc, err := html.Parse(response.Body)
	this.errMsg("smsid.Payuterus: Html parsing error: %s", err)

	var matcher func(*html.Node) (keep, exit bool)

	matcher = func(n *html.Node) (keep, exit bool) {
		if n.Type == html.TextNode && strings.TrimSpace(n.Data) != "" {
			keep = true
		}
		if n.DataAtom == atom.P {
			exit = true
		}

		return
	}

	nodes := this.traverse(doc, matcher)
	for i, _ := range nodes {
		//fmt.Println(node)

		if strings.TrimSpace(nodes[i].Data) == "Hasil Penjumlahan" {
			result = nodes[i+1]
		}
	}

	this.verbose.Info("Captcha operation...")
	this.verbose.Info("Captcha question:")

	text := strings.TrimSpace(result.Data)

	this.verbose.Warn("[?] %s", text)
	this.verbose.Info("Captcha answer:")

	text = strings.Replace(text, "=", "", -1)
	text = strings.TrimSpace(text)
	textArr := strings.Split(text, " ")
	this.resolveOperation(textArr)

	this.verbose.Warn("[*] %s", this.captcha)

}

//
//
//
func (this *Payuterus) traverse(doc *html.Node, matcher func(*html.Node) (bool, bool)) (nodes []*html.Node) {

	var keep, exit bool
	var f func(*html.Node)
	f = func(n *html.Node) {
		keep, exit = matcher(n)
		if keep {
			nodes = append(nodes, n)
		}
		if exit {
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return
}

//
//
//
func (this *Payuterus) resolveOperation(strArr []string) {

	var op int
	var sec int
	var err error

	op, err = strconv.Atoi(strArr[0])
	this.errMsg("smsid.Payuterus: Resolve operation error: %s", err)

	sec, err = strconv.Atoi(strArr[2])
	this.errMsg("smsid.Payuterus: Resolve operation error: %s", err)

	switch strArr[1] {
	case "+":
		op += sec
		break
	case "-":
		op -= sec
		break
	case ":", "/":
		op /= sec
		break
	case "*", "x", "X":
		op *= sec
		break
	}

	this.captcha = strconv.Itoa(op)

}

func (this *Payuterus) errMsg(text string, err error) {
	if err != nil {
		this.err = errors.New(fmt.Sprintf(
			text, err,
		))
	}
}
