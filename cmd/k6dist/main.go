// Package main contains the main function for k6dist.
package main

import (
	"log/slog"
	"os"

	"github.com/grafana/k6dist/cmd"
	sloglogrus "github.com/samber/slog-logrus/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var (
	appname = "k6dist" //nolint:unused
	version = "dev"
)

func initLogging() *slog.LevelVar {
	levelVar := new(slog.LevelVar)

	logrus.SetLevel(logrus.DebugLevel)

	logger := slog.New(sloglogrus.Option{Level: levelVar}.NewLogrusHandler())

	slog.SetDefault(logger)

	return levelVar
}

func main() {
	root, err := newCmd(os.Args[1:], initLogging()) //nolint:forbidigo
	checkErr(err)
	checkErr(root.Execute())
}

func newCmd(args []string, levelVar *slog.LevelVar) (*cobra.Command, error) {
	root := cmd.New(levelVar)

	root.Version = version

	args, err := cmd.AddGitHubArgs(args, root.Flags())
	if err != nil {
		return nil, err
	}

	slog.Debug("arguments", "args", args)

	root.SetArgs(args)

	return root, nil
}

func checkErr(err error) {
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1) //nolint:forbidigo
	}
}
