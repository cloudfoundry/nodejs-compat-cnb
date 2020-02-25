package compat_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/cloudfoundry/nodejs-compat-cnb/compat"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testPackageJSONParser(t *testing.T, context spec.G, it spec.S) {
	var Expect = NewWithT(t).Expect

	context("Parse", func() {
		var (
			path   string
			parser compat.PackageJSONParser
		)

		context("when there is no heroku pre build script and no heroku post build script", func() {
			it.Before(func() {
				file, err := ioutil.TempFile("", "package.json")
				Expect(err).NotTo(HaveOccurred())
				defer file.Close()

				_, err = file.WriteString(`{
				"engines": {
					"node": "1.2.3"
				}
			}`)
				Expect(err).NotTo(HaveOccurred())

				path = file.Name()

				parser = compat.NewPackageJSONParser()
			})

			it.After(func() {
				Expect(os.RemoveAll(path)).To(Succeed())
			})

			it("returns false", func() {
				scriptsExist, err := parser.Parse(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(scriptsExist).To(BeFalse())
			})

		})

		context("when heroku-prebuild exists", func() {
			it.Before(func() {
				file, err := ioutil.TempFile("", "package.json")
				Expect(err).NotTo(HaveOccurred())
				defer file.Close()

				_, err = file.WriteString(`{
				scripts": {
					"heroku-prebuild": "echo whatever"
				},
				"engines": {
					"node": "1.2.3"
				}
			}`)
				Expect(err).NotTo(HaveOccurred())

				path = file.Name()

				parser = compat.NewPackageJSONParser()
			})

			it.After(func() {
				Expect(os.RemoveAll(path)).To(Succeed())
			})

			it("returns true", func() {
				scriptsExist, err := parser.Parse(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(scriptsExist).To(BeTrue())
			})
		})

		context("when heroku-postbuild exists", func() {
			it.Before(func() {
				file, err := ioutil.TempFile("", "package.json")
				Expect(err).NotTo(HaveOccurred())
				defer file.Close()

				_, err = file.WriteString(`{
				scripts": {
					"heroku-postbuild": "echo whatever"
				},
				"engines": {
					"node": "1.2.3"
				}
			}`)
				Expect(err).NotTo(HaveOccurred())

				path = file.Name()

				parser = compat.NewPackageJSONParser()
			})

			it.After(func() {
				Expect(os.RemoveAll(path)).To(Succeed())
			})

			it("returns true", func() {
				scriptsExist, err := parser.Parse(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(scriptsExist).To(BeTrue())
			})
		})

		// context("failure cases", func() {
		// 	context("when the package.json file does not exist", func() {
		// 		it("returns an error", func() {
		// 			_, err := parser.Parse("/missing/file")
		// 			Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
		// 		})
		// 	})
		//
		// 	context("when the package.json contents are malformed", func() {
		// 		it.Before(func() {
		// 			err := ioutil.WriteFile(path, []byte("%%%"), 0644)
		// 			Expect(err).NotTo(HaveOccurred())
		// 		})
		//
		// 		it("returns an error", func() {
		// 			_, err := parser.Parse(path)
		// 			Expect(err).To(MatchError(ContainSubstring("invalid character")))
		// 		})
		// 	})
		// })
	})
}
