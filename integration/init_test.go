package integration_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cloudfoundry/dagger"
	"github.com/cloudfoundry/packit/pexec"
	"github.com/paketo-buildpacks/occam"
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

	nodeCompatURI = fmt.Sprintf("%s.tgz", nodeCompatURI)
	defer dagger.DeleteBuildpack(nodeCompatURI)

	SetDefaultEventuallyTimeout(5 * time.Second)

	suite := spec.New("Integration", spec.Report(report.Terminal{}))
	suite("EnvVars", testEnvVars)
	suite("Logging", testLogging)
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

func GetGitVersion() (string, error) {
	gitExec := pexec.NewExecutable("git")
	stdout := bytes.NewBuffer(nil)
	err := gitExec.Execute(pexec.Execution{
		Args:   []string{"describe", "--abbrev=0", "--tags"},
		Stdout: stdout,
	})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(strings.TrimPrefix(stdout.String(), "v")), nil
}
