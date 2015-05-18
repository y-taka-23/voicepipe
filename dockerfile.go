package main

import (
	"fmt"
	"strings"
)

type Statement interface {
	statement()
}

type Dockerfile struct {
	Statements []Statement
}

type From struct {
	Image  string
	Tag    string
	Digest string
}

type Maintainer struct {
	Name string
}

type Run struct {
	Tokens []string
}

type Cmd struct {
	Tokens []string
}

type Label struct {
	Labels map[string]string
}

type Expose struct {
	Ports []int
}

type Env struct {
	Variables map[string]string
}

type Add struct {
	Sources     []string
	Destination string
}

type Copy struct {
	Sources     []string
	Destination string
}

type Entrypoint struct {
	Tokens []string
}

type Volume struct {
	Points []string
}

type User struct {
	Name string
}

type Workdir struct {
	Path string
}

type Onbuild struct {
	Statement Statement
}

func (x From) statement()       {}
func (x Maintainer) statement() {}
func (x Run) statement()        {}
func (x Cmd) statement()        {}
func (x Label) statement()      {}
func (x Expose) statement()     {}
func (x Env) statement()        {}
func (x Add) statement()        {}
func (x Copy) statement()       {}
func (x Entrypoint) statement() {}
func (x Volume) statement()     {}
func (x User) statement()       {}
func (x Workdir) statement()    {}
func (x Onbuild) statement()    {}

func (x From) String() string {
	if len(x.Digest) != 0 {
		return fmt.Sprintf("FROM %s@%s", x.Image, x.Digest)
	} else if len(x.Tag) != 0 {
		return fmt.Sprintf("FROM %s:%s", x.Image, x.Tag)
	} else {
		return fmt.Sprintf("FROM %s", x.Image)
	}
}

func (x Maintainer) String() string {
	return fmt.Sprintf("MAINTAINER %s", x.Name)
}

func (x Run) String() string {
	return fmt.Sprintf("RUN [\"%s\"]", strings.Join(x.Tokens, "\", \""))
}

func (x Cmd) String() string {
	return fmt.Sprintf("CMD [\"%s\"]", strings.Join(x.Tokens, "\", \""))
}

func (x Label) String() string {
	res := make([]byte, 0)
	res = append(res, "LABEL"...)
	for k, v := range x.Labels {
		res = append(res, fmt.Sprintf(" %s=%s", k, v)...)
	}
	return string(res)
}

func (x Expose) String() string {
	res := make([]byte, 0)
	res = append(res, "EXPOSE"...)
	for _, p := range x.Ports {
		res = append(res, fmt.Sprintf(" %d", p)...)
	}
	return string(res)
}

func (x Env) String() string {
	res := make([]byte, 0)
	res = append(res, "ENV"...)
	for k, v := range x.Variables {
		res = append(res, fmt.Sprintf(" %s=%s", k, v)...)
	}
	return string(res)
}

func (x Add) String() string {
	return fmt.Sprintf(
		"ADD [\"%s\", \"%s\"]",
		strings.Join(x.Sources, "\", \""),
		x.Destination,
	)
}

func (x Copy) String() string {
	return fmt.Sprintf(
		"COPY [\"%s\", \"%s\"]",
		strings.Join(x.Sources, "\", \""),
		x.Destination,
	)
}

func (x Entrypoint) String() string {
	return fmt.Sprintf("ENTRYPOINT [\"%s\"]", strings.Join(x.Tokens, "\", \""))
}

func (x Volume) String() string {
	return fmt.Sprintf("VOLUME %s", strings.Join(x.Points, " "))
}

func (x User) String() string {
	return fmt.Sprintf("USER %s", x.Name)
}

func (x Workdir) String() string {
	return fmt.Sprintf("WORKDIR %s", x.Path)
}

func (x Onbuild) String() string {
	return fmt.Sprintf("ONBUILD %v", x.Statement)
}

// TODO: appending here is not so effective
func (df *Dockerfile) Marshal() []byte {
	res := make([]byte, 0)
	for _, st := range df.Statements {
		res = append(res, fmt.Sprintln(st)...)
	}
	return res
}
