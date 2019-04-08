package main

import (
	"flag"
	"os"

	"github.com/linkchain/common/util/log"
	"github.com/robfig/cron"
	"hum.tan/server/server"
	"hum.tan/server/server/client"
)

func main() {
	logLevel := flag.Int("loglevel", 3, "log level")

	//init log
	log.Root().SetHandler(
		log.LvlFilterHandler(log.Lvl(*logLevel),
			log.StreamHandler(os.Stdout, log.TerminalFormat(true))))

	go func() {
		c := cron.New()
		spec := "0 0/1 * * * ?"
		_ = c.AddFunc(spec, func() {
			client.Sync()
		})
		log.Info("Start sync block and transaction info")
		c.Start()
	}()

	s := server.NewServer()
	s.Start()
}
