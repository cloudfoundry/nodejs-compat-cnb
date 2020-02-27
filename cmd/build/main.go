package main

import (
	"github.com/cloudfoundry/nodejs-compat-cnb/compat"
	"github.com/cloudfoundry/packit"
)

func main() {
	packageJSONParser := compat.NewPackageJSONParser()

	packit.Build(compat.Build(packageJSONParser))
}
