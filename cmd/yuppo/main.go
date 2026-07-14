package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"yuppo/internal/config"
	"yuppo/internal/debouncer"
	"yuppo/internal/gitignore"
	"yuppo/internal/gitops"
	"yuppo/internal/watcher"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	matcher, err := gitignore.Load(cfg.WatchPath)
	if err != nil {
		log.Fatalf("gitignore: %v", err)
	}

	d := debouncer.New(cfg.DebounceTimeout)
	defer d.Stop()

	commit := func() {
		d.Trigger(func() {
			if err := gitops.AutoCommit(cfg.WatchPath, cfg.CommitMessage, cfg.AutoPush); err != nil {
				log.Printf("commit error: %v", err)
			}
		})
	}

	w, err := watcher.New(cfg.WatchPath, matcher, commit)
	if err != nil {
		log.Fatalf("watcher: %v", err)
	}
	defer w.Close()

	log.Printf("yuppo: watching %s", cfg.WatchPath)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("yuppo: shutting down")
}
