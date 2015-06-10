package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

func (d *Directive) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var in = &intermediateDirective{}
	if err := unmarshal(in); err != nil {
		return err
	}
	d.Repository = in.Repository
	d.ImageDirectives = make([]*ImageDirective, 0)
	for _, i := range in.Images {
		params := make(map[string]string)
		for _, p := range i.Parameters {
			params[p.Name] = p.Value
		}
		d.ImageDirectives = append(
			d.ImageDirectives,
			&ImageDirective{Tag: i.Tag, Parameters: params},
		)
	}
	return nil
}

func NewDirective(path string) (*Directive, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var d = &Directive{}
	err = yaml.Unmarshal(buf, d)
	if err != nil {
		return nil, err
	}
	return d, nil
}
