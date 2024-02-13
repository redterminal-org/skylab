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
		if query == "" {
			fmt.Print("20 text/gemini\r\nDidn't find you anywhere in my list. If you have an account already, then you can log in with a single-use password:\n")
			fmt.Println("=> /?login Login")
			return
		} else if query == "login" {
			fmt.Print("11 Enter your single-use password:\r\n")
			return
		} else {
			if err := tools.CheckTempPassword(query); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Something went wrong trying to add your cert:", err)
				return
			}
			fmt.Print("30 /\r\n")
		}
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
		fmt.Print("30 /\r\n")
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
	fmt.Printf("=> inbox 📪 Inbox%s\n", inboxCount)
	fmt.Println("=> compose/new ✍️ Compose")
//	fmt.Printf("=> archive 🗄️ Archive%s\n", archiveCount)
	fmt.Println("=> archive 🗄️ Archive")
//	fmt.Printf("=> sent 📤 Sent%s\n", sentCount)
	fmt.Println("=> sent ☞ Sent")
	fmt.Printf("=> drafts 📝 Drafts%s\n", draftsCount)
//	fmt.Printf("=> trash 🗑️ Trash%s\n", trashCount)
	fmt.Println("=> trash 🗑️ Trash")
	fmt.Println("=> tags 🏷️ Tags")
	fmt.Println("=> settings ⚙️ Settings")
	if !cfg.AutoFetch {
		fmt.Println("=> /?fetch 🔃 Check for New Mail")
	}
}
