package deploy

import (
	"github.com/takutakahashi/cluster-setup/pkg/config"
	"github.com/takutakahashi/cluster-setup/pkg/server"
)

func Execute(cfg *config.Config) error {
	for _, node := range cfg.Nodes {

		server.ParseConfig(node)
		for _, t := range node.Targets {
			if err := server.Transport(t); err != nil {
				return err
			}
			if err := server.ExecuteMitamae(t); err != nil {
				return err
			}
		}

	}
	return nil
}
