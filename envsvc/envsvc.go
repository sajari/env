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

// Parse is equivalent to ParseWithExitFn(os.Exit).
func Parse() {
	ParseWithExitFn(os.Exit)
}

// ParseWithExitFn registers command line arguments -env-check, -env-dump, -env-dump-yaml
// and calls env.Parse().  Any errors are written to stderr and then exitFn(1) is called.
//
// Registered flags:
// -env-dump: skips parsing step and writes each env.Var to stderr, calls exitFn(0) when done.
// -env-dump-yaml: skips parsing steps and write each env.Var to stderr in YAML format, calls
// os.Exit(0) when done.
// -env-check: calls exitFn(0) if env.Parse() succeeds without error.
func ParseWithExitFn(exitFn func(int)) {
	envCheck := flag.Bool("env-check", false, "check env variables")
	envDump := flag.Bool("env-dump", false, "dump env variables")
	envDumpYAML := flag.Bool("env-dump-yaml", false, "dump env variables in YAML format")
	envDumpJSON := flag.Bool("env-dump-json", false, "dump env variables in JSON format")
	envDumpCUE := flag.Bool("env-dump-cue", false, "dump env variables as CUE schema")

	envPodYAML := flag.String("env-pod-spec", "", "path to pod YAML to read env from")
	envPodYAMLContainerName := flag.String("env-pod-spec-container-name", "", "extract env from container `name` in the pod spec (required if more than one container is in the pod)")

	flag.Parse()

	var g Getter = osLookup{}

	if *envPodYAML != "" {
		f, err := os.Open(*envPodYAML)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not open pod spec %q: %v", *envPodYAML, err)
		}
		g, err = PodENVLookup(f, *envPodYAMLContainerName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read pod spec: %v", err)
		}
	}

	var outWriter io.Writer = os.Stderr

	if *envDumpJSON {
		fmt.Fprintf(outWriter, "{\n")
		first := true
		env.Visit(func(v *env.Var) {
			if !first {
				fmt.Fprintf(outWriter, ",\n")
			}
			first = false
			fmt.Fprintf(outWriter, "    %q: %q", v.Name, get(g, v.Name))
		})
		fmt.Fprintf(outWriter, "\n}\n")
		exitFn(0)
	}

	if *envDumpYAML {
		env.Visit(func(v *env.Var) {
			fmt.Fprintf(outWriter, "- name: %v\n  value: %q\n", v.Name, get(g, v.Name))
		})
		exitFn(0)
	}

	if *envDumpCUE {
		fmt.Fprintf(outWriter, "package %s\n\n", env.CmdName())
		fmt.Fprintf(outWriter, "#Env: [string]: string")
		env.Visit(func(v *env.Var) {
			// Insert newlines between fields to avoid cue fmt issues
			fmt.Fprintf(outWriter, "\n\n#Env: \"%v\": string", v.Name)
		})
		fmt.Fprintln(outWriter, "")
		exitFn(0)
	}

	if *envDump {
		first := true
		env.Visit(func(v *env.Var) {
			if !first {
				fmt.Fprintf(outWriter, "\n")
			}
			first = false
			fmt.Fprintf(outWriter, "# %v\nexport %v=%q\n", v.Usage, v.Name, get(g, v.Name))
		})
		exitFn(0)
	}

	if err := env.Parse(); err != nil {
		if es, ok := err.(env.Errors); ok {
			for _, e := range es {
				fmt.Fprintln(outWriter, e)
			}
		} else {
			fmt.Fprintln(outWriter, err)
		}
		exitFn(1)
	}

	if *envCheck {
		exitFn(0)
	}
}
