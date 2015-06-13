package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

type VoicePipe struct {
	RootDir   string
	Directive *Directive
	Stdout    io.Writer
	Stderr    io.Writer
}

func NewVoicePipe(path string, stdout, stderr io.Writer) (*VoicePipe, error) {
	d, err := NewDirective(path)
	if err != nil {
		return nil, err
	}
	return &VoicePipe{
		RootDir:   path,
		Directive: d,
		Stdout:    stdout,
		Stderr:    stderr,
	}, nil
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

func (vp *VoicePipe) Setup(id ImageDirective, rs []os.FileInfo) error {
	dir := path.Join(vp.RootDir, ".voicepipe", id.Tag)
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0775); err != nil {
		return err
	}
	for _, fi := range rs {
		if fi.Name() != "Dockerfile" {
			src := path.Join(vp.RootDir, fi.Name())
			tgt := path.Join(dir, fi.Name())
			if err := os.Link(src, tgt); err != nil {
				return err
			}
			continue
		}
		buf, err := ioutil.ReadFile(path.Join(vp.RootDir, fi.Name()))
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

func (vp *VoicePipe) SetupAll() error {
	rs, err := vp.Resources()
	if err != nil {
		return err
	}
	for _, id := range vp.Directive.ImageDirectives {
		if err := vp.Setup(*id, rs); err != nil {
			return err
		}
	}
	return nil
}

func (vp *VoicePipe) Build(id ImageDirective) error {
	dir := path.Join(vp.RootDir, ".voicepipe", id.Tag)
	tag := vp.Directive.Repository + ":" + id.Tag
	cmd := exec.Command("docker", "build", "--rm", "-t", tag, dir)
	cmd.Stdout = vp.Stdout
	cmd.Stderr = vp.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (vp *VoicePipe) BuildAll() error {
	for _, id := range vp.Directive.ImageDirectives {
		if err := vp.Build(*id); err != nil {
			return err
		}
	}
	return nil
}

func (vp *VoicePipe) List() {
	fmt.Fprint(vp.Stdout, "REPOSITORY:\n")
	fmt.Fprintf(vp.Stdout, "   %s\n", vp.Directive.Repository)
	fmt.Fprint(vp.Stdout, "TAGS:\n")
	for _, id := range vp.Directive.ImageDirectives {
		fmt.Fprintf(vp.Stdout, "   %s\n", id.Tag)
	}
}

func (vp *VoicePipe) CleanAll() error {
	dir := path.Join(vp.RootDir, ".voicepipe")
	if _, err := os.Stat(dir); err == nil {
		// the directory does not exist
		return nil
	}
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return nil
}
