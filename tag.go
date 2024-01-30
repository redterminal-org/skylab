package main

import (
    "fmt"
	"strings"
	"os"
	"unicode"
	"git.sr.ht/~satchlj/skylab/tools"
)

func main() {
	query := os.Getenv("QUERY_STRING")
	path := os.Getenv("PATH_INFO")
	msgid := strings.TrimPrefix(path, "/tag/")
	filename := fmt.Sprintf("../STORE/%s.msfn", msgid)
	validTag := true
	for _, char := range query {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' && char != '-' {
			validTag = false
		}
	}
	_, err := os.Stat(filename)
	if err != nil || validTag != true {
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
		if tagged, _ := tools.HasTag(msgid, query); !tagged {
			if err := tools.Tag(msgid, query); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Printf("Error: Could not tag as %s.\n", query)
				fmt.Println("=> inbox 📪 Inbox")
			} else {
				fmt.Printf("30 /tags?%s\r\n", query)
				return
			}
		} else {
			if err := tools.Untag(msgid, query); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Printf("Error: Could not tag as %s.\n", query)
				fmt.Println("=> inbox 📪 Inbox")
			} else {
				fmt.Printf("30 /mail?%s\r\n", msgid)
			}
		}
	} else {
		if query == "" {
			fmt.Print("20 text/gemini\r\n")
			options, err := tools.ListTags()
			if err != nil {
				fmt.Println("error: ", err)
			}
			for _, option := range options {
				var actionName string
				if tagged, _ := tools.HasTag(msgid, option); !tagged {
					actionName = "Assign"
				} else {
					actionName = "Remove"
				}
				fmt.Printf("=> /tag/%s?%s 🏷️ %s tag %s\n", msgid, option, actionName, option)
			}
			fmt.Printf("=> /tag/%s?new 🏷️ Assign new tag\n", msgid)
			return
		} else if query == "new" {
			fmt.Print("10 enter name for new tag (you may use letters, digits, underscores, and dashes):\r\n")
		} else {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("the tag", query, "is reserved.")
			fmt.Println("=> tags 🏷️ tags")
		}
	}
}
