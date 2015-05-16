package main

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
	// to be inplemented
}
