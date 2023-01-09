package deploy

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/takutakahashi/cluster-setup/pkg/config"
)

func Execute(cfg *config.Config) error {
	logrus.Info(cfg)
	return fmt.Errorf("not implemented")
}
