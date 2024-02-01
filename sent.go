package main

import (
    "fmt"
//	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
)

func main() {
    fmt.Print("20 text/gemini\r\n")
	fmt.Println("# Sent 📤")
	if err := tools.WalkTag("Sent", true); err != nil {
		fmt.Println("Error:", err)
        return
	}
}
