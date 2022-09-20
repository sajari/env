package envsvc

import (
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type pod struct {
	Spec struct {
		Containers []struct {
			Name string `yaml:"name"`
			Env  []struct {
				Name  string `yaml:"name"`
				Value string `yaml:"value"`
			} `yaml:"env"`
		} `yaml:"containers"`
	} `yaml:"spec"`
}

type Getter interface {
	Get(string) (string, bool)
}

// PodENVLookup creates a lookup Getter from the given Pod YAML.
// Set name if there is more than one container in the pod.
func PodENVLookup(r io.Reader, name string) (Getter, error) {
	dec := yaml.NewDecoder(r)
	p := &pod{}
	if err := dec.Decode(p); err != nil {
		return nil, err
	}

	if len(p.Spec.Containers) == 0 {
		return nil, errors.New("no containers in pod spec")
	}

	if name == "" {
		if len(p.Spec.Containers) != 1 {
			return nil, fmt.Errorf("name empty but %d containers in pod spec, must set name", len(p.Spec.Containers))
		}
		return lookupFromContainerEnv(p, 0), nil
	}

	for i, c := range p.Spec.Containers {
		if c.Name == name {
			return lookupFromContainerEnv(p, i), nil
		}
	}
	return nil, fmt.Errorf("no container for name %q", name)
}

func lookupFromContainerEnv(p *pod, n int) *lookup {
	m := make(map[string]string)
	for _, x := range p.Spec.Containers[n].Env {
		m[x.Name] = x.Value
	}
	return &lookup{
		g: osLookup{},
		m: m,
	}
}

type lookup struct {
	g Getter
	m map[string]string
}

func (l lookup) Get(x string) (string, bool) {
	z, ok := l.g.Get(x)
	if ok {
		return z, ok
	}

	z, ok = l.m[x]
	return z, ok
}

type osLookup struct{}

func (osLookup) Get(x string) (string, bool) { return os.LookupEnv(x) }

func get(g Getter, x string) string {
	z, _ := g.Get(x)
	return z
}
