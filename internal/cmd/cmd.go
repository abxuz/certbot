package cmd

import (
	"certbot/internal/certbot"
	"certbot/internal/config"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Cmd struct {
	cobra.Command

	config string
}

func NewCmd() *Cmd {
	c := &Cmd{
		Command: cobra.Command{
			Use: filepath.Base(os.Args[0]),
		},
	}

	c.Flags().StringVarP(&c.config, "config", "c", "config.yaml", "config file path")
	c.MarkFlagFilename("config")

	c.Command.Run = c.Run
	return c
}

func (c *Cmd) Run(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfigFromFile(c.config)
	if err != nil {
		c.PrintErrln(err)
		os.Exit(1)
	}

	bot := &certbot.CertBot{
		Cfg: cfg,
	}

	if err := bot.Run(); err != nil {
		c.PrintErrln(err)
		os.Exit(1)
	}
}
