package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "ddx"
	app.Usage = "Analysis of all considered diagnoses across all episodes of House M.D."
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:  "parse",
			Usage: "parse the subtitle files and seed the database",
			Action: func(c *cli.Context) {
				fmt.Println("Parsing subtitle files and seeding database with possible fragments")
			},
		},
	}
	app.Run(os.Args)
}
