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
	"fmt"

	"github.com/shikingram/adctl/cmd/require"
	"github.com/shikingram/adctl/pkg/action"
	"github.com/spf13/cobra"
)

const unInstallDesc = `This command takes a release name and uninstalls the release.`

func newUnInstallCmd(cfg *action.Configuration) *cobra.Command {
	client := action.NewUnInstall(cfg)
	// cmd represents the uninstall command
	var cmd = &cobra.Command{
		Use:        "uninstall",
		Short:      "uninstall application",
		Aliases:    []string{"del", "delete", "un"},
		SuggestFor: []string{"remove", "rm"},
		Args:       require.MinimumNArgs(1),
		Long:       unInstallDesc,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return require.Environment()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			for i := 0; i < len(args); i++ {

				err := client.Run(args[i])
				if err != nil {
					return err
				}

				fmt.Printf("release \"%s\" uninstalled\n", args[i])
			}
			return nil
		},
	}

	// f := cmd.Flags()
	// f.BoolVar(&client.DryRun, "dry-run", false, "simulate a uninstall")
	// f.DurationVar(&client.Timeout, "timeout", 300*time.Second, "time to wait for any individual Kubernetes operation (like Jobs for hooks)")

	return cmd
}
