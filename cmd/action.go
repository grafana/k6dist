package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// AddGitHubArgs adds GitHub input parameters as arguments.
func AddGitHubArgs(args []string) []string {
	if !isGitHubAction() {
		return args
	}

	args = ghinput(args, "distro_name", "--distro-name")
	args = ghinput(args, "distro_version", "--distro-version")
	args = ghinput(args, "executable", "--executable")
	args = ghinput(args, "archive", "--archive")
	args = ghinput(args, "docker", "--docker")
	args = ghinput(args, "docker_template", "--docker-template")
	args = ghinput(args, "notes", "--notes")
	args = ghinput(args, "notes_template", "--notes-template")
	args = ghinput(args, "notes_latest", "--notes-latest")
	args = ghinput(args, "readme", "--readme")
	args = ghinput(args, "license", "--license")
	args = ghinput(args, "platform", "--platform")

	if getenv("INPUT_VERBOSE", "false") == "true" {
		args = append(args, "--verbose")
	}

	if getenv("INPUT_QUIET", "false") == "true" {
		args = append(args, "--quiet")
	}

	if in := getenv("INPUT_IN", ""); len(in) > 0 {
		args = append(args, in)
	}

	return args
}

//nolint:forbidigo
func isGitHubAction() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

//nolint:forbidigo
func getenv(name string, defval string) string {
	value, found := os.LookupEnv(name)
	if found {
		return value
	}

	return defval
}

func ghinput(args []string, name string, flag string) []string {
	val := getenv("INPUT_"+strings.ToUpper(name), "")
	if len(val) > 0 {
		args = append(args, flag, val)
	}

	return args
}

//nolint:forbidigo
func emitOutput(changed bool, version string) error {
	ghOutput := getenv("GITHUB_OUTPUT", "")
	if len(ghOutput) == 0 {
		return nil
	}

	file, err := os.Create(filepath.Clean(ghOutput))
	if err != nil {
		return err
	}

	slog.Debug("Emit changed", "changed", changed)

	_, err = fmt.Fprintf(file, "changed=%t\n", changed)
	if err != nil {
		return err
	}

	slog.Debug("Emit version", "version", version)

	_, err = fmt.Fprintf(file, "version=%s\n", version)
	if err != nil {
		return err
	}

	return file.Close()
}
