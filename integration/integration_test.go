package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/cloudfoundry/dagger"
	"github.com/cloudfoundry/occam"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/cloudfoundry/occam/matchers"
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

	SetDefaultEventuallyTimeout(5 * time.Second)

	spec.Run(t, "Integration", testIntegration, spec.Report(report.Terminal{}), spec.Parallel())
}

func ContainerLogs(id string) func() string {
	docker := occam.NewDocker()

	return func() string {
		logs, _ := docker.Container.Logs.Execute(id)
		return logs.String()
	}
}

func testIntegration(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	context("when heroku-postbuild and heroku-prebuild scripts are in package.json", func() {
		var (
			image     occam.Image
			container occam.Container

			name string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
		})

		it("rewrites them as preinstall and postinstall scripts", func() {
			var err error
			image, _, err = pack.Build.
				WithBuildpacks(nodeCompatURI).
				WithNoPull().
				Execute(name, filepath.Join("testdata", "pre_post_commands"))
			Expect(err).NotTo(HaveOccurred())

			container, err = docker.Container.Run.
				WithCommand("/workspace/server.sh").
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(BeAvailable(), ContainerLogs(container.ID))

			var content []byte
			Eventually(func() error {
				response, err := http.Get(fmt.Sprintf("http://localhost:%s", container.HostPort()))
				if err != nil {
					return err
				}
				defer response.Body.Close()

				content, err = ioutil.ReadAll(response.Body)
				if err != nil {
					return err
				}

				return nil

			}).Should(Succeed())

			var packageJSON struct {
				Scripts struct {
					Preinstall  string
					Postinstall string
				}
			}
			Expect(json.Unmarshal(content, &packageJSON)).To(Succeed())
			Expect(packageJSON.Scripts.Preinstall).To(Equal("heroku-prebuild && preinstall"))
			Expect(packageJSON.Scripts.Postinstall).To(Equal("postinstall && heroku-postbuild"))
		})
	})

	context("when $VCAP_APPLICATION is defined", func() {
		var (
			image     occam.Image
			container occam.Container

			name string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
		})

		it("assigns $MEMORY_AVAILABLE from that JSON object", func() {
			var err error
			image, _, err = pack.Build.
				WithBuildpacks(nodeCompatURI).
				WithNoPull().
				WithEnv(map[string]string{"VCAP_APPLICATION": `{"limits": {"mem": "some-memory-limit"}}`}).
				Execute(name, filepath.Join("testdata", "env_vars"))
			Expect(err).NotTo(HaveOccurred())

			container, err = docker.Container.Run.
				WithEnv(map[string]string{
					"PORT":             "8080",
					"VCAP_APPLICATION": `{"limits": {"mem": "some-memory-limit"}}`,
				}).
				WithCommand("/workspace/server.sh").
				Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			Eventually(container).Should(BeAvailable(), ContainerLogs(container.ID))

			var body string
			Eventually(func() error {
				response, err := http.Get(fmt.Sprintf("http://localhost:%s", container.HostPort()))
				if err != nil {
					return err
				}
				defer response.Body.Close()

				content, err := ioutil.ReadAll(response.Body)
				if err != nil {
					return err
				}

				body = string(content)
				return nil
			}).Should(Succeed())

			Expect(body).To(ContainSubstring("MEMORY_AVAILABLE=some-memory-limit"))
			Expect(body).To(ContainSubstring("NODE_MODULES_CACHE=true"))
			Expect(body).To(ContainSubstring("WEB_MEMORY=512"))
			Expect(body).To(ContainSubstring("WEB_CONCURRENCY=1"))
		})
	})
}
