package main

import (
	"fmt"
	"os"
	//	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
	"path/filepath"
)

func main() {
	query := os.Getenv("QUERY_STRING")
	if query != "" {
		filename := fmt.Sprintf("../STORE/%s.msfn", query)
		_, err := os.Stat(filename)
		if err != nil {
			fmt.Println("30 /\r\n")
			return
		}
		if err := tools.UntagAll(query); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("Error: Could not remove from somewhere")
			fmt.Println(err)
			fmt.Println("=> inbox 📪 Inbox")
			return
		} else {
			filename := filepath.Join("..", "STORE", query+".msfn")
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
