package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"

	diff "github.com/khulnasoft/gojsondiff"
	"github.com/khulnasoft/gojsondiff/formatter"
)

func main() {
	app := cli.NewApp()
	app.Name = "jd"
	app.Usage = "JSON Diff"
	app.Version = "0.0.2"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "format, f",
			Value:  "ascii",
			Usage:  "Diff Output Format (ascii, delta)",
			EnvVar: "DIFF_FORMAT",
		},
		cli.BoolFlag{
			Name:   "coloring, c",
			Usage:  "Enable coloring in the ASCII mode (not available in the delta mode)",
			EnvVar: "COLORING",
		},
		cli.BoolFlag{
			Name:   "quiet, q",
			Usage:  "Suppress output, if no differences are found",
			EnvVar: "QUIET",
		},
	}

	app.Action = func(c *cli.Context) error {
		if len(c.Args()) < 2 {
			fmt.Println("Not enough arguments.")
			fmt.Printf("Usage: %s json_file another_json_file\n", app.Name)
			os.Exit(1)
		}

		aFilePath := c.Args()[0]
		bFilePath := c.Args()[1]

		// Prepare your JSON string as `[]byte`, not `string`
		aString, err := ioutil.ReadFile(aFilePath)
		if err != nil {
			fmt.Printf("Failed to open file '%s': %s\n", aFilePath, err.Error())
			os.Exit(2)
		}

		// Another JSON string
		bString, err := ioutil.ReadFile(bFilePath)
		if err != nil {
			fmt.Printf("Failed to open file '%s': %s\n", bFilePath, err.Error())
			os.Exit(2)
		}

		// Then, compare them
		differ := diff.New()
		d, err := differ.Compare(aString, bString)
		if err != nil {
			fmt.Printf("Failed to unmarshal file: %s\n", err.Error())
			os.Exit(3)
		}

		// Output the result
		if d.Modified() || !c.Bool("quiet") {
			format := c.String("format")
			var diffString string
			if format == "ascii" {
				var aJson map[string]interface{}
				json.Unmarshal(aString, &aJson)

				config := formatter.AsciiFormatterConfig{
					ShowArrayIndex: true,
					Coloring:       c.Bool("coloring"),
				}

				formatter := formatter.NewAsciiFormatter(aJson, config)
				diffString, err = formatter.Format(d)
				if err != nil {
					// No error can occur
				}
			} else if format == "delta" {
				formatter := formatter.NewDeltaFormatter()
				diffString, err = formatter.Format(d)
				if err != nil {
					// No error can occur
				}
			} else {
				fmt.Printf("Unknown Foramt %s\n", format)
				os.Exit(4)
			}

			fmt.Print(diffString)
			return cli.NewExitError("", 1)
		}
		return nil
	}

	app.Run(os.Args)
}
