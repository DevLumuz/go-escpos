//go:build windows

package escpos

import (
	"bytes"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	winspool             = syscall.NewLazyDLL("winspool.drv")
	procOpenPrinterW     = winspool.NewProc("OpenPrinterW")
	procClosePrinter     = winspool.NewProc("ClosePrinter")
	procStartDocPrinterW = winspool.NewProc("StartDocPrinterW")
	procEndDocPrinter    = winspool.NewProc("EndDocPrinter")
	procStartPagePrinter = winspool.NewProc("StartPagePrinter")
	procEndPagePrinter   = winspool.NewProc("EndPagePrinter")
	procWritePrinter     = winspool.NewProc("WritePrinter")
	procEnumPrintersW    = winspool.NewProc("EnumPrintersW")
)

const (
	PRINTER_ENUM_LOCAL       = 0x00000002
	PRINTER_ENUM_CONNECTIONS = 0x00000004
)

type PRINTER_INFO_4 struct {
	PrinterName *uint16
	ServerName  *uint16
	Attributes  uint32
}

type DOC_INFO_1 struct {
	DocName    *uint16
	OutputFile *uint16
	Datatype   *uint16
}

// WindowsPrinter implements io.ReadWriteCloser for Windows printers
type WindowsPrinter struct {
	name       string
	handle     syscall.Handle
	buffer     bytes.Buffer
	jobStarted bool
}

// NewWindowsPrinter creates a new printer connection to a Windows printer by name
func NewWindowsPrinter(name string) (Printer, error) {
	wp := &WindowsPrinter{
		name: name,
	}

	// Open printer
	namePtr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return Printer{}, fmt.Errorf("invalid printer name: %w", err)
	}

	ret, _, err := procOpenPrinterW.Call(
		uintptr(unsafe.Pointer(namePtr)),
		uintptr(unsafe.Pointer(&wp.handle)),
		0,
	)
	if ret == 0 {
		return Printer{}, fmt.Errorf("failed to open printer %q: %w", name, err)
	}

	// Start document
	if err := wp.startDoc(); err != nil {
		wp.Close()
		return Printer{}, err
	}

	return NewPrinter(wp), nil
}

func (wp *WindowsPrinter) startDoc() error {
	docName, err := syscall.UTF16PtrFromString("ESC/POS Document")
	if err != nil {
		return fmt.Errorf("failed to create doc name: %w", err)
	}

	datatype, err := syscall.UTF16PtrFromString("RAW")
	if err != nil {
		return fmt.Errorf("failed to create datatype: %w", err)
	}

	docInfo := DOC_INFO_1{
		DocName:    docName,
		OutputFile: nil,
		Datatype:   datatype,
	}

	ret, _, err := procStartDocPrinterW.Call(
		uintptr(wp.handle),
		1,
		uintptr(unsafe.Pointer(&docInfo)),
	)
	if ret == 0 {
		return fmt.Errorf("failed to start document: %w", err)
	}

	// Start page
	ret, _, err = procStartPagePrinter.Call(uintptr(wp.handle))
	if ret == 0 {
		return fmt.Errorf("failed to start page: %w", err)
	}

	wp.jobStarted = true
	return nil
}

// Write writes data to the printer
func (wp *WindowsPrinter) Write(p []byte) (int, error) {
	if !wp.jobStarted {
		return 0, fmt.Errorf("print job not started")
	}

	var written uint32
	ret, _, err := procWritePrinter.Call(
		uintptr(wp.handle),
		uintptr(unsafe.Pointer(&p[0])),
		uintptr(len(p)),
		uintptr(unsafe.Pointer(&written)),
	)
	if ret == 0 {
		return int(written), fmt.Errorf("failed to write to printer: %w", err)
	}

	// Store in buffer for Bytes() method
	wp.buffer.Write(p)

	return int(written), nil
}

// Read is not supported for Windows printers
func (wp *WindowsPrinter) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("read operation not supported on Windows printers")
}

// Close closes the printer connection
func (wp *WindowsPrinter) Close() error {
	if wp.jobStarted {
		procEndPagePrinter.Call(uintptr(wp.handle))
		procEndDocPrinter.Call(uintptr(wp.handle))
		wp.jobStarted = false
	}

	if wp.handle != 0 {
		ret, _, err := procClosePrinter.Call(uintptr(wp.handle))
		if ret == 0 {
			return fmt.Errorf("failed to close printer: %w", err)
		}
		wp.handle = 0
	}

	return nil
}

// Bytes returns all bytes that have been written to the printer
func (wp *WindowsPrinter) Bytes() []byte {
	return wp.buffer.Bytes()
}

// GetInstalledPrinters returns a list of all installed printers on the system
func GetInstalledPrinters() ([]string, error) {
	flags := PRINTER_ENUM_LOCAL | PRINTER_ENUM_CONNECTIONS
	var needed, returned uint32

	// First call to get buffer size
	procEnumPrintersW.Call(
		uintptr(flags),
		0,
		4, // PRINTER_INFO_4
		0,
		0,
		uintptr(unsafe.Pointer(&needed)),
		uintptr(unsafe.Pointer(&returned)),
	)

	if needed == 0 {
		return []string{}, nil
	}

	// Allocate buffer and call again
	buf := make([]byte, needed)
	ret, _, err := procEnumPrintersW.Call(
		uintptr(flags),
		0,
		4, // PRINTER_INFO_4
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(needed),
		uintptr(unsafe.Pointer(&needed)),
		uintptr(unsafe.Pointer(&returned)),
	)

	if ret == 0 {
		return nil, fmt.Errorf("failed to enumerate printers: %w", err)
	}

	// Parse printer names
	printers := make([]string, 0, returned)
	info := (*[1 << 20]PRINTER_INFO_4)(unsafe.Pointer(&buf[0]))[:returned:returned]

	for i := uint32(0); i < returned; i++ {
		if info[i].PrinterName != nil {
			name := syscall.UTF16ToString((*[1 << 10]uint16)(unsafe.Pointer(info[i].PrinterName))[:])
			printers = append(printers, name)
		}
	}

	return printers, nil
}
