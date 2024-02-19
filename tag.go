package main

import (
    "fmt"
	"strings"
	"os"
	"unicode"
	"path/filepath"
	"git.sr.ht/~satchlj/skylab/tools"
	"git.sr.ht/~satchlj/skylab/config"
)

func tag() {
	username, _ := tools.Authenticate()
	query := os.Getenv("QUERY_STRING")
	if strings.Contains(query, "..") {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("That query is not allowed")
		return
	}
	sitecfg, err := config.ExportSiteConfig()
	if err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Couldn't export site-wide config", err)
		return
	}
	skylabPath := sitecfg.SkylabPath
	path := os.Getenv("PATH_INFO")
	if strings.Contains(path, "..") {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("That query is not allowed")
	}
	var msgid string
	if skylabPath != "" {
		skylabPath = "/"+skylabPath
		msgid = strings.TrimPrefix(path, skylabPath+"/tag/")
	} else {
		msgid = strings.TrimPrefix(path, "/tag/")
	}
	filename := filepath.Join("..", "users", username, "STORE", msgid+".msfn")
	validTag := true
	for _, char := range query {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' && char != '-' {
			validTag = false
		}
	}
	if _, err := os.Stat(filename); err != nil || validTag != true {
		fmt.Print("20 text/gemini\r\n")
		fmt.Printf("Tag %s is not allowed.\n", query)
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
			if err := tools.Tag(msgid, query, username, true); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Printf("Error: Could not tag as %s.\n", query)
				fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
			} else {
				fmt.Printf("30 %s/tag/%s\r\n", skylabPath, msgid)
				return
			}
		} else {
			if err := tools.Untag(msgid, query, username, true); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Printf("Error: Could not tag as %s.\n", query)
				fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
			} else {
				fmt.Printf("30 %s/tag/%s\r\n", skylabPath, msgid)
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
				fmt.Printf("=> %s/tag/%s?%s 🏷️ %s tag %s\n", skylabPath, msgid, option, actionName, option)
			}
			fmt.Printf("=> %s/tag/%s?new 🏷️ Assign new tag\n", skylabPath, msgid)
			fmt.Printf("=> %s/mail?%s ◎ View message\n", skylabPath, msgid)
			return
		} else if query == "new" {
			fmt.Print("10 Enter name for new tag (you may use letters, digits, underscores, and dashes):\r\n")
		} else {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("The tag", query, "is reserved.")
			fmt.Printf("=> %s/tags 🏷️ tags", skylabPath)
		}
	}
}
