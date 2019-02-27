module github.com/cloudfoundry/nodejs-compat-cnb

require (
	github.com/buildpack/libbuildpack v1.10.0
	github.com/cloudfoundry/dagger v0.0.0
	github.com/cloudfoundry/libcfbuildpack v1.44.0
	github.com/cloudfoundry/npm-cnb v0.0.4 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/onsi/gomega v1.4.3
	github.com/pkg/errors v0.8.1
	github.com/sclevine/spec v1.2.0
)

replace github.com/cloudfoundry/libcfbuildpack => /Users/pivotal/workspace/libcfbuildpack

replace github.com/cloudfoundry/dagger => /Users/pivotal/workspace/dagger
