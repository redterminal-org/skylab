package main

import (
    "fmt"
	"strings"
	"time"
	"os"
    "path/filepath"
//	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
)

func archive() {
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
		if archived, _ := tools.HasTag(query, "Archive", username); !archived {
			if err := tools.Tag(query, "Archive", username, true); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Error: Could not move to Archive")
				fmt.Println("=> inbox 📪 Inbox")
			} else {
				time.Sleep(50 * time.Millisecond)
				if err := tools.Untag(query, "Inbox", username, true); err != nil {
					fmt.Print("20 text/gemini\r\n")
					fmt.Println("Error: Could not remove from Inbox", err)
					fmt.Println("=> inbox 📪 Inbox")
					return
				}
				fmt.Print("30 /\r\n")
			}
		} else {
			if err := tools.Tag(query, "Inbox", username, true); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Error: Could not move to Inbox")
				fmt.Println("=> inbox 📪 Inbox")
			} else {
				time.Sleep(50 * time.Millisecond)
				if err := tools.Untag(query, "Archive", username, true); err != nil {
					fmt.Print("20 text/gemini\r\n")
					fmt.Println("Error: Could not remove from Archive", err)
					fmt.Println("=> inbox 📪 Inbox")
					return
				}
				fmt.Print("30 /inbox\r\n")
			}
		}
	} else {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("# Archive 🗄")
		if err := tools.WalkTag("Archive", username, true); err != nil {
			fmt.Println("Error:", err)
			return
		}

	}
}
