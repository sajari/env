package envsvc

import (
	"fmt"
	"io"
	"net/http"

	"code.sajari.com/env"
)

// Handler returns the env HTTP Handler.
//
// This is only needed to install the handler in a non-standard location.
func Handler() http.Handler {
	return http.HandlerFunc(envHandler)
}

func envHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	vs := r.URL.Query()
	if _, ok := vs["short"]; ok {
		shortHandler(w)
		return
	}
	detailHandler(w)
}

func shortHandler(w io.Writer) {
	fmt.Fprintf(w, "{\n")
	first := true
	env.Visit(func(v *env.Var) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "    %q: %q", v.Name, v.Value.String())
	})
	fmt.Fprintf(w, "\n}\n")
}

func detailHandler(w io.Writer) {
	fmt.Fprintf(w, "{\n    %q: [\n", "env")
	first := true
	env.Visit(func(v *env.Var) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "        {\n")
		fmt.Fprintf(w, "            %q: %q,\n", "name", v.Name)
		fmt.Fprintf(w, "            %q: %q,\n", "usage", v.Usage)
		fmt.Fprintf(w, "            %q: %q\n", "value", v.Value.String())
		fmt.Fprintf(w, "        }")
	})
	fmt.Fprintf(w, "\n    ]\n}\n")
}
