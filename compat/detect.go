package compat

import (
	"os"
	"path/filepath"

	"github.com/cloudfoundry/packit"
)

const (
	PlanDependency = "compat"
)

//go:generate faux --interface PrePostParser --output fakes/prepost_parser.go
type PrePostParser interface {
	ContainsScripts(path string) (scriptsExist bool, err error)
}

func Detect(packageJSONParser PrePostParser) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {

		hasHerokuScripts, err := packageJSONParser.ContainsScripts(filepath.Join(context.WorkingDir, "package.json"))
		if err != nil {
			return packit.DetectResult{}, err
		}

		_, set := os.LookupEnv("VCAP_APPLICATION")
		if !hasHerokuScripts && !set {
			return packit.DetectResult{}, packit.Fail
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: PlanDependency},
				},
				Requires: []packit.BuildPlanRequirement{
					{Name: PlanDependency},
				},
			},
		}, nil
	}
}
