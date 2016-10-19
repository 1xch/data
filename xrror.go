package data

import (
	"fmt"
)

type drror struct {
	base string
	vals []interface{}
}

func (d *drror) Error() string {
	return fmt.Sprintf("%s", fmt.Sprintf(d.base, d.vals...))
}

func (d *drror) Out(vals ...interface{}) *drror {
	d.vals = vals
	return d
}

func Drror(base string) *drror {
	return &drror{base: base}
}
