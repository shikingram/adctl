/*
Copyright © 2022 Kingram <kingram@163.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/shikingram/adctl/cmd/require"
	"github.com/shikingram/adctl/pkg/action"
	"github.com/shikingram/adctl/pkg/chart/loader"
	"github.com/shikingram/adctl/pkg/cli/values"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const upgradeDesc = "This command upgrades a release to a new version of a chart"

func newUpgradeCmd(cfg *action.Configuration) *cobra.Command {
	client := action.NewUpgrade(cfg)
	valueOpts := &values.Options{}

	// upgradeCmd represents the upgrade command
	var cmd = &cobra.Command{
		Use:   "upgrade [NAME] [CHART]",
		Short: "upgrade application",
		Long:  upgradeDesc,
		Args:  require.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return require.Environment()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpgrade(args, client, valueOpts)
		},
	}
	addUpgradeFlags(cmd, cmd.Flags(), client, valueOpts)
	return cmd
}

func addUpgradeFlags(cmd *cobra.Command, f *pflag.FlagSet, client *action.Upgrade, valueOpts *values.Options) {
	// f.BoolVar(&client.UseReleaseName, "release-name", false, "use release name in the output-dir path.")
	f.BoolVar(&client.Force, "force", false, "force resource updates through a replacement strategy")
	addValueOptionsFlags(f, valueOpts)
}

func runUpgrade(args []string, client *action.Upgrade, valueOpts *values.Options) error {
	name, cp, err := client.NameAndChart(args)
	if err != nil {
		return err
	}

	// validate name
	if !client.ValidateName(name) {
		return errors.New("must specify already deployed name")
	}

	client.ReleaseName = name
	vals, err := valueOpts.MergeValues()
	if err != nil {
		return err
	}

	charts, err := loader.LoadChart(cp)
	if err != nil {
		return err
	}

	// Create context and prepare the handle of SIGTERM
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Handle SIGTERM
	cSignal := make(chan os.Signal, 1)
	signal.Notify(cSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cSignal
		fmt.Fprintf(os.Stdout, "Release %s has been cancelled.\n", args[0])
		cancel()
	}()

	return client.RunWithContext(ctx, charts, vals)
}
