package k6dist

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed Dockerfile.tpl
var defaultDockerTemplate string

func createDockerfile(executable string, tmplfile string, output string) error {
	var tmplsrc string

	if len(tmplfile) != 0 {
		bin, err := os.ReadFile(filepath.Clean(tmplfile)) //nolint:forbidigo
		if err != nil {
			return err
		}

		tmplsrc = string(bin)
	} else {
		tmplsrc = defaultDockerTemplate
	}

	dir := filepath.Dir(executable)

	str, err := expandTemplate("Dockerfile", tmplsrc, map[string]string{"Executable": filepath.Base(executable)})
	if err != nil {
		return err
	}

	name := output
	if len(name) == 0 {
		name = filepath.Join(dir, "Dockerfile")
	}

	return os.WriteFile(name, []byte(str), 0o600) //nolint:forbidigo
}
