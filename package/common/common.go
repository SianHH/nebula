package common

import "fmt"

type helpError struct {
	s string
}

func (he *helpError) Error() string {
	return he.s
}

func newHelpErrorf(s string, v ...interface{}) error {
	return &helpError{s: fmt.Sprintf(s, v...)}
}

func mustFlagString(name string, val *string) error {
	if *val == "" {
		return newHelpErrorf("-%s is required", name)
	}
	return nil
}
