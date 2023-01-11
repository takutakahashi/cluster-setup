package server

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/takutakahashi/cluster-setup/pkg/config"
	"gopkg.in/yaml.v3"
)

type Server struct {
	Version string
	Node    config.Node
	Host    string
	Admin   bool
	Secret  config.Secret
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
	if out, err := s.Execute([]string{"mitamae", "local", "cluster-setup/default.rb", "--dry-run"}, true); err != nil {
		logrus.Errorf("%s", out)
		return err
	} else {
		logrus.Errorf("%s", out)
	}
	return nil
}

func (s Server) ParseConfig() error {
	// TODO: this is mock
	if err := exec.Command("rm", "-rf", "dist").Run(); err != nil {
		return err
	}
	w, err := os.Create("assets/rootfs/etc/rancher/k3s/config.yaml")
	if err != nil {
		return err
	}
	fm := template.FuncMap{
		"toYaml": toYAML,
	}
	tpl, err := template.New("config.yaml").Funcs(fm).ParseFiles("assets/templates/etc/rancher/k3s/config.yaml")
	if err != nil {
		return err
	}
	if err := tpl.Execute(w, s); err != nil {
		return err
	}
	w.Close()
	w, err = os.Create("assets/bin/install_k3s.sh")
	if err != nil {
		return err
	}
	tpl, err = template.ParseFiles("assets/templates/bin/install_k3s.sh")
	if err != nil {
		return err
	}
	if err := tpl.Execute(w, s); err != nil {
		return err
	}
	if err := exec.Command("cp", "-rf", "assets", "dist").Run(); err != nil {
		return err
	}
	return nil
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

func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}
