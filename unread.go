package main

import (
    "fmt"
	"os"
	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
    "path/filepath"
)

func unread() {
	username, _ := tools.Authenticate()
	query := os.Getenv("QUERY_STRING")
	sitecfg, err := config.ExportSiteConfig()
	if err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Couldn't export site-wide config", err)
		return
	}
	skylabPath := sitecfg.SkylabPath
	if query != "" {
		filename := filepath.Join("..", "users", username, "STORE", query+".msfn")
		_, err := os.Stat(filename)
		if err != nil {
			fmt.Printf("30 /%s\r\n", skylabPath)
			return
		}
		if err := tools.Tag(query, "Unread", username, true); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("Error: Could not mark as unread")
			fmt.Printf("=> /%s/inbox 📪 Inbox\n", skylabPath)
		} else {
			fmt.Printf("30 /%s\r\n", skylabPath)
		}
	} else {
		fmt.Printf("30 /%s\r\n", skylabPath)
	}
}
