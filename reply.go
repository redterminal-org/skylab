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
	query := os.Getenv("QUERY_STRING")
	path := os.Getenv("PATH_INFO")
	msgid := strings.TrimPrefix(path, "/reply/")
	filename := filepath.Join("..", "users", username, "STORE", msgid+".msfn")
	cfg, err := config.ExportConfig(username)
	if err != nil {
		fmt.Println("Couldn't get the config info")
		return
	}
	message, err := tools.ParseSkymail(filename)
	if err != nil {
		fmt.Print("30 /\r\n")
		return
	}
	if query != "" {
		address, err := send.GetAddressFromCert(cfg.CertFilename)
		if err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println(err)
			fmt.Println("=> inbox 📪 Inbox")
			return
		}
		decodedQuery, err := url.QueryUnescape(query)
		if err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println(err)
			fmt.Println("=> inbox 📪 Inbox")
			return
		}
		replyMessage := tools.Skymail{
			Senders:		[]string{address},
			Recipients:		[]string{message.Senders[0]},
			Timestamps:		[]time.Time{time.Now()},
			Subject:		message.Subject,
			MessageBody:	decodedQuery,
		}
		if err := send.SendSkymail(replyMessage, cfg.CertFilename); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println(err)
			fmt.Println("=> inbox 📪 Inbox")
			return
		}
		replyMsgid := tools.GenerateID(replyMessage.Timestamps[0], tools.IsolateAddress(replyMessage.Senders[0]))
		filename := filepath.Join("..", "users", username, "STORE", replyMsgid+".msfn")
		if err := ioutil.WriteFile(filename, []byte(send.SkymailToString(replyMessage)), 0644); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println(err)
			fmt.Println("=> inbox 📪 Inbox")
			return
		}
		if err := tools.Tag(replyMsgid, "Sent", username); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("Error tagging as sent:", err)
			fmt.Println("=> inbox 📪 Inbox")
			return
		}
		fmt.Printf("30 /mail?%s", replyMsgid)
	} else {
		fmt.Printf("10 Reply to %s:\r\n", tools.FormatSender(message.Senders[0]))
	}
}
