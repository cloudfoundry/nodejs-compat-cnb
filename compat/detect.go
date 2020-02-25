package compat

import "github.com/cloudfoundry/packit"

const (
	PlanDependency = "compat"
)

//go:generate faux --interface PrePostParser --output fakes/prepost_parser.go
type PrePostParser interface {
	Parse(path string) (scriptsExist bool, err error)
}

func Detect(packageJSONParser PrePostParser) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		return packit.DetectResult{}, packit.Fail
	}
}
