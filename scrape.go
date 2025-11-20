package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/FurqanSoftware/blink/site"
	"github.com/spf13/cobra"
)

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape scrapes programming language references",
	Run: func(cmd *cobra.Command, args []string) {
		exitCh := make(chan struct{})
		errCh := make(chan error)

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		s := site.Get(args[0])
		go func() {
			err := site.Scraper{}.Run(ctx, s)
			if err != nil {
				errCh <- err
			}
			close(exitCh)
		}()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		select {
		case err := <-errCh:
			log.Fatal(err)

		case <-exitCh:

		case sig := <-sigCh:
			log.Printf("Received %s; exiting", sig)
			cancel()
			select {
			case <-exitCh:
			case <-sigCh:
			}
		}

		log.Print("Bye")
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
}
