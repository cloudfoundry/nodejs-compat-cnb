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
	packageFile := filepath.Join(context.Application.Root, "package.json")
	file, err := os.Open(packageFile)
	if os.IsNotExist(err) {
		return context.Fail(), nil
	} else if err != nil {
		return context.Fail(), errors.Wrap(err, "failed to open package.json")
	}
	defer file.Close()

	pkgJSON := compat.PackageJSON{}
	if err := json.NewDecoder(file).Decode(&pkgJSON); err != nil {
		return context.Fail(), err
	}

	if pkgJSON.Scripts.HerokuPreBuild == "" && pkgJSON.Scripts.HerokuPostBuild == "" {
		return context.Fail(), nil
	}

	return context.Pass(buildplan.Plan{
		Requires: []buildplan.Required {
			{Name: compat.Dependency},
		},
		Provides: []buildplan.Provided {
			{Name: compat.Dependency},
		},
	})
}
