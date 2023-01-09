package deploy

import (
	"github.com/takutakahashi/cluster-setup/pkg/config"
	"github.com/takutakahashi/cluster-setup/pkg/server"
)

func Execute(cfg *config.Config) error {
	for _, node := range cfg.Nodes {
		for _, t := range node.Targets {

			s := server.Server{
				Host:  t,
				Admin: true,
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
