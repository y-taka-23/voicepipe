package main

type Statement interface {
	statement()
}

type Dockerfile []*Statement

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

type Lable struct {
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
	Statement *Statement
}

func (x *From) statement()       {}
func (x *Maintainer) statement() {}
func (x *Run) statement()        {}
func (x *Cmd) statement()        {}
func (x *Lable) statement()      {}
func (x *Expose) statement()     {}
func (x *Env) statement()        {}
func (x *Add) statement()        {}
func (x *Copy) statement()       {}
func (x *Entrypoint) statement() {}
func (x *Volume) statement()     {}
func (x *User) statement()       {}
func (x *Workdir) statement()    {}
func (x *Onbuild) statement()    {}
