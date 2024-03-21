package cmd

import (
	"certbot/internal/certbot"
	"certbot/internal/config"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var configPath string

	c := &cobra.Command{
		Use: filepath.Base(os.Args[0]),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadConfigFromFile(configPath)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			bot := &certbot.CertBot{
				Cfg: cfg,
			}

			if err := bot.Run(); err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		},
	}

	c.Flags().StringVarP(&configPath, "config", "c", "config.yaml", "config file path")
	c.MarkFlagFilename("config")
	return c
}
