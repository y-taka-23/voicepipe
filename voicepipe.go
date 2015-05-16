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

func Resources(root string) ([]os.FileInfo, error) {
	rs := make([]os.FileInfo, 0)
	fis, err := ioutil.ReadDir(root)
	if err != nil {
		return rs, err
	}
	for _, fi := range fis {
		n := fi.Name()
		if n != ".voicepipe" && n != "voicepipe.yml" && n != "Dockerfile" {
			rs = append(rs, fi)
		}
	}
	return rs, nil
}

func SetupWorkingDir(d Directive, root string) error {
	rs, err := Resources(root)
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
		for _, fi := range rs {
			err = os.Symlink(root+"/"+fi.Name(), dir+"/"+fi.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

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
