package main

import (
	"fmt"
	"strings"
)

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
		res = append(res, fmt.Sprintf(" \"%s\"=\"%s\"", k, v)...)
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
		res = append(res, fmt.Sprintf(" %s=\"%s\"", k, v)...)
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
func (df *Dockerfile) marshal() []byte {
	res := make([]byte, 0)
	for _, st := range df.Statements {
		res = append(res, fmt.Sprintln(st)...)
	}
	return res
}
