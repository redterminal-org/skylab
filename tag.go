package main

import (
    "fmt"
	"strings"
	"os"
	"git.sr.ht/~satchlj/skylab/tools"
)

func main() {
	query := os.Getenv("QUERY_STRING")
	path := os.Getenv("PATH_INFO")
	msgid := strings.TrimPrefix(path, "/tag/")
	filename := fmt.Sprintf("../STORE/%s.msfn", msgid)
	_, err := os.Stat(filename)
	if err != nil {
		fmt.Println("30 /\r\n")
		return
	}
	reservedTags := []string{"", "Inbox", "Archive", "Trash", "Sent", "Unread", "Drafts", "new"}
	var reserved bool
	for _, tag := range reservedTags {
		if query == tag {
			reserved = true
			break
		}
	}
	if !reserved {
		if err := tools.Tag(msgid, query); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Printf("Error: Could not tag as %s.\n", query)
			fmt.Println("=> inbox 📪 Inbox")
		} else {
			fmt.Printf("30 /tags?%s\r\n", query)
			return
		}
	} else {
		if query == "" {
			fmt.Print("20 text/gemini\r\n")
			options, err := tools.ListTags()
			if err != nil {
				fmt.Println("Error: ", err)
			}
			for _, option := range options {
				fmt.Printf("=> /tag/%s?%s 🏷️ Assign tag %s\n", msgid, option, option)
			}
			fmt.Printf("=> /tag/%s?new 🏷️ Assign new tag\n", msgid)
			return
		} else if query == "new" {
			fmt.Print("10 Enter name for new tag (you may use letters, digits, underscores, and dashes):\r\n")
		} else {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("The tag", query, "is reserved.")
			fmt.Println("=> tags 🏷️ Tags")
		}
	}
}
