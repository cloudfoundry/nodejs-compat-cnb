package compat

import (
	"io"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
)

type LogEmitter struct {
	scribe.Logger
}

func NewLogEmitter(buffer io.Writer) LogEmitter {
	return LogEmitter{
		Logger: scribe.NewLogger(buffer),
	}
}

func (l LogEmitter) Title(bpInfo packit.BuildpackInfo) {
	l.Logger.Title("%s %s", bpInfo.Name, bpInfo.Version)
}

func (l LogEmitter) ExplainMemoryAvailable() {
	l.Logger.Break()
	l.Logger.Subprocess("Writing profile.d/0_memory_available.sh")
	l.Logger.Action("Calculates available memory based on memory limits declared in $VCAP_APPLICATION.")
	l.Logger.Action("Made available in the $MEMORY_AVAILABLE environment variable.")
}

func (l LogEmitter) RewritePackageJSON(pre, post bool) {
	l.Logger.Process("Executing build process")
	l.Logger.Subprocess("Detected Heroku build scripts")

	if pre {
		l.Logger.Action("Prepending \"scripts.heroku-prebuild\" on \"scripts.preinstall\"")
	}

	if post {
		l.Logger.Action("Appending \"scripts.heroku-postbuild\" on \"scripts.postinstall\"")
	}

	l.Logger.Action("Rewriting package.json")
	l.Logger.Break()
}
