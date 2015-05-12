package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type ImageDirective struct {
	Tag        string
	Parameters map[string]string
}

type Directive struct {
	Repository      string
	ImageDirectives []*ImageDirective
}

type intermediateDirective struct {
	Repository string
	Images     []struct {
		Tag        string
		Parameters []struct {
			Name  string
			Value string
		}
	}
}

func main() {
	buf, err := ioutil.ReadFile("voicepipe.yml")
	if err != nil {
		log.Println(err)
		return
	}

	in := &intermediateDirective{}
	err = yaml.Unmarshal(buf, in)
	if err != nil {
		log.Println(err)
		return
	}

	d := &Directive{
		Repository:      in.Repository,
		ImageDirectives: make([]*ImageDirective, 0),
	}
	for _, i := range in.Images {
		params := make(map[string]string)
		for _, p := range i.Parameters {
			params[p.Name] = params[p.Value]
		}
		d.ImageDirectives = append(
			d.ImageDirectives,
			&ImageDirective{Tag: i.Tag, Parameters: params},
		)
	}
	fmt.Println(d)
}
