package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type VoicePipe struct {
	RootDir   string
	Directive *Directive
}

func NewVoicePipe(path string) (*VoicePipe, error) {
	d, err := NewDirective(path)
	if err != nil {
		return nil, err
	}
	return &VoicePipe{RootDir: path, Directive: d}, nil
}

func (vp *VoicePipe) Resources() ([]os.FileInfo, error) {
	rs := make([]os.FileInfo, 0)
	fis, err := ioutil.ReadDir(vp.RootDir)
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

func (vp *VoicePipe) SetupWorkDirFor(id ImageDirective, rs []os.FileInfo) error {
	dir := vp.RootDir + "/.voicepipe/" + id.Tag
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0775); err != nil {
		return err
	}
	for _, fi := range rs {
		if fi.Name() != "Dockerfile" {
			if err := os.Link(vp.RootDir+"/"+fi.Name(), dir+"/"+fi.Name()); err != nil {
				return err
			}
			continue
		}
		buf, err := ioutil.ReadFile(vp.RootDir + "/" + fi.Name())
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
		if err := ioutil.WriteFile(fi.Name(), df.Marshal(), 775); err != nil {
			return err
		}
	}
	return nil
}

func (vp *VoicePipe) SetupWorkingDir() error {
	rs, err := vp.Resources()
	if err != nil {
		return err
	}
	for _, id := range vp.Directive.ImageDirectives {
		if err := vp.SetupWorkDirFor(*id, rs); err != nil {
			return err
		}
	}
	return nil
}

func (vp *VoicePipe) BuildImages(stdout, stderr io.Writer) error {
	for _, id := range vp.Directive.ImageDirectives {
		dir := vp.RootDir + "/.voicepipe/" + id.Tag
		tag := vp.Directive.Repository + ":" + id.Tag
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
	err := vp.SetupWorkingDir()
	if err != nil {
		return err
	}
	stdout := bufio.NewWriter(os.Stdout)
	stderr := bufio.NewWriter(os.Stderr)
	err = vp.BuildImages(stdout, stderr)
	if err != nil {
		return err
	}
	return nil
}
