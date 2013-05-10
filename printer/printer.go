// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package printer

import (
	"syscall"
)

type DOC_INFO_1 struct {
	DocName    *uint16
	OutputFile *uint16
	Datatype   *uint16
}

//sys	GetDefaultPrinter(buf *uint16, bufN *uint32) (err error) = winspool.GetDefaultPrinterW
//sys	ClosePrinter(h syscall.Handle) (err error) = winspool.ClosePrinter
//sys	OpenPrinter(name *uint16, h *syscall.Handle, defaults uintptr) (err error) = winspool.OpenPrinterW
//sys	StartDocPrinter(h syscall.Handle, level uint32, docinfo *DOC_INFO_1) (err error) = winspool.StartDocPrinterW
//sys	EndDocPrinter(h syscall.Handle) (err error) = winspool.EndDocPrinter
//sys	WritePrinter(h syscall.Handle, buf *byte, bufN uint32, written *uint32) (err error) = winspool.WritePrinter
//sys	StartPagePrinter(h syscall.Handle) (err error) = winspool.StartPagePrinter
//sys	EndPagePrinter(h syscall.Handle) (err error) = winspool.EndPagePrinter

func Default() (string, error) {
	b := make([]uint16, 3)
	n := uint32(len(b))
	err := GetDefaultPrinter(&b[0], &n)
	if err != nil {
		if err != syscall.ERROR_INSUFFICIENT_BUFFER {
			return "", err
		}
		b = make([]uint16, n)
		err = GetDefaultPrinter(&b[0], &n)
		if err != nil {
			return "", err
		}
	}
	return syscall.UTF16ToString(b), nil
}

type Printer struct {
	h syscall.Handle
}

func Open(name string) (*Printer, error) {
	var p Printer
	// TODO: implement pDefault parameter
	err := OpenPrinter(&(syscall.StringToUTF16(name))[0], &p.h, 0)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Printer) StartDocument(name, datatype string) error {
	d := DOC_INFO_1{
		DocName:    &(syscall.StringToUTF16(name))[0],
		OutputFile: nil,
		Datatype:   &(syscall.StringToUTF16(datatype))[0],
	}
	return StartDocPrinter(p.h, 1, &d)
}

func (p *Printer) Write(b []byte) (int, error) {
	var written uint32
	err := WritePrinter(p.h, &b[0], uint32(len(b)), &written)
	if err != nil {
		return 0, err
	}
	return int(written), nil
}

func (p *Printer) EndDocument() error {
	return EndDocPrinter(p.h)
}

func (p *Printer) StartPage() error {
	return StartPagePrinter(p.h)
}

func (p *Printer) EndPage() error {
	return EndPagePrinter(p.h)
}

func (p *Printer) Close() error {
	return ClosePrinter(p.h)
}
