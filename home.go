package main

import (
    "fmt"
	"os"
	"git.sr.ht/~satchlj/skylab/fetch"
	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
)

func main() {
	query := os.Getenv("QUERY_STRING")
    cfg := config.ExportConfig()
    fmt.Print("20 text/gemini\r\n")
	fmt.Println("# Misfin ✉️")
	if cfg.AutoFetch || query == "fetch" {
		err := fetch.FetchGembox(cfg.GeminiFetchAddress, cfg.CertFilename)
		if err != nil {
			fmt.Println("Error fetching new mail:")
			fmt.Println(err)
		}
	}
	if err := tools.WalkTag("Inbox", false, 4); err != nil {
		fmt.Println("Error:", err)
        return
	}
	inboxCount, err := tools.CountTags("Inbox")
	if err != nil {
        fmt.Println("Error counting Inbox:")
        fmt.Println(err)
    }
	archiveCount, err := tools.CountTags("Archive")
	if err != nil {
        fmt.Println("Error counting Archive:")
        fmt.Println(err)
    }
	sentCount, err := tools.CountTags("Sent")
	if err != nil {
        fmt.Println("Error counting Sent:")
        fmt.Println(err)
    }
	draftsCount, err := tools.CountTags("Drafts")
	if err != nil {
        fmt.Println("Error counting Drafts:")
        fmt.Println(err)
	}
	trashCount, err := tools.CountTags("Trash")
	if err != nil {
		fmt.Println("Error counting Trash:")
		fmt.Println(err)
	}
	fmt.Printf("=> inbox 📪 Inbox%s\n", inboxCount)
	fmt.Printf("=> archive 🗄️ Archive%s\n", archiveCount)
	fmt.Printf("=> sent 📤 Sent%s\n", sentCount)
	fmt.Printf("=> drafts 📝 Drafts%s\n", draftsCount)
	fmt.Printf("=> trash 🗑️ Trash%s\n", trashCount)
	fmt.Println("=> tags 🏷️ Tags")
	if !cfg.AutoFetch {
		fmt.Println("=> /?fetch 🔃 Check for New Mail")
	}

}
