package k6dist

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/grafana/k6dist/internal/registry"
	"github.com/grafana/k6foundry"
)

var initialVersion = semver.MustParse("v0.1.0") //nolint:gochecknoglobals

func detectChange(reg registry.Registry, latest string) (bool, *semver.Version, error) {
	if len(latest) == 0 {
		return true, initialVersion, nil
	}

	contents, err := os.ReadFile(filepath.Clean(latest)) //nolint:forbidigo
	if err != nil {
		return false, nil, err
	}

	found, version, modules, err := parseNotes(contents)
	if err != nil {
		return false, nil, err
	}

	if !found {
		return true, initialVersion, nil
	}

	bumped := reg.AddLatest(modules, version)

	return !bumped.Equal(version), bumped, nil
}

func newBuilder(ctx context.Context, modules registry.Modules) (k6foundry.Builder, error) {
	cgo := "0"
	if modules.Cgo() {
		cgo = "1"
	}

	return k6foundry.NewNativeBuilder(
		ctx,
		k6foundry.NativeBuilderOpts{
			Logger: slog.Default(),
			Stdout: os.Stdout, //nolint:forbidigo
			Stderr: os.Stderr, //nolint:forbidigo
			GoOpts: k6foundry.GoOpts{
				CopyGoEnv: true,
				Env:       map[string]string{"GOWORK": "off", "CGO_ENABLED": cgo},
			},
		},
	)
}

// Build builds k6 binaries and archives based on opts parameter.
func Build(ctx context.Context, opts *Options) (bool, *semver.Version, error) {
	opts.setDefaults()

	registry, err := registry.LoadRegistry(ctx, opts.Registry)
	if err != nil {
		return false, nil, err
	}

	changed, version, err := detectChange(registry, opts.NotesLatest)
	if err != nil {
		return false, nil, err
	}

	if !changed {
		return false, nil, nil
	}

	if opts.Version != nil {
		version = opts.Version
	}

	notes, err := expandNotes(opts.Name, version, registry, opts.NotesTemplate)
	if err != nil {
		return false, nil, err
	}

	filename, err := expandAsTargetPath("notes", opts.Notes, newInstsanceData(opts.Name, version, &Platform{}))
	if err != nil {
		return false, nil, err
	}

	if err := os.WriteFile(filename, []byte(notes), 0o600); err != nil { //nolint:forbidigo
		return false, nil, err
	}

	modules := registry.ToModules()
	k6Version, mods := modules.ToFoundry()

	builder, err := newBuilder(ctx, modules)
	if err != nil {
		return false, nil, err
	}

	for _, platform := range opts.Platforms {
		data := newInstsanceData(opts.Name, version, platform)

		filename, err := expandAsTargetPath("executable", opts.Executable, data)
		if err != nil {
			return false, nil, err
		}

		err = buildExecutable(ctx, builder, platform, k6Version, mods, filename)
		if err != nil {
			return false, nil, err
		}

		err = createDockerfile(filename, opts.DockerfileTemplate, opts.Dockerfile)
		if err != nil {
			return false, nil, err
		}

		archive, err := expandAsTargetPath("archive", opts.Archive, data)
		if err != nil {
			return false, nil, err
		}

		err = buildArchive(archive, filename, opts.Readme, opts.License)
		if err != nil {
			return false, nil, err
		}
	}

	return true, version, nil
}

func buildExecutable(
	ctx context.Context,
	builder k6foundry.Builder,
	platform *Platform,
	k6Version string,
	mods []k6foundry.Module,
	filename string,
) error {
	file, err := os.Create(filepath.Clean(filename)) //nolint:forbidigo
	if err != nil {
		return err
	}

	err = builder.Build(
		ctx,
		k6foundry.NewPlatform(platform.OS, platform.Arch),
		k6Version,
		mods,
		[]string{`-ldflags`, `-s -w`},
		file,
	)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return os.Chmod(filename, 0o700) //nolint:forbidigo,gosec
}

func buildArchive(archive string, exe, readme, license string) error {
	files := []string{exe}

	if len(readme) > 0 {
		files = append(files, readme)
	}

	if len(license) > 0 {
		files = append(files, license)
	}

	if strings.ToLower(filepath.Ext(archive)) == ".zip" {
		return zipFor(archive, files...)
	}

	return tgzFor(archive, files...)
}

func tgzFor(archive string, files ...string) error {
	tgz, err := os.Create(filepath.Clean(archive)) //nolint:forbidigo
	if err != nil {
		return err
	}
	defer tgz.Close() //nolint:errcheck

	gzipWriter := gzip.NewWriter(tgz)
	defer gzipWriter.Close() //nolint:errcheck

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close() //nolint:errcheck

	for _, file := range files {
		if err := addToTar(tarWriter, file); err != nil {
			return err
		}
	}

	return nil
}

func zipFor(archive string, files ...string) error {
	out, err := os.Create(filepath.Clean(archive)) //nolint:forbidigo
	if err != nil {
		return err
	}
	defer out.Close() //nolint:errcheck

	zipWriter := zip.NewWriter(out)
	defer zipWriter.Close() //nolint:errcheck

	for _, file := range files {
		if err := addToZip(zipWriter, file); err != nil {
			return err
		}
	}

	return nil
}

func addToTar(archive *tar.Writer, path string) error {
	file, err := os.Open(filepath.Clean(path)) //nolint:forbidigo
	if err != nil {
		return err
	}

	defer file.Close() //nolint:errcheck

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	name := filepath.Base(path)

	header := &tar.Header{ //nolint:exhaustruct
		Name:    name,
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}

	err = archive.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(archive, file)
	if err != nil {
		return err
	}

	return nil
}

func addToZip(archive *zip.Writer, path string) error {
	file, err := os.Open(filepath.Clean(path)) //nolint:forbidigo
	if err != nil {
		return err
	}

	defer file.Close() //nolint:errcheck

	_, err = file.Stat()
	if err != nil {
		return err
	}

	name := filepath.Base(path)

	writer, err := archive.Create(name)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)

	return err
}
