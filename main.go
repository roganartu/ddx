package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
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
				parse()
			},
		},
	}
	app.Run(os.Args)
}

func parse() {
	fmt.Println("Parsing subtitle files and seeding database with possible fragments")

	files, err := ioutil.ReadDir("subtitles")
	if err != nil {
		fmt.Printf("There was an error listing the files in the subtitles directory: ", err)
	}

	for _, f := range files {
		// Extract the title and episode from the filename
		metadata := regexp.MustCompile("([0-9])x([0-9]{2}) - (.+)\\.en\\.srt").FindStringSubmatch(f.Name())
		season, _ := strconv.Atoi(metadata[1])
		episode, _ := strconv.Atoi(metadata[2])
		title := metadata[3]

		fmt.Printf("Parsing Season %d, Episode %d: %s\n", season, episode, title)
	}
}
