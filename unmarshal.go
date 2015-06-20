package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func logicalLines(src []byte) [][]byte {
	src = regexp.MustCompile("\\\\\n").ReplaceAll(src, []byte(""))
	lines := [][]byte{}
	head := 0
	for i, b := range src {
		if b == '\n' {
			lines = append(lines, src[head:i])
			head = i + 1
		}
	}
	if head < len(src) {
		lines = append(lines, src[head:])
	}
	return lines
}

func trimComment(line []byte) []byte {
	for i, b := range line {
		if b == '#' {
			return line[:i]
		}
	}
	return line
}

func parseJSONArray(s string) ([]string, error) {
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

func parseLine(line []byte) (Statement, error) {
	s := string(line)
	i := strings.IndexAny(s, " \t")
	if i < 0 {
		return nil, errors.New("missing argument")
	}
	instr := strings.ToUpper(s[:i])
	body := line[i:]
	switch instr {
	case "FROM":
		return parseFrom(body)
	case "MAINTAINER":
		return parseMaintainer(body)
	case "RUN":
		return parseRun(body)
	case "CMD":
		return parseCmd(body)
	case "LABEL":
		return parseLabel(body)
	case "EXPOSE":
		return parseExpose(body)
	case "ENV":
		return parseEnv(body)
	case "ADD":
		return parseAdd(body)
	case "COPY":
		return parseCopy(body)
	case "ENTRYPOINT":
		return parseEntrypoint(body)
	case "VOLUME":
		return parseVolume(body)
	case "USER":
		return parseUser(body)
	case "WORKDIR":
		return parseWorkdir(body)
	case "ONBUILD":
		return parseOnbuild(body)
	}
	return nil, fmt.Errorf("illegal instruction '%s'", instr)
}

func parseFrom(body []byte) (*From, error) {
	s := strings.TrimSpace(string(body))
	if args := strings.Split(s, "@"); len(args) >= 2 {
		return &From{Image: args[0], Digest: args[1]}, nil
	}
	if args := strings.Split(s, ":"); len(args) >= 2 {
		return &From{Image: args[0], Tag: args[1]}, nil
	}
	return &From{Image: s}, nil
}

func parseMaintainer(body []byte) (*Maintainer, error) {
	s := strings.TrimSpace(string(body))
	return &Maintainer{Name: s}, nil
}

func parseRun(body []byte) (*Run, error) {
	ts, err := parseJSONArray(string(body))
	if err == nil {
		return &Run{Tokens: ts}, nil
	}
	ts = []string{"/bin/sh", "-c"}
	for _, t := range strings.Fields(string(body)) {
		ts = append(ts, t)
	}
	return &Run{Tokens: ts}, nil
}

func parseCmd(body []byte) (*Cmd, error) {
	ts, err := parseJSONArray(string(body))
	if err == nil {
		return &Cmd{Tokens: ts}, nil
	}
	ts = []string{"/bin/sh", "-c"}
	for _, t := range strings.Fields(string(body)) {
		ts = append(ts, t)
	}
	return &Cmd{Tokens: ts}, nil
}

func fetchKey(src []byte) (string, []byte, error) {
	for i, b := range src {
		if b == '=' {
			return strings.Trim(string(src[:i]), " \""), src[i+1:], nil
		}
	}
	return "", nil, errors.New("missing '='")
}

func fetchValue(src []byte) (string, []byte, error) {
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

func parseLabel(body []byte) (*Label, error) {
	if !strings.Contains(string(body), "=") {
		return &Label{Labels: map[string]string{}}, nil
	}
	k, t, err := fetchKey(body)
	if err != nil {
		return nil, err
	}
	v, tail, err := fetchValue(t)
	if err != nil {
		return nil, err
	}
	l, err := parseLabel(tail)
	if err != nil {
		return nil, err
	}
	l.Labels[k] = v
	return l, nil
}

func parseExpose(body []byte) (*Expose, error) {
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

func parseMultiEnv(body []byte) (*Env, error) {
	if !strings.Contains(string(body), "=") {
		return &Env{Variables: map[string]string{}}, nil
	}
	k, t, err := fetchKey(body)
	if err != nil {
		return nil, err
	}
	v, tail, err := fetchValue(t)
	if err != nil {
		return nil, err
	}
	e, err := parseMultiEnv(tail)
	if err != nil {
		return nil, err
	}
	e.Variables[k] = v
	return e, nil
}

func parseSingleEnv(body []byte) (*Env, error) {
	s := strings.TrimSpace(string(body))
	i := strings.IndexAny(s, " \t")
	if i < 0 || i >= len(s)-1 {
		return nil, errors.New("missing value")
	}
	k := strings.TrimSpace(s[:i])
	v := strings.TrimSpace(s[i+1:])
	return &Env{Variables: map[string]string{k: v}}, nil
}

func parseEnv(body []byte) (*Env, error) {
	if strings.Contains(string(body), "=") {
		return parseMultiEnv(body)
	} else {
		return parseSingleEnv(body)
	}
}

func parseAdd(body []byte) (*Add, error) {
	fs, err := parseJSONArray(string(body))
	if err != nil {
		fs = strings.Fields(string(body))
	}
	if len(fs) == 0 {
		return nil, errors.New("no destination directory")
	}
	return &Add{Sources: fs[:len(fs)-1], Destination: fs[len(fs)-1]}, nil
}

func parseCopy(body []byte) (*Copy, error) {
	fs, err := parseJSONArray(string(body))
	if err != nil {
		fs = strings.Fields(string(body))
	}
	if len(fs) == 0 {
		return nil, errors.New("no destination directory")
	}
	return &Copy{Sources: fs[:len(fs)-1], Destination: fs[len(fs)-1]}, nil
}

func parseEntrypoint(body []byte) (*Entrypoint, error) {
	ts, err := parseJSONArray(string(body))
	if err == nil {
		return &Entrypoint{Tokens: ts}, nil
	}
	ts = []string{"/bin/sh", "-c"}
	for _, t := range strings.Fields(string(body)) {
		ts = append(ts, t)
	}
	return &Entrypoint{Tokens: ts}, nil
}

func parseVolume(body []byte) (*Volume, error) {
	ps, err := parseJSONArray(string(body))
	if err != nil {
		return &Volume{Points: strings.Fields(string(body))}, nil
	}
	return &Volume{Points: ps}, nil
}

func parseUser(body []byte) (*User, error) {
	s := strings.TrimSpace(string(body))
	return &User{Name: s}, nil
}

func parseWorkdir(body []byte) (*Workdir, error) {
	s := strings.TrimSpace(string(body))
	return &Workdir{Path: s}, nil
}

func parseOnbuild(body []byte) (*Onbuild, error) {
	st, err := parseLine(body)
	if err != nil {
		return nil, err
	}
	return &Onbuild{Statement: st}, nil
}

func unmarshal(src []byte) (*Dockerfile, error) {
	sts := []Statement{}
	lines := logicalLines(src)
	for _, l := range lines {
		content := trimComment(l)
		if len(strings.TrimSpace(string(content))) != 0 {
			st, err := parseLine(content)
			if err != nil {
				return nil, err
			}
			sts = append(sts, st)
		}
	}
	return &Dockerfile{Statements: sts}, nil
}
