package main

import (
	"log"

	_ "github.com/FurqanSoftware/blink/site/com.cppreference"
	_ "github.com/FurqanSoftware/blink/site/dev.golang.pkg"
	_ "github.com/FurqanSoftware/blink/site/org.python.docs"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "blink",
	Short: "Blink is a programming language documentation and reference scraping and mirroring tool",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}
