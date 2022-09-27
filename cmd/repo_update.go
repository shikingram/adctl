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
	"sync"

	"github.com/pkg/errors"
	"github.com/shikingram/adctl/cmd/require"
	"github.com/shikingram/adctl/pkg/getter"
	"github.com/shikingram/adctl/pkg/repo"
	"github.com/spf13/cobra"
)

const updateDesc = `
Update gets the latest information about charts from the respective chart repositories.
Information is cached locally.

You can optionally specify a list of repositories you want to update.
	$ adctl repo update <repo_name> ...
To update all the repositories, use 'adctl repo update'.
`

var errNoRepositories = errors.New("no repositories found. You must add one before updating")

type repoUpdateOptions struct {
	update               func([]*repo.ChartRepository, bool) error
	repoFile             string
	repoCache            string
	names                []string
	failOnRepoUpdateFail bool
}

func newRepoUpdateCmd() *cobra.Command {
	o := &repoUpdateOptions{update: updateCharts}

	cmd := &cobra.Command{
		Use:     "update [REPO1 [REPO2 ...]]",
		Aliases: []string{"up"},
		Short:   "update information of available charts locally from chart repositories",
		Long:    updateDesc,
		Args:    require.MinimumNArgs(0),
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

	f := cmd.Flags()

	f.BoolVar(&o.failOnRepoUpdateFail, "fail-on-repo-update-fail", false, "update fails if any of the repository updates fail")

	return cmd
}

func isNotExist(err error) bool {
	return os.IsNotExist(errors.Cause(err))
}

func (o *repoUpdateOptions) run() error {
	f, err := repo.LoadFile(o.repoFile)
	switch {
	case isNotExist(err):
		return errNoRepositories
	case err != nil:
		return errors.Wrapf(err, "failed loading file: %s", o.repoFile)
	case len(f.Repositories) == 0:
		return errNoRepositories
	}

	var repos []*repo.ChartRepository
	updateAllRepos := len(o.names) == 0

	if !updateAllRepos {
		if err := checkRequestedRepos(o.names, f.Repositories); err != nil {
			return err
		}
	}

	for _, cfg := range f.Repositories {
		if updateAllRepos || isRepoRequested(cfg.Name, o.names) {
			r, err := repo.NewChartRepository(cfg, getter.All(settings))
			if err != nil {
				return err
			}
			if o.repoCache != "" {
				r.CachePath = o.repoCache
			}
			repos = append(repos, r)
		}
	}

	return o.update(repos, o.failOnRepoUpdateFail)
}

func updateCharts(repos []*repo.ChartRepository, failOnRepoUpdateFail bool) error {
	fmt.Println("Hang tight while we grab the latest from your chart repositories...")
	var wg sync.WaitGroup
	var repoFailList []string
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				fmt.Printf("...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
				repoFailList = append(repoFailList, re.Config.URL)
			} else {
				fmt.Printf("...Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()

	if len(repoFailList) > 0 && failOnRepoUpdateFail {
		return errors.New(fmt.Sprintf("Failed to update the following repositories: %s",
			repoFailList))
	}

	fmt.Println("Update Complete. ⎈Have Fun!⎈")
	return nil
}

func checkRequestedRepos(requestedRepos []string, validRepos []*repo.Entry) error {
	for _, requestedRepo := range requestedRepos {
		found := false
		for _, repo := range validRepos {
			if requestedRepo == repo.Name {
				found = true
				break
			}
		}
		if !found {
			return errors.Errorf("no repositories found matching '%s'.  Nothing will be updated", requestedRepo)
		}
	}
	return nil
}

func isRepoRequested(repoName string, requestedRepos []string) bool {
	for _, requestedRepo := range requestedRepos {
		if repoName == requestedRepo {
			return true
		}
	}
	return false
}
