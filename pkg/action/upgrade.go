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
	Force       bool
}

func NewUpgrade(cfg *Configuration) *Upgrade {
	return &Upgrade{cfg: cfg}
}

func (i *Upgrade) NameAndChart(args []string) (string, string, error) {

	if len(args) > 2 {
		return args[0], args[1], errors.Errorf("expected at most two arguments, unexpected arguments: %v", strings.Join(args[1:], ", "))
	}

	if len(args) == 2 {
		return args[0], args[1], nil
	}

	if i.ReleaseName != "" {
		return i.ReleaseName, args[0], nil
	}
	return "", "", errors.New("must specify name")
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
	return deploy.UpgradeWithContext(ctx, i.ReleaseName,i.Force)
}

func (i *Upgrade) ValidateName(name string) bool {
	num, err := deploy.CheckReleaseDeploy(name)
	return err == nil && num > 0
}
