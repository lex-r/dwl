package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/lex-r/dwl/config"
	"github.com/lex-r/dwl/downloader"
	"github.com/lex-r/dwl/saver/fs"
)

func main() {
	cfg := config.NewConfig()
	err := cfg.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	svr := fs.NewSaver(cwd)

	dwl, err := downloader.NewDownloader(svr,
		downloader.WithTimeout(cfg.Timeout),
		downloader.WithThreads(cfg.Threads),
		downloader.WithUserAgent(cfg.UserAgent),
	)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	defer func() {
		signal.Stop(signalCh)
		cancel()
	}()
	go func() {
		select {
		case <-signalCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	dwl.Download(ctx, cfg.Links)
}
