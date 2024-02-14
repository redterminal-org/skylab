package main

import (
    "fmt"
//	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
)

func drafts() {
	username, _ := tools.Authenticate()
    fmt.Print("20 text/gemini\r\n")
	fmt.Println("# Drafts 📝")
	if err := tools.WalkTag("Drafts", username, false); err != nil {
		fmt.Println("Error:", err)
        return
	}
}
