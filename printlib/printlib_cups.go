package printlib

// +build linux

// #cgo LDFLAGS: -ldl
// #include "linux_cups.h"
import "C"

import (
	"errors"
	"unsafe"
)

var linux_cups_init_ok int = int(C.linux_cups_init_dl())

type CupsPrinter struct {
	name     string
	instance *string
}

func (p CupsPrinter) Name() string {
	if p.instance == nil {
		return p.name
	} else {
		return p.name + "/" + *(p.instance)
	}
}

func ListPrinters() ([]Printer, error) {
	if linux_cups_init_ok == 0 {
		return nil, errors.New("no library")
	}

	var dests *C.cups_dest_t = nil
	ndests := C.linux_cups_get_dests(&dests)
	defer C.linux_cups_free_dests(ndests, dests)

	printers := make([]Printer, ndests)

	var pdest uintptr = uintptr(unsafe.Pointer(dests))
	for i := 0; i < int(ndests); i++ {
		curpdest := (*C.cups_dest_t)(unsafe.Pointer(pdest))

		var p CupsPrinter
		p.name = C.GoString(curpdest.name)
		p.instance = nil
		if curpdest.instance != nil {
			instance := C.GoString(curpdest.instance)
			p.instance = &instance
		}
		printers[i] = p

		pdest = pdest + unsafe.Sizeof(*curpdest)
	}

	return printers, nil
}

func GetDefaultPrinter() (string, error) {
	if linux_cups_init_ok == 0 {
		return "", errors.New("no library")
	}

	var dests *C.cups_dest_t = nil
	ndests := C.linux_cups_get_dests(&dests)
	defer C.linux_cups_free_dests(ndests, dests)
	defaultdest := C.linux_cups_get_dest(nil, nil, ndests, dests)
	if defaultdest == nil {
		return "", errors.New("no default")
	}
	result := C.GoString(defaultdest.name)
	// is it actually possible for this to return an instance?
	if defaultdest.instance != nil {
		result = result + "/" + C.GoString(defaultdest.instance)
	}
	return result, nil
}

func Print(printer, data string) (bool, error) {
	return true, nil
}

// vim: ft=go
