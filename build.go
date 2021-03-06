package gogenerate

import (
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
)

//go:generate faux --interface BuildProcess --output fakes/build_process.go
type BuildProcess interface {
	Execute(workingDir string) error
}

func Build(buildProcess BuildProcess, logs scribe.Logger) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logs.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)
		err := buildProcess.Execute(context.WorkingDir)
		if err != nil {
			return packit.BuildResult{}, err
		}
		return packit.BuildResult{}, nil
	}
}
