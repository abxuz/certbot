package main

import (
	"certbot/internal/cmd"

	_ "certbot/provider/bdns"

	_ "certbot/reciever/bvhost"
)

func main() {
	cmd.NewCmd().Execute()
}
