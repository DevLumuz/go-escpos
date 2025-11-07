//go:build windows

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/DevLumuz/go-escpos"
)

func main() {
	list := flag.Bool("list", false, "List installed printers")
	printerName := flag.String("printer", "", "Printer name to use")
	testType := flag.String("test", "receipt", "Test type: receipt, barcode, image, all")
	flag.Parse()

	// List printers
	if *list {
		printers, err := escpos.GetInstalledPrinters()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Installed printers:")
		for i, name := range printers {
			fmt.Printf("%d. %s\n", i+1, name)
		}
		return
	}

	// Get printer name
	if *printerName == "" {
		printers, err := escpos.GetInstalledPrinters()
		if err != nil {
			log.Fatal(err)
		}
		if len(printers) == 0 {
			log.Fatal("No printers found")
		}
		*printerName = printers[0]
		fmt.Printf("Using default printer: %s\n", *printerName)
	}

	// Connect to printer
	printer, err := escpos.NewWindowsPrinter(*printerName)
	if err != nil {
		log.Fatal(err)
	}
	defer printer.Close()

	// Run tests
	switch *testType {
	case "receipt":
		testReceipt(printer)
	case "barcode":
		testBarcode(printer)
	case "format":
		testFormatting(printer)
	case "all":
		testReceipt(printer)
		printer.FeedLines(2)
		testBarcode(printer)
		printer.FeedLines(2)
		testFormatting(printer)
	default:
		log.Fatalf("Unknown test type: %s", *testType)
	}

	fmt.Println("Test completed successfully!")
}

func testReceipt(printer escpos.Printer) {
	printer.Initialize()

	// Header
	printer.Justify(escpos.CenterJustify)
	printer.SetBold(true)
	printer.SetCharacterSize(1, 1)
	printer.Println("TEST RECEIPT")
	printer.SetBold(false)
	printer.SetCharacterSize(0, 0)
	printer.Println("Windows Native Printing")
	printer.Println("================================")
	printer.LF()

	// Items
	printer.Justify(escpos.LeftJustify)
	printer.Println("Item                    Price")
	printer.Println("--------------------------------")
	printer.Println("Coffee                  $3.00")
	printer.Println("Sandwich               $12.00")
	printer.Println("Cookie                  $2.50")
	printer.Println("--------------------------------")

	// Total
	printer.SetBold(true)
	printer.Println("TOTAL:                 $17.50")
	printer.SetBold(false)
	printer.LF()

	// Footer
	printer.Justify(escpos.CenterJustify)
	printer.Println("Thank you!")
	printer.Println("Visit us again")

	printer.FeedLines(3)
	printer.Cut()
}

func testBarcode(printer escpos.Printer) {
	printer.Initialize()

	printer.Justify(escpos.CenterJustify)
	printer.SetBold(true)
	printer.Println("BARCODE TEST")
	printer.SetBold(false)
	printer.LF()

	// Test different barcode types
	printer.SetHRIPosition(escpos.HRIBelow)
	printer.SetBarCodeHeight(50)

	printer.Println("CODE39:")
	printer.PrintBarCode(escpos.BcCODE39, "TEST123")
	printer.LF()

	printer.Println("CODE128:")
	printer.PrintBarCode(escpos.BcCODE123, "ABC123")
	printer.LF()

	printer.FeedLines(3)
	printer.Cut()
}

func testFormatting(printer escpos.Printer) {
	printer.Initialize()

	printer.Justify(escpos.CenterJustify)
	printer.SetBold(true)
	printer.Println("FORMATTING TEST")
	printer.SetBold(false)
	printer.LF()

	printer.Justify(escpos.LeftJustify)

	// Bold
	printer.Println("Normal text")
	printer.SetBold(true)
	printer.Println("Bold text")
	printer.SetBold(false)
	printer.LF()

	// Fonts
	printer.SetFont(escpos.FontA)
	printer.Println("Font A (default)")
	printer.SetFont(escpos.FontB)
	printer.Println("Font B (smaller)")
	printer.SetFont(escpos.FontA)
	printer.LF()

	// Sizes
	printer.SetCharacterSize(0, 0)
	printer.Println("Normal size")
	printer.SetCharacterSize(1, 0)
	printer.Println("Double width")
	printer.SetCharacterSize(0, 1)
	printer.Println("Double height")
	printer.SetCharacterSize(1, 1)
	printer.Println("Double both")
	printer.SetCharacterSize(0, 0)
	printer.LF()

	// Justification
	printer.Justify(escpos.LeftJustify)
	printer.Println("Left aligned")
	printer.Justify(escpos.CenterJustify)
	printer.Println("Center aligned")
	printer.Justify(escpos.RightJustify)
	printer.Println("Right aligned")

	printer.FeedLines(3)
	printer.Cut()
}
