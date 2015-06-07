package main

import (
	"strings"
	"testing"
)

func TestFromString(t *testing.T) {
	cases := []struct {
		in   From
		want string
	}{
		{From{Image: "ubuntu"}, "FROM ubuntu"},
		{From{Image: "ubuntu", Tag: "14.04"}, "FROM ubuntu:14.04"},
		{From{Image: "ubuntu"}, "FROM ubuntu"},
	}
	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("%v.String() == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestMaintainerString(t *testing.T) {
	in := Maintainer{Name: "John Doe"}
	want := "MAINTAINER John Doe"
	got := in.String()
	if got != want {
		t.Errorf("%v.String() == %q, want %q", in, got, want)
	}
}

func TestRunString(t *testing.T) {
	in := Run{Tokens: []string{"/bin/rm", "foo"}}
	want := "RUN [\"/bin/rm\", \"foo\"]"
	got := in.String()
	if got != want {
		t.Errorf("%v.String() == %q, want %q", in, got, want)
	}
}

func TestCmdString(t *testing.T) {
	in := Cmd{Tokens: []string{"/bin/rm", "foo"}}
	want := "CMD [\"/bin/rm\", \"foo\"]"
	got := in.String()
	if got != want {
		t.Errorf("%v.String() == %q, want %q", in, got, want)
	}
}

func TestLabelString(t *testing.T) {
	in := Label{Labels: map[string]string{"foo": "bar", "fizz": "buzz"}}
	wants := []string{"\"foo\"=\"bar\"", "\"fizz\"=\"buzz\""}
	got := in.String()
	for _, w := range wants {
		if !strings.Contains(got, w) {
			t.Errorf("%v.String() should contain %q, but doesn't", in, w)
		}
	}
	if strings.Count(got, "=") != len(wants) {
		t.Errorf("%v.String() should contain exact %d key-value pairs, but doesn't", in, len(wants))
	}
}

func TestExposeString(t *testing.T) {
	in := Expose{Ports: []int{22, 80}}
	want := "EXPOSE 22 80"
	got := in.String()
	if got != want {
		t.Errorf("%v.String() == %q, want %q", in, got, want)
	}
}

func TestEnvString(t *testing.T) {
	in := Env{Variables: map[string]string{"FOO": "bar", "FIZZ": "buzz"}}
	wants := []string{"FOO=\"bar\"", "FIZZ=\"buzz\""}
	got := in.String()
	for _, w := range wants {
		if !strings.Contains(got, w) {
			t.Errorf("%v.String() should contain %q, but doesn't", in, w)
		}
	}
	if strings.Count(got, "=") != len(wants) {
		t.Errorf("%v.String() should contain exact %d key-value pairs, but doesn't", in, len(wants))
	}
}

func TestAddString(t *testing.T) {
	in := Add{Sources: []string{"/s1", "/s2"}, Destination: "/dest"}
	want := "ADD [\"/s1\", \"/s2\", \"/dest\"]"
	got := in.String()
	if got != want {
		t.Errorf("%v.String() == %q, want %q", in, got, want)
	}
}

func TestCopyString(t *testing.T) {
	in := Copy{Sources: []string{"/s1", "/s2"}, Destination: "/dest"}
	want := "COPY [\"/s1\", \"/s2\", \"/dest\"]"
	got := in.String()
	if got != want {
		t.Errorf("%v.String() == %q, want %q", in, got, want)
	}
}

func TestEntrypoinString(t *testing.T) {
	in := Entrypoint{Tokens: []string{"/bin/rm", "foo"}}
	want := "ENTRYPOINT [\"/bin/rm\", \"foo\"]"
	got := in.String()
	if got != want {
		t.Errorf("%v.String() == %q, want %q", in, got, want)
	}
}

func TestVolumeString(t *testing.T) {
	in := Volume{Points: []string{"/opt", "/etc"}}
	want := "VOLUME /opt /etc"
	got := in.String()
	if got != want {
		t.Errorf("%v.String() == %q, want %q", in, got, want)
	}
}

func TestUserString(t *testing.T) {
	in := User{Name: "root"}
	want := "USER root"
	got := in.String()
	if got != want {
		t.Errorf("%v.String() == %q, want %q", in, got, want)
	}
}

func TestWorkdirString(t *testing.T) {
	in := Workdir{Path: "/home/foo"}
	want := "WORKDIR /home/foo"
	got := in.String()
	if got != want {
		t.Errorf("%v.String() == %q, want %q", in, got, want)
	}
}
