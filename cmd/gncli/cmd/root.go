package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/prmsrswt/gonotify/cmd/gncli/client"
	"github.com/prmsrswt/gonotify/cmd/gncli/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gncli",
	Short:   "gncli - a commanf-line client for GoNotify",
	Version: "v0.1.0",
	// Run:     list,
}

func getClient(baseURL, token string) *client.Client {
	c, err := client.NewClient(baseURL, token)
	if err != nil {
		fmt.Println("Error initialising API client")
		os.Exit(2)
	}

	return c
}

func getConfig() *config.Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	baseDir := path.Join(homeDir, ".gonotify")

	err = os.MkdirAll(baseDir, 0755)
	if err != nil {
		panic(err)
	}

	c := &config.Config{Path: path.Join(baseDir, "config.json")}

	err = c.Load()
	if err != nil {
		panic(err)
	}

	return c
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
