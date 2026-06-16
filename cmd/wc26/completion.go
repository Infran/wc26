package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion script for the specified shell.
Load it with:

  Bash:   source <(wc26 completion bash)
  Zsh:    source <(wc26 completion zsh)
  Fish:   wc26 completion fish | source
  PowerShell: wc26 completion powershell | Out-String | Invoke-Expression

To persist, redirect to a file and add to your shell config.`,
	Args: cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return fmt.Errorf("unsupported shell: %s (use: bash, zsh, fish, powershell)", args[0])
		}
	},
}
