package compat

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/packit"
	"github.com/cloudfoundry/packit/scribe"
)

const (
	MemoryAvailableScript = `if which jq > /dev/null; then
	export MEMORY_AVAILABLE="$(echo $VCAP_APPLICATION | jq -r .limits.mem)"
fi`
)

//go:generate faux --interface ScriptRewriter --output fakes/script_rewriter.go
type ScriptRewriter interface {
	RewriteInstallScripts(path string) error
}

func Build(packageJSONParser ScriptRewriter, logEmitter LogEmitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {

		logEmitter.Title(context.BuildpackInfo)

		compatLayer, err := context.Layers.Get("compat", packit.LaunchLayer)
		if err != nil {
			return packit.BuildResult{}, err
		}

		if err = compatLayer.Reset(); err != nil {
			return packit.BuildResult{}, err
		}

		err = packageJSONParser.RewriteInstallScripts(filepath.Join(context.WorkingDir, "package.json"))
		if err != nil {
			return packit.BuildResult{}, err
		}

		logEmitter.Process("Configuring environment")

		compatLayer.SharedEnv.Override("NODE_MODULES_CACHE", "true")
		compatLayer.SharedEnv.Override("WEB_MEMORY", "512")
		compatLayer.SharedEnv.Override("WEB_CONCURRENCY", "1")

		logEmitter.Subprocess("%s", scribe.NewFormattedMapFromEnvironment(compatLayer.SharedEnv))

		if _, ok := os.LookupEnv("VCAP_APPLICATION"); ok {
			profileDPath := filepath.Join(compatLayer.Path, "profile.d")

			err = os.MkdirAll(profileDPath, os.ModePerm)
			if err != nil {
				return packit.BuildResult{}, err
			}

			logEmitter.ExplainMemoryAvailable()

			err = ioutil.WriteFile(filepath.Join(profileDPath, "0_memory_available.sh"), []byte(MemoryAvailableScript), 0644)
			if err != nil {
				return packit.BuildResult{}, err
			}
		}

		return packit.BuildResult{
			Plan:   context.Plan,
			Layers: []packit.Layer{compatLayer},
		}, nil
	}
}
