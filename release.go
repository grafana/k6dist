package k6dist

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"

	"github.com/grafana/k6dist/internal/registry"
)

//go:embed NOTES.md.tpl
var defaultNotesTemplate string

const (
	notesFooterBegin = "<!--```json"
	notesFooterEnd   = "```-->"
)

func notesFooter(registry registry.Registry) (string, error) {
	var buff bytes.Buffer

	buff.WriteString(notesFooterBegin)
	buff.WriteRune('\n')

	encoder := json.NewEncoder(&buff)

	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(registry.ToModules()); err != nil {
		return "", err
	}

	buff.WriteString(notesFooterEnd)
	buff.WriteRune('\n')

	return buff.String(), nil
}

func expandNotes(name, version string, reg registry.Registry, tmplfile string) (string, error) {
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

var reModules = regexp.MustCompile("(?ms:^" + notesFooterBegin + "(?P<modules>.*)" + notesFooterEnd + ")")

func parseNotes(notes []byte) (bool, registry.Modules, error) {
	match := reModules.FindSubmatch(notes)

	if match == nil {
		return false, nil, nil
	}

	var modules registry.Modules

	if err := json.Unmarshal(match[reModules.SubexpIndex("modules")], &modules); err != nil {
		return false, nil, err
	}

	return true, modules, nil
}
