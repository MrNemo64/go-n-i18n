package util

type WarningsCollector struct {
	warnings []error
}

func NewWarningsCollector() *WarningsCollector {
	return &WarningsCollector{warnings: make([]error, 0)}
}

func (wc *WarningsCollector) AddWarning(err error) {
	wc.warnings = append(wc.warnings, err)
}

func (wc *WarningsCollector) Warnings() []error {
	return wc.warnings
}

func (wc *WarningsCollector) IsEmpty() bool {
	return len(wc.warnings) == 0
}

func (wc *WarningsCollector) Clear() {
	wc.warnings = make([]error, 0)
}
