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
	var (
		Expect = NewWithT(t).Expect

		path   string
		parser compat.PackageJSONParser
	)

	it.Before(func() {
		file, err := ioutil.TempFile("", "package.json")
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		path = file.Name()

		parser = compat.NewPackageJSONParser()
	})

	it.After(func() {
		Expect(os.RemoveAll(path)).To(Succeed())
	})

	context("ContainsScripts", func() {

		context("when there is no heroku pre build script and no heroku post build script", func() {
			it.Before(func() {
				err := ioutil.WriteFile(path, []byte(`{
				"engines": {
					"node": "1.2.3"
				}
			}`), 0666)
				Expect(err).NotTo(HaveOccurred())

			})

			it("returns false", func() {
				scriptsExist, err := parser.ContainsScripts(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(scriptsExist).To(BeFalse())
			})

		})

		context("when heroku-prebuild exists", func() {
			it.Before(func() {
				err := ioutil.WriteFile(path, []byte(`{
				"scripts": {
					"heroku-prebuild": "echo whatever"
				},
				"engines": {
					"node": "1.2.3"
				}
			}`), 0666)
				Expect(err).NotTo(HaveOccurred())
			})

			it("returns true", func() {
				scriptsExist, err := parser.ContainsScripts(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(scriptsExist).To(BeTrue())
			})
		})

		context("when heroku-postbuild exists", func() {
			it.Before(func() {
				err := ioutil.WriteFile(path, []byte(`{
				"scripts": {
					"heroku-postbuild": "echo whatever"
				},
				"engines": {
					"node": "1.2.3"
				}
			}`), 0666)
				Expect(err).NotTo(HaveOccurred())

			})

			it("returns true", func() {
				scriptsExist, err := parser.ContainsScripts(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(scriptsExist).To(BeTrue())
			})
		})

		context("failure cases", func() {
			context("when the package.json file does not exist", func() {
				it("returns an error", func() {
					_, err := parser.ContainsScripts("/missing/file")
					Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
				})
			})

			context("when the package.json contents are malformed", func() {
				it.Before(func() {
					err := ioutil.WriteFile(path, []byte("%%%"), 0644)
					Expect(err).NotTo(HaveOccurred())
				})

				it("returns an error", func() {
					_, err := parser.ContainsScripts(path)
					Expect(err).To(MatchError(ContainSubstring("invalid character")))
				})
			})
		})
	})

	context("RewriteInstallScripts", func() {

		context("Pre Install", func() {

			context("when preinstall script exists and heroku-prebuild does not exist", func() {
				it.Before(func() {
					err := ioutil.WriteFile(path, []byte(`{
				"scripts": {
					"preinstall": "echo whatever"
				}
			}`), 0666)
					Expect(err).NotTo(HaveOccurred())
				})

				it("should still have the preinstall script in the package.json", func() {
					err := parser.RewriteInstallScripts(path)
					Expect(err).NotTo(HaveOccurred())

					content, err := ioutil.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					Expect(string(content)).To(MatchJSON(`{
"scripts": {
	"preinstall": "echo whatever"
	}
}`))
				})
			})

			context("when preinstall script does not exist and heroku-prebuild exists", func() {
				it.Before(func() {
					err := ioutil.WriteFile(path, []byte(`{
				"scripts": {
					"heroku-prebuild": "echo whatever"
				}
			}`), 0666)
					Expect(err).NotTo(HaveOccurred())
				})

				it("should still have the preinstall script in the package.json", func() {
					err := parser.RewriteInstallScripts(path)
					Expect(err).NotTo(HaveOccurred())

					content, err := ioutil.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					Expect(string(content)).To(MatchJSON(`{
"scripts": {
	"heroku-prebuild": "echo whatever",
	"preinstall": "echo whatever"
	}
}`))
				})
			})

			context("when preinstall script exists and heroku-prebuild exists", func() {
				it.Before(func() {
					err := ioutil.WriteFile(path, []byte(`{
				"scripts": {
					"preinstall": "echo preinstall",
					"heroku-prebuild": "echo whatever"
				}
			}`), 0666)
					Expect(err).NotTo(HaveOccurred())
				})

				it("should still have the preinstall script in the package.json", func() {
					err := parser.RewriteInstallScripts(path)
					Expect(err).NotTo(HaveOccurred())

					content, err := ioutil.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					Expect(string(content)).To(MatchJSON(`{
"scripts": {
	"heroku-prebuild": "echo whatever",
	"preinstall": "echo whatever && echo preinstall"
	}
}`))
				})
			})

			context("when preinstall script does not exist and heroku-prebuild does not exist", func() {
				it.Before(func() {
					err := ioutil.WriteFile(path, []byte(`{
				"scripts": {
					"some-script": "echo do something"
				}
			}`), 0666)
					Expect(err).NotTo(HaveOccurred())
				})

				it("should still have the preinstall script in the package.json", func() {
					err := parser.RewriteInstallScripts(path)
					Expect(err).NotTo(HaveOccurred())

					content, err := ioutil.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					Expect(string(content)).To(MatchJSON(`{
"scripts": {
	"some-script": "echo do something"
	}
}`))
				})
			})

		})

		context("Post Install", func() {
			context("when postinstall script exists and heroku-postbuild does not exist", func() {
				it.Before(func() {
					err := ioutil.WriteFile(path, []byte(`{
				"scripts": {
					"postinstall": "echo whatever"
				}
			}`), 0666)
					Expect(err).NotTo(HaveOccurred())
				})

				it("should still have the postinstall script in the package.json", func() {
					err := parser.RewriteInstallScripts(path)
					Expect(err).NotTo(HaveOccurred())

					content, err := ioutil.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					Expect(string(content)).To(MatchJSON(`{
"scripts": {
	"postinstall": "echo whatever"
	}
}`))
				})
			})

			context("when postinstall script does not exist and heroku-postbuild exists", func() {
				it.Before(func() {
					err := ioutil.WriteFile(path, []byte(`{
				"scripts": {
					"heroku-postbuild": "echo whatever"
				}
			}`), 0666)
					Expect(err).NotTo(HaveOccurred())
				})

				it("should still have the postinstall script in the package.json", func() {
					err := parser.RewriteInstallScripts(path)
					Expect(err).NotTo(HaveOccurred())

					content, err := ioutil.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					Expect(string(content)).To(MatchJSON(`{
"scripts": {
	"heroku-postbuild": "echo whatever",
	"postinstall": "echo whatever"
	}
}`))
				})
			})

			context("when postinstall script exists and heroku-postbuild exists", func() {
				it.Before(func() {
					err := ioutil.WriteFile(path, []byte(`{
				"scripts": {
					"postinstall": "echo postinstall",
					"heroku-postbuild": "echo whatever"
				}
			}`), 0666)
					Expect(err).NotTo(HaveOccurred())
				})

				it("should still have the postinstall script in the package.json", func() {
					err := parser.RewriteInstallScripts(path)
					Expect(err).NotTo(HaveOccurred())

					content, err := ioutil.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					Expect(string(content)).To(MatchJSON(`{
"scripts": {
	"heroku-postbuild": "echo whatever",
	"postinstall": "echo postinstall && echo whatever"
	}
}`))
				})
			})

			context("when postinstall script does not exist and heroku-postbuild does not exist", func() {
				it.Before(func() {
					err := ioutil.WriteFile(path, []byte(`{
				"scripts": {
					"some-script": "echo do something"
				}
			}`), 0666)
					Expect(err).NotTo(HaveOccurred())
				})

				it("should still have the postinstall script in the package.json", func() {
					err := parser.RewriteInstallScripts(path)
					Expect(err).NotTo(HaveOccurred())

					content, err := ioutil.ReadFile(path)
					Expect(err).NotTo(HaveOccurred())

					Expect(string(content)).To(MatchJSON(`{
"scripts": {
	"some-script": "echo do something"
	}
}`))
				})
			})
		})

		context("failure cases", func() {
			context("when the package.json file cannot be opened", func() {
				it("returns an error", func() {
					err := parser.RewriteInstallScripts("/missing/some-file")
					Expect(err).To(MatchError(ContainSubstring("no such file or directory")))
				})
			})

			context("when the json is malformed", func() {
				it.Before(func() {
					err := ioutil.WriteFile(path, []byte(`%%%`), 0666)
					Expect(err).NotTo(HaveOccurred())
				})

				it("returns an error", func() {
					err := parser.RewriteInstallScripts(path)
					Expect(err).To(MatchError(ContainSubstring("invalid character")))
				})
			})

			context("when the underlying json struct is malformed", func() {
				it.Before(func() {
					err := ioutil.WriteFile(path, []byte(`{
				"scripts" : "string"
			}`), 0666)
					Expect(err).NotTo(HaveOccurred())
				})

				it("returns an error", func() {
					err := parser.RewriteInstallScripts(path)
					Expect(err).To(MatchError(ContainSubstring("cannot unmarshal")))
				})
			})
		})
	})
}
