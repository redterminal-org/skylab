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
	if os.Getenv("TLS_CLIENT_HASH") == "" {
		if query != "login" {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("# Welcome to the Skylab Misfin Client")
			fmt.Println("=> /?login Login")
		} else {
			fmt.Print("60 Please use a client certificate\r\n")
		}
		return
	}
	username, err := tools.Authenticate()
	if err != nil {
		fmt.Print("20 text/gemini\r\nError with somethin, prolly not ur fault :)")
		return
	}
	if username == "" {
		fmt.Print("20 text/gemini\r\nDidn't find you anywhere in my list, perhaps you're not supposed to be here?\n")
		return
	}
	if query == "login" {
		fmt.Print("30 /\r\n")
		return
	}
    cfg, err := config.ExportConfig(username)
	if err != nil {
		fmt.Println("Couldn't get the config info:", err)
		return
	}
	if cfg.AutoFetch || query == "fetch" {
		err := fetch.FetchGembox(cfg.GeminiFetchAddress, cfg.CertFilename, username)
		if err != nil {
    		fmt.Print("20 text/gemini\r\n")
			fmt.Println("Error fetching new mail:", err)
			fmt.Println("Cert filename is:", cfg.CertFilename)
			return
		} else {
			fmt.Print("30 /\r\n")
			return
		}
	}
    fmt.Print("20 text/gemini\r\n")
	fmt.Println("# Misfin ✉️")
	if err := tools.WalkTag("Inbox", username, false, 4); err != nil {
		fmt.Println("Error:", err)
        return
	}
	inboxCount, err := tools.CountTags("Inbox", username)
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
	fmt.Printf("=> inbox 📪 Inbox%s\n", inboxCount)
//	fmt.Printf("=> archive 🗄️ Archive%s\n", archiveCount)
	fmt.Println("=> archive 🗄️ Archive")
//	fmt.Printf("=> sent 📤 Sent%s\n", sentCount)
	fmt.Println("=> sent ☞ Sent")
	fmt.Printf("=> drafts 📝 Drafts%s\n", draftsCount)
//	fmt.Printf("=> trash 🗑️ Trash%s\n", trashCount)
	fmt.Println("=> trash 🗑️ Trash")
	fmt.Println("=> tags 🏷️ Tags")
	if !cfg.AutoFetch {
		fmt.Println("=> /?fetch 🔃 Check for New Mail")
	}

}
