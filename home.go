package main

import (
    "fmt"
	"os"
	"git.sr.ht/~satchlj/skylab/fetch"
	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
)

func home() {
	query := os.Getenv("QUERY_STRING")
	username, _ := tools.Authenticate()
    cfg, err := config.ExportConfig(username)
	if err != nil {
		fmt.Println("Couldn't get the config info:", err)
		return
	}
	sitecfg, err := config.ExportSiteConfig()
	if err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Couldn't export site-wide config", err)
		return
	}
	skylabPath := sitecfg.SkylabPath
	var fetchGembox func(username string) error
	if cfg.GMAP {
		fetchGembox = fetch.FetchGemboxGMAP
	} else {
		fetchGembox = fetch.FetchGembox
	}
	if cfg.AutoFetch || query == "fetch" {
		if err := fetchGembox(username); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("Error fetching new mail:", err)
			fmt.Println("Cert filename is:", cfg.CertFilename)
			return
		}
	}

	if query == "fetch" {
		fmt.Printf("30 /%s\r\n", skylabPath)
	}
	fmt.Print("20 text/gemini\r\n")
	if cfg.HeaderText != "" {
		fmt.Println("# ", cfg.HeaderText)
	} else {
		fmt.Println("# Misfin ✉️")
	}
	display := cfg.DisplayNumber
	if err := tools.WalkTag("Inbox", username, false, display); err != nil {
		fmt.Println("Error:", err)
        return
	}
	inboxCount, err := tools.CountTags("Unread", username)
	if err != nil {
        fmt.Println("Error counting Inbox:")
        fmt.Println(err)
    }
//	archiveCount, err := tools.CountTags("Archive")
//	if err != nil {
//        fmt.Println("Error counting Archive:")
//        fmt.Println(err)
//    }
//	sentCount, err := tools.CountTags("Sent")
//	if err != nil {
//        fmt.Println("Error counting Sent:")
//        fmt.Println(err)
//    }
	draftsCount, err := tools.CountTags("Drafts", username)
	if err != nil {
        fmt.Println("Error counting Drafts:")
        fmt.Println(err)
	}
//	trashCount, err := tools.CountTags("Trash")
//	if err != nil {
//		fmt.Println("Error counting Trash:")
//		fmt.Println(err)
//	}
	fmt.Printf("=> %s/inbox 📪 Inbox%s\n", skylabPath, inboxCount)
	fmt.Printf("=> %s/compose/new ✍️ Compose\n", skylabPath)
//	fmt.Printf("=> archive 🗄️ Archive%s\n", archiveCount)
	fmt.Printf("=> %s/archive 🗄️ Archive\n", skylabPath)
//	fmt.Printf("=> sent 📤 Sent%s\n", sentCount)
	fmt.Printf("=> %s/sent ☞ Sent\n", skylabPath)
	fmt.Printf("=> %s/drafts 📝 Drafts%s\n", skylabPath, draftsCount)
//	fmt.Printf("=> trash 🗑️ Trash%s\n", trashCount)
	fmt.Printf("=> %s/trash 🗑️ Trash\n", skylabPath)
	fmt.Printf("=> %s/tags 🏷️ Tags\n", skylabPath)
	fmt.Printf("=> %s/settings ⚙️ Settings\n", skylabPath)
	if !cfg.AutoFetch {
		fmt.Printf("=> %s?fetch 🔃 Check for New Mail", skylabPath)
	}
}
