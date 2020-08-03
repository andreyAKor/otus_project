package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andreyAKor/otus_project/internal/app"
	"github.com/andreyAKor/otus_project/internal/cache"
	"github.com/andreyAKor/otus_project/internal/configs"
	"github.com/andreyAKor/otus_project/internal/http/client"
	"github.com/andreyAKor/otus_project/internal/http/server"
	"github.com/andreyAKor/otus_project/internal/image"
	"github.com/andreyAKor/otus_project/internal/logging"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "image_previewer",
	Short: "Image previewer service application",
	Long:  "The image previewer service is the most simplified service for storing resized images from another web-resources and storing this resized images in own file lru-cache.",
	RunE:  run,
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&cfgFile, "config", "", "config file")
	if err := cobra.MarkFlagRequired(pf, "config"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//nolint:funlen
func run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init config
	c := new(configs.Config)
	if err := c.Init(cfgFile); err != nil {
		return errors.Wrap(err, "init config failed")
	}

	// Init logger
	l := logging.New(c.Logging.File, c.Logging.Level)
	if err := l.Init(); err != nil {
		return errors.Wrap(err, "init logging failed")
	}

	// Init cache
	cache, err := cache.New(c.Cache.Capacity)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	// Init image
	image, err := image.New(c.Image.MaxWidth, c.Image.MaxHeight)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	// Init http-client
	client, err := client.New(c.Client.Timeout, c.Client.BodyLimit)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize http-client")
	}

	// Init http-server
	server, err := server.New(client, image, cache, c.HTTP.Host, c.HTTP.Port, c.HTTP.BodyLimit)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize http-server")
	}

	// Init and run app
	a, err := app.New(server)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize app")
	}
	if err := a.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("app runnign fail")
	}

	// Graceful shutdown
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)
	<-interruptCh

	log.Info().Msg("Stopping...")

	if err := server.Close(); err != nil {
		log.Fatal().Err(err).Msg("http-server closing fail")
	}
	if err := cache.Close(); err != nil {
		log.Fatal().Err(err).Msg("cache closing fail")
	}
	if err := a.Close(); err != nil {
		log.Fatal().Err(err).Msg("app closing fail")
	}

	log.Info().Msg("Stopped")

	if err := l.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}
