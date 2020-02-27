package compat

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/packit"
)

var (
	MemoryAvailableScript = `if which jq > /dev/null; then
	export MEMORY_AVAILABLE="$(echo $VCAP_APPLICATION | jq -r .limits.mem)"
fi`
)

//go:generate faux --interface ScriptRewriter --output fakes/script_rewriter.go
type ScriptRewriter interface {
	RewriteInstallScripts(path string) error
}

func Build(packageJSONParser ScriptRewriter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {

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

		compatLayer.SharedEnv.Override("NODE_MODULES_CACHE", "true")
		compatLayer.SharedEnv.Override("WEB_MEMORY", "512")
		compatLayer.SharedEnv.Override("WEB_CONCURRENCY", "1")

		if _, ok := os.LookupEnv("VCAP_APPLICATION"); ok {
			profileDPath := filepath.Join(compatLayer.Path, "profile.d")

			err = os.MkdirAll(profileDPath, os.ModePerm)
			if err != nil {
				return packit.BuildResult{}, err
			}

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
