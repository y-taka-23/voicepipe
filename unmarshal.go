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
	body := strings.TrimSpace(s)
	if !strings.HasPrefix(body, "[") || !strings.HasSuffix(body, "]") {
		return nil, errors.New("unmatched '[' and ']'")
	}
	args := strings.Split(body[1:len(body)-1], ",")
	for i, a := range args {
		arg := strings.TrimSpace(a)
		if !strings.HasPrefix(arg, "\"") ||
			!strings.HasSuffix(arg, "\"") ||
			len(arg) <= 1 {
			return nil, errors.New("unmatched '\"'")
		}
		args[i] = arg[1 : len(arg)-1]
	}
	return args, nil
}

func ParseLine(line []byte) (Statement, error) {
	s := string(line)
	i := strings.IndexAny(s, " \t")
	if i < 0 {
		return nil, errors.New("missing argument")
	}
	instr := strings.ToUpper(s[:i])
	body := line[i:]
	switch instr {
	case "FROM":
		return ParseFrom(body)
	case "MAINTAINER":
		return ParseMaintainer(body)
	case "RUN":
		return ParseRun(body)
	case "CMD":
		return ParseCmd(body)
	case "LABEL":
		return ParseLabel(body)
	case "EXPOSE":
		return ParseExpose(body)
	case "ENV":
		return ParseEnv(body)
	case "ADD":
		return ParseAdd(body)
	case "COPY":
		return ParseCopy(body)
	case "ENTRYPOINT":
		return ParseEntrypoint(body)
	case "VOLUME":
		return ParseVolume(body)
	case "USER":
		return ParseUser(body)
	case "WORKDIR":
		return ParseWorkdir(body)
	case "Onbuild":
		return ParseOnbuild(body)
	}
	return nil, errors.New("illegal instruction")
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
	ts, err := ParseJSONArray(string(body))
	if err == nil {
		return &Run{Tokens: ts}, nil
	}
	ts = []string{"/bin/sh", "-c"}
	for _, t := range strings.Fields(string(body)) {
		ts = append(ts, t)
	}
	return &Run{Tokens: ts}, nil
}

func ParseCmd(body []byte) (*Cmd, error) {
	ts, err := ParseJSONArray(string(body))
	if err == nil {
		return &Cmd{Tokens: ts}, nil
	}
	ts = []string{"/bin/sh", "-c"}
	for _, t := range strings.Fields(string(body)) {
		ts = append(ts, t)
	}
	return &Cmd{Tokens: ts}, nil
}

func FetchKey(src []byte) (string, []byte, error) {
	for i, b := range src {
		if b == '=' {
			return strings.Trim(string(src[:i]), " \""), src[i+1:], nil
		}
	}
	return "", nil, errors.New("missing '='")
}

func FetchValue(src []byte) (string, []byte, error) {
	if len(src) == 0 {
		return "", nil, errors.New("missing value")
	}
	if src[0] == '"' {
		for i := 1; i < len(src); i++ {
			if src[i] == '"' {
				return string(src[1:i]), src[i+1:], nil
			}
		}
		return "", nil, errors.New("unmatched '\"'")
	} else {
		for i, b := range src {
			if b == ' ' {
				return string(src[:i]), src[i+1:], nil
			}
		}
		return string(src), []byte{}, nil
	}
}

func ParseLabel(body []byte) (*Label, error) {
	if !strings.Contains(string(body), "=") {
		return &Label{Labels: map[string]string{}}, nil
	}
	k, t, err := FetchKey(body)
	if err != nil {
		return nil, err
	}
	v, tail, err := FetchValue(t)
	if err != nil {
		return nil, err
	}
	l, err := ParseLabel(tail)
	if err != nil {
		return nil, err
	}
	l.Labels[k] = v
	return l, nil
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

func ParseMultiEnv(body []byte) (*Env, error) {
	if !strings.Contains(string(body), "=") {
		return &Env{Variables: map[string]string{}}, nil
	}
	k, t, err := FetchKey(body)
	if err != nil {
		return nil, err
	}
	v, tail, err := FetchValue(t)
	if err != nil {
		return nil, err
	}
	e, err := ParseMultiEnv(tail)
	if err != nil {
		return nil, err
	}
	e.Variables[k] = v
	return e, nil
}

func ParseSingleEnv(body []byte) (*Env, error) {
	s := strings.TrimSpace(string(body))
	i := strings.IndexAny(s, " \t")
	if i < 0 || i >= len(s)-1 {
		return nil, errors.New("missing value")
	}
	k := strings.TrimSpace(s[:i])
	v := strings.TrimSpace(s[i+1:])
	return &Env{Variables: map[string]string{k: v}}, nil
}

func ParseEnv(body []byte) (*Env, error) {
	if strings.Contains(string(body), "=") {
		return ParseMultiEnv(body)
	} else {
		return ParseSingleEnv(body)
	}
}

func ParseAdd(body []byte) (*Add, error) {
	fs, err := ParseJSONArray(string(body))
	if err != nil {
		fs = strings.Fields(string(body))
	}
	if len(fs) == 0 {
		return nil, errors.New("no destination directory")
	}
	return &Add{Sources: fs[:len(fs)-1], Destination: fs[len(fs)-1]}, nil
}

func ParseCopy(body []byte) (*Copy, error) {
	fs, err := ParseJSONArray(string(body))
	if err != nil {
		fs = strings.Fields(string(body))
	}
	if len(fs) == 0 {
		return nil, errors.New("no destination directory")
	}
	return &Copy{Sources: fs[:len(fs)-1], Destination: fs[len(fs)-1]}, nil
}

func ParseEntrypoint(body []byte) (*Entrypoint, error) {
	ts, err := ParseJSONArray(string(body))
	if err == nil {
		return &Entrypoint{Tokens: ts}, nil
	}
	ts = []string{"/bin/sh", "-c"}
	for _, t := range strings.Fields(string(body)) {
		ts = append(ts, t)
	}
	return &Entrypoint{Tokens: ts}, nil
}

func ParseVolume(body []byte) (*Volume, error) {
	ps, err := ParseJSONArray(string(body))
	if err != nil {
		return &Volume{Points: strings.Fields(string(body))}, nil
	}
	return &Volume{Points: ps}, nil
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
