package main

import (
    "fmt"
	"strings"
	"os"
	"path/filepath"
	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
)

func trash() {
	username, _ := tools.Authenticate()
	query := os.Getenv("QUERY_STRING")
	if strings.Contains(query, "..") {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("That query is not allowed")
		return
	}
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
		if trashed, _ := tools.HasTag(query, "Trash", username); !trashed {
			if err := tools.Tag(query, "Trash", username, true); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Error: Could not move to Trash:", err)
				fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
				return
			}
			fmt.Printf("30 /%s\r\n", skylabPath)
		} else {
			if err := tools.Untag(query, "Trash", username, true); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Error: Could not remove from Trash:", err)
				fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
				return
			}
			fmt.Printf("30 %s/trash\r\n", skylabPath)
		}
	} else {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("# Trash 🗑️")
		if err := tools.WalkTag("Trash", username, true); err != nil {
			fmt.Println(err)
			return
		}

	}
}
