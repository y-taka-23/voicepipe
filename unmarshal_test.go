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
		if got.Image != c.want.Image {
			t.Errorf(
				"ParseFrom(%q).Image == %q, want %q",
				c.in,
				got.Image,
				c.want.Image,
			)
		}
		if got.Tag != c.want.Tag {
			t.Errorf(
				"ParseFrom(%q).Tag == %q, want %q",
				c.in,
				got.Tag,
				c.want.Tag,
			)
		}
		if got.Digest != c.want.Digest {
			t.Errorf(
				"ParseFrom(%q).Digest == %q, want %q",
				c.in,
				got.Digest,
				c.want.Digest,
			)
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
		if got.Name != c.want.Name {
			t.Errorf(
				"ParseMaintainer(%q).Name == %q, want %q",
				c.in,
				got.Name,
				c.want.Name,
			)
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
		if got.Name != c.want.Name {
			t.Errorf(
				"ParseUser(%q).Name == %q, want %q",
				c.in,
				got.Name,
				c.want.Name,
			)
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
		if got.Path != c.want.Path {
			t.Errorf(
				"ParseWorkdir(%q).Path == %q, want %q",
				c.in,
				got.Path,
				c.want.Path,
			)
		}
	}
}
