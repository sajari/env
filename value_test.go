package env_test

import (
	"errors"
	"fmt"
	"strconv"
)

// positiveInteger is a custom environment variable type
type positiveInteger int

func (p *positiveInteger) String() string { return fmt.Sprintf("%d", *p) }

func (p *positiveInteger) Set(in string) error {
	n, err := strconv.Atoi(in)
	if err != nil {
		return err
	}
	if n < 0 {
		return errors.New("must be >= 0")
	}
	*p = positiveInteger(n)
	return nil
}
