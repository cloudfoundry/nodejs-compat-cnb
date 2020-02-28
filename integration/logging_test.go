package integration_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/occam"
	"github.com/sclevine/spec"

	. "github.com/cloudfoundry/occam/matchers"
	. "github.com/onsi/gomega"
)

func testLogging(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		pack   occam.Pack
		docker occam.Docker
		image  occam.Image

		name             string
		buildpackVersion string
	)

	it.Before(func() {
		pack = occam.NewPack().WithNoColor()
		docker = occam.NewDocker()

		var err error
		name, err = occam.RandomName()
		Expect(err).NotTo(HaveOccurred())

		buildpackVersion, err = GetGitVersion()
		Expect(err).ToNot(HaveOccurred())
	})

	it.After(func() {
		Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
		Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
	})

	context("building an image", func() {
		it("log output is consistent with CNB style", func() {
			var err error
			var buildLogs fmt.Stringer

			image, buildLogs, err = pack.Build.
				WithBuildpacks(nodeCompatURI).
				WithNoPull().
				WithEnv(map[string]string{"VCAP_APPLICATION": `{"limits": {"mem": "1024"}}`}).
				Execute(name, filepath.Join("testdata", "env_vars"))
			Expect(err).NotTo(HaveOccurred())

			Expect(buildLogs).To(ContainLines(
				fmt.Sprintf("Node.js Compat Buildpack %s", buildpackVersion),
				"  Configuring environment",
				`    NODE_MODULES_CACHE -> "true"`,
				`    WEB_CONCURRENCY    -> "1"`,
				`    WEB_MEMORY         -> "512"`,
				"",
				"    Writing profile.d/0_memory_available.sh",
				"      Calculates available memory based on memory limits declared in $VCAP_APPLICATION.",
				"      Made available in the $MEMORY_AVAILABLE environment variable.",
			))
		})
	})

	context("when package.json needs to be rewritten", func() {
		it("log output is consistent with CNB style", func() {
			var err error
			var buildLogs fmt.Stringer

			image, buildLogs, err = pack.Build.
				WithBuildpacks(nodeCompatURI).
				WithNoPull().
				Execute(name, filepath.Join("testdata", "pre_post_commands"))
			Expect(err).NotTo(HaveOccurred())

			Expect(buildLogs).To(ContainLines(
				fmt.Sprintf("Node.js Compat Buildpack %s", buildpackVersion),
				"  Executing build process",
				"    Detected Heroku build scripts",
				"      Prepending \"scripts.heroku-prebuild\" on \"scripts.preinstall\"",
				"      Appending \"scripts.heroku-postbuild\" on \"scripts.postinstall\"",
				"      Rewriting package.json",
				"",
				"  Configuring environment",
				`    NODE_MODULES_CACHE -> "true"`,
				`    WEB_CONCURRENCY    -> "1"`,
				`    WEB_MEMORY         -> "512"`,
			))
		})
	})
}
