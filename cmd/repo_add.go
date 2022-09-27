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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"github.com/shikingram/adctl/pkg/getter"
	"golang.org/x/term"
	"sigs.k8s.io/yaml"

	"github.com/pkg/errors"
	"github.com/shikingram/adctl/cmd/require"
	"github.com/shikingram/adctl/pkg/repo"
	"github.com/spf13/cobra"
)

type repoAddOptions struct {
	name                 string
	url                  string
	username             string
	password             string
	passwordFromStdinOpt bool
	passCredentialsAll   bool
	forceUpdate          bool

	certFile              string
	keyFile               string
	caFile                string
	insecureSkipTLSverify bool

	repoFile  string
	repoCache string
}

func newRepoAddCmd() *cobra.Command {
	o := &repoAddOptions{}

	cmd := &cobra.Command{
		Use:   "add [NAME] [URL]",
		Short: "add a chart repository",
		Args:  require.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.name = args[0]
			o.url = args[1]
			o.repoFile = settings.RepositoryConfig
			o.repoCache = settings.RepositoryCache

			return o.run()
		},
	}
	f := cmd.Flags()
	f.StringVar(&o.username, "username", "", "chart repository username")
	f.StringVar(&o.password, "password", "", "chart repository password")
	return cmd
}

func (o *repoAddOptions) run() error {

	err := os.MkdirAll(filepath.Dir(o.repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Acquire a file lock for process synchronization
	repoFileExt := filepath.Ext(o.repoFile)
	var lockPath string
	if len(repoFileExt) > 0 && len(repoFileExt) < len(o.repoFile) {
		lockPath = strings.Replace(o.repoFile, repoFileExt, ".lock", 1)
	} else {
		lockPath = o.repoFile + ".lock"
	}
	fileLock := flock.New(lockPath)
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(o.repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	if o.username != "" && o.password == "" {
		if o.passwordFromStdinOpt {
			passwordFromStdin, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			password := strings.TrimSuffix(string(passwordFromStdin), "\n")
			password = strings.TrimSuffix(password, "\r")
			o.password = password
		} else {
			fd := int(os.Stdin.Fd())
			password, err := term.ReadPassword(fd)
			if err != nil {
				return err
			}
			o.password = string(password)
		}
	}

	c := repo.Entry{
		Name:                  o.name,
		URL:                   o.url,
		Username:              o.username,
		Password:              o.password,
		PassCredentialsAll:    o.passCredentialsAll,
		CertFile:              o.certFile,
		KeyFile:               o.keyFile,
		CAFile:                o.caFile,
		InsecureSkipTLSverify: o.insecureSkipTLSverify,
	}

	// If the repo exists do one of two things:
	// 1. If the configuration for the name is the same continue without error
	// 2. When the config is different require --force-update
	if !o.forceUpdate && f.Has(o.name) {
		existing := f.Get(o.name)
		if c != *existing {

			// The input coming in for the name is different from what is already
			// configured. Return an error.
			return errors.Errorf("repository name (%s) already exists, please specify a different name", o.name)
		}

		// The add is idempotent so do nothing
		fmt.Printf("%q already exists with the same configuration, skipping\n", o.name)
		return nil
	}

	r, err := repo.NewChartRepository(&c, getter.All(settings))
	if err != nil {
		return err
	}

	if o.repoCache != "" {
		r.CachePath = o.repoCache
	}
	if _, err := r.DownloadIndexFile(); err != nil {
		return errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", o.url)
	}

	f.Update(&c)

	if err := f.WriteFile(o.repoFile, 0644); err != nil {
		return err
	}
	fmt.Printf("%q has been added to your repositories\n", o.name)
	return nil
}
