package compat_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/nodejs-compat-cnb/compat"
	"github.com/cloudfoundry/nodejs-compat-cnb/compat/fakes"
	"github.com/cloudfoundry/packit"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		packageJSONParser *fakes.ScriptRewriter
		logEmitter        compat.LogEmitter
		workingDir        string
		layersDir         string
		build             packit.BuildFunc
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "workingDir")
		Expect(err).NotTo(HaveOccurred())

		layersDir, err = ioutil.TempDir("", "layersDir")
		Expect(err).NotTo(HaveOccurred())

		packageJSONParser = &fakes.ScriptRewriter{}
		packageJSONParser.RewriteInstallScriptsCall.Returns.Error = nil

		logEmitter = compat.NewLogEmitter(ioutil.Discard)

		build = compat.Build(packageJSONParser, logEmitter)
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
		Expect(os.RemoveAll(layersDir)).To(Succeed())
	})

	context("when VCAP_APPLICATION is not set", func() {
		it("calls the build process and does not write a profile.d script", func() {
			result, err := build(packit.BuildContext{
				WorkingDir: workingDir,
				Layers:     packit.Layers{Path: layersDir},
				Plan: packit.BuildpackPlan{
					Entries: []packit.BuildpackPlanEntry{
						{Name: "compat"},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(packit.BuildResult{
				Plan: packit.BuildpackPlan{
					Entries: []packit.BuildpackPlanEntry{
						{Name: "compat"},
					},
				},
				Layers: []packit.Layer{
					{
						Name: "compat",
						Path: filepath.Join(layersDir, "compat"),
						SharedEnv: packit.Environment{
							"NODE_MODULES_CACHE.override": "true",
							"WEB_MEMORY.override":         "512",
							"WEB_CONCURRENCY.override":    "1",
						},
						BuildEnv:  packit.Environment{},
						LaunchEnv: packit.Environment{},
						Build:     false,
						Launch:    true,
						Cache:     false,
					},
				},
			}))

			Expect(packageJSONParser.RewriteInstallScriptsCall.Receives.Path).To(Equal(filepath.Join(workingDir, "package.json")))

			Expect(filepath.Join(layersDir, "compat", "profile.d", "0_memory_available.sh")).NotTo(BeAnExistingFile())
		})
	})

	context("when VCAP_APPLICATION is set", func() {
		it.Before(func() {
			Expect(os.Setenv("VCAP_APPLICATION", "some-value")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("VCAP_APPLICATION")).To(Succeed())
		})

		it("calls the build process and writes a profile.d script", func() {
			result, err := build(packit.BuildContext{
				WorkingDir: workingDir,
				Layers:     packit.Layers{Path: layersDir},
				Plan: packit.BuildpackPlan{
					Entries: []packit.BuildpackPlanEntry{
						{Name: "compat"},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(packit.BuildResult{
				Plan: packit.BuildpackPlan{
					Entries: []packit.BuildpackPlanEntry{
						{Name: "compat"},
					},
				},
				Layers: []packit.Layer{
					{
						Name: "compat",
						Path: filepath.Join(layersDir, "compat"),
						SharedEnv: packit.Environment{
							"NODE_MODULES_CACHE.override": "true",
							"WEB_MEMORY.override":         "512",
							"WEB_CONCURRENCY.override":    "1",
						},
						BuildEnv:  packit.Environment{},
						LaunchEnv: packit.Environment{},
						Build:     false,
						Launch:    true,
						Cache:     false,
					},
				},
			}))

			Expect(packageJSONParser.RewriteInstallScriptsCall.Receives.Path).To(Equal(filepath.Join(workingDir, "package.json")))

			Expect(filepath.Join(layersDir, "compat", "profile.d", "0_memory_available.sh")).To(BeAnExistingFile())
		})
	})

	context("failure cases", func() {
		context("when you cannot get the compat layer", func() {
			it.Before(func() {
				Expect(os.Chmod(layersDir, 0000)).To(Succeed())
			})

			it.After(func() {
				Expect(os.Chmod(layersDir, os.ModePerm)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := build(packit.BuildContext{
					WorkingDir: workingDir,
					Layers:     packit.Layers{Path: layersDir},
					Plan: packit.BuildpackPlan{
						Entries: []packit.BuildpackPlanEntry{
							{Name: "compat"},
						},
					},
				})
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})
		})

		context("when the layer cannot be reset", func() {
			it.Before(func() {
				Expect(os.MkdirAll(filepath.Join(layersDir, "compat", "profile.d"), os.ModePerm)).To(Succeed())
				Expect(os.Chmod(filepath.Join(layersDir, "compat"), 0000)).To(Succeed())
			})

			it.After(func() {
				Expect(os.Chmod(filepath.Join(layersDir, "compat"), os.ModePerm)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := build(packit.BuildContext{
					WorkingDir: workingDir,
					Layers:     packit.Layers{Path: layersDir},
					Plan: packit.BuildpackPlan{
						Entries: []packit.BuildpackPlanEntry{
							{Name: "compat"},
						},
					},
				})
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})
		})

		context("when RewriteInstallScripts fails", func() {
			it.Before(func() {
				packageJSONParser.RewriteInstallScriptsCall.Returns.Error = errors.New("failed to rewrite scripts")
			})

			it("returns an error", func() {
				_, err := build(packit.BuildContext{
					WorkingDir: workingDir,
					Layers:     packit.Layers{Path: layersDir},
					Plan: packit.BuildpackPlan{
						Entries: []packit.BuildpackPlanEntry{
							{Name: "compat"},
						},
					},
				})
				Expect(err).To(MatchError(ContainSubstring("failed to rewrite scripts")))
			})
		})

	})
}
