package main

import (
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/FurqanSoftware/blink/pipe"
	"github.com/adrg/frontmatter"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve serves scraped programming language references",
	Run: func(cmd *cobra.Command, args []string) {
		exitCh := make(chan struct{})
		errCh := make(chan error)

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "serve/style.css")
		})
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			f, err := os.Open(filepath.Join("out", r.URL.Path, "index.html"))
			if errors.Is(err, os.ErrNotExist) {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
			tpl, err := template.ParseGlob("serve/*.html")
			catch(err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			meta := pipe.Meta{}
			body, err := frontmatter.Parse(f, &meta)
			catch(err)
			err = tpl.Execute(w, struct {
				Meta pipe.Meta
				Body template.HTML
			}{
				Meta: meta,
				Body: template.HTML(body),
			})
			catch(err)
		})

		go func() {
			err := http.ListenAndServe(":8080", nil)
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
		}

		log.Print("Bye")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
