package compat

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpack/libbuildpack/application"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/mitchellh/mapstructure"
)

const (
	Dependency = "compat"

	MemoryAvailableScript = `if which jq > /dev/null; then
	export MEMORY_AVAILABLE="$(echo $VCAP_APPLICATION | jq .limits.mem)"
fi`
)

type Contributor struct {
	app   application.Application
	layer layers.Layer
}

type Scripts struct {
	PreInstall      string `mapstructure:"preinstall"`
	PostInstall     string `mapstructure:"postinstall"`
	HerokuPreBuild  string `mapstructure:"heroku-prebuild" json:"heroku-prebuild"`
	HerokuPostBuild string `mapstructure:"heroku-postbuild" json:"heroku-postbuild"`
}

type PackageJSON struct {
	Scripts Scripts `json:"scripts"`
}

func NewContributor(context build.Build) (Contributor, bool, error) {
	if !context.Plans.Has(Dependency) {
		return Contributor{}, false, nil
	}

	return Contributor{
		app:   context.Application,
		layer: context.Layers.Layer(Dependency),
	}, true, nil
}

func (c Contributor) Contribute() error {
	packagePath := filepath.Join(c.app.Root, "package.json")
	if exists, err := helper.FileExists(packagePath); err != nil {
		return err
	} else if !exists {
		return errors.New("package.json does not exist")
	}

	file, err := os.OpenFile(packagePath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	pkgJSON := map[string]interface{}{}
	if err := json.NewDecoder(file).Decode(&pkgJSON); err != nil {
		return err
	}

	scriptsMap, ok := pkgJSON["scripts"].(map[string]interface{})
	if !ok {
		scriptsMap = map[string]interface{}{}
	}

	var scripts Scripts
	if err := mapstructure.Decode(scriptsMap, &scripts); err != nil {
		return err
	}

	if scripts.HerokuPreBuild != "" {
		final := scripts.HerokuPreBuild
		if scripts.PreInstall != "" {
			final = fmt.Sprintf("%s && %s", final, scripts.PreInstall)
		}

		scriptsMap["preinstall"] = final
	}

	if scripts.HerokuPostBuild != "" {
		final := scripts.HerokuPostBuild
		if scripts.PostInstall != "" {
			final = fmt.Sprintf("%s && %s", scripts.PostInstall, final)
		}

		scriptsMap["postinstall"] = final
	}

	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	enc := json.NewEncoder(file)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(pkgJSON); err != nil {
		return err
	}

	if _, ok := os.LookupEnv("VCAP_APPLICATION"); ok {
		if err := c.layer.WriteProfile("0_memory_available.sh", MemoryAvailableScript); err != nil {
			return err
		}

		if err := c.layer.WriteMetadata(nil, layers.Launch); err != nil {
			return err
		}
	}

	return nil
}
