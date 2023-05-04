package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/drone/envsubst"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

const (
	// varsubRegex is the regular expression used to validate
	// the var names before substitution
	varsubRegex             = "^[_[:alpha:]][_[:alpha:][:digit:]]*$"
	substituteAnnotationKey = "kustomize.toolkit.fluxcd.io/substitute"
	DisabledValue           = "disabled"
)

type Substitution struct {
	Values map[string]string `yaml:"values" json:"values"`
}

func main() {
	config := new(Substitution)
	fn := func(items []*yaml.RNode) ([]*yaml.RNode, error) {
		vars := config.Values
		if vars != nil {
			r, _ := regexp.Compile(varsubRegex)
			for v := range vars {
				if !r.MatchString(v) {
					return nil, fmt.Errorf("'%s' var name is invalid, must match '%s'", v, varsubRegex)
				}
			}
			for _, item := range items {
				err := substitution(config.Values, item)
				if err != nil {
					return nil, err
				}
			}
		}
		return items, nil
	}
	p := framework.SimpleProcessor{Config: config, Filter: kio.FilterFunc(fn)}
	cmd := command.Build(p, command.StandaloneDisabled, false)
	command.AddGenerateDockerfile(cmd)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func substitution(vars map[string]string, res *yaml.RNode) error {
	if res.GetLabels()[substituteAnnotationKey] == DisabledValue || res.GetAnnotations()[substituteAnnotationKey] == DisabledValue {
		return nil
	}

	z, err := res.MarshalJSON()
	if err != nil {
		return fmt.Errorf("error converting manifest: %w", err)
	}
	// Run substitution
	output, err := envsubst.Eval(string(z), func(s string) string {
		return vars[s]
	})

	if err != nil {
		return fmt.Errorf("variable substitution failed: %w", err)
	}

	err = res.UnmarshalJSON([]byte(output))
	if err != nil {
		return fmt.Errorf("UnmarshalJSON: %w %s", err, output)
	}

	return nil
}
