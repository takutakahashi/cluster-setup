package server

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Host  string
	Admin bool
}

func (s Server) Transport() error {
	logrus.Infof("transport files to %s", s.Host)
	if out, err := s.Execute([]string{"rm", "-rf", "cluster-setup"}, true); err != nil {
		logrus.Error(out)
		return err
	}
	c := exec.Command(
		"scp",
		"-r",
		"dist/",
		fmt.Sprintf("%s:cluster-setup", s.Host),
	)
	out, err := c.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "output: %s", out)
	}
	return nil
}

func (s Server) ExecuteMitamae() error {
	if err := s.setupMitamae(); err != nil {
		return err
	}
	return nil
}

func (s Server) ParseConfig() error {
	// TODO: this is mock
	if err := exec.Command("rm", "-rf", "dist").Run(); err != nil {
		return nil
	}
	return exec.Command("cp", "-rf", "assets", "dist").Run()
}

func (s Server) setupMitamae() error {
	out, err := s.Execute([]string{"bash", "cluster-setup/bin/install_mitamae.sh"}, true)
	logrus.Infof("out: %s", out)
	return err
}

func (s Server) Execute(params []string, sudo bool) ([]byte, error) {
	if sudo && !s.Admin {
		return nil, fmt.Errorf("can't exec as sudo")
	}
	in := []string{}
	in = append(in, s.Host)
	if sudo {
		in = append(in, "sudo")
	}
	in = append(in, params...)
	return exec.Command("ssh", in...).Output()
}
