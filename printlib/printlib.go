package printlib

type Printer interface {
	Name() string
}

type DummyPrinter struct {
	name string
}

func (p DummyPrinter) Name() string {
	return p.name
}

func ListPrinters() ([]Printer, error) {
	dummyPrinters := make([]Printer, 1)
	var dummy DummyPrinter
	dummy.name = "Dummy Printer"
	dummyPrinters[0] = dummy
	return dummyPrinters, nil
}

// vim: ft=go
