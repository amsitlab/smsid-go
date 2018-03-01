package smsid

import (
	"errors"
	"fmt"
	"os"
)

type Factory struct {
	Manager
	Verbose
	adapter     map[string]Adapter
	err         error
	initialized bool
}

func (this *Factory) Initialize() {
	if nil == this.adapter {
		this.adapter = make(map[string]Adapter)
	}
	if nil == this.Verbose {
		this.Verbose = new(NilVerbose)
	}

	this.initialized = true
	//register
	this.SetAdapter("payuterus", new(Payuterus))
}

//
//
//
func (this *Factory) IsInitialized() bool {
	return this.initialized
}

//
//
//
func (this *Factory) Terminate() {
	if nil != this.err {
		fmt.Printf("smsid.Factory: %s", this.err)
		os.Exit(1)
	}

}

//
//
//
func (this *Factory) SetAdapter(tag string, adapter Adapter) Manager {
	if nil == this.adapter {
		this.errMsg("%s", errors.New("smsid.Factory: Missing initialize"))
		return nil
	}
	this.adapter[tag] = adapter
	return this
}

//
//
//
func (this *Factory) Adapter(tag string) Adapter {
	if ok := this.adapter[tag]; ok == nil {
		this.errMsg("%s", errors.New("smsid.Factory: Unknown adapter"))
		return nil
	}
	return this.adapter[tag]
}

//
//
//
func (this *Factory) Send(adapterTag, phone, message string) (stat Status) {
	adapt := this.Adapter(adapterTag)
	adapt.SetVerbose(this.Verbose)
	adapt.Initialize()
	defer adapt.Terminate()
	stat = adapt.Send(phone, message)
	return
}

// ================================= PRIVATE ==================================== //

//
//
//
func (this *Factory) errMsg(format string, err error) {
	if nil != err {
		this.err = errors.New(fmt.Sprintf(
			format, err,
		))
	}
}
