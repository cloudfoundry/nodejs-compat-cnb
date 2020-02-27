package compat

import (
	"encoding/json"
	"fmt"
	"os"
)

type PackageJSONParser struct {
}

func NewPackageJSONParser() PackageJSONParser {
	return PackageJSONParser{}
}

func (p PackageJSONParser) ContainsScripts(path string) (bool, error) {

	file, err := os.Open(path)
	if err != nil {
		return false, err
	}

	var pkg struct {
		Scripts struct {
			Prebuild  string `json:"heroku-prebuild"`
			Postbuild string `json:"heroku-postbuild"`
		} `json:"scripts"`
	}

	err = json.NewDecoder(file).Decode(&pkg)
	if err != nil {
		return false, err
	}

	return pkg.Scripts.Prebuild != "" || pkg.Scripts.Postbuild != "", nil
}

func (p PackageJSONParser) RewriteInstallScripts(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var pkg map[string]json.RawMessage
	err = json.NewDecoder(file).Decode(&pkg)
	if err != nil {
		return err
	}

	s, ok := pkg["scripts"]
	if !ok {
		return nil
	}

	var scripts map[string]string
	err = json.Unmarshal(s, &scripts)
	if err != nil {
		return err
	}

	prebuildScript, prebuildSet := scripts["heroku-prebuild"]
	postbuildScript, postbuildSet := scripts["heroku-postbuild"]
	if !prebuildSet && !postbuildSet {
		return nil
	}

	if prebuildSet {
		script, ok := scripts["preinstall"]
		if ok {
			scripts["preinstall"] = fmt.Sprintf("%s && %s", prebuildScript, script)
		} else {
			scripts["preinstall"] = prebuildScript
		}
	}

	if postbuildSet {
		script, ok := scripts["postinstall"]
		if ok {
			scripts["postinstall"] = fmt.Sprintf("%s && %s", script, postbuildScript)
		} else {
			scripts["postinstall"] = postbuildScript
		}
	}

	pkg["scripts"], err = json.Marshal(scripts)
	if err != nil {
		return err
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(&pkg)
	if err != nil {
		return err
	}

	return nil
}
