// Package env provides a standardised interface to environment variables,
// including parsing, validation and export checks.
package env // import "code.sajari.com/env"

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// Var represents the state of a variable.
type Var struct {
	Name  string // name
	Usage string // help message
	Value Value  // value as set
}

// Value is the interface to the dynamic value stored in Var.
type Value interface {
	// String is a string representation of the stored value.
	String() string

	// Set assigns a new value to a stored value from a string
	// representation.
	Set(string) error
}

type stringValue string

func newStringValue(x string, p *string) *stringValue {
	*p = x
	return (*stringValue)(p)
}

func (v *stringValue) Set(x string) error {
	*v = stringValue(x)
	return nil
}

func (v *stringValue) String() string {
	return string(*v)
}

type intValue int

func newIntValue(x int, p *int) *intValue {
	*p = x
	return (*intValue)(p)
}

func (v *intValue) Set(x string) error {
	n, err := strconv.Atoi(x)
	*v = intValue(n)
	if err != nil {
		if ne, ok := err.(*strconv.NumError); ok {
			return errors.New("invalid integer " + strconv.Quote(ne.Num) + ": " + ne.Err.Error())
		}
	}
	return err
}

func (v *intValue) String() string {
	return strconv.Itoa(int(*v))
}

type int64Value int64

func newInt64Value(x int64, p *int64) *int64Value {
	*p = x
	return (*int64Value)(p)
}

func (v *int64Value) Set(x string) error {
	n, err := strconv.ParseInt(x, 10, 64)
	*v = int64Value(n)
	if err != nil {
		if ne, ok := err.(*strconv.NumError); ok {
			return errors.New("invalid 64-bit integer " + strconv.Quote(ne.Num) + ": " + ne.Err.Error())
		}
	}
	return err
}

func (v *int64Value) String() string {
	return strconv.FormatInt(int64(*v), 10)
}

type durationValue time.Duration

func newDurationValue(x time.Duration, p *time.Duration) *durationValue {
	*p = x
	return (*durationValue)(p)
}

func (v *durationValue) Set(x string) error {
	d, err := time.ParseDuration(x)
	*v = durationValue(d)
	return err
}

func (v *durationValue) String() string {
	return time.Duration(*v).String()
}

type boolValue bool

func newBoolValue(x bool, p *bool) *boolValue {
	*p = x
	return (*boolValue)(p)
}

func (v *boolValue) Set(x string) error {
	b, err := strconv.ParseBool(x)
	*v = boolValue(b)
	if err != nil {
		if ne, ok := err.(*strconv.NumError); ok {
			return errors.New("invalid bool " + strconv.Quote(ne.Num) + ": " + ne.Err.Error())
		}
	}
	return err
}

func (v *boolValue) String() string {
	return strconv.FormatBool(bool(*v))
}

type urlValue url.URL

func newURLValue(x url.URL, p *url.URL) *urlValue {
	*p = x
	return (*urlValue)(p)
}

func (v *urlValue) Set(x string) error {
	if x == "" {
		return errors.New("empty")
	}
	u, err := url.Parse(x)
	if err != nil {
		return err
	}
	*v = urlValue(*u)
	return err
}

func (v *urlValue) String() string {
	u := url.URL(*v)
	return u.String()
}

// NewVarSet creates a new variable set with given name.
//
// If name is non-empty, then all variables will have a strings.ToUpper(name)+"_"
// prefix.
func NewVarSet(name string) *VarSet {
	return &VarSet{
		name:   name,
		prefix: strings.Replace(strings.ToUpper(name), "-", "_", -1),
	}
}

// VarSet contains a set of variables.
type VarSet struct {
	name   string
	prefix string

	vars []*Var
}

// Var defines a variable with the specified name and usage string.
func (v *VarSet) Var(value Value, name, usage string) {
	var prefix string
	if v.prefix != "" {
		prefix = v.prefix + "_"
	}
	x := &Var{Value: value, Name: prefix + name, Usage: usage}
	v.vars = append(v.vars, x)
}

// Name is the name of the variable set.
func (v *VarSet) Name() string {
	return v.name
}

// Prefix applied to all variables when they are created.
func (v *VarSet) Prefix() string {
	return v.prefix
}

// Visit visits the variables in the order in which they were defined, calling fn for each.
func (v *VarSet) Visit(fn func(v *Var)) {
	for _, x := range v.vars {
		fn(x)
	}
}

// String defines a string variable with specified name, usage string and validation checks.
// The return value is the address of a string variable that stores the value of the variable.
func (v *VarSet) String(name string, usage string) *string {
	p := new(string)
	v.Var(newStringValue("", p), name, usage)
	return p
}

// StringRequired defines a required string variable with specified name and usage string.
// The return value is the address of a string variable that stores the value of the variable.
func (v *VarSet) StringRequired(name string, usage string) *string {
	p := new(string)
	v.Var(checkedValue{
		fn:    isNonEmpty,
		Value: newStringValue("", p),
	}, name, usage)
	return p
}

// Int defines an int variable with specified name, usage string and validation checks.
// The return value is the address of an int variable that stores the value of the variable.
func (v *VarSet) Int(name string, usage string) *int {
	p := new(int)
	v.Var(newIntValue(0, p), name, usage)
	return p
}

// Int64 defines an int64 variable with specified name, usage string and validation checks.
// The return value is the address of an int variable that stores the value of the variable.
func (v *VarSet) Int64(name string, usage string) *int64 {
	p := new(int64)
	v.Var(newInt64Value(0, p), name, usage)
	return p
}

// Bool defines a bool variable with specified name, usage string and validation checks.
// The return value is the address of a bool variable that stores the value of the variable.
func (v *VarSet) Bool(name string, usage string) *bool {
	p := new(bool)
	v.Var(newBoolValue(false, p), name, usage)
	return p
}

// Duration defines a time.Duration variable with specified name, usage string and validation checks.
// The return value is the address of a time.Duration variable that stores the value of the variable.
func (v *VarSet) Duration(name string, usage string) *time.Duration {
	p := new(time.Duration)
	v.Var(newDurationValue(time.Duration(0), p), name, usage)
	return p
}

// BindAddr defines a string variable with specified name, usage string validated as a
// bind address (host:port).
// The return value is the address of a string variable that stores the value of the variable.
func (v *VarSet) BindAddr(name, usage string) *string {
	p := new(string)
	v.Var(checkedValue{
		fn:    isBindAddr,
		Value: newStringValue("", p),
	}, name, usage)
	return p
}

// DialAddr defines a string variable with specified name, usage string validated as a
// dial address (host:port).
// The return value is the address of a string variable that stores the value of the variable.
func (v *VarSet) DialAddr(name, usage string) *string {
	p := new(string)
	v.Var(checkedValue{
		fn:    isDialAddr,
		Value: newStringValue("", p),
	}, name, usage)
	return p
}

// URL defines a string variable with specified name, usage string validated as a URL.
// The return value is the address of a URL variable that stores the value of the variable.
func (v *VarSet) URL(name, usage string) *url.URL {
	p := new(url.URL)
	v.Var(newURLValue(url.URL{}, p), name, usage)
	return p
}

// Path defines a string variable with specified name, usage string validated as a local path.
// The return value is the address of a string variable that stores the value of the variable.
func (v *VarSet) Path(name, usage string) *string {
	p := new(string)
	v.Var(checkedValue{
		fn:    isPath,
		Value: newStringValue("", p),
	}, name, usage)
	return p
}

// Errors is returned from Parse.
type Errors []error

// Error implements error.
func (me Errors) Error() string {
	n := 0
	msg := ""
	for _, e := range me {
		if e != nil {
			if n == 0 {
				msg = e.Error()
			}
			n++
		}
	}

	switch n {
	case 0:
		return "(0 errors)"
	case 1:
		return msg
	case 2:
		return fmt.Sprintf("%v (and 1 other error)", msg)
	}
	return fmt.Sprintf("%v (and %d other errors)", msg, n)
}

// Getter defines the Get method.
type Getter interface {
	// Get retrieves an evironment variable.
	Get(string) (string, bool)
}

type osLookup struct{}

func (osLookup) Get(x string) (string, bool) { return os.LookupEnv(x) }

// Parse parses variables from the environment provided by
// the Getter.
func (v *VarSet) Parse(g Getter) error {
	var errs []error

	for _, x := range v.vars {
		z, ok := g.Get(x.Name)
		if !ok {
			errs = append(errs, fmt.Errorf("missing env %v", x.Name))
			continue
		}

		if err := x.Value.Set(z); err != nil {
			errs = append(errs, fmt.Errorf("could not set env %v: %v", x.Name, err))
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return Errors(errs)
}

// CmdVar is the default variable set used for command-line based applications.
// The name of the variable set (and hence all variable prefixes) is given
// by CmdName.
var CmdVar = NewVarSet(CmdName())

// CmdName is used to create the default variable set name.
var CmdName = func() string {
	return path.Base(os.Args[0])
}

// String defines a string variable with specified name, usage string and validation checks.
// The return value is the address of a string variable that stores the value of the variable.
func String(name, usage string) *string {
	return CmdVar.String(name, usage)
}

// StringRequired defines a required string variable with specified name and usage string..
// The return value is the address of a string variable that stores the value of the variable.
func StringRequired(name, usage string) *string {
	return CmdVar.StringRequired(name, usage)
}

// BindAddr defines a string variable with specified name, usage string validated as a
// bind address (host:port).
// The return value is the address of a string variable that stores the value of the variable.
func BindAddr(name, usage string) *string {
	return CmdVar.BindAddr(name, usage)
}

// DialAddr defines a string variable with specified name, usage string validated as a
// dial address (host:port).
// The return value is the address of a string variable that stores the value of the variable.
func DialAddr(name, usage string) *string {
	return CmdVar.DialAddr(name, usage)
}

// URL defines a string variable with specified name, usage string validated as a URL.
// The return value is the address of a URL variable that stores the value of the variable.
func URL(name, usage string) *url.URL {
	return CmdVar.URL(name, usage)
}

// Path defines a string variable with specified name, usage string validated as a
// local path.
// The return value is the address of a string variable that stores the value of the variable.
func Path(name, usage string) *string {
	return CmdVar.Path(name, usage)
}

// Int defines an int variable with specified name and usage string.
// The return value is the address of an int variable that stores the value of the variable.
func Int(name string, usage string) *int {
	return CmdVar.Int(name, usage)
}

// Int64 defines an int64 variable with specified name and usage string.
// The return value is the address of an int variable that stores the value of the variable.
func Int64(name string, usage string) *int64 {
	return CmdVar.Int64(name, usage)
}

// Bool defines a bool variable with specified name and usage string.
// The return value is the address of a bool variable that stores the value of the variable.
func Bool(name string, usage string) *bool {
	return CmdVar.Bool(name, usage)
}

// Duration defines a time.Duration variable with specified name, usage string and validation checks.
// The return value is the address of a time.Duration variable that stores the value of the variable.
func Duration(name string, usage string) *time.Duration {
	return CmdVar.Duration(name, usage)
}

// Visit visits the variables in the order in which they were defined, calling fn for each.
func Visit(fn func(*Var)) {
	CmdVar.Visit(fn)
}

// Parse parses variables from the process environment.
func Parse() error {
	return CmdVar.Parse(osLookup{})
}
