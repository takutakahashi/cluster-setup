package deploy

import (
	"github.com/takutakahashi/cluster-setup/pkg/config"
	"github.com/takutakahashi/cluster-setup/pkg/server"
)

func Execute(cfg *config.Config) error {
	for _, node := range cfg.Nodes {
		for _, t := range node.Targets {

			s := server.Server{
				Version: cfg.Version,
				Host:    t,
				Admin:   true,
				Node:    node,
				Secret:  cfg.Secret,
			}
			if err := s.ParseConfig(); err != nil {
				return err
			}
			if err := s.Transport(); err != nil {
				return err
			}
			if err := s.ExecuteMitamae(); err != nil {
				return err
			}
		}

	}
	return nil
}
