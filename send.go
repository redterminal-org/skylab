package main

import (
    "fmt"
	"strings"
	"os"
	"io/ioutil"
    "path/filepath"
	"time"
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
	message, err := tools.ParseSkymail(filename)
	if err != nil {
		fmt.Print("30 /\r\n")
		return
	}
	message.MessageBody = strings.ReplaceAll(message.MessageBody, "\n  \n", "\n\n")
	message.Timestamps[0] = time.Now().UTC()
	if err := send.SendSkymail(message, cfg.CertFilename); err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println(err)
		fmt.Println("=> inbox 📪 Inbox")
		return
	}
	if err := ioutil.WriteFile(filename, []byte(send.SkymailToString(message)), 0644); err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println(err)
		fmt.Println("=> inbox 📪 Inbox")
		return
	}
	if err := tools.Tag(msgid, "Sent", username); err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Error tagging as sent:", err)
		fmt.Println("=> inbox 📪 Inbox")
		return
	}
	if err := tools.Untag(msgid, "Drafts", username); err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Error untagging as draft:", err)
		fmt.Println("=> inbox 📪 Inbox")
		return
	}
	fmt.Printf("30 /mail?%s\r\n", msgid)
}
