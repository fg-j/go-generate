package gogenerate_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	gogenerate "github.com/paketo-buildpacks/go-generate"
	"github.com/paketo-buildpacks/go-generate/fakes"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		layersDir    string
		workingDir   string
		logs         *bytes.Buffer
		buildProcess *fakes.BuildProcess

		build packit.BuildFunc
	)

	it.Before(func() {
		var err error
		layersDir, err = ioutil.TempDir("", "layers-dir")
		Expect(err).NotTo(HaveOccurred())

		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		logs = bytes.NewBuffer(nil)

		buildProcess = &fakes.BuildProcess{}

		build = gogenerate.Build(
			buildProcess,
			scribe.NewLogger(logs),
		)
	})

	it.After(func() {
		Expect(os.RemoveAll(layersDir)).To(Succeed())
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	it("builds", func() {
		result, err := build(packit.BuildContext{
			Layers:     packit.Layers{Path: layersDir},
			WorkingDir: workingDir,
			BuildpackInfo: packit.BuildpackInfo{
				Name:    "Some Buildpack",
				Version: "some-version",
			},
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(result).To(Equal(packit.BuildResult{}))

		Expect(buildProcess.ExecuteCall.Receives.WorkingDir).To(Equal(workingDir))

		Expect(logs.String()).To(ContainSubstring("Some Buildpack some-version"))
	})

	context("failure cases", func() {

		context("build process fails to execute", func() {
			it.Before(func() {
				buildProcess.ExecuteCall.Stub = nil
				buildProcess.ExecuteCall.Returns.Error = errors.New("build process failed to execute")
			})

			it("returns an error", func() {
				_, err := build(packit.BuildContext{
					Layers:     packit.Layers{Path: layersDir},
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(ContainSubstring("build process failed to execute")))
			})
		})
	})
}
