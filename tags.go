package main

import (
    "fmt"
	"os"
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
