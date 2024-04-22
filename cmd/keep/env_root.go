package main

import (
	"github.com/spf13/cobra"
)

type EnvRoot struct {
	PrintRequirementsOnly bool
	Requirements          string
	Force                 bool
	DryRun                bool
}

func parseEnvRoot(cmd *cobra.Command, _ []string) (EnvRoot, error) {
	var err error

	env := EnvRoot{}

	env.PrintRequirementsOnly, err = cmd.Flags().GetBool("print-requirements-only")
	if err != nil {
		return env, err
	}

	env.Requirements, err = cmd.Flags().GetString("requirements")
	if err != nil {
		return env, err
	}

	env.Force, err = cmd.Flags().GetBool("force")
	if err != nil {
		return env, err
	}

	env.DryRun, err = cmd.Flags().GetBool("dry-run")
	if err != nil {
		return env, err
	}
	return env, nil
}
