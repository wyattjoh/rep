package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/gobeam/stringy"
)

var questions = []*survey.Question{
	{
		Name:     "issue",
		Prompt:   &survey.Input{Message: "GitHub Issue #"},
		Validate: survey.Required,
	},
	{
		Name:     "description",
		Prompt:   &survey.Input{Message: "GitHub Issue Description"},
		Validate: survey.Required,
	},
	{
		Name:   "repository",
		Prompt: &survey.Input{Message: "GitHub Reproduction Repository"},
		Transform: func(ans interface{}) interface{} {
			repository, ok := ans.(string)
			if !ok {
				return ans
			}

			if strings.HasSuffix(repository, ".git") {
				return repository
			}

			return fmt.Sprintf("%s.git", repository)
		},
		Validate: func(ans interface{}) error {
			repository, ok := ans.(string)
			if !ok {
				return errors.New("Value is not a string")
			}

			if !strings.HasPrefix(repository, "https://github.com") {
				return errors.New("Expecter GitHub HTTP URL")
			}

			if strings.Count(repository, "/") != 4 {
				return errors.New("Unexpected format")
			}

			return nil
		},
	},
}

func main() {
	config, err := LoadOrCreateConfig()
	if err != nil {
		if !errors.Is(err, terminal.InterruptErr) {
			fmt.Fprintln(os.Stderr, err.Error())
		}

		os.Exit(1)
	}

	// the answers will be written to this struct
	answers := struct {
		Issue       int32  `survey:"issue"`
		Description string `survey:"description"`
		Repository  string `survey:"repository"`
	}{}

	if err := survey.Ask(questions, &answers); err != nil {
		if !errors.Is(err, terminal.InterruptErr) {
			fmt.Fprintln(os.Stderr, err.Error())
		}

		os.Exit(1)
	}

	// Format the new folder name.
	description := stringy.New(answers.Description).SnakeCase().ToLower()
	folder := fmt.Sprintf("%d_%s", answers.Issue, description)

	path := filepath.Join(config.Directory, folder)

	if err := clone(answers.Repository, path); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := install(path); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("GitHub Reproduction at %s\n\nCloned to %s\n", answers.Repository, path)
}

func clone(repository, path string) error {
	cmd := exec.Command("git", "clone", repository, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func install(path string) error {
	// check if the package.json exists before installing.
	if _, err := os.Stat(filepath.Join(path, "package.json")); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	cmd := exec.Command("pnpm", "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Dir = path

	return cmd.Run()
}
