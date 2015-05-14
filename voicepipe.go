package main

import (
	"fmt"
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

func SetupWorkingDir(d Directive, root string) error {
	fis, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}
	for _, id := range d.ImageDirectives {
		dir := root + "/.voicepipe/" + id.Tag
		err = os.RemoveAll(dir)
		if err != nil {
			return err
		}
		err = os.MkdirAll(dir, 0775)
		if err != nil {
			return err
		}
		for _, fi := range fis {
			err = os.Symlink(root+"/"+fi.Name(), dir+"/"+fi.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}
	root += "/example" // just for debug

	buf, err := ioutil.ReadFile(root + "/voicepipe.yml")
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

	err = SetupWorkingDir(*d, root)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("SUCCESS")
}
