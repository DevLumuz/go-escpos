//go:build !windows

package escpos

import "fmt"

// NewWindowsPrinter creates a new printer connection to a Windows printer by name.
// This function is only available on Windows systems.
func NewWindowsPrinter(name string) (Printer, error) {
	return Printer{}, fmt.Errorf("Windows printer support is not available on this platform")
}

// GetInstalledPrinters returns a list of all installed printers on the system.
// This function is only available on Windows systems.
func GetInstalledPrinters() ([]string, error) {
	return nil, fmt.Errorf("printer enumeration is not available on this platform")
}
