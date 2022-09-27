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
	"io"
	"os"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/pkg/errors"
	"github.com/shikingram/adctl/cmd/require"
	"github.com/shikingram/adctl/pkg/cli/output"
	"github.com/shikingram/adctl/pkg/repo"
	"github.com/spf13/cobra"
)

func newRepoListCmd() *cobra.Command {
	var outfmt output.Format
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list chart repositories",
		Args:    require.NoArgs,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return compListRepos(toComplete, args), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := repo.LoadFile(settings.RepositoryConfig)
			if isNotExist(err) || (len(f.Repositories) == 0) {
				return errors.New("no repositories to show")
			}

			return outfmt.Write(os.Stdout, &repoListWriter{f.Repositories})
		},
	}

	bindOutputFlag(cmd, &outfmt)

	return cmd
}

type repositoryElement struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type repoListWriter struct {
	repos []*repo.Entry
}

func (r *repoListWriter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("NAME", "URL")
	for _, re := range r.repos {
		table.AddRow(re.Name, re.URL)
	}
	return output.EncodeTable(out, table)
}

func (r *repoListWriter) WriteJSON(out io.Writer) error {
	return r.encodeByFormat(out, output.JSON)
}

func (r *repoListWriter) WriteYAML(out io.Writer) error {
	return r.encodeByFormat(out, output.YAML)
}

func (r *repoListWriter) encodeByFormat(out io.Writer, format output.Format) error {
	// Initialize the array so no results returns an empty array instead of null
	repolist := make([]repositoryElement, 0, len(r.repos))

	for _, re := range r.repos {
		repolist = append(repolist, repositoryElement{Name: re.Name, URL: re.URL})
	}

	switch format {
	case output.JSON:
		return output.EncodeJSON(out, repolist)
	case output.YAML:
		return output.EncodeYAML(out, repolist)
	}

	// Because this is a non-exported function and only called internally by
	// WriteJSON and WriteYAML, we shouldn't get invalid types
	return nil
}

// Returns all repos from repos, except those with names matching ignoredRepoNames
// Inspired by https://stackoverflow.com/a/28701031/893211
func filterRepos(repos []*repo.Entry, ignoredRepoNames []string) []*repo.Entry {
	// if ignoredRepoNames is nil, just return repo
	if ignoredRepoNames == nil {
		return repos
	}

	filteredRepos := make([]*repo.Entry, 0)

	ignored := make(map[string]bool, len(ignoredRepoNames))
	for _, repoName := range ignoredRepoNames {
		ignored[repoName] = true
	}

	for _, repo := range repos {
		if _, removed := ignored[repo.Name]; !removed {
			filteredRepos = append(filteredRepos, repo)
		}
	}

	return filteredRepos
}

// Provide dynamic auto-completion for repo names
func compListRepos(prefix string, ignoredRepoNames []string) []string {
	var rNames []string

	f, err := repo.LoadFile(settings.RepositoryConfig)
	if err == nil && len(f.Repositories) > 0 {
		filteredRepos := filterRepos(f.Repositories, ignoredRepoNames)
		for _, repo := range filteredRepos {
			if strings.HasPrefix(repo.Name, prefix) {
				rNames = append(rNames, fmt.Sprintf("%s\t%s", repo.Name, repo.URL))
			}
		}
	}
	return rNames
}
