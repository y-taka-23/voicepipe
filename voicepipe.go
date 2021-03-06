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

func newVoicePipe(root string, stdout, stderr io.Writer) (*VoicePipe, error) {
	d, err := newDirective(path.Join(root, "voicepipe.yml"))
	if err != nil {
		return nil, err
	}
	return &VoicePipe{
		RootDir:   root,
		Directive: d,
		Stdout:    stdout,
		Stderr:    stderr,
	}, nil
}

func (vp *VoicePipe) resources() ([]os.FileInfo, error) {
	rs := make([]os.FileInfo, 0)
	fis, err := ioutil.ReadDir(vp.RootDir)
	if err != nil {
		return rs, err
	}
	for _, fi := range fis {
		n := fi.Name()
		if n != ".voicepipe" && n != "voicepipe.yml" {
			rs = append(rs, fi)
		}
	}
	return rs, nil
}

func (vp *VoicePipe) setup(id ImageDirective, rs []os.FileInfo) error {
	dir := path.Join(vp.RootDir, ".voicepipe", id.Tag)
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0775); err != nil {
		return err
	}
	for _, fi := range rs {
		src := path.Join(vp.RootDir, fi.Name())
		if fi.Name() != "Dockerfile" {
			tgt := path.Join(dir, fi.Name())
			if err := os.Link(src, tgt); err != nil {
				return err
			}
			continue
		}
		buf, err := ioutil.ReadFile(src)
		if err != nil {
			return err
		}
		df, err := unmarshal(buf)
		if err != nil {
			return err
		}
		// TODO: copying structures costs a lot
		for k, v := range id.Parameters {
			df = replaceEnv(*df, k, v)
		}
		tgt := path.Join(dir, fi.Name())
		if err := ioutil.WriteFile(tgt, df.marshal(), 775); err != nil {
			return err
		}
	}
	return nil
}

func (vp *VoicePipe) setupAll() error {
	rs, err := vp.resources()
	if err != nil {
		return err
	}
	for _, id := range vp.Directive.ImageDirectives {
		if err := vp.setup(*id, rs); err != nil {
			return err
		}
	}
	return nil
}

func (vp *VoicePipe) build(id ImageDirective) error {
	dir := path.Join(vp.RootDir, ".voicepipe", id.Tag)
	fullName := vp.Directive.Repository + ":" + id.Tag
	cmd := exec.Command("docker", "build", "--rm", "-t", fullName, dir)
	cmd.Stdout = vp.Stdout
	cmd.Stderr = vp.Stderr
	vp.showHeader(vp.Directive.Repository, id.Tag)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (vp *VoicePipe) buildAll() error {
	for _, id := range vp.Directive.ImageDirectives {
		if err := vp.build(*id); err != nil {
			return err
		}
	}
	return nil
}

func (vp *VoicePipe) buildLatest() error {
	repo := vp.Directive.Repository
	cmd := exec.Command("docker", "build", "--rm", "-t", repo, vp.RootDir)
	cmd.Stdout = vp.Stdout
	cmd.Stderr = vp.Stderr
	vp.showHeader(repo, "latest")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (vp *VoicePipe) showHeader(repo, tag string) {
	fmt.Fprintln(vp.Stdout, "-------------------------------------------------------------------------------")
	fmt.Fprintf(vp.Stdout, "[VoicePipe] Building image %s:%s\n", repo, tag)
	fmt.Fprintln(vp.Stdout, "-------------------------------------------------------------------------------")
}

func (vp *VoicePipe) list() {
	fmt.Fprint(vp.Stdout, "REPOSITORY:\n")
	fmt.Fprintf(vp.Stdout, "   %s\n", vp.Directive.Repository)
	fmt.Fprint(vp.Stdout, "\n")
	fmt.Fprint(vp.Stdout, "TAGS:\n")
	for _, id := range vp.Directive.ImageDirectives {
		fmt.Fprintf(vp.Stdout, "   %s\t%s\n", id.Tag, id.Description)
	}
	fmt.Fprint(vp.Stdout, "\n")
}

func (vp *VoicePipe) cleanAll() error {
	dir := path.Join(vp.RootDir, ".voicepipe")
	if _, err := os.Stat(dir); err != nil {
		// the directory does not exist
		return nil
	}
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return nil
}
