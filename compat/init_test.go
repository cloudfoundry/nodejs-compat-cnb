package compat_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitCompat(t *testing.T) {
	suite := spec.New("compat", spec.Report(report.Terminal{}))
	suite("Build", testBuild)
	suite("Detect", testDetect)
	suite("LogEmitter", testLogEmitter)
	suite("PackageJSONParser", testPackageJSONParser)
	suite.Run(t)

}
