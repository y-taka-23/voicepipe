package main

import (
	"testing"
)

func initialize() Dockerfile {
	return Dockerfile{
		Statements: []Statement{
			From{Image: "ubuntu", Tag: "14.04"},
			Maintainer{Name: "John Doe"},
			Run{Tokens: []string{"/bin/rm", "foo"}},
			Cmd{Tokens: []string{"/usr/rm", "foo"}},
			Label{Labels: map[string]string{"foo": "bar", "fizz": "buzz"}},
			Expose{Ports: []int{22, 80}},
			Env{Variables: map[string]string{"FOO": "bar", "FIZZ": "buzz"}},
			Add{Sources: []string{"/init", "/config"}, Destination: "/opt"},
			Copy{Sources: []string{"/init", "/config"}, Destination: "/opt"},
			Entrypoint{Tokens: []string{"/usr/rm", "foo"}},
			Volume{Points: []string{"/etc", "/opt"}},
			User{Name: "root"},
			Workdir{Path: "/home/foo"},
			Onbuild{Statement: Run{Tokens: []string{"/bin/rm", "foo"}}},
		},
	}
}

func TestReplaceEnv(t *testing.T) {
	cases := []struct {
		key   string
		value string
	}{
		{"FOO", "xxx"},
		{"HOGE", "xxx"},
	}
	for _, c := range cases {
		df := ReplaceEnv(initialize(), c.key, c.value)
		for _, st := range df.Statements {
			if x, ok := st.(*Env); ok {
				if x.Variables[c.key] != c.value {
					t.Errorf("value of %q is %q, want %q ", c.key, x.Variables[c.key], c.value)
				}
			}
		}
	}
}
