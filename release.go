package k6dist

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Masterminds/semver/v3"
	"github.com/grafana/k6dist/internal/registry"
)

//go:embed NOTES.md.tpl
var defaultNotesTemplate string

const (
	notesFooterBegin = "<!--```json"
	notesFooterEnd   = "```-->"
)

type footerData struct {
	Version string           `json:"version,omitempty"`
	Modules registry.Modules `json:"modules,omitempty"`
}

func notesFooter(version *semver.Version, registry registry.Registry) (string, error) {
	var buff bytes.Buffer

	buff.WriteString(notesFooterBegin)
	buff.WriteRune('\n')

	encoder := json.NewEncoder(&buff)

	encoder.SetEscapeHTML(false)

	data := &footerData{Version: version.Original(), Modules: registry.ToModules()}

	if err := encoder.Encode(data); err != nil {
		return "", err
	}

	buff.WriteString(notesFooterEnd)
	buff.WriteRune('\n')

	return buff.String(), nil
}

func expandNotes(name string, version *semver.Version, reg registry.Registry, tmplfile string) (string, error) {
	var tmplsrc string

	if len(tmplfile) != 0 {
		bin, err := os.ReadFile(filepath.Clean(tmplfile)) //nolint:forbidigo
		if err != nil {
			return "", err
		}

		tmplsrc = string(bin)
	} else {
		tmplsrc = defaultNotesTemplate
	}

	data, err := newReleaseData(name, version, reg)
	if err != nil {
		return "", err
	}

	return expandTemplate("notes", tmplsrc, data)
}

var reModules = regexp.MustCompile("(?ms:^" + notesFooterBegin + "(?P<state>.*)" + notesFooterEnd + ")")

func parseNotes(notes []byte) (bool, *semver.Version, registry.Modules, error) {
	match := reModules.FindSubmatch(notes)

	if match == nil {
		return false, nil, nil, nil
	}

	var data footerData

	if err := json.Unmarshal(match[reModules.SubexpIndex("state")], &data); err != nil {
		return false, nil, nil, err
	}

	version, err := semver.NewVersion(data.Version)
	if err != nil {
		return false, nil, nil, err
	}

	return true, version, data.Modules, nil
}
