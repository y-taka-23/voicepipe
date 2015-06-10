package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

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

func BuildImages(d Directive, root string, stdout, stderr io.Writer) error {
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

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}
	root += "/example" // just for debug

	d, err := NewDirective(root + "/voicepipe.yml")
	if err != nil {
		log.Println(err)
		return
	}

	err = SetupWorkingDir(*d, root)
	if err != nil {
		log.Println(err)
		return
	}

	stdout := bufio.NewWriter(os.Stdout)
	stderr := bufio.NewWriter(os.Stderr)
	err = BuildImages(*d, root, stdout, stderr)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("SUCCESS")
}
