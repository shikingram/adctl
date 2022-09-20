package cli

import (
	"os"

	"github.com/spf13/pflag"
)

type EnvSettings struct {
	namespace string
}

func New() *EnvSettings {
	env := &EnvSettings{
		namespace: os.Getenv("HELM_NAMESPACE"),
	}
	return env
}

func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&s.namespace, "namespace", "n", s.namespace, "namespace scope for this request")
}
