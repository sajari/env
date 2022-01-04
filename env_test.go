package env_test

import (
	"testing"

	"code.sajari.com/env"
)

type testGetter map[string]string

func (g testGetter) Get(x string) (string, bool) {
	v, ok := g[x]
	return v, ok
}

func TestAll(t *testing.T) {
	env.ResetForTesting()

	env.Bool("BOOL", "bool test")
	env.Int("INT", "int test")
	env.BindAddr("LISTEN", "bindaddr test")
	env.DialAddr("ADDR", "dialaddr test")
	env.URL("URL", "URL test")
	env.String("STRING", "string test")
	env.Duration("TIMEOUT", "timeout test")
	env.Float32("FLOAT32", "float32 test")
	env.Float64("FLOAT64", "float64 test")

	tg := testGetter{
		"TEST_BOOL":    "true",
		"TEST_INT":     "1",
		"TEST_LISTEN":  ":1234",
		"TEST_ADDR":    "localhost:1234",
		"TEST_URL":     "http://localhost:1234/api",
		"TEST_STRING":  "name",
		"TEST_TIMEOUT": "1m1s",
		"TEST_FLOAT32": "1.23",
		"TEST_FLOAT64": "1.24",
	}

	if err := env.CmdVar.Parse(tg); err != nil {
		t.Errorf("unexpected error from Parse: %v", err)
	}

	checked := make(map[string]bool, len(tg))
	env.Visit(func(v *env.Var) {
		if _, ok := checked[v.Name]; ok {
			t.Errorf("already seen var: %v", v.Name)
		}
		checked[v.Name] = true

		if s := v.Value.String(); s != tg[v.Name] {
			t.Errorf("v.Value.String() = %q, expected %q", s, tg[v.Name])
		}
	})
}
