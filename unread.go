package main

import (
    "fmt"
	"os"
//	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
    "path/filepath"
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
		if err := tools.Tag(query, "Unread", username); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("Error: Could not mark as unread")
			fmt.Println("=> inbox 📪 Inbox")
		} else {
			fmt.Print("30 /\r\n")
		}
	} else {
		fmt.Print("30 /\r\n")
	}
}
