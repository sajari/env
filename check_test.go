package env_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"code.sajari.com/env"
)

func TestInt(t *testing.T) {
	tests := []struct {
		in      string
		out     int
		wantErr bool
	}{
		// Valid
		{"1234", 1234, false},
		{"0", 0, false},

		// Invalid
		{"", 0, true},
		{" ", 0, true},
		{"a", 0, true},
		{"12.3", 0, true},
	}

	env.ResetForTesting()
	intValue := env.Int("INT", "int test")
	name := "TEST_INT"

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			os.Setenv(name, tt.in)

			if err := env.Parse(); (err != nil) != tt.wantErr {
				t.Errorf("env.Parse() = %v, wantErr %v", err, tt.wantErr)
			}

			if *intValue != tt.out {
				t.Errorf(" = %d, expected %d", *intValue, tt.out)
			}
		})
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		in      string
		out     bool
		wantErr bool
	}{
		// Valid, true
		{"1", true, false},
		{"T", true, false},
		{"TRUE", true, false},
		{"true", true, false},

		// Valid, false
		{"0", false, false},
		{"F", false, false},
		{"FALSE", false, false},
		{"false", false, false},

		// Invalid
		{"", false, true},
		{" ", false, true},
		{"2", false, true},
		{"a", false, true},
		{"12.3", false, true},
	}

	env.ResetForTesting()
	boolValue := env.Bool("BOOL", "int test")
	name := "TEST_BOOL"

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			os.Setenv(name, tt.in)

			if err := env.Parse(); (err != nil) != tt.wantErr {
				t.Errorf("env.Parse() = %v, wantErr %v", err, tt.wantErr)
			}

			if *boolValue != tt.out {
				t.Errorf(" = %v, expected %v", *boolValue, tt.out)
			}
		})
	}
}

func TestBindAddr(t *testing.T) {
	tests := []struct {
		in      string
		wantErr bool
	}{
		// Valid
		{":1234", false},
		{"localhost:1234", false},
		{"192.168.0.1:1234", false},

		// Invalid
		{"", true},
		{":", true},
		{"192.168.0.1:", true},
		{"localhost:", true},
	}

	env.ResetForTesting()
	prefix := env.CmdVar.Name()
	_ = env.BindAddr("BIND", "bind address test")
	name := strings.ToUpper(prefix) + "_BIND"

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			os.Setenv(name, tt.in)

			if err := env.Parse(); (err != nil) != tt.wantErr {
				t.Errorf("env.Parse() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsDialAddr(t *testing.T) {
	tests := []struct {
		in      string
		wantErr bool
	}{
		// Valid
		{"localhost:1234", false},
		{"192.168.0.1:1234", false},
		{"sajari.com:1234", false},

		// Invalid
		{"", true},
		{":", true},
		{":1234", true},
		{"192.168.0.1:", true},
		{"localhost:", true},
	}

	env.ResetForTesting()
	prefix := env.CmdVar.Name()
	_ = env.DialAddr("ADDR", "dial address test")
	name := strings.ToUpper(prefix) + "_ADDR"

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			os.Setenv(name, tt.in)

			if err := env.Parse(); (err != nil) != tt.wantErr {
				t.Errorf("env.Parse() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsPath(t *testing.T) {
	env.ResetForTesting()

	tmpFile, err := ioutil.TempFile("", "IsPath")
	if err != nil {
		t.Fatalf("could not create temporary file: %v", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	prefix := env.CmdVar.Name()
	_ = env.Path("PATH", "path test")
	name := strings.ToUpper(prefix) + "_PATH"

	os.Setenv(name, tmpFile.Name())
	if err := env.Parse(); err != nil {
		t.Errorf("env.Parse() = %v, expected nil error", err)
	}

	os.Setenv(name, "filedoesnotexist.txt")
	if err := env.Parse(); err == nil {
		t.Error("env.Parse() should return en error")
	}
}

func TestMissingGetter(t *testing.T) {
	tg := testGetter{}

	vs := env.NewVarSet("")
	vs.Bool("MISSING", "missing test")

	if err := vs.Parse(tg); err == nil {
		t.Errorf("expected error for missing var")
	}
}
