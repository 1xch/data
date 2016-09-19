package data

import (
	"fmt"
)

type drror struct {
	base string
	vals []interface{}
}

func (x *drror) Error() string {
	return fmt.Sprintf("%s", fmt.Sprintf(x.base, x.vals...))
}

func (x *drror) Out(vals ...interface{}) *drror {
	x.vals = vals
	return x
}

func Drror(base string) *drror {
	return &drror{base: base}
}
