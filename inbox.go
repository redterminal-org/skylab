package main

import (
    "fmt"
	"git.sr.ht/~satchlj/skylab/fetch"
	"git.sr.ht/~satchlj/skylab/config"
)

func main() {
    cfg := config.ExportConfig()
    fmt.Print("20 text/gemini\r\n")
    err := fetch.FetchGembox(cfg.GeminiFetchAddress, cfg.CertFilename)
    if err != nil {
        fmt.Println("Error fetching new mail:")
        fmt.Println(err)
    }
    fmt.Println("=> vviw7643ftwnw4fn ◉ Unread Email Subject - sender0@host.org")
    fmt.Println("Example text from the beginning of the most recently received email...")
}
