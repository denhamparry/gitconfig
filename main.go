package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var logLevel string

func main() {
	var rootCmd = &cobra.Command{
		Use: "gitconfig-cli",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logLevel, _ = cmd.Flags().GetString("log")
		},
	}

	var setupGitsignCmd = &cobra.Command{
		Use:   "setup-gitsign",
		Short: "Setup git signing configuration",
		Run: func(cmd *cobra.Command, args []string) {
			email, _ := cmd.Flags().GetString("email")
			setupGitsign(email)
		},
	}

	var clearGitsignCmd = &cobra.Command{
		Use:   "clear-gitsign",
		Short: "Clear git signing configuration",
		Run: func(cmd *cobra.Command, args []string) {
			setupGitLocalCleanup()
		},
	}

	setupGitsignCmd.Flags().StringP("email", "e", "", "Email address (optional)")

	rootCmd.PersistentFlags().StringVar(&logLevel, "log", "", "Set log level")
	rootCmd.AddCommand(setupGitsignCmd)
	rootCmd.AddCommand(clearGitsignCmd)
	rootCmd.Execute()
}

func setupGitsign(email string) {
	err := checkDirectoryForGit()
	if err != nil {
		logError(err)
		return
	}
	err = setupGitLocalCleanup()
	if err != nil {
		logError(err)
		return
	}
	err = setupGitLocalUser(email)
	if err != nil {
		logError(err)
		return
	}
}

func checkDirectoryForGit() error {
	_, err := os.Stat(".git")
	if os.IsNotExist(err) {
		return fmt.Errorf("this directory does not contain a git repository")
	}
	return nil
}

func setupGitLocalCleanup() error {
	configs := []string{
		"user.email",
		"commit.gpgsign",
		"tag.gpgsign",
		"gpg.x509.program",
		"gpg.format",
		"gitsign.connectorID",
		"user.signingkey",
		"gpg.ssh.program",
	}

	for _, config := range configs {
		err := runGitUnset(config)
		if err != nil {
			if !strings.Contains(fmt.Sprintf("%s", err), "exit status 5") {
				return fmt.Errorf("%s", err)
			}
		}
	}

	return nil
}

func setupGitLocalUser(email string) error {

	// If email is not provided, prompt user for email
	if email == "" {

		fmt.Println("Enter an email address:")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}
		email = input
	}

	// Check that email is a valid email address
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		log.Fatalf("Invalid email address")
	}

	err := setupGitConfigLocalUser(email)
	if err != nil {
		log.Fatalf("Failed to set up git config: %v", err)
	}

	return nil
}

func setupGitConfigLocalUser(email string) error {
	configurations := map[string]string{
		"commit.gpgsign":      "true",
		"tag.gpgsign":         "true",
		"gpg.x509.program":    "gitsign",
		"gpg.format":          "x509",
		"gitsign.connectorID": "https://accounts.google.com",
	}

	configurations["user.email"] = email

	for key, value := range configurations {
		err := runGitConfig(key, value)
		if err != nil {
			return fmt.Errorf("%s", err)
		}
	}

	return nil
}

func runGitConfig(key, value string) error {
	cmd := exec.Command("git", "config", "--local", key, value)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run git config for %s: %v", key, err)
	}
	return nil
}

func runGitUnset(key string) error {
	cmd := exec.Command("git", "config", "--local", "--unset", key)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		if !strings.Contains(fmt.Sprintf("%s", err), "exit status 5") {
			return fmt.Errorf("failed to unset git config key %s: %v", key, err)
		}
	}
	return nil
}

func logError(err error) {
	if logLevel == "error" {
		fmt.Println(err)
	}
}
