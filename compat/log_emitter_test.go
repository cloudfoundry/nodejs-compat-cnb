package compat_test

import (
	"bytes"
	"testing"

	"github.com/cloudfoundry/nodejs-compat-cnb/compat"
	"github.com/cloudfoundry/packit"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testLogEmitter(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		logEmitter compat.LogEmitter
		buffer     *bytes.Buffer
	)

	it.Before(func() {
		buffer = bytes.NewBuffer(nil)
		logEmitter = compat.NewLogEmitter(buffer)
	})

	context("Title", func() {
		it("logs the buildpack title", func() {
			logEmitter.Title(packit.BuildpackInfo{
				Name:    "Node.js Compat Buildpack",
				Version: "0.0.0",
			})

			Expect(buffer).To(ContainSubstring("Node.js Compat Buildpack 0.0.0"))
		})
	})

	context("ExplainMemoryAvailable", func() {
		it("logs an explanation of the memory available script", func() {
			logEmitter.ExplainMemoryAvailable()

			Expect(buffer.String()).To(Equal(`
    Writing profile.d/0_memory_available.sh
      Calculates available memory based on memory limits declared in $VCAP_APPLICATION.
      Made available in the $MEMORY_AVAILABLE environment variable.
`))
		})
	})

	context("RewritePackageJSON", func() {
		it("logs an explanation of rewriting package.json", func() {
			logEmitter.RewritePackageJSON(true, true)

			Expect(buffer.String()).To(Equal(`  Executing build process
    Detected Heroku build scripts
      Prepending "scripts.heroku-prebuild" on "scripts.preinstall"
      Appending "scripts.heroku-postbuild" on "scripts.postinstall"
      Rewriting package.json

`))
		})

		context("when there is no prebuild", func() {
			it("omits that line", func() {
				logEmitter.RewritePackageJSON(false, true)

				Expect(buffer.String()).To(Equal(`  Executing build process
    Detected Heroku build scripts
      Appending "scripts.heroku-postbuild" on "scripts.postinstall"
      Rewriting package.json

`))
			})
		})

		context("when there is no postbuild", func() {
			it("omits that line", func() {
				logEmitter.RewritePackageJSON(true, false)

				Expect(buffer.String()).To(Equal(`  Executing build process
    Detected Heroku build scripts
      Prepending "scripts.heroku-prebuild" on "scripts.preinstall"
      Rewriting package.json

`))
			})
		})
	})
}
