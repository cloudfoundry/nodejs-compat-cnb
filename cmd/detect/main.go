package main

import (
	"os"

	"github.com/cloudfoundry/nodejs-compat-cnb/compat"
	"github.com/paketo-buildpacks/packit"
)

func main() {
	logEmitter := compat.NewLogEmitter(os.Stdout)
	packageJSONParser := compat.NewPackageJSONParser(logEmitter)

	packit.Detect(compat.Detect(packageJSONParser))
}
