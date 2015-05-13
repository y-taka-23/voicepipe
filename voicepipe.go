package main

import (
	//	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
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
			params[p.Name] = params[p.Value]
		}
		d.ImageDirectives = append(
			d.ImageDirectives,
			&ImageDirective{Tag: i.Tag, Parameters: params},
		)
	}

	return nil
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}
	wd += "/example" // just for debug

	buf, err := ioutil.ReadFile(wd + "/voicepipe.yml")
	if err != nil {
		log.Println(err)
		return
	}

	var d = &Directive{}
	err = yaml.Unmarshal(buf, d)
	if err != nil {
		log.Println(err)
		return
	}

	fis, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Println(err)
		return
	}

	for _, id := range d.ImageDirectives {
		dir := wd + "/.voicepipe/" + id.Tag
		err = os.RemoveAll(dir)
		if err != nil {
			log.Println(err)
			return
		}
		err = os.MkdirAll(dir, 0775)
		if err != nil {
			log.Println(err)
			return
		}
		for _, fi := range fis {
			err = os.Symlink(wd+"/"+fi.Name(), dir+"/"+fi.Name())
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
