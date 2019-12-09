package compat_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/buildpackplan"

	"github.com/cloudfoundry/nodejs-compat-cnb/compat"

	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitCompat(t *testing.T) {
	spec.Run(t, "Compat", testCompat, spec.Report(report.Terminal{}))
}

func testCompat(t *testing.T, context spec.G, it spec.S) {
	var factory *test.BuildFactory

	it.Before(func() {
		RegisterTestingT(t)
		factory = test.NewBuildFactory(t)
	})

	context("NewContributor", func() {
		context("when a buildplan exists with the dependency", func() {
			it("indicates that it will contribute", func() {
				factory.AddPlan(buildpackplan.Plan{Name: compat.Dependency})

				_, willContribute, err := compat.NewContributor(factory.Build)
				Expect(err).NotTo(HaveOccurred())
				Expect(willContribute).To(BeTrue())
			})
		})

		context("when a buildplan does not exist with the dependency", func() {
			it("indicates that it will not contribute", func() {
				factory.AddPlan(buildpackplan.Plan{})

				_, willContribute, err := compat.NewContributor(factory.Build)
				Expect(err).NotTo(HaveOccurred())
				Expect(willContribute).To(BeFalse())
			})
		})
	})

	context("Contribute", func() {
		var contributor compat.Contributor

		it.Before(func() {
			factory.AddPlan(buildpackplan.Plan{Name: compat.Dependency})

			var (
				err            error
				willContribute bool
			)

			contributor, willContribute, err = compat.NewContributor(factory.Build)
			Expect(err).NotTo(HaveOccurred())
			Expect(willContribute).To(BeTrue())

			err = ioutil.WriteFile(filepath.Join(factory.Build.Application.Root, "package.json"), []byte(`{
					"name": "simple_app",
					"version": "0.0.0",
					"description": "hello, world",
					"main": "server.js",
					"engines": {
						"node": "8.x"
					},
					"author": "",
					"license": "BSD-2-Clause",
					"repository": {
						"type" : "git",
						"url" : "http://github.com/cloudfoundry/nodejs-buildpack.git"
					}
				}`), 0644)
			Expect(err).NotTo(HaveOccurred())
		})

		it("leaves the package.json unmodified", func() {
			Expect(contributor.Contribute()).To(Succeed())

			contents, err := ioutil.ReadFile(filepath.Join(factory.Build.Application.Root, "package.json"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(contents)).To(MatchJSON(`{
				"name": "simple_app",
				"version": "0.0.0",
				"description": "hello, world",
				"main": "server.js",
				"engines": {
					"node": "8.x"
				},
				"author": "",
				"license": "BSD-2-Clause",
				"repository": {
					"type" : "git",
					"url" : "http://github.com/cloudfoundry/nodejs-buildpack.git"
				}
			}`))
		})

		it("does not write a 0_memory_available.sh profile.d script", func() {
			Expect(contributor.Contribute()).To(Succeed())

			layer := factory.Build.Layers.Layer(compat.Dependency)
			Expect(filepath.Join(layer.Root, "profile.d", "0_memory_available.sh")).NotTo(BeARegularFile())
		})

		it("sets heroku-specific environment variables", func() {
			Expect(contributor.Contribute()).To(Succeed())

			layer := factory.Build.Layers.Layer(compat.Dependency)
			Expect(layer).To(test.HaveOverrideSharedEnvironment("NODE_MODULES_CACHE", "true"))
			Expect(layer).To(test.HaveOverrideSharedEnvironment("WEB_MEMORY", "512"))
			Expect(layer).To(test.HaveOverrideSharedEnvironment("WEB_CONCURRENCY", "1"))
		})

		context("when the package.json contains heroku build hooks", func() {
			it.Before(func() {
				err := ioutil.WriteFile(filepath.Join(factory.Build.Application.Root, "package.json"), []byte(`{
					"name": "simple_app",
					"version": "0.0.0",
					"description": "hello, world",
					"main": "server.js",
					"engines": {
						"node": "8.x"
					},
					"scripts": {
						"heroku-prebuild": "heroku-prebuild",
						"preinstall": "preinstall",
						"postinstall": "postinstall",
						"heroku-postbuild": "heroku-postbuild"
					},
					"author": "",
					"license": "BSD-2-Clause",
					"repository": {
						"type" : "git",
						"url" : "http://github.com/cloudfoundry/nodejs-buildpack.git"
					}
				}`), 0644)
				Expect(err).NotTo(HaveOccurred())
			})

			it("rewrites package.json", func() {
				Expect(contributor.Contribute()).To(Succeed())

				contents, err := ioutil.ReadFile(filepath.Join(factory.Build.Application.Root, "package.json"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(contents)).To(MatchJSON(`{
					"name": "simple_app",
					"version": "0.0.0",
					"description": "hello, world",
					"main": "server.js",
					"engines": {
						"node": "8.x"
					},
					"scripts": {
						"heroku-prebuild": "heroku-prebuild",
						"preinstall": "heroku-prebuild && preinstall",
						"postinstall": "postinstall && heroku-postbuild",
						"heroku-postbuild": "heroku-postbuild"
					},
					"author": "",
					"license": "BSD-2-Clause",
					"repository": {
						"type" : "git",
						"url" : "http://github.com/cloudfoundry/nodejs-buildpack.git"
					}
				}`))
			})
		})

		context("when $VCAP_APPLICATION is assigned", func() {
			it.Before(func() {
				Expect(os.Setenv("VCAP_APPLICATION", `{}`)).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("VCAP_APPLICATION")).To(Succeed())
			})

			it("writes a 0_memory_available.sh profile.d script", func() {
				Expect(contributor.Contribute()).To(Succeed())

				layer := factory.Build.Layers.Layer(compat.Dependency)
				contents, err := ioutil.ReadFile(filepath.Join(layer.Root, "profile.d", "0_memory_available.sh"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(contents)).To(Equal(`if which jq > /dev/null; then
	export MEMORY_AVAILABLE="$(echo $VCAP_APPLICATION | jq -r .limits.mem)"
fi`))
			})
		})
	})
}
