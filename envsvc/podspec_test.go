package envsvc_test

import (
	"strings"
	"testing"

	"code.sajari.com/env/envsvc"
)

const singleInput = `apiVersion: v1
kind: Pod
metadata:
  name: podname
  labels:
	 name: podname
spec:
  containers:
	 - name: one
		image: one
		ports:
		  - name: grpc
			 protocol: TCP
			 containerPort: 50551
		resources:
		  requests:
			 cpu: 100m
			 memory: 50Mi
		  limits:
			 cpu: "1"
			 memory: 50Mi
		env:
		  - name: ONE_BOOLEAN
			 value: "false"
		  - name: ONE_STRING
			 value: "string value"
		  - name: ONE_INTEGER
			 value: "1"`

func TestPodENVLookup(t *testing.T) {
	t.Run("single_container", func(t *testing.T) {
		x, err := envsvc.PodENVLookup(strings.NewReader(singleInput), "")
		must(t, err)

		expectValue(t, x, "ONE_BOOLEAN", "false")
		expectValue(t, x, "ONE_STRING", "string value")
		expectValue(t, x, "ONE_INTEGER", "1")

		// Using the name should be fine (for 1 container also).
		x, err = envsvc.PodENVLookup(strings.NewReader(singleInput), "one")
		must(t, err)

		expectValue(t, x, "ONE_BOOLEAN", "false")
		expectValue(t, x, "ONE_STRING", "string value")
		expectValue(t, x, "ONE_INTEGER", "1")
	})

	t.Run("multiple_containers", func(t *testing.T) {
		input := `apiVersion: v1
kind: Pod
metadata:
	name: podname
	labels:
		name: podname
spec:
	containers:
		- name: one
		image: one
		ports:
			- name: grpc
				protocol: TCP
				containerPort: 50551
		resources:
			requests:
				cpu: 100m
				memory: 50Mi
			limits:
				cpu: "1"
				memory: 50Mi
		env:
			- name: ONE_BOOLEAN
				value: "false"
			- name: ONE_STRING
				value: "string value"
			- name: ONE_INTEGER
				value: "1"
		- name: two
		image: two
		ports:
			- name: grpc
				protocol: TCP
				containerPort: 50552
		env:
			- name: TWO_BOOLEAN
				value: "true"
			- name: TWO_STRING
				value: "string value two"
			- name: TWO_INTEGER
				value: "2"`

		x, err := envsvc.PodENVLookup(strings.NewReader(input), "two")
		must(t, err)

		expectValue(t, x, "TWO_BOOLEAN", "true")
		expectValue(t, x, "TWO_STRING", "string value two")
		expectValue(t, x, "TWO_INTEGER", "2")
	})

}

func expectValue(t *testing.T, g envsvc.Getter, name, value string) {
	t.Helper()

	x, ok := g.Get(name)
	if !ok {
		t.Errorf("expected value for %q", name)
		return
	}
	if x != value {
		t.Errorf("g.Get(%q) = %q, want %q", name, x, value)
	}
}

func must(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
