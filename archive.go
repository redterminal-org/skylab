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
		if err := tools.Tag(query, "Archive"); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("Error: Could not move to Archive")
			fmt.Println("=> inbox 📪 Inbox")
		} else {
			if err := tools.Untag(query, "Inbox"); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Error: Could not remove from Inbox")
				fmt.Println("=> inbox 📪 Inbox")
				return
			}
			fmt.Print("30 /\r\n")
		}
	} else {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("# Archive 🗄")
		if err := tools.WalkTag("Archive", true); err != nil {
			fmt.Println("Error:", err)
			return
		}

	}
}
