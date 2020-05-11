package integration_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testPackageScripts(t *testing.T, context spec.G, it spec.S) {
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
}
