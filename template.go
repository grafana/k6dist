package k6dist

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/semver/v3"
	sprig "github.com/go-task/slim-sprig/v3"

	"github.com/grafana/k6dist/internal/registry"
)

type instanceData struct {
	Name    string
	Version string
	OS      string
	Arch    string
	ExeExt  string
	ZipExt  string
}

type releaseData struct {
	Name     string
	Version  string
	Registry registry.Registry
	Footer   string
}

func newInstsanceData(name string, version *semver.Version, platform *Platform) *instanceData {
	exe := ""
	zip := ".tar.gz"

	if platform.OS == "windows" {
		exe = ".exe"
		zip = ".zip"
	}

	return &instanceData{
		Name:    name,
		Version: version.Original(),
		OS:      platform.OS,
		Arch:    platform.Arch,
		ExeExt:  exe,
		ZipExt:  zip,
	}
}

func newReleaseData(name string, version *semver.Version, reg registry.Registry) (*releaseData, error) {
	footer, err := notesFooter(version, reg)
	if err != nil {
		return nil, err
	}

	return &releaseData{Name: name, Version: version.Original(), Registry: reg, Footer: footer}, nil
}

func expandTemplate(name string, tmplsrc string, data interface{}) (string, error) {
	tmpl, err := template.New(name).Funcs(sprig.TxtFuncMap()).Parse(tmplsrc)
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer

	if err := tmpl.Execute(&buff, data); err != nil {
		return "", err
	}

	return buff.String(), nil
}

func expandAsTargetPath(name string, tmplsrc string, data interface{}) (string, error) {
	filename, err := expandTemplate(name, tmplsrc, data)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0o700); err != nil { //nolint:forbidigo
		return "", err
	}

	return filename, nil
}
