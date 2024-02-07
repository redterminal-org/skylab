package main

import (
	"fmt"
	"strings"
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
		if err := tools.UntagAll(query, username); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("Error: Could not remove from somewhere")
			fmt.Println(err)
			fmt.Println("=> inbox 📪 Inbox")
			return
		} else {
			filename := filepath.Join("..", "users", username, "STORE", query+".msfn")
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("File was already gone.")
				fmt.Println("=> inbox 📪 Inbox")
				return
			}

			err := os.Remove(filename)
			if err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Could not remove file.")
				fmt.Println("=> inbox 📪 Inbox")
				return
			}

			fmt.Print("30 /trash\r\n")
		}
	} else {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("# Empty Graveyard 🪦")
		fmt.Println("There is no trace left of your deleted messages. I'd return home if I were you:")
		fmt.Println("=> / 🏠 Home")
	}
}
