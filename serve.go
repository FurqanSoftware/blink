package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/FurqanSoftware/blink/pipe"
	"github.com/adrg/frontmatter"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve serves scraped programming language references",
	Run: func(cmd *cobra.Command, args []string) {
		errCh := make(chan error)

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

		srv := &http.Server{
			Addr:    fmt.Sprintf(":8080"),
			Handler: nil,
		}
		go func() {
			log.Printf("Listening on %s", srv.Addr)
			err := srv.ListenAndServe()
			if err != nil {
				errCh <- err
			}
		}()

		sigch := make(chan os.Signal, 2)
		signal.Notify(sigch, os.Interrupt)

		select {
		case err := <-errCh:
			log.Fatal(err)

		case sig := <-sigch:
			log.Printf("Received %s; exiting", sig)
		}

		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(7*time.Second))
		defer cancel()
		err := srv.Shutdown(ctx)
		catch(err)

		log.Print("Bye")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
