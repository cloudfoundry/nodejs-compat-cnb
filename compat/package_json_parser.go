package compat

import (
	"encoding/json"
	"os"
)

type PackageJSONParser struct {
}

func NewPackageJSONParser() PackageJSONParser {
	return PackageJSONParser{}
}

func (p PackageJSONParser) Parse(path string) (bool, error) {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	var pkg struct {
		Scripts struct {
			Prebuild  string `json:"heroku-prebuild"`
			Postbuild string `json:"heroku-postbuild"`
		} `json:"scripts"`
	}

	err = json.NewDecoder(file).Decode(&pkg)
	if err != nil {
		panic(err)
	}

	if pkg.Scripts.Prebuild != "" || pkg.Scripts.Postbuild != "" {
		return true, nil
	}

	return false, nil
}
