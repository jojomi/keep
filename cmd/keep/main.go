package main

import (
	"bufio"
	"fmt"
	"github.com/djherbis/times"
	"github.com/jojomi/keep"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	rootCmd := cobra.Command{
		Use: "keep",
		Run: runRoot,
	}

	flags := rootCmd.PersistentFlags()
	flags.StringP("requirements", "r", "10 last, 14 days, 12 weeks, 12 months, 12 years", "keep config")
	flags.Bool("print-requirements-only", false, "print perceived requirements")
	flags.BoolP("dry-run", "n", false, "don't actually delete files, but show which would be deleted")
	flags.BoolP("force", "f", false, "don't ask questions, just do it")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	env, err := parseEnvRoot(cmd, args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	now := time.Now()
	reqs := keep.NewRequirementsFromString(env.Requirements)
	fmt.Println(reqs)

	if env.PrintRequirementsOnly {
		os.Exit(0)
	}

	jh := keep.NewDefaultJailhouse[keep.File]()

	// file selection
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error reading current directory: %v", err)
	}
	files, err := os.ReadDir(wd)
	if err != nil {
		log.Fatalf("Error reading directory contents: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		t, err := times.Stat(filename)
		if err != nil {
			log.Fatal(err.Error())
		}

		jh.AddElements(keep.File{
			Filename: filename,
			Time:     t.BirthTime(),
		})
	}

	// apply requirements to find which files to keep and which to delete
	jh.ApplyRequirementsForDate(*reqs, now)

	k := jh.KeptElements()
	fmt.Printf("\nKeeping %d files:\n", len(k))
	for _, keepElement := range k {
		tags := keepElement.GetTags()
		tagStrings := make([]string, len(tags))
		for i, tag := range tags {
			tagStrings[i] = tag.String()
		}
		fmt.Printf("%s [%s]\n", keepElement.TimeResource.Filename, strings.Join(tagStrings, ", "))
	}

	r := jh.FreeElements()
	if len(r) > 0 {
		fmt.Printf("\nRemoving %d files:\n", len(r))
		for _, keepElement := range r {
			fmt.Println(keepElement.TimeResource.Filename)
		}

		doRemove := env.Force
		if !env.Force {
			reader := bufio.NewReader(os.Stdin)

			fmt.Printf("\nRemove %d files as listed above? (Y/n) ", len(r))
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				return
			}

			input = input[:len(input)-1] // Remove newline character from the input

			doRemove = input == "y" || input == "j" || input == ""
		}
		if doRemove {
			for _, keepElement := range r {
				f := keepElement.TimeResource.Filename
				if env.DryRun {
					fmt.Printf("[DRY-RUN] Would be deleting %s...\n", f)
				} else {
					err = os.Remove(f)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(4)
					}
				}
			}

			if !env.DryRun {
				fmt.Printf("deleted %d files\n", len(r))
			}
		}
	}
}
