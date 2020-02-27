package integration_test

import (
	"testing"
	"time"

	"github.com/cloudfoundry/dagger"
	"github.com/cloudfoundry/occam"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var nodeCompatURI string

func TestIntegration(t *testing.T) {
	var Expect = NewWithT(t).Expect

	bpDir, err := dagger.FindBPRoot()
	Expect(err).NotTo(HaveOccurred())

	nodeCompatURI, err = dagger.PackageBuildpack(bpDir)
	Expect(err).NotTo(HaveOccurred())
	defer dagger.DeleteBuildpack(nodeCompatURI)

	SetDefaultEventuallyTimeout(5 * time.Second)

	suite := spec.New("Integration", spec.Report(report.Terminal{}), spec.Parallel())
	suite("EnvVars", testEnvVars)
	suite("PackageScripts", testPackageScripts)
	suite.Run(t)
}

func ContainerLogs(id string) func() string {
	docker := occam.NewDocker()

	return func() string {
		logs, _ := docker.Container.Logs.Execute(id)
		return logs.String()
	}
}
