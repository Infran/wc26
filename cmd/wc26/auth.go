package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/Infran/wc26/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands (register, login, logout, status)",
}

var registerCmd = &cobra.Command{
	Use:   "register <name> <email> <password>",
	Short: "Register a new user and save JWT token",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := newAPIClient()
		resp, err := cl.Register(args[0], args[1], args[2])
		if err != nil {
			return fmt.Errorf("registration failed: %w", err)
		}

		viper.Set("auth.token", resp.Token)
		viper.Set("auth.email", resp.User.Email)
		if err := config.Set("auth.token", resp.Token); err != nil {
			return err
		}
		if err := config.Set("auth.email", resp.User.Email); err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "Registered and logged in as: %s (%s)\n", resp.User.Name, resp.User.Email)
		return nil
	},
}

var loginCmd = &cobra.Command{
	Use:   "login <email> <password>",
	Short: "Login and save JWT token",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := newAPIClient()
		resp, err := cl.Login(args[0], args[1])
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		viper.Set("auth.token", resp.Token)
		viper.Set("auth.email", resp.User.Email)
		if err := config.Set("auth.token", resp.Token); err != nil {
			return err
		}
		if err := config.Set("auth.email", resp.User.Email); err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "Logged in as: %s (%s)\n", resp.User.Name, resp.User.Email)
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear saved authentication token",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set("auth.token", "")
		viper.Set("auth.email", "")
		if err := config.Set("auth.token", ""); err != nil {
			return err
		}
		if err := config.Set("auth.email", ""); err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "Logged out. Token cleared.")
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := viper.GetString("auth.token")
		email := viper.GetString("auth.email")
		baseURL := viper.GetString("api.base_url")

		fmt.Fprintf(os.Stderr, "API URL: %s\n", baseURL)
		if token == "" {
			fmt.Fprintln(os.Stderr, "Status: not authenticated")
			fmt.Fprintln(os.Stderr, "Run 'wc26 auth login <email> <password>' to authenticate")
		} else {
			fmt.Fprintf(os.Stderr, "Status: authenticated\n")
			fmt.Fprintf(os.Stderr, "Email: %s\n", email)
			fmt.Fprintf(os.Stderr, "Token: %s…\n", token[:min(len(token), 20)])
		}
		return nil
	},
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Display JWT token (prompts for password to validate first)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		email := viper.GetString("auth.email")
		if email == "" {
			return fmt.Errorf("not logged in. Run 'wc26 auth login <email> <password>' first")
		}

		fmt.Fprintf(os.Stderr, "Password for %s: ", email)
		password, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return fmt.Errorf("reading password: %w", err)
		}
		if len(password) == 0 {
			return fmt.Errorf("password cannot be empty")
		}

		cl := newAPIClient()
		resp, err := cl.Login(email, string(password))
		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		fmt.Println(resp.Token)
		return nil
	},
}

func init() {
	authCmd.AddCommand(registerCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
	authCmd.AddCommand(tokenCmd)
}
