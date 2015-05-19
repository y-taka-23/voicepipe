package main

import (
	"strings"
)

func LogicalLines(src string) []string {
	s := strings.Replace(src, "\\\n", " ", -1)
	return strings.Split(s, "\n")
}

func Tokenize(line string) []string {
	s := strings.SplitN(line, "#", 1)[0]
	return strings.Fields(s)
}

func (df *Dockerfile) Unmarshal(src []byte) error {
	return nil
}
