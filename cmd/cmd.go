package cmd

import (
	"certbot/internal/certbot"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var config string

	c := &cobra.Command{
		Use: filepath.Base(os.Args[0]),
		Run: func(cmd *cobra.Command, args []string) {
			cb := new(certbot.CertBot)
			if err := cb.Init(config); err != nil {
				log.Fatal(err)
			}
			cb.Serve()
		},
	}

	c.Flags().StringVarP(&config, "config", "c", "config.yaml", "config file path")
	c.MarkFlagFilename("config")
	return c
}
