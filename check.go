package env

import (
	"errors"
	"net"
	"os"
)

// checkedValue wraps a Value and runs fn on any values passed to Set
// before calling the underlying Value.Set.
type checkedValue struct {
	fn func(string) error

	Value
}

func (v checkedValue) Set(x string) error {
	if err := v.fn(x); err != nil {
		return err
	}
	return v.Value.Set(x)
}

// isNonEmpty checks if x is a non-empty string.
func isNonEmpty(x string) error {
	if x == "" {
		return errors.New("empty")
	}
	return nil
}

// isBindAddr checks if x is a valid bind address.
//
// A valid bind addresses is of the form host:port,
// and port must be non-empty.
func isBindAddr(x string) error {
	_, port, err := net.SplitHostPort(x)
	if err != nil {
		return err
	}
	if port == "" {
		return errors.New("empty port")
	}
	return nil
}

// isDialAddr checks if x is a valid bind address.
//
// A valid bind addresses is of the form host:port,
// and port must be non-empty.
func isDialAddr(x string) error {
	host, port, err := net.SplitHostPort(x)
	if err != nil {
		return err
	}
	if host == "" {
		return errors.New("empty host")
	}
	if port == "" {
		return errors.New("empty port")
	}
	return nil
}

// isPath checks if x is a valid path.
func isPath(x string) error {
	_, err := os.Stat(x)
	return err
}
