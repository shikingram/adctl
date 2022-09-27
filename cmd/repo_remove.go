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
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/shikingram/adctl/cmd/require"
	"github.com/shikingram/adctl/pkg/adctlpath"
	"github.com/shikingram/adctl/pkg/repo"
	"github.com/spf13/cobra"
)

type repoRemoveOptions struct {
	names     []string
	repoFile  string
	repoCache string
}

func newRepoRemoveCmd() *cobra.Command {
	o := &repoRemoveOptions{}

	cmd := &cobra.Command{
		Use:     "remove [REPO1 [REPO2 ...]]",
		Aliases: []string{"rm"},
		Short:   "remove one or more chart repositories",
		Args:    require.MinimumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return compListRepos(toComplete, args), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			o.repoFile = settings.RepositoryConfig
			o.repoCache = settings.RepositoryCache
			o.names = args
			return o.run()
		},
	}
	return cmd
}

func (o *repoRemoveOptions) run() error {
	r, err := repo.LoadFile(o.repoFile)
	if isNotExist(err) || len(r.Repositories) == 0 {
		return errors.New("no repositories configured")
	}

	for _, name := range o.names {
		if !r.Remove(name) {
			return errors.Errorf("no repo named %q found", name)
		}
		if err := r.WriteFile(o.repoFile, 0644); err != nil {
			return err
		}

		if err := removeRepoCache(o.repoCache, name); err != nil {
			return err
		}
		fmt.Printf("%q has been removed from your repositories\n", name)
	}

	return nil
}

func removeRepoCache(root, name string) error {
	idx := filepath.Join(root, adctlpath.CacheChartsFile(name))
	if _, err := os.Stat(idx); err == nil {
		os.Remove(idx)
	}

	idx = filepath.Join(root, adctlpath.CacheIndexFile(name))
	if _, err := os.Stat(idx); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return errors.Wrapf(err, "can't remove index file %s", idx)
	}
	return os.Remove(idx)
}
