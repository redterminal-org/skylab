package main

import (
    "fmt"
	"os"
	"git.sr.ht/~satchlj/skylab/fetch"
	"git.sr.ht/~satchlj/skylab/config"
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
    cfg, err := config.ExportConfig(username)
	if err != nil {
		fmt.Println("Couldn't get the config info")
		return
	}
    fmt.Print("20 text/gemini\r\n")
	fmt.Println("# Inbox ✉️")
	if cfg.AutoFetch {
		err := fetch.FetchGembox(cfg.GeminiFetchAddress, cfg.CertFilename, username)
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
