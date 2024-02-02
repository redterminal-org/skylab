package main

import (
    "fmt"
	"strings"
	"os"
	"io/ioutil"
	"git.sr.ht/~satchlj/skylab/tools"
	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/send"
)

func main() {
	msgid := os.Getenv("QUERY_STRING")
	filename := fmt.Sprintf("../STORE/%s.msfn", msgid)
	cfg := config.ExportConfig()
	message, err := tools.ParseSkymail(filename)
	if err != nil {
		fmt.Print("30 /\r\n")
		return
	}
	message.MessageBody = strings.ReplaceAll(message.MessageBody, "\n  \n", "\n\n")
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
	if err := tools.Tag(msgid, "Sent"); err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Error tagging as sent:", err)
		fmt.Println("=> inbox 📪 Inbox")
		return
	}
	if err := tools.Untag(msgid, "Drafts"); err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Error untagging as draft:", err)
		fmt.Println("=> inbox 📪 Inbox")
		return
	}
	fmt.Printf("30 /mail?%s\r\n", msgid)
}
