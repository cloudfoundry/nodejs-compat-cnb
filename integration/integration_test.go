package integration

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/dagger"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var (
	bpDir         string
	nodeCompatURI string
)

func TestIntegration(t *testing.T) {
	var Expect = NewWithT(t).Expect

	var err error
	bpDir, err = dagger.FindBPRoot()
	Expect(err).NotTo(HaveOccurred())

	nodeCompatURI, err = dagger.PackageBuildpack(bpDir)
	Expect(err).NotTo(HaveOccurred())
	defer dagger.DeleteBuildpack(nodeCompatURI)

	spec.Run(t, "Integration", testIntegration, spec.Report(report.Terminal{}), spec.Parallel())
}

func testIntegration(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually
	)

	context("when heroku-postbuild and heroku-prebuild scripts are in package.json", func() {
		it("rewrites them as preinstall and postinstall scripts", func() {
			app, err := dagger.NewPack(
				filepath.Join("testdata", "pre_post_commands"),
				dagger.RandomImage(),
				dagger.SetBuildpacks(nodeCompatURI),
			).Build()
			Expect(err).NotTo(HaveOccurred())
			defer app.Destroy()

			Expect(app.StartWithCommand("/workspace/server.sh")).To(Succeed())

			var body string
			Eventually(func() error {
				var err error
				body, _, err = app.HTTPGet("/")
				return err
			}, "5s").Should(Succeed())

			var packageJSON struct {
				Scripts struct {
					Preinstall  string
					Postinstall string
				}
			}
			Expect(json.Unmarshal([]byte(body), &packageJSON)).To(Succeed())
			Expect(packageJSON.Scripts.Preinstall).To(Equal("heroku-prebuild && preinstall"))
			Expect(packageJSON.Scripts.Postinstall).To(Equal("postinstall && heroku-postbuild"))
		})
	})

	context("when $VCAP_APPLICATION is defined", func() {
		it("assigns $MEMORY_AVAILABLE from that JSON object", func() {
			app, err := dagger.NewPack(
				filepath.Join("testdata", "env_vars"),
				dagger.RandomImage(),
				dagger.SetBuildpacks(nodeCompatURI),
				dagger.SetEnv(map[string]string{"VCAP_APPLICATION": `{"limits": {"mem": "some-memory-limit"}}`}),
			).Build()
			Expect(err).NotTo(HaveOccurred())
			defer app.Destroy()

			app.Env["VCAP_APPLICATION"] = `{"limits": {"mem": "some-memory-limit"}}`
			Expect(app.StartWithCommand("/workspace/server.sh")).To(Succeed())

			var body string
			Eventually(func() error {
				var err error
				body, _, err = app.HTTPGet("/")
				return err
			}, "5s").Should(Succeed())
			Expect(body).To(ContainSubstring("MEMORY_AVAILABLE=some-memory-limit"))
			Expect(body).To(ContainSubstring("NODE_MODULES_CACHE=true"))
			Expect(body).To(ContainSubstring("WEB_MEMORY=512"))
			Expect(body).To(ContainSubstring("WEB_CONCURRENCY=1"))
		})
	})
}
