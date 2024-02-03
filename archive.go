package main

import (
    "fmt"
	"os"
    "path/filepath"
//	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
)

func main() {
	if os.Getenv("TLS_CLIENT_HASH") == "" {
		fmt.Print("30 /\r\n")
		return
	}
	username, err := tools.Authenticate()
	if err != nil {
		fmt.Print("20 text/gemini\r\nError with somethin, prolly not ur fault :)")
		return
	}
	if username == "" {
		fmt.Print("30 /\r\n")
		return
	}
	query := os.Getenv("QUERY_STRING")
	if query != "" {
		filename := filepath.Join("..", "users", username, "STORE", query+".msfn")
		_, err := os.Stat(filename)
		if err != nil {
			fmt.Println("30 /\r\n")
			return
		}
		if archived, _ := tools.HasTag(query, "Archive", username); !archived {
			if err := tools.Tag(query, "Archive", username); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Error: Could not move to Archive")
				fmt.Println("=> inbox 📪 Inbox")
			} else {
				if err := tools.Untag(query, "Inbox", username); err != nil {
					fmt.Print("20 text/gemini\r\n")
					fmt.Println("Error: Could not remove from Inbox")
					fmt.Println("=> inbox 📪 Inbox")
					return
				}
				fmt.Print("30 /\r\n")
			}
		} else {
			if err := tools.Tag(query, "Inbox", username); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Error: Could not move to Inbox")
				fmt.Println("=> inbox 📪 Inbox")
			} else {
				if err := tools.Untag(query, "Archive", username); err != nil {
					fmt.Print("20 text/gemini\r\n")
					fmt.Println("Error: Could not remove from Archive")
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
