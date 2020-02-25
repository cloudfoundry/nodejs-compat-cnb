package integration_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/occam"
	"github.com/sclevine/spec"

	. "github.com/cloudfoundry/occam/matchers"
	. "github.com/onsi/gomega"
)

func testEnvVars(t *testing.T, context spec.G, it spec.S) {
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
				WithEnv(map[string]string{"VCAP_APPLICATION": `{"limits": {"mem": "2048"}}`}).
				Execute(name, filepath.Join("testdata", "env_vars"))
			Expect(err).NotTo(HaveOccurred())

			container, err = docker.Container.Run.
				WithEnv(map[string]string{
					"PORT":             "8080",
					"VCAP_APPLICATION": `{"limits": {"mem": "2048"}}`,
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

			Expect(body).To(ContainSubstring("MEMORY_AVAILABLE=2048"))
			Expect(body).To(ContainSubstring("NODE_MODULES_CACHE=true"))
			Expect(body).To(ContainSubstring("WEB_MEMORY=512"))
			Expect(body).To(ContainSubstring("WEB_CONCURRENCY=1"))
		})
	})
}
