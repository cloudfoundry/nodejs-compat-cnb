package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/nodejs-compat-cnb/compat"

	"github.com/cloudfoundry/libcfbuildpack/build"
)

func main() {
	context, err := build.DefaultBuild()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create a default build context: %s", err)
		os.Exit(100)
	}

	code, err := runBuild(context)
	if err != nil {
		context.Logger.Info(err.Error())
	}

	os.Exit(code)
}

func runBuild(context build.Build) (int, error) {
	context.Logger.FirstLine(context.Logger.PrettyIdentity(context.Buildpack))

	compatContributor, willContribute, err := compat.NewContributor(context)
	if err != nil {
		return context.Failure(102), err
	}

	if willContribute {
		if err := compatContributor.Contribute(); err != nil {
			return context.Failure(103), err
		}
	}

	return context.Success(context.BuildPlan)
}
