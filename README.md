# go-escpos

[![Go Report Card](https://goreportcard.com/badge/github.com/DevLumuz/go-escpos)](https://goreportcard.com/report/github.com/DevLumuz/go-escpos)
![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)
[![GoDoc](https://godoc.org/github.com/DevLumuz/go-escpos?status.svg)](https://godoc.org/github.com/DevLumuz/go-escpos)

ESC/POS Thermal Printer library for Go with **native Windows support**.

This fork adds Windows API integration for direct printing to installed printers without network connectivity.

## Features

- ✅ Network printing (TCP/IP)
- ✅ Serial/USB printing (any `io.ReadWriteCloser`)
- ✅ **Windows native printing** (new!)
- ✅ Text formatting (bold, underline, fonts, sizes)
- ✅ Barcodes (UPC, EAN, Code39, Code93, Code128, etc.)
- ✅ Image printing (8-bit and 24-bit)
- ✅ Paper control (cut, feed)
- ✅ Printer status queries
- ✅ No external dependencies

## Installation

```bash
go get github.com/DevLumuz/go-escpos
```

## Usage

### Windows Native Printing

Print directly to installed Windows printers:

```go
package main

import (
	"log"
	"github.com/DevLumuz/go-escpos"
)

func main() {
	// List available printers
	printers, err := escpos.GetInstalledPrinters()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to first printer
	printer, err := escpos.NewWindowsPrinter(printers[0])
	if err != nil {
		log.Fatal(err)
	}
	defer printer.Close()

	// Print receipt
	printer.Initialize()
	
	printer.Justify(escpos.CenterJustify)
	printer.SetBold(true)
	printer.Println("MY STORE")
	printer.SetBold(false)
	printer.LF()
	
	printer.Justify(escpos.LeftJustify)
	printer.Println("Item 1............$10.00")
	printer.Println("Item 2............$15.00")
	printer.Println("----------------------")
	printer.SetBold(true)
	printer.Println("TOTAL.............$25.00")
	printer.SetBold(false)
	
	printer.FeedLines(3)
	printer.Cut()
}
```

### Network Printing (TCP/IP)

```go
package main

import (
	"net"
	"github.com/DevLumuz/go-escpos"
)

func main() {
	conn, _ := net.Dial("tcp", "192.168.1.100:9100")
	defer conn.Close()

	printer := escpos.NewPrinter(conn)
	printer.Println("Hello World!")
	printer.Cut()
}
```

### Serial/USB Printing

```go
printer := escpos.NewPrinter(anyIOReadWriteCloser)
printer.Println("Hello!")
printer.Cut()
```

## Demo Utility

The program in `./cmd/printhis/` is a demo utility to demonstrate some basic printing use cases.

## Testing

What? Did I hear you ask for testing? You think we make useless mocks that only tests our assumptions about the hoin printer instead of REAL **HONEST** ***GOOD*** boots on the ground testing.

Run `go run ./cmd/test-printer/` to print out our test program.

Really, how are we supposed to tests without a firmware dump? Total incongruity.

Also the test program assumes some things will work line printing and the such, cause how can we test functions without that. It'd be obvious if nothing prints. The goal is to test all the extra functions like horizontal tabbing, justifications, images, etc.

## Windows API

### GetInstalledPrinters()

List all installed printers on the system.

```go
printers, err := escpos.GetInstalledPrinters()
// Returns: []string{"Printer 1", "Printer 2", ...}
```

### NewWindowsPrinter(name)

Connect to a Windows printer by name.

```go
printer, err := escpos.NewWindowsPrinter("Your Printer Name")
defer printer.Close()
```

### Debug: Capture Bytes

Access raw bytes sent to printer for debugging.

```go
wp := printer.dst.(*escpos.WindowsPrinter)
bytes := wp.Bytes()
fmt.Printf("Sent %d bytes\n", len(bytes))
```

## Common Functions

```go
// Text formatting
printer.SetBold(true)
printer.SetFont(escpos.FontB)
printer.SetCharacterSize(2, 2)

// Alignment
printer.Justify(escpos.LeftJustify)
printer.Justify(escpos.CenterJustify)
printer.Justify(escpos.RightJustify)

// Paper control
printer.Feed(50)
printer.FeedLines(3)
printer.Cut()

// Barcodes
printer.SetHRIPosition(escpos.HRIBelow)
printer.PrintBarCode(escpos.BcCODE39, "123456")

// Sound
printer.Beep(3, 5)
```

## Testing

Test Windows printing:

```bash
# List printers
go run ./cmd/test-windows -list

# Run default test (receipt)
go run ./cmd/test-windows

# Run all tests
go run ./cmd/test-windows -test all

# Specific tests
go run ./cmd/test-windows -test receipt
go run ./cmd/test-windows -test barcode
go run ./cmd/test-windows -test format

# Specify printer
go run ./cmd/test-windows -printer "Your Printer Name" -test all
```

Test with real hardware (network/USB):

```bash
go run ./cmd/test-printer/ 192.168.1.100:9100
```

## Implementation

Windows support uses `winspool.drv` API via syscall:
- `OpenPrinterW` - Open printer handle
- `StartDocPrinterW` - Start print job
- `WritePrinter` - Send RAW data
- `EnumPrintersW` - List printers

No external dependencies required.

## Requirements

- **Windows**: Windows 7 or later
- **Go**: 1.16+
- **Printer**: Must support RAW mode and ESC/POS commands

## Troubleshooting

**"failed to open printer"**
- Use `GetInstalledPrinters()` to get exact names (case-sensitive)

**Nothing prints**
- Call `printer.Cut()` at the end
- Check printer has paper and is online

**Garbled output**
- Printer must support ESC/POS commands

## Fork Information

Fork of [joeyak/go-escpos](https://github.com/joeyak/go-escpos) with Windows native printing support.

**Changes:**
- Windows API integration (`printer_windows.go`)
- `GetInstalledPrinters()` function
- `WindowsPrinter.Bytes()` for debugging
- Windows test utility (`cmd/test-windows`)

## License

MIT License - See [license.md](license.md)

Original work Copyright 2022 joeyak  
Windows support Copyright 2025 DevLumuz
