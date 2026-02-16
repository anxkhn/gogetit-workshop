package cmd

import (
	"github.com/anxkhn/gogetit-workshop/internal/config"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "gogetit",
	Short: "A concurrent file downloader and website metadata scraper",
	Long: `GoGetIt is a powerful CLI tool for downloading files concurrently
and scraping website metadata. It features a beautiful TUI progress
bar and supports configuration files.`,
	Example: `  # Download a single file
  gogetit download https://example.com/file.zip

  # Download multiple files concurrently
  gogetit download https://example.com/file1.zip https://example.com/file2.zip

  # Scrape website metadata
  gogetit scrape https://example.com`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gogetit.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		config.SetConfigFile(cfgFile)
	}
	config.Load()
}
