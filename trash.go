package main

import (
    "fmt"
	"os"
//	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
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
		if trashed, _ := tools.HasTag(query, "Trash"); !trashed {
			if err := tools.Tag(query, "Trash"); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Error: Could not move to Trash")
				fmt.Println("=> inbox 📪 Inbox")
				return
			}
			fmt.Print("30 /\r\n")
		} else {
			if err := tools.Untag(query, "Trash"); err != nil {
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
		if err := tools.WalkTag("Trash", true); err != nil {
			fmt.Println(err)
			return
		}

	}
}
