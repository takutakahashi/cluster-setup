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
	return fmt.Errorf("not implemented")
}

func (s Server) ParseConfig() error {
	return fmt.Errorf("not implemented")
}

func (s Server) setupMitamae() error {
	_, err := s.Execute("/usr/local/bin/install_mitamae.sh", true)
	return err
}

func (s Server) Execute(execPath string, sudo bool) ([]byte, error) {
	if sudo && !s.Admin {
		return nil, fmt.Errorf("can't exec as sudo")
	}
	var c *exec.Cmd
	if sudo {
		c = exec.Command(
			"ssh",
			s.Host,
			"sudo",
			execPath)
	} else {
		c = exec.Command(
			"ssh",
			s.Host,
			execPath)
	}
	return c.Output()
}
