package compat_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/nodejs-compat-cnb/compat"
	"github.com/cloudfoundry/nodejs-compat-cnb/compat/fakes"
	"github.com/paketo-buildpacks/packit"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		packageJSONParser *fakes.PrePostParser
		workingDir        string
		detect            packit.DetectFunc
	)
	it.Before(func() {
		var err error

		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(workingDir, "package.json"), []byte{}, 0644)
		Expect(err).NotTo(HaveOccurred())

		packageJSONParser = &fakes.PrePostParser{}
		packageJSONParser.ContainsScriptsCall.Returns.ScriptsExist = false

		detect = compat.Detect(packageJSONParser)
	})

	context("when $VCAP_APPLICATION is set", func() {
		it.Before(func() {
			err := os.Setenv("VCAP_APPLICATION", "some-value")
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			err := os.Unsetenv("VCAP_APPLICATION")
			Expect(err).NotTo(HaveOccurred())
		})

		context("when heroku pre and post build scripts exist", func() {
			it.Before(func() {
				packageJSONParser.ContainsScriptsCall.Returns.ScriptsExist = true
			})

			it("passes and adds compat dependency to buildplan", func() {
				result, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Plan).To(Equal(packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{
						{Name: compat.PlanDependency},
					},
					Requires: []packit.BuildPlanRequirement{
						{Name: compat.PlanDependency},
					},
				}))
				Expect(packageJSONParser.ContainsScriptsCall.Receives.Path).To(Equal(filepath.Join(workingDir, "package.json")))
			})
		})

		context("when heroku pre and post build scripts do not exist", func() {
			it.Before(func() {
				packageJSONParser.ContainsScriptsCall.Returns.ScriptsExist = false
			})

			it("passes and adds compat dependency to buildplan", func() {
				result, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Plan).To(Equal(packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{
						{Name: compat.PlanDependency},
					},
					Requires: []packit.BuildPlanRequirement{
						{Name: compat.PlanDependency},
					},
				}))
				Expect(packageJSONParser.ContainsScriptsCall.Receives.Path).To(Equal(filepath.Join(workingDir, "package.json")))
			})
		})
	})

	context("when $VCAP_APPLICATION is not set", func() {

		context("when heroku pre and post build scripts exist", func() {
			it.Before(func() {
				packageJSONParser.ContainsScriptsCall.Returns.ScriptsExist = true
			})

			it("passes and adds compat dependency to buildplan", func() {
				result, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Plan).To(Equal(packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{
						{Name: compat.PlanDependency},
					},
					Requires: []packit.BuildPlanRequirement{
						{Name: compat.PlanDependency},
					},
				}))
				Expect(packageJSONParser.ContainsScriptsCall.Receives.Path).To(Equal(filepath.Join(workingDir, "package.json")))
			})
		})

		context("when heroku pre and post build scripts do not exist", func() {
			it.Before(func() {
				packageJSONParser.ContainsScriptsCall.Returns.ScriptsExist = false
			})

			it("fails detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(packit.Fail))
				Expect(packageJSONParser.ContainsScriptsCall.Receives.Path).To(Equal(filepath.Join(workingDir, "package.json")))
			})
		})
	})

	context("failure cases", func() {
		context("when checking for scripts in package.json fails", func() {
			it.Before(func() {
				packageJSONParser.ContainsScriptsCall.Returns.Err = errors.New("failed to check for pre and post build scripts")
			})
			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError("failed to check for pre and post build scripts"))
				Expect(packageJSONParser.ContainsScriptsCall.Receives.Path).To(Equal(filepath.Join(workingDir, "package.json")))
			})
		})
	})

}
