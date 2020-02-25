package compat_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/nodejs-compat-cnb/compat"
	"github.com/cloudfoundry/nodejs-compat-cnb/compat/fakes"
	"github.com/cloudfoundry/packit"
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
	//WHEN PACKAGE.JSON EXISTS BUT CONTAINS NO PRE-BUILD, POST-BUILD and $VCAP_APPLICATION is not found
	it.Before(func() {
		var err error

		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(workingDir, "package.json"), []byte{}, 0644)
		Expect(err).NotTo(HaveOccurred())

		packageJSONParser = &fakes.PrePostParser{}
		packageJSONParser.ParseCall.Returns.ScriptsExist = false

		detect = compat.Detect(packageJSONParser)
	})

	it("fails", func() {

		_, err := detect(packit.DetectContext{
			WorkingDir: workingDir,
		})
		Expect(err).To(MatchError(packit.Fail))
	})

}
