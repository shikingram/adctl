package cli

import (
	"os"

	"github.com/shikingram/adctl/pkg/adctlpath"
)

type EnvSettings struct {
	RepositoryCache  string
	RepositoryConfig string
}

func New() *EnvSettings {
	env := &EnvSettings{
		RepositoryConfig: envOr("ADCTL_REPOSITORY_CONFIG", adctlpath.ConfigPath("repositories.yaml")),
		RepositoryCache:  envOr("ADCTL_REPOSITORY_CACHE", adctlpath.CachePath("repository")),
	}
	return env
}

func envOr(name, def string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return def
}
