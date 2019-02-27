package integration

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/dagger"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
	spec.Run(t, "Integration", testIntegration, spec.Report(report.Terminal{}))
}

func testIntegration(t *testing.T, when spec.G, it spec.S) {
	var (
		bp     string
		nodeBP string
		npmBP  string
	)

	it.Before(func() {
		RegisterTestingT(t)

		var err error

		err = dagger.BuildCFLinuxFS3()
		Expect(err).ToNot(HaveOccurred())

		bp, err = dagger.PackageBuildpack()
		Expect(err).ToNot(HaveOccurred())

		nodeBP, err = dagger.GetLatestBuildpack("nodejs-cnb")
		Expect(err).ToNot(HaveOccurred())

		npmBP, err = dagger.GetLatestBuildpack("npm-cnb")
		Expect(err).ToNot(HaveOccurred())
	})

	when("when heroku-postbuild and heroku-prebuild scripts are in package.json", func() {
		it("rewrites them as preinstall and postinstall scripts", func() {
			app, err := dagger.PackBuild(filepath.Join("testdata", "pre_post_commands"), bp, nodeBP, npmBP)
			Expect(err).ToNot(HaveOccurred())
			defer app.Destroy()

			Expect(app.Start()).To(Succeed())

			body, _, err := app.HTTPGet("/")
			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(Equal("Text: heroku-prebuild\npreinstall\npostinstall\nheroku-postbuild\n"))
		})
	})
}
