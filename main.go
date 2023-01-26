package main

import (
	"github.com/urfave/cli/v2"
	"github.com/ytake/kcr/log"
	"os"
)

func main() {
	l := log.NewLogger()
	defer l.Provider.Sync()
	app := &cli.App{
		Commands: []*cli.Command{}}
	app.Name = `kcr`
	if err := app.Run(os.Args); err != nil {
		l.RuntimeFatalError("kcr command error", err)
	}
}
