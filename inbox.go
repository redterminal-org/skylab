package main

import (
	"fmt"
	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/fetch"
	"git.sr.ht/~satchlj/skylab/tools"
)

func inbox() {
	username, _ := tools.Authenticate()
	cfg, err := config.ExportConfig(username)
	if err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Couldn't get the config info")
		return
	}
	fmt.Print("20 text/gemini\r\n")
	fmt.Println("# Inbox ✉️")
	var fetchGembox func(username string) error
	if cfg.GMAP {
		fetchGembox = fetch.GemboxGMAP
	} else {
		fetchGembox = fetch.Gembox
	}
	if cfg.AutoFetch {
		err := fetchGembox(username)
		if err != nil {
			fmt.Println("Error fetching new mail:")
			fmt.Println(err)
		}
	}
	if err := tools.WalkTag("Inbox", username, true); err != nil {
		fmt.Println("Error:", err)
		return
	}
}
