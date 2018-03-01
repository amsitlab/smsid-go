package smsid

import (
	"syscall"
	"unsafe"
)

type Status int32

const (
	Failed Status = iota
	Success
)

type Adapter interface {
	Initialize()
	IsInitialized() bool
	Terminate()
	Send(phone, text string) Status
	SetVerbose(verb Verbose)
}

type Manager interface {
	SetAdapter(tag string, adapt Adapter) Manager

	Adapter(tag string) Adapter

	Initialize()

	IsInitialized() bool

	Terminate()

	Send(adaptTag, phone, message string) (stat Status)
}

// Struct for give terminal size
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	YPixel uint16
}

// Getting size of terminal
func terminalSize() (int, int) {

	ws := &winsize{}
	ret, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)
	if int(ret) == -1 {
		panic(err)
	}

	return int(ws.Col), int(ws.Row)
}
