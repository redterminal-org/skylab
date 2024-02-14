package main

import (
    "fmt"
	"git.sr.ht/~satchlj/skylab/tools"
)

func sent() {
	username, _ := tools.Authenticate()
    fmt.Print("20 text/gemini\r\n")
	fmt.Println("# Sent 📤")
	if err := tools.WalkTag("Sent", username, true); err != nil {
		fmt.Println("Error:", err)
        return
	}
}
