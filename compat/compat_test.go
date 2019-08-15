package compat_test

import (
	"github.com/cloudfoundry/libcfbuildpack/buildpackplan"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/nodejs-compat-cnb/compat"

	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitCompat(t *testing.T) {
	spec.Run(t, "Compat", testCompat, spec.Report(report.Terminal{}))
}

func testCompat(t *testing.T, when spec.G, it spec.S) {
	var (
		factory *test.BuildFactory
	)

	it.Before(func() {
		RegisterTestingT(t)
		factory = test.NewBuildFactory(t)
	})

	when("NewContributor", func() {
		it("returns true if a build plan exists with the dep", func() {
			factory.AddPlan(buildpackplan.Plan{Name: compat.Dependency})

			_, willContribute, err := compat.NewContributor(factory.Build)
			Expect(err).NotTo(HaveOccurred())
			Expect(willContribute).To(BeTrue())
		})
	})

	when("Contribute", func() {
		var (
			contributor    compat.Contributor
			willContribute bool
			appRoot        string
			err            error
		)
		it.Before(func() {
			appRoot = factory.Build.Application.Root
			factory.AddPlan(buildpackplan.Plan{Name: compat.Dependency})

			test.CopyFile(t, filepath.Join("testdata", "package.json"), filepath.Join(appRoot, "package.json"))

			contributor, willContribute, err = compat.NewContributor(factory.Build)
			Expect(err).NotTo(HaveOccurred())
			Expect(willContribute).To(BeTrue())
		})

		it("rewrites package.json", func() {
			Expect(contributor.Contribute()).To(Succeed())

			contents, err := ioutil.ReadFile(filepath.Join(appRoot, "package.json"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(contents)).To(ContainSubstring(`"preinstall":"heroku-prebuild && preinstall"`))
			Expect(string(contents)).To(ContainSubstring(`"postinstall":"postinstall && heroku-postbuild"`))
		})
	})
}
