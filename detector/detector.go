package detector

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/nodejs-compat-cnb/compat"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/pkg/errors"
)

type Detector struct{}

func (d *Detector) RunDetect(context detect.Detect) (int, error) {
	file, err := os.Open(filepath.Join(context.Application.Root, "package.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return context.Fail(), nil
		}

		return context.Fail(), errors.Wrap(err, "failed to open package.json")
	}
	defer file.Close()

	var packageJSON compat.PackageJSON
	if err := json.NewDecoder(file).Decode(&packageJSON); err != nil {
		return context.Fail(), err
	}

	_, ok := os.LookupEnv("VCAP_APPLICATION")
	if packageJSON.Scripts.HerokuPreBuild == "" && packageJSON.Scripts.HerokuPostBuild == "" && !ok {
		return context.Fail(), nil
	}

	return context.Pass(buildplan.Plan{
		Requires: []buildplan.Required{
			{Name: compat.Dependency},
		},
		Provides: []buildplan.Provided{
			{Name: compat.Dependency},
		},
	})
}
