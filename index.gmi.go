package main

import (
    "fmt"
	"git.sr.ht/~satchlj/skylab/fetch"
	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
)

func main() {
    cfg := config.ExportConfig()
    fmt.Print("20 text/gemini\r\n")
	fmt.Println("# Misfin ✉️")
    err := fetch.FetchGembox(cfg.GeminiFetchAddress, cfg.CertFilename)
    if err != nil {
        fmt.Println("Error fetching new mail:")
        fmt.Println(err)
    }
	if err := tools.WalkTag("Inbox"); err != nil {
		fmt.Println("Error:", err)
        return
	}
}
