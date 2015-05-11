package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Image struct {
	Tag        string
	Parameters map[string]string
}

type Directive struct {
	Repository string
	Images     []*Image
}

type unmarshalledYaml struct {
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

	uy := &unmarshalledYaml{}
	err = yaml.Unmarshal(buf, uy)
	if err != nil {
		log.Println(err)
		return
	}

	d := &Directive{
		Repository: uy.Repository,
		Images:     make([]*Image, 0),
	}
	for _, i := range uy.Images {
		params := make(map[string]string)
		for _, p := range i.Parameters {
			params[p.Name] = params[p.Value]
		}
		d.Images = append(d.Images, &Image{Tag: i.Tag, Parameters: params})
	}
	fmt.Println(d)
}
