package detector_test

import (
	"path/filepath"
	"testing"

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

func testDetector(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	it("fails without a package.json", func() {
		f := test.NewDetectFactory(t)

		d := detector.Detector{}
		code, err := d.RunDetect(f.Detect)
		Expect(err).NotTo(HaveOccurred())
		Expect(code).To(Equal(detect.FailStatusCode))
	})

	it("passes when package.json contains heroku-post-build script", func() {
		f := test.NewDetectFactory(t)
		test.CopyFile(t, filepath.Join("testdata", "package.json"), filepath.Join(f.Detect.Application.Root, "package.json"))
		runDetectAndExpectBuildplan(f, buildplan.BuildPlan{
			"compat": buildplan.Dependency{},
		})
	})
}

func runDetectAndExpectBuildplan(factory *test.DetectFactory, buildplan buildplan.BuildPlan) {
	d := detector.Detector{}
	code, err := d.RunDetect(factory.Detect)
	Expect(err).NotTo(HaveOccurred())

	Expect(code).To(Equal(detect.PassStatusCode))

	Expect(factory.Output).To(Equal(buildplan))
}
