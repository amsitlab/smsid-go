package smsid

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

// ========================== Verbose ===========================

type Verbose interface {
	Start()

	NewLine(yes bool)

	SetPrefix(prefix string)

	Info(format string, args ...interface{})

	Warn(format string, args ...interface{})

	Success(format string, args ...interface{})

	Default(format string, arg ...interface{})
}

// ========================= HighlightVerbose ====================

type HighlightVerbose struct {
	call    func(fg, bg color.Attribute, format string, args []interface{})
	prefix  string
	newline bool
}

func (this *HighlightVerbose) Start() {
	color.HiCyan("Starting")
}

//
//
//
func (this *HighlightVerbose) SetPrefix(prefix string) {
	this.prefix = prefix
}

//
//
//
func (this *HighlightVerbose) NewLine(yes bool) {
	this.newline = yes
}

//
//
//
func (this *HighlightVerbose) Info(format string, args ...interface{}) {
	this.create()
	this.call(color.FgHiYellow, color.BgBlack, format, args)
}

//
//
//
func (this *HighlightVerbose) Warn(format string, args ...interface{}) {
	this.create()
	this.call(color.FgHiRed, color.BgBlack, format, args)
}

//
//
//
func (this *HighlightVerbose) Success(format string, args ...interface{}) {
	this.create()
	this.call(color.FgHiCyan, color.BgBlack, format, args)
}

//
//
//
func (this *HighlightVerbose) Default(format string, args ...interface{}) {
	this.create()
	this.call(color.FgWhite, color.BgBlack, format, args)
}

//
//
//

func (this *HighlightVerbose) create() {
	if !this.newline {
		time.Sleep(2 * time.Second)
		fmt.Printf("\033[1A")
		fmt.Printf("\033[K")
	}
	// terminalSize() has defined in interface.go
	w, _ := terminalSize()

	var spaces func(*color.Color, string)
	spaces = func(c *color.Color, text string) {
		w -= len(text)
		c.Printf("%s", strings.Repeat(" ", w))

	}

	this.call = func(fg, bg color.Attribute, text string, args []interface{}) {

		text = this.prefix + text
		c := color.New(fg, bg)
		switch len(args) {
		case 0:
			c.Printf(text)

			spaces(c, text)
			fmt.Print("\n")
			break
		case 1:
			c.Printf(text, args[0])
			spaces(c, fmt.Sprintf(text, args[0]))
			fmt.Print("\n")
			break
		case 2:
			c.Printf(text, args[0], args[1])
			spaces(c, fmt.Sprintf(text, args[0], args[1]))
			fmt.Print("\n")
			break
		case 3:
			c.Printf(text, args[0], args[1], args[2])
			spaces(c, fmt.Sprintf(text, args[0], args[1], args[2]))
			fmt.Print("\n")
			break
		default:
			panic("[smsid.HighlightVerbose] Too long Atgument(s)")
		}

	}

}

//
//
//

// =================== DisableVerbose =====================

type NilVerbose struct{}

func (this *NilVerbose) Start() {}

func (this *NilVerbose) SetPrefix(prefix string) {}

func (this *NilVerbose) NewLine(yes bool) {}

func (this *NilVerbose) Info(format string, args ...interface{}) {}

func (this *NilVerbose) Warn(format string, args ...interface{}) {}

func (this *NilVerbose) Success(format string, args ...interface{}) {}

func (this *NilVerbose) Default(format string, args ...interface{}) {}
