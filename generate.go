package gogenerate

import (
	"bytes"
	"strings"
	"time"

	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/pexec"
)

//go:generate faux --interface Executable --output fakes/executable.go
type Executable interface {
	Execute(pexec.Execution) error
}

type Generate struct {
	executable Executable
	logs       LogEmitter
	clock      chronos.Clock
}

func NewGenerate(executable Executable, logs LogEmitter, clock chronos.Clock) Generate {
	return Generate{
		executable: executable,
		logs:       logs,
		clock:      clock,
	}
}

func (g Generate) Execute(workingDir string) error {
	buffer := bytes.NewBuffer(nil)
	args := []string{"generate", "./..."}

	g.logs.Process("Executing build process")
	g.logs.Subprocess("Running 'go %s'", strings.Join(args, " "))

	duration, err := g.clock.Measure(func() error {
		return g.executable.Execute(pexec.Execution{
			Args:   args,
			Dir:    workingDir,
			Stdout: buffer,
			Stderr: buffer,
		})
	})
	if err != nil {
		g.logs.Action("Failed after %s", duration.Round(time.Millisecond))
		g.logs.Detail(buffer.String())

		return err
	}

	g.logs.Action("Completed in %s", duration.Round(time.Millisecond))
	g.logs.Break()

	return nil
}