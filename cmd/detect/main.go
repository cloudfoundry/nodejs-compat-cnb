package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/nodejs-compat-cnb/detector"
)

func main() {
	context, err := detect.DefaultDetect()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create a default detection context: %s", err)
		os.Exit(100)
	}

	d := detector.Detector{}
	code, err := d.RunDetect(context)
	if err != nil {
		context.Logger.Info(err.Error())
	}

	os.Exit(code)
}
