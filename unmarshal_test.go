package main

import (
	"io/ioutil"
	"testing"
)

func TestLogicalLines(t *testing.T) {
	cases := []struct {
		in   []byte
		want [][]byte
	}{
		{[]byte("abcdef"), [][]byte{[]byte("abcdef")}},
		{[]byte("abc\\\ndef"), [][]byte{[]byte("abcdef")}},
		{[]byte("abc\ndef"), [][]byte{[]byte("abc"), []byte("def")}},
	}
	for _, c := range cases {
		got := logicalLines(c.in)
		for i, bs := range got {
			for j, b := range bs {
				if b != c.want[i][j] {
					t.Errorf("logicalLines(%q) == %q, want %q", c.in, got, c.want)
				}
			}
		}
	}
}

func TestTrimComment(t *testing.T) {
	cases := []struct {
		in   []byte
		want []byte
	}{
		{[]byte("foo#bar"), []byte("foo")},
		{[]byte("foobar"), []byte("foobar")},
		{[]byte("#foo#bar"), []byte("")},
	}
	for _, c := range cases {
		got := trimComment(c.in)
		for i, b := range got {
			if b != c.want[i] {
				t.Errorf("trimComment(%q) == %q, want %q", c.in, got, c.want)
			}
		}
	}
}

func TestParseFrom(t *testing.T) {
	cases := []struct {
		in   []byte
		want *From
	}{
		{[]byte("ubuntu"), &From{Image: "ubuntu"}},
		{[]byte(" ubuntu"), &From{Image: "ubuntu"}},
		{[]byte("ubuntu "), &From{Image: "ubuntu"}},
		{[]byte("ubuntu:14.04"), &From{Image: "ubuntu", Tag: "14.04"}},
		{[]byte(" ubuntu:14.04"), &From{Image: "ubuntu", Tag: "14.04"}},
		{[]byte("ubuntu:14.04 "), &From{Image: "ubuntu", Tag: "14.04"}},
		{[]byte("ubuntu@12345"), &From{Image: "ubuntu", Digest: "12345"}},
		{[]byte(" ubuntu@12345"), &From{Image: "ubuntu", Digest: "12345"}},
		{[]byte("ubuntu@12345 "), &From{Image: "ubuntu", Digest: "12345"}},
	}
	for _, c := range cases {
		got, _ := parseFrom(c.in)
		if *got != *c.want {
			t.Errorf("parseFrom(%q) == %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseMaintainer(t *testing.T) {
	cases := []struct {
		in   []byte
		want *Maintainer
	}{
		{[]byte("JohnDoe"), &Maintainer{Name: "JohnDoe"}},
		{[]byte("John Doe"), &Maintainer{Name: "John Doe"}},
		{[]byte(" John Doe"), &Maintainer{Name: "John Doe"}},
		{[]byte("John Doe "), &Maintainer{Name: "John Doe"}},
	}
	for _, c := range cases {
		got, _ := parseMaintainer(c.in)
		if *got != *c.want {
			t.Errorf("parseMaintainer(%q) == %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseRun(t *testing.T) {
	cases := []struct {
		in   []byte
		want *Run
	}{
		{[]byte("/bin/ls"), &Run{Tokens: []string{"/bin/sh", "-c", "/bin/ls"}}},
		{[]byte(" /bin/ls"), &Run{Tokens: []string{"/bin/sh", "-c", "/bin/ls"}}},
		{[]byte("/bin/ls "), &Run{Tokens: []string{"/bin/sh", "-c", "/bin/ls"}}},
		{[]byte("/bin/rm foo"), &Run{Tokens: []string{"/bin/sh", "-c", "/bin/rm", "foo"}}},
		{[]byte(" /bin/rm foo"), &Run{Tokens: []string{"/bin/sh", "-c", "/bin/rm", "foo"}}},
		{[]byte("/bin/rm foo "), &Run{Tokens: []string{"/bin/sh", "-c", "/bin/rm", "foo"}}},
		{[]byte("[\"/bin/ls\"]"), &Run{Tokens: []string{"/bin/ls"}}},
		{[]byte(" [ \"/bin/ls\" ] "), &Run{Tokens: []string{"/bin/ls"}}},
		{[]byte("[\"/bin/rm\",\"foo\"]"), &Run{Tokens: []string{"/bin/rm", "foo"}}},
		{[]byte(" [ \"/bin/rm\", \"foo\" ] "), &Run{Tokens: []string{"/bin/rm", "foo"}}},
	}
	for _, c := range cases {
		got, _ := parseRun(c.in)
		for i, tok := range got.Tokens {
			if tok != c.want.Tokens[i] {
				t.Errorf("parseRun(%q) == %v, want %v", c.in, got, c.want)
			}
		}
	}
}

func TestParseCmd(t *testing.T) {
	cases := []struct {
		in   []byte
		want *Cmd
	}{
		{[]byte("/bin/ls"), &Cmd{Tokens: []string{"/bin/sh", "-c", "/bin/ls"}}},
		{[]byte(" /bin/ls"), &Cmd{Tokens: []string{"/bin/sh", "-c", "/bin/ls"}}},
		{[]byte("/bin/ls "), &Cmd{Tokens: []string{"/bin/sh", "-c", "/bin/ls"}}},
		{[]byte("/bin/rm foo"), &Cmd{Tokens: []string{"/bin/sh", "-c", "/bin/rm", "foo"}}},
		{[]byte(" /bin/rm foo"), &Cmd{Tokens: []string{"/bin/sh", "-c", "/bin/rm", "foo"}}},
		{[]byte("/bin/rm foo "), &Cmd{Tokens: []string{"/bin/sh", "-c", "/bin/rm", "foo"}}},
		{[]byte("[\"/bin/ls\"]"), &Cmd{Tokens: []string{"/bin/ls"}}},
		{[]byte(" [ \"/bin/ls\" ] "), &Cmd{Tokens: []string{"/bin/ls"}}},
		{[]byte("[\"/bin/rm\",\"foo\"]"), &Cmd{Tokens: []string{"/bin/rm", "foo"}}},
		{[]byte(" [ \"/bin/rm\", \"foo\" ] "), &Cmd{Tokens: []string{"/bin/rm", "foo"}}},
	}
	for _, c := range cases {
		got, _ := parseCmd(c.in)
		for i, tok := range got.Tokens {
			if tok != c.want.Tokens[i] {
				t.Errorf("parseCmd(%q) == %v, want %v", c.in, got, c.want)
			}
		}
	}
}

func TestParseLabel(t *testing.T) {
	cases := []struct {
		in   []byte
		want *Label
	}{
		{[]byte("foo=bar"), &Label{Labels: map[string]string{"foo": "bar"}}},
		{[]byte(" foo=bar"), &Label{Labels: map[string]string{"foo": "bar"}}},
		{[]byte("foo=bar "), &Label{Labels: map[string]string{"foo": "bar"}}},
		{[]byte("foo=bar fizz=buzz"), &Label{Labels: map[string]string{"foo": "bar", "fizz": "buzz"}}},
		{[]byte(" foo=bar fizz=buzz"), &Label{Labels: map[string]string{"foo": "bar", "fizz": "buzz"}}},
		{[]byte("foo=bar fizz=buzz "), &Label{Labels: map[string]string{"foo": "bar", "fizz": "buzz"}}},
		{[]byte("foo=\"b a r\""), &Label{Labels: map[string]string{"foo": "b a r"}}},
		{[]byte("\"f o o\"=bar"), &Label{Labels: map[string]string{"f o o": "bar"}}},
		{[]byte("\"f o o\"=\"b a r\""), &Label{Labels: map[string]string{"f o o": "b a r"}}},
	}
	for _, c := range cases {
		got, _ := parseLabel(c.in)
		for k, v := range got.Labels {
			if v != c.want.Labels[k] {
				t.Errorf("parseLabel(%q) == %v, want %v", c.in, got, c.want)
			}
		}
	}
}

func TestParseExpose(t *testing.T) {
	cases := []struct {
		in   []byte
		want *Expose
	}{
		{[]byte("22"), &Expose{Ports: []int{22}}},
		{[]byte(" 22"), &Expose{Ports: []int{22}}},
		{[]byte("22 "), &Expose{Ports: []int{22}}},
		{[]byte("22 80"), &Expose{Ports: []int{22, 80}}},
		{[]byte(" 22 80"), &Expose{Ports: []int{22, 80}}},
		{[]byte("22 80 "), &Expose{Ports: []int{22, 80}}},
	}
	for _, c := range cases {
		got, _ := parseExpose(c.in)
		for i, p := range got.Ports {
			if p != c.want.Ports[i] {
				t.Errorf("parseExpose(%q) == %v, want %v", c.in, got, c.want)
			}
		}
	}
}

func TestParseEnv(t *testing.T) {
	cases := []struct {
		in   []byte
		want *Env
	}{
		{[]byte("FOO bar"), &Env{Variables: map[string]string{"FOO": "bar"}}},
		{[]byte(" FOO bar"), &Env{Variables: map[string]string{"FOO": "bar"}}},
		{[]byte("FOO bar "), &Env{Variables: map[string]string{"FOO": "bar"}}},
		{[]byte("FOO bar fizz buzz"), &Env{Variables: map[string]string{"FOO": "bar fizz buzz"}}},
		{[]byte("FOO=bar"), &Env{Variables: map[string]string{"FOO": "bar"}}},
		{[]byte(" FOO=bar"), &Env{Variables: map[string]string{"FOO": "bar"}}},
		{[]byte("FOO=bar "), &Env{Variables: map[string]string{"FOO": "bar"}}},
		{[]byte("FOO=bar FIZZ=buzz"), &Env{Variables: map[string]string{"FOO": "bar", "FIZZ": "buzz"}}},
		{[]byte("FOO=\"b a r\""), &Env{Variables: map[string]string{"FOO": "b a r"}}},
		{[]byte("FOO=bar FIZZ=\"b u z z\""), &Env{Variables: map[string]string{"FOO": "bar", "FIZZ": "b u z z"}}},
		{[]byte("FOO=\"b a r\" FIZZ=\"b u z z\""), &Env{Variables: map[string]string{"FOO": "b a r", "FIZZ": "b u z z"}}},
	}
	for _, c := range cases {
		got, _ := parseEnv(c.in)
		for k, v := range got.Variables {
			if c.want.Variables[k] != v {
				t.Errorf("parseEnv(%q) == %v, want %v", c.in, got, c.want)
			}
		}
	}
}

func TestParseAdd(t *testing.T) {
	cases := []struct {
		in   []byte
		want *Add
	}{
		{[]byte("/src /dest"), &Add{Sources: []string{"/src"}, Destination: "/dest"}},
		{[]byte(" /src /dest"), &Add{Sources: []string{"/src"}, Destination: "/dest"}},
		{[]byte("/src /dest "), &Add{Sources: []string{"/src"}, Destination: "/dest"}},
		{[]byte("[\"/src\",\"/dest\"]"), &Add{Sources: []string{"/src"}, Destination: "/dest"}},
		{[]byte(" [ \"/src\", \"/dest\" ] "), &Add{Sources: []string{"/src"}, Destination: "/dest"}},
		{[]byte("/s1 /s2 /dest"), &Add{Sources: []string{"/s1", "/s2"}, Destination: "/dest"}},
		{[]byte("[\"/s1\",\"/s2\",\"/dest\"]"), &Add{Sources: []string{"/s1", "/s2"}, Destination: "/dest"}},
	}
	for _, c := range cases {
		got, _ := parseAdd(c.in)
		for i, s := range got.Sources {
			if s != c.want.Sources[i] {
				t.Errorf("parseAdd(%q) == %v, want %v", c.in, got, c.want)
			}
		}
		if got.Destination != c.want.Destination {
			t.Errorf("parseAdd(%q) == %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseCopy(t *testing.T) {
	cases := []struct {
		in   []byte
		want *Copy
	}{
		{[]byte("/src /dest"), &Copy{Sources: []string{"/src"}, Destination: "/dest"}},
		{[]byte(" /src /dest"), &Copy{Sources: []string{"/src"}, Destination: "/dest"}},
		{[]byte("/src /dest "), &Copy{Sources: []string{"/src"}, Destination: "/dest"}},
		{[]byte("[\"/src\",\"/dest\"]"), &Copy{Sources: []string{"/src"}, Destination: "/dest"}},
		{[]byte(" [ \"/src\", \"/dest\" ] "), &Copy{Sources: []string{"/src"}, Destination: "/dest"}},
		{[]byte("/s1 /s2 /dest"), &Copy{Sources: []string{"/s1", "/s2"}, Destination: "/dest"}},
		{[]byte("[\"/s1\",\"/s2\",\"/dest\"]"), &Copy{Sources: []string{"/s1", "/s2"}, Destination: "/dest"}},
	}
	for _, c := range cases {
		got, _ := parseCopy(c.in)
		for i, s := range got.Sources {
			if s != c.want.Sources[i] {
				t.Errorf("parseCopy(%q) == %v, want %v", c.in, got, c.want)
			}
		}
		if got.Destination != c.want.Destination {
			t.Errorf("parseCopy(%q) == %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseEntrypoint(t *testing.T) {
	cases := []struct {
		in   []byte
		want *Entrypoint
	}{
		{[]byte("/bin/ls"), &Entrypoint{Tokens: []string{"/bin/sh", "-c", "/bin/ls"}}},
		{[]byte(" /bin/ls"), &Entrypoint{Tokens: []string{"/bin/sh", "-c", "/bin/ls"}}},
		{[]byte("/bin/ls "), &Entrypoint{Tokens: []string{"/bin/sh", "-c", "/bin/ls"}}},
		{[]byte("/bin/rm foo"), &Entrypoint{Tokens: []string{"/bin/sh", "-c", "/bin/rm", "foo"}}},
		{[]byte(" /bin/rm foo"), &Entrypoint{Tokens: []string{"/bin/sh", "-c", "/bin/rm", "foo"}}},
		{[]byte("/bin/rm foo "), &Entrypoint{Tokens: []string{"/bin/sh", "-c", "/bin/rm", "foo"}}},
		{[]byte("[\"/bin/ls\"]"), &Entrypoint{Tokens: []string{"/bin/ls"}}},
		{[]byte(" [ \"/bin/ls\" ] "), &Entrypoint{Tokens: []string{"/bin/ls"}}},
		{[]byte("[\"/bin/rm\",\"foo\"]"), &Entrypoint{Tokens: []string{"/bin/rm", "foo"}}},
		{[]byte(" [ \"/bin/rm\", \"foo\" ] "), &Entrypoint{Tokens: []string{"/bin/rm", "foo"}}},
	}
	for _, c := range cases {
		got, _ := parseEntrypoint(c.in)
		for i, tok := range got.Tokens {
			if tok != c.want.Tokens[i] {
				t.Errorf("parseEntrypoint(%q) == %v, want %v", c.in, got, c.want)
			}
		}
	}
}

func TestParseVolume(t *testing.T) {
	cases := []struct {
		in   []byte
		want *Volume
	}{
		{[]byte("/opt"), &Volume{Points: []string{"/opt"}}},
		{[]byte(" /opt"), &Volume{Points: []string{"/opt"}}},
		{[]byte("/opt "), &Volume{Points: []string{"/opt"}}},
		{[]byte("/opt /etc"), &Volume{Points: []string{"/opt", "/etc"}}},
		{[]byte(" /opt /etc"), &Volume{Points: []string{"/opt", "/etc"}}},
		{[]byte("/opt /etc "), &Volume{Points: []string{"/opt", "/etc"}}},
		{[]byte("[\"/opt\"]"), &Volume{Points: []string{"/opt"}}},
		{[]byte(" [ \"/opt\" ] "), &Volume{Points: []string{"/opt"}}},
		{[]byte("[\"/opt\",\"/etc\"]"), &Volume{Points: []string{"/opt", "/etc"}}},
		{[]byte(" [ \"/opt\", \"/etc\" ] "), &Volume{Points: []string{"/opt", "/etc"}}},
	}
	for _, c := range cases {
		got, _ := parseVolume(c.in)
		for i, p := range got.Points {
			if p != c.want.Points[i] {
				t.Errorf("parseVolume(%q) == %v, want %v", c.in, got, c.want)
			}
		}
	}
}

func TestParseUser(t *testing.T) {
	cases := []struct {
		in   []byte
		want *User
	}{
		{[]byte("root"), &User{Name: "root"}},
		{[]byte(" root"), &User{Name: "root"}},
		{[]byte("root "), &User{Name: "root"}},
	}
	for _, c := range cases {
		got, _ := parseUser(c.in)
		if *got != *c.want {
			t.Errorf("parseUser(%q) == %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseWorkdir(t *testing.T) {
	cases := []struct {
		in   []byte
		want *Workdir
	}{
		{[]byte("/home/foo"), &Workdir{Path: "/home/foo"}},
		{[]byte(" /home/foo"), &Workdir{Path: "/home/foo"}},
		{[]byte("/home/foo "), &Workdir{Path: "/home/foo"}},
	}
	for _, c := range cases {
		got, _ := parseWorkdir(c.in)
		if *got != *c.want {
			t.Errorf("parseWorkdir(%q) == %v, want %v", c.in, got, c.want)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	in, _ := ioutil.ReadFile("stub/DockerfileTest")
	_, err := unmarshal(in)
	if err != nil {
		t.Errorf("unmarshal(%q) returns %s", in, err)
	}
	// TODO test the contents
}
