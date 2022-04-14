package module

import "fmt"

type DuplicatedComponentErr struct{ component string }

func newDuplicatedComponentErr(component string) *DuplicatedComponentErr {
	return &DuplicatedComponentErr{component: component}
}

func (d *DuplicatedComponentErr) Error() string {
	return fmt.Sprintf("component %q already present on module", d.component)
}
