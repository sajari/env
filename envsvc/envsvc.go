// Package envsvc provides convenience methods for using env with
// services.
//
// It exposes these variables via HTTP at /debug/env in JSON format.
package envsvc

import (
	"flag"
	"fmt"
	"io"
	"os"

	"code.sajari.com/env"
)

// Parse registers command line arguments -env-check, -env-dump, -env-dump-yaml
// and calls env.Parse().  Any errors are written to stderr and then os.Exit(1) is called.
//
// Registered flags:
// -env-dump: skips parsing step and writes each env.Var to stderr, calls os.Exit(0) when done.
// -env-dump-yaml: skips parsing steps and write each env.Var to stderr in YAML format, calls
// os.Exit(0) when done.
// -env-check: calls os.Exit(0) if env.Parse() succeeds without error.
func Parse() {
	envCheck := flag.Bool("env-check", false, "check env variables")
	envDump := flag.Bool("env-dump", false, "dump env variables")
	envDumpYAML := flag.Bool("env-dump-yaml", false, "dump env variables in YAML format")
	envDumpJSON := flag.Bool("env-dump-json", false, "dump env variables in JSON format")

	flag.Parse()

	var outWriter io.Writer = os.Stderr

	if *envDumpJSON {
		fmt.Fprintf(outWriter, "{\n")
		first := true
		env.Visit(func(v *env.Var) {
			if !first {
				fmt.Fprintf(outWriter, ",\n")
			}
			first = false
			fmt.Fprintf(outWriter, "    %q: %q", v.Name, os.Getenv(v.Name))
		})
		fmt.Fprintf(outWriter, "\n}\n")
		os.Exit(0)
	}

	if *envDumpYAML {
		env.Visit(func(v *env.Var) {
			fmt.Fprintf(outWriter, "- name: %v\n  value: %q\n", v.Name, os.Getenv(v.Name))
		})
		os.Exit(0)
	}

	if *envDump {
		first := true
		env.Visit(func(v *env.Var) {
			if !first {
				fmt.Fprintf(outWriter, "\n")
			}
			first = false
			fmt.Fprintf(outWriter, "# %v\nexport %v=%q\n", v.Usage, v.Name, os.Getenv(v.Name))
		})
		os.Exit(0)
	}

	if err := env.Parse(); err != nil {
		if es, ok := err.(env.Errors); ok {
			for _, e := range es {
				fmt.Fprintln(outWriter, e)
			}
		} else {
			fmt.Fprintln(outWriter, err)
		}
		os.Exit(1)
	}

	if *envCheck {
		os.Exit(0)
	}
}
