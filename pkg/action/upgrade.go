package action

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/shikingram/adctl/pkg/chart"
	"github.com/shikingram/adctl/pkg/chartutil"
	"github.com/shikingram/adctl/pkg/deploy"
)

type Upgrade struct {
	cfg         *Configuration
	ReleaseName string
}

func NewUpgrade(cfg *Configuration) *Upgrade {
	return &Upgrade{cfg: cfg}
}

func (i *Upgrade) NameAndChart(args []string) (string, error) {

	if len(args) > 1 {
		return args[0], errors.Errorf("expected at most one arguments, unexpected arguments: %v", strings.Join(args[1:], ", "))
	}

	if len(args) == 1 {
		return args[0], nil
	}

	if i.ReleaseName != "" {
		return i.ReleaseName, nil
	}
	return args[0], errors.New("must specify name")
}

func (i *Upgrade) Run(ch *chart.Chart, vals chartutil.Values) error {
	ctx := context.Background()
	return i.RunWithContext(ctx, ch, vals)
}

func (i *Upgrade) RunWithContext(ctx context.Context, ch *chart.Chart, vals chartutil.Values) error {
	options := chartutil.ReleaseOptions{
		Name:     i.ReleaseName,
		Revision: 1,
	}
	valuesToRender, err := chartutil.ToRenderValues(ch, options, vals)
	if err != nil {
		return err
	}
	err = i.cfg.renderResources(ch, valuesToRender, i.ReleaseName)
	if err != nil {
		return err
	}
	return deploy.UpgradeWithContext(ctx, i.ReleaseName)
}

func (i *Upgrade) ValidateName(name string) bool {
	num, err := deploy.CheckReleaseDeploy(name)
	return err == nil && num > 0
}
