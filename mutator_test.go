package main

import (
	"testing"
)

func initialize() Dockerfile {
	return Dockerfile{
		Statements: []Statement{
			&From{Image: "ubuntu", Tag: "14.04"},
			&Maintainer{Name: "John Doe"},
			&Run{Tokens: []string{"/bin/rm", "foo"}},
			&Cmd{Tokens: []string{"/usr/rm", "foo"}},
			&Label{Labels: map[string]string{"foo": "bar", "fizz": "buzz"}},
			&Expose{Ports: []int{22, 80}},
			&Env{Variables: map[string]string{"FOO": "bar", "FIZZ": "buzz"}},
			&Add{Sources: []string{"/init", "/config"}, Destination: "/opt"},
			&Copy{Sources: []string{"/init", "/config"}, Destination: "/opt"},
			&Entrypoint{Tokens: []string{"/usr/rm", "foo"}},
			&Volume{Points: []string{"/etc", "/opt"}},
			&User{Name: "root"},
			&Workdir{Path: "/home/foo"},
			&Onbuild{Statement: Run{Tokens: []string{"/bin/rm", "foo"}}},
		},
	}
}

func TestReplaceEnv(t *testing.T) {
	cases := []struct {
		key   string
		value string
		want  map[string]string
	}{
		{"FOO", "newbar", map[string]string{"FOO": "newbar", "FIZZ": "buzz"}},
		{"UNKNOWN", "dummy", map[string]string{"FOO": "bar", "FIZZ": "buzz"}},
	}
	in := initialize()
	idx := 6
	for _, c := range cases {
		df := replaceEnv(in, c.key, c.value)
		env, ok := df.Statements[idx].(Env)
		if !ok {
			t.Errorf("the %dth statement is not an instance of Env", idx)
		}
		for k, v := range env.Variables {
			if v != c.want[k] {
				t.Errorf("the value of %s is %s, want %s", k, v, c.want[k])
			}
		}
	}
}
