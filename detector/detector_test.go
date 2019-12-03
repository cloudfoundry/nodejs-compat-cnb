package detector_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/nodejs-compat-cnb/compat"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/nodejs-compat-cnb/detector"

	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitDetector(t *testing.T) {
	spec.Run(t, "Detector", testDetector, spec.Report(report.Terminal{}))
}

func testDetector(t *testing.T, context spec.G, it spec.S) {
	var (
		factory *test.DetectFactory

		d detector.Detector
	)

	it.Before(func() {
		RegisterTestingT(t)

		factory = test.NewDetectFactory(t)
		err := ioutil.WriteFile(filepath.Join(factory.Detect.Application.Root, "package.json"), []byte(`{
			"name": "simple_app",
			"version": "0.0.0",
			"description": "some app",
			"main": "server.js",
			"author": "",
			"license": "",
			"repository": {
				"type": "git",
				"url": ""
			},
			"engines": {
				"node": "~10"
			}
		}`), 0644)
		Expect(err).NotTo(HaveOccurred())

		d = detector.Detector{}
	})

	it("fails", func() {
		code, err := d.RunDetect(factory.Detect)
		Expect(err).NotTo(HaveOccurred())
		Expect(code).To(Equal(detect.FailStatusCode))
	})

	context("when the package.json file is missing", func() {
		it.Before(func() {
			Expect(os.Remove(filepath.Join(factory.Detect.Application.Root, "package.json"))).To(Succeed())
		})

		it("fails", func() {
			code, err := d.RunDetect(factory.Detect)
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(detect.FailStatusCode))
		})
	})

	context("when package.json contains heroku-postbuild script", func() {
		it.Before(func() {
			err := ioutil.WriteFile(filepath.Join(factory.Detect.Application.Root, "package.json"), []byte(`{
				"name": "simple_app",
				"version": "0.0.0",
				"description": "some app",
				"main": "server.js",
				"scripts": {
					"heroku-postbuild": "echo whatever"
				},
				"author": "",
				"license": "",
				"repository": {
					"type": "git",
					"url": ""
				},
				"engines": {
					"node": "~10"
				}
			}`), 0644)
			Expect(err).NotTo(HaveOccurred())
		})

		it("passes", func() {
			code, err := d.RunDetect(factory.Detect)
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(detect.PassStatusCode))

			Expect(factory.Plans.Plan).To(Equal(buildplan.Plan{
				Requires: []buildplan.Required{
					{Name: compat.Dependency},
				},
				Provides: []buildplan.Provided{
					{Name: compat.Dependency},
				},
			}))
		})
	})

	context("when package.json contains heroku-prebuild script", func() {
		it.Before(func() {
			err := ioutil.WriteFile(filepath.Join(factory.Detect.Application.Root, "package.json"), []byte(`{
				"name": "simple_app",
				"version": "0.0.0",
				"description": "some app",
				"main": "server.js",
				"scripts": {
					"heroku-prebuild": "echo whatever"
				},
				"author": "",
				"license": "",
				"repository": {
					"type": "git",
					"url": ""
				},
				"engines": {
					"node": "~10"
				}
			}`), 0644)
			Expect(err).NotTo(HaveOccurred())
		})

		it("passes", func() {
			code, err := d.RunDetect(factory.Detect)
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(detect.PassStatusCode))

			Expect(factory.Plans.Plan).To(Equal(buildplan.Plan{
				Requires: []buildplan.Required{
					{Name: compat.Dependency},
				},
				Provides: []buildplan.Provided{
					{Name: compat.Dependency},
				},
			}))
		})
	})

	context("when $VCAP_APPLICATION is assigned", func() {
		it.Before(func() {
			Expect(os.Setenv("VCAP_APPLICATION", `{}`)).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("VCAP_APPLICATION")).To(Succeed())
		})

		it("passes", func() {
			code, err := d.RunDetect(factory.Detect)
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(detect.PassStatusCode))

			Expect(factory.Plans.Plan).To(Equal(buildplan.Plan{
				Requires: []buildplan.Required{
					{Name: compat.Dependency},
				},
				Provides: []buildplan.Provided{
					{Name: compat.Dependency},
				},
			}))
		})
	})
}
