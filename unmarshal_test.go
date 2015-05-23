package main

import (
	"testing"
)

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
