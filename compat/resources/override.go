package resources

import (
	"fmt"
	"github.com/buildpack/libbuildpack/buildplan"
	"strings"
)

type OverrideDependency struct {
	CfStacks []string `json:"cf_stacks"`
	File     string   `json:"file"`
	Name     string   `json:"name"`
	Sha256   string   `json:"sha256"`
	URI      string   `json:"uri"`
	Version  string   `json:"version"`
}

type Override struct {
	Nodejs struct {
		Dependencies []OverrideDependency `json:"dependencies"`
	} `json:"nodejs"`
}

func Convert(dependency OverrideDependency) buildplan.Dependency {
	var result = buildplan.Dependency{
		Version: dependency.Version,
	}

	metadata :=  make(map[string]interface{})
	metadata["stacks"] = fmt.Sprintf("org.cloudfoundry.stacks.%s", dependency.CfStacks)
	metadata["name"] = strings.Title(dependency.Name)
	metadata["sha256"] = dependency.Sha256
	metadata["uri"] = dependency.URI

	result.Metadata = metadata

	return result
}
