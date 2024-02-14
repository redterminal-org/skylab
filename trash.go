package main

import (
    "fmt"
	"strings"
	"os"
	"path/filepath"
//	"git.sr.ht/~satchlj/skylab/config"
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
	if query != "" {
		filename := filepath.Join("..", "users", username, "STORE", query+".msfn")
		_, err := os.Stat(filename)
		if err != nil {
			fmt.Println("30 /\r\n")
			return
		}
		if trashed, _ := tools.HasTag(query, "Trash", username); !trashed {
			if err := tools.Tag(query, "Trash", username); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Error: Could not move to Trash")
				fmt.Println("=> inbox 📪 Inbox")
				return
			}
			fmt.Print("30 /\r\n")
		} else {
			if err := tools.Untag(query, "Trash", username); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Error: Could not remove from Trash")
				fmt.Println("=> inbox 📪 Inbox")
				return
			}
			fmt.Print("30 /trash\r\n")
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
