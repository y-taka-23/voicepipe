package main

import (
	"testing"
)

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
		got, _ := ParseFrom(c.in)
		if *got != *c.want {
			t.Errorf("ParseFrom(%q) == %v, want %v", c.in, got, c.want)
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
		got, _ := ParseMaintainer(c.in)
		if *got != *c.want {
			t.Errorf("ParseMaintainer(%q) == %v, want %v", c.in, got, c.want)
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
		got, _ := ParseExpose(c.in)
		for i, p := range got.Ports {
			if p != c.want.Ports[i] {
				t.Errorf("ParseExpose(%q) == %v, want %v", c.in, got, c.want)
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
		got, _ := ParseVolume(c.in)
		for i, p := range got.Points {
			if p != c.want.Points[i] {
				t.Errorf("ParseVolume(%q) == %v, want %v", c.in, got, c.want)
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
		got, _ := ParseUser(c.in)
		if *got != *c.want {
			t.Errorf("ParseUser(%q) == %v, want %v", c.in, got, c.want)
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
		got, _ := ParseWorkdir(c.in)
		if *got != *c.want {
			t.Errorf("ParseWorkdir(%q) == %v, want %v", c.in, got, c.want)
		}
	}
}
