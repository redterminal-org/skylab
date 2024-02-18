package main

import (
    "fmt"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path/filepath"
	"net/url"
	"git.sr.ht/~satchlj/skylab/tools"
	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/send"
)

func reply() {
	username, _ := tools.Authenticate()
	query := os.Getenv("QUERY_STRING")
	path := os.Getenv("PATH_INFO")
	if strings.Contains(path, "..") {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("That path is not allowed")
		return
	}
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
	skylabPath := sitecfg.SkylabPath
	var msgid string
	if skylabPath != "" {
		msgid = strings.TrimPrefix(path, "/"+skylabPath+"/reply/")
	} else {
		msgid = strings.TrimPrefix(path, "/reply/")
	}
	filename := filepath.Join("..", "users", username, "STORE", msgid+".msfn")
	message, err := tools.ParseSkymail(filename)
	if err != nil {
		fmt.Printf("30 /%s\r\n", skylabPath)
		return
	}
	if query != "" {
		address, err := send.GetAddressFromCert(cfg.CertFilename)
		if err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println(err)
			fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
			return
		}
		decodedQuery, err := url.QueryUnescape(query)
		if err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println(err)
			fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
			return
		}
		replyMessage := tools.Skymail{
			Senders:		[]string{address},
			Recipients:		[]string{message.Senders[0]},
			Timestamps:		[]time.Time{time.Now().UTC()},
			Subject:		message.Subject,
			MessageBody:	decodedQuery,
		}
		if err := send.SendSkymail(replyMessage, cfg.CertFilename); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println(err)
			fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
			return
		}
		replyMsgid := tools.GenerateID(replyMessage.Timestamps[0], tools.IsolateAddress(replyMessage.Senders[0]))
		filename := filepath.Join("..", "users", username, "STORE", replyMsgid+".msfn")
		if err := ioutil.WriteFile(filename, []byte(send.SkymailToString(replyMessage)), 0644); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println(err)
			fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
			return
		}
		if err := tools.Tag(replyMsgid, "Sent", username, false); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("Error tagging as sent:", err)
			fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
			return
		}
		fmt.Printf("30 /%s/mail?%s\r\n", skylabPath, replyMsgid)
	} else {
		fmt.Printf("10 Reply to %s:\r\n", tools.FormatSender(message.Senders[0]))
	}
}
