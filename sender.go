package main

import (
	"fmt"
	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/send"
	"git.sr.ht/~satchlj/skylab/tools"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func sender() {
	username, _ := tools.Authenticate()
	msgid := os.Getenv("QUERY_STRING")
	if strings.Contains(msgid, "..") {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("That query is not allowed")
		return
	}
	filename := filepath.Join("..", "users", username, "STORE", msgid+".msfn")
	cfg, err := config.ExportConfig(username)
	if err != nil {
		fmt.Println("Couldn't get the config info")
		return
	}
	sitecfg, err := config.ExportSiteConfig()
	if err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Couldn't export site-wide config", err)
		return
	}
	skylabPath := appPrefix(sitecfg.SkylabPath)
	message, err := tools.ParseSkymail(filename)
	if err != nil {
		fmt.Printf("30 %s\r\n", appPath(skylabPath, ""))
		return
	}
	message.MessageBody = strings.ReplaceAll(message.MessageBody, "\n  \n", "\n\n")
	message.Timestamps[0] = time.Now().UTC()
	if err := send.Skymail(message, cfg.CertFilename); err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println(err)
		fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
		return
	}
	if err := os.WriteFile(filename, []byte(send.SkymailToString(message)), 0644); err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println(err)
		fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
		return
	}
	if err := tools.Tag(msgid, "Sent", username, false); err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Error tagging as sent:", err)
		fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
		return
	}
	if err := tools.Untag(msgid, "Drafts", username, false); err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Error untagging as draft:", err)
		fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
		return
	}
	fmt.Printf("30 %s?%s\r\n", appPath(skylabPath, "mail"), msgid)
}
