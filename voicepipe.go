package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type VoicePipe struct {
}

func NewVoicePipe() *VoicePipe {
	return &VoicePipe{}
}

func (vp *VoicePipe) Resources(root string) ([]os.FileInfo, error) {
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

func (vp *VoicePipe) SetupWorkingDir(d Directive, root string) error {
	rs, err := vp.Resources(root)
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
			if fi.Name() != "Dockerfile" {
				err = os.Link(root+"/"+fi.Name(), dir+"/"+fi.Name())
				if err != nil {
					return err
				}
			}
			buf, err := ioutil.ReadFile(root + "/Dockerfile")
			if err != nil {
				return err
			}
			df, err := Unmarshal(buf)
			if err != nil {
				return err
			}
			// TODO: copying structures costs a lot
			for k, v := range id.Parameters {
				df = ReplaceEnv(*df, k, v)
			}
			err = ioutil.WriteFile(dir+"/Dockerfile", df.Marshal(), 775)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (vp *VoicePipe) BuildImages(d Directive, root string, stdout, stderr io.Writer) error {
	for _, id := range d.ImageDirectives {
		dir := root + "/.voicepipe/" + id.Tag
		tag := d.Repository + ":" + id.Tag
		cmd := exec.Command("docker", "build", "--rm", "-t", tag, dir)
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (vp *VoicePipe) Run() error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	root += "/example" // just for debug

	d, err := NewDirective(root + "/voicepipe.yml")
	if err != nil {
		return err
	}

	err = vp.SetupWorkingDir(*d, root)
	if err != nil {
		return err
	}

	stdout := bufio.NewWriter(os.Stdout)
	stderr := bufio.NewWriter(os.Stderr)
	err = vp.BuildImages(*d, root, stdout, stderr)
	if err != nil {
		return err
	}

	return nil
}
