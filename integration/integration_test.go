package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/dagger"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var (
	bpDir, nodeCompatURI, nodeURI, npmURI string
)

func TestIntegration(t *testing.T) {
	var err error

	Expect := NewWithT(t).Expect
	bpDir, err = dagger.FindBPRoot()
	Expect(err).NotTo(HaveOccurred())
	nodeCompatURI, err = dagger.PackageBuildpack(bpDir)
	Expect(err).ToNot(HaveOccurred())
	defer os.RemoveAll(nodeCompatURI)

	nodeURI, err = dagger.GetLatestBuildpack("nodejs-cnb")
	Expect(err).ToNot(HaveOccurred())
	defer os.RemoveAll(nodeURI)

	npmURI, err = dagger.GetLatestBuildpack("npm-cnb")
	Expect(err).ToNot(HaveOccurred())
	defer os.RemoveAll(npmURI)

	spec.Run(t, "Integration", testIntegration, spec.Report(report.Terminal{}))
}

func testIntegration(t *testing.T, when spec.G, it spec.S) {
	var Expect func(interface{}, ...interface{}) GomegaAssertion
	it.Before(func() {
		Expect = NewWithT(t).Expect
	})

	when("when heroku-postbuild and heroku-prebuild scripts are in package.json", func() {
		it("rewrites them as preinstall and postinstall scripts", func() {
			app, err := dagger.PackBuild(filepath.Join("testdata", "pre_post_commands"), nodeCompatURI, nodeURI, npmURI)
			Expect(err).ToNot(HaveOccurred())
			defer app.Destroy()

			Expect(app.Start()).To(Succeed())

			body, _, err := app.HTTPGet("/")
			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(Equal("Text: heroku-prebuild\npreinstall\npostinstall\nheroku-postbuild\n"))
		})
	})
}
