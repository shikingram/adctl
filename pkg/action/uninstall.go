package action

import (
	"time"

	"github.com/shikingram/adctl/pkg/deploy"
)

type UnInstall struct {
	cfg *Configuration

	ReleaseName   string
	DryRun        bool
	CleanInstance bool
	Timeout       time.Duration
}

func NewUnInstall(cfg *Configuration) *UnInstall {
	return &UnInstall{cfg: cfg}
}

func (i *UnInstall) Run(name string) error {
	return deploy.UnInstall(name, i.CleanInstance)
}
