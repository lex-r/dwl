package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/lex-r/dwl/downloader"
)

// Config represents config for application.
type Config struct {
	Threads   int
	Timeout   time.Duration
	UserAgent string
	Links     []string
	fs        *flag.FlagSet
}

// NewConfig create new instance of Config.
func NewConfig() *Config {
	cfg := &Config{
		fs: flag.NewFlagSet(os.Args[0], flag.ExitOnError),
	}

	cfg.fs.IntVar(&cfg.Threads, "threads", downloader.DefaultThreadOpt, "number of threads")
	cfg.fs.DurationVar(&cfg.Timeout, "timeout", downloader.DefaultTimeoutOpt, "timeout for downloading")
	cfg.fs.StringVar(&cfg.UserAgent, "user-agent", downloader.DefaultUserAgent(), "send user-agent to server")

	return cfg
}

// Parse takes args, parses it and validates.
func (c *Config) Parse(args []string) error {
	err := c.fs.Parse(args)
	if err != nil {
		return err
	}

	c.Links = c.fs.Args()

	if err := c.validate(); err != nil {
		return err
	}

	return nil
}

func (c *Config) validate() error {
	if c.Threads < 1 {
		return fmt.Errorf("invalid value for flag -threads")
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("invalid value for flag -timeout")
	}

	return nil
}
