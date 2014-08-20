package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type episode struct {
	seasonNum  int
	episodeNum int
	title      string
	lines      []*line
}

type line struct {
	order   int
	start   int
	end     int
	content string
}

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

	episodes := parseSubtitleFiles()

	fmt.Printf("Episode Count %d\n", len(episodes))

	// Read list of common words into memory
	byte_contents, err := ioutil.ReadFile("common_25000_lowercase.txt")
	if err != nil {
		fmt.Printf("Error opening list of common words: %s\n", err)
	}

	common := make(map[string]bool)
	for _, s := range strings.Split(string(byte_contents), "\n") {
		common[s] = true
	}

	// Check how effective the common list even is
	whitespace := regexp.MustCompile("\\s")
	skipped, total := 0, 0
	for _, ep := range episodes {
		fmt.Printf("Parsing Season %d, Episode %d: %s\n", ep.seasonNum, ep.episodeNum, ep.title)
		for _, line := range ep.lines {
			words := whitespace.Split(line.content, -1)
			for _, word := range words {
				total += 1
				if common[word] {
					skipped += 1
				}
			}
		}
	}
	fmt.Printf("Total words: %d\nSkipped words: %d\nOther: %d\n", total, skipped, total-skipped)
}

func timestampToMilliseconds(str string) (int, int) {
	r := regexp.MustCompile("([0-9]*):([0-9]*):([0-9]*),([0-9]*)")
	matches := r.FindAllStringSubmatch(str, -1)

	// Start time in milliseconds
	hours, _ := strconv.Atoi(matches[0][0])
	minutes, _ := strconv.Atoi(matches[0][1])
	seconds, _ := strconv.Atoi(matches[0][2])
	milliseconds, _ := strconv.Atoi(matches[0][3])
	start := hours*3600000 + minutes*60000 + seconds*1000 + milliseconds

	// End time in milliseconds
	hours, _ = strconv.Atoi(matches[0][0])
	minutes, _ = strconv.Atoi(matches[0][1])
	seconds, _ = strconv.Atoi(matches[0][2])
	milliseconds, _ = strconv.Atoi(matches[0][3])
	end := hours*3600000 + minutes*60000 + seconds*1000 + milliseconds

	return start, end
}

func parseSubtitleFiles() []*episode {
	files, err := ioutil.ReadDir("subtitles")
	if err != nil {
		fmt.Printf("There was an error listing the files in the subtitles directory: ", err)
	}

	empty := regexp.MustCompile("^\\s$")
	exclude := regexp.MustCompile("([.?,\\-=/\"\\(\\)\\[\\]!0-9♪]|</?i>|)")

	episodes := []*episode{}
	for _, f := range files {
		// Extract the title and episode from the filename
		metadata := regexp.MustCompile("([0-9])x([0-9]{2}) - (.+)\\.en\\.srt").FindStringSubmatch(f.Name())
		ep := &episode{
			title: metadata[3],
		}
		ep.seasonNum, _ = strconv.Atoi(metadata[1])
		ep.episodeNum, _ = strconv.Atoi(metadata[2])

		fmt.Printf("Reading Season %d, Episode %d: %s\n", ep.seasonNum, ep.episodeNum, ep.title)

		// Get the file contents
		byte_contents, err := ioutil.ReadFile("subtitles/" + f.Name())
		if err != nil {
			fmt.Printf("There was an error reading the subtitle file: %s\n", err)
		}
		contents := string(byte_contents)

		// Subtitle files are structured nicely using newlines
		segments := regexp.MustCompile("\n\n+").Split(contents, -1)
		for _, segment := range segments {
			rows := strings.Split(segment, "\n")
			l := &line{}
			l.order, err = strconv.Atoi(rows[0])
			l.start, l.end = timestampToMilliseconds(rows[1])
			l.content = ""
			for i := 2; i < len(rows); i++ {
				l.content += " " + strings.ToLower(exclude.ReplaceAllString(rows[i], ""))
			}

			// Empty subtitles hold no information
			if empty.MatchString(l.content) {
				continue
			}

			ep.lines = append(ep.lines, l)
		}
		episodes = append(episodes, ep)
	}
	return episodes
}
