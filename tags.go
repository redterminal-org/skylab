package main

import (
    "fmt"
	"os"
//	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
)

func tags() {
	username, _ := tools.Authenticate()
	query := os.Getenv("QUERY_STRING")
	options, err := tools.ListTags(username)
	if query != "" {
		realTag := false
		for _, option := range options {
			if query == option {
				realTag = true
				break
			}
		}
		if !realTag {
			fmt.Print("30 tags\r\n")
			return
		}
		if err != nil {
			fmt.Println("error: ", err)
		}
		fmt.Print("20 text/gemini\r\n")
		fmt.Printf("# Tag: %s 🏷️\n", query)
		if err := tools.WalkTag(query, username, true); err != nil {
			fmt.Println("Error:", err)
			return
		}
	} else {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("# Tags 🏷️")
		for _, option := range options {
			tagCount, _ := tools.CountTags(option, username)
			fmt.Printf("=> /tags?%s 🏷️ %s%s\n", option, option, tagCount)
		}
	}
}
