package main

import (
	"errors"
	"strconv"
	"strings"
)

func LogicalLines([]byte) [][]byte {
	return nil
}

func ParseJSONArray(s string) ([]string, error) {
	body := TrimSpace(s)
	if !strings.HasPrefix(body, "[") || !strings.HasSurfix(body, "]") {
		return nil, errors.New("unmatched '[' and ']'")
	}
	args := strings.Split(body[1:len(body)-1], ",")
	for i, a := range args {
		arg := strings.TrimSpace(a)
		if !strings.HasPrefix(arg, "\"") ||
			!strings.HasSurfix(arg, "\"") ||
			len(arg) <= 1 {
			return nil, errors.New("unmatched '\"'")
		}
		args[i] = arg[1 : len(arg)-1]
	}
	return args, nil
}

func Parse(args []byte) (*Statement, error) {
	return nil, nil
}

func ParseFrom(body []byte) (*From, error) {
	s := strings.TrimSpace(string(body))
	if args := strings.Split(s, "@"); len(args) >= 2 {
		return &From{Image: args[0], Digest: args[1]}, nil
	}
	if args := strings.Split(s, ":"); len(args) >= 2 {
		return &From{Image: args[0], Tag: args[1]}, nil
	}
	return &From{Image: s}, nil
}

func ParseMaintainer(body []byte) (*Maintainer, error) {
	s := strings.TrimSpace(string(body))
	return &Maintainer{Name: s}, nil
}

func ParseRun(body []byte) (*Run, error) {
	return nil, nil
}

func ParseCmd(body []byte) (*Cmd, error) {
	return nil, nil
}

func ParseLable(body []byte) (*Label, error) {
	return nil, nil
}

func ParseExpose(body []byte) (*Expose, error) {
	args := strings.Fields(string(body))
	var ps = make([]int, len(args))
	for i, a := range args {
		p, err := strconv.Atoi(a)
		if err != nil {
			return nil, errors.New("illegal port number")
		}
		ps[i] = p
	}
	return &Expose{Ports: ps}, nil
}

func ParseEnv(body []byte) (*Env, error) {
	return nil, nil
}

func ParseAdd(body []byte) (*Add, error) {
	return nil, nil
}

func ParseCopy(body []byte) (*Copy, error) {
	return nil, nil
}

func ParseEntrypoint(body []byte) (*Entrypoint, error) {
	return nil, nil
}

func ParseVolume(body []byte) (*Volume, error) {
	return nil, nil
}

func ParseUser(body []byte) (*User, error) {
	s := strings.TrimSpace(string(body))
	return &User{Name: s}, nil
}

func ParseWorkdir(body []byte) (*Workdir, error) {
	s := strings.TrimSpace(string(body))
	return &Workdir{Path: s}, nil
}

func ParseOnbuild(body []byte) (*Onbuild, error) {
	return nil, nil
}

func (df *Dockerfile) Unmarshal(src []byte) error {
	return nil
}
