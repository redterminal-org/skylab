package main

import (
    "fmt"
	"strings"
	"os"
	"unicode"
	"path/filepath"
	"git.sr.ht/~satchlj/skylab/tools"
)

func main() {
	if os.Getenv("TLS_CLIENT_HASH") == "" {
		fmt.Print("30 /\r\n")
		return
	}
	username, err := tools.Authenticate()
	if err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Error with somethin, prolly not ur fault :)")
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
	path := os.Getenv("PATH_INFO")
	if strings.Contains(path, "..") {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("That query is not allowed")
	}
	msgid := strings.TrimPrefix(path, "/tag/")
	filename := filepath.Join("..", "users", username, "STORE", msgid+".msfn")
	validTag := true
	for _, char := range query {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' && char != '-' {
			validTag = false
		}
	}
	if _, err := os.Stat(filename); err != nil || validTag != true {
		fmt.Print("30 /\r\n")
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
		if tagged, _ := tools.HasTag(msgid, query, username); !tagged {
			if err := tools.Tag(msgid, query, username); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Printf("Error: Could not tag as %s.\n", query)
				fmt.Println("=> inbox 📪 Inbox")
			} else {
				fmt.Printf("30 /tags?%s\r\n", query)
				return
			}
		} else {
			if err := tools.Untag(msgid, query, username); err != nil {
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
			options, err := tools.ListTags(username)
			if err != nil {
				fmt.Println("error: ", err)
			}
			for _, option := range options {
				var actionName string
				if tagged, _ := tools.HasTag(msgid, option, username); !tagged {
					actionName = "Assign"
				} else {
					actionName = "Remove"
				}
				fmt.Printf("=> /tag/%s?%s 🏷️ %s tag %s\n", msgid, option, actionName, option)
			}
			fmt.Printf("=> /tag/%s?new 🏷️ Assign new tag\n", msgid)
			return
		} else if query == "new" {
			fmt.Print("10 Enter name for new tag (you may use letters, digits, underscores, and dashes):\r\n")
		} else {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("the tag", query, "is reserved.")
			fmt.Println("=> tags 🏷️ tags")
		}
	}
}
