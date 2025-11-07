# Windows Printer Test Utility

Command-line utility to test Windows native printing functionality.

## Usage

### List installed printers

```bash
go run main.go -list
```

### Run default test (receipt)

```bash
go run main.go
```

### Run specific test

```bash
# Receipt test
go run main.go -test receipt

# Barcode test
go run main.go -test barcode

# Formatting test
go run main.go -test format

# All tests
go run main.go -test all
```

### Specify printer

```bash
go run main.go -printer "Your Printer Name" -test all
```

## Build

```bash
go build -o test-windows.exe
```

## Test Types

- **receipt**: Prints a sample receipt with items and total
- **barcode**: Tests different barcode formats (CODE39, CODE128)
- **format**: Tests text formatting (bold, fonts, sizes, alignment)
- **all**: Runs all tests

## Requirements

- Windows OS
- At least one printer installed
- Printer must support RAW/ESC-POS commands
