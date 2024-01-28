package main

import (
    "fmt"
	"os"
//	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
)

func main() {
	query := os.Getenv("QUERY_STRING")
	if query != "" {
		filename := fmt.Sprintf("../STORE/%s.msfn", query)
		message, err := tools.ParseSkymail(filename)
		if err != nil {
			fmt.Println("30 /\r\n")
			return
		}
		fmt.Print("20 text/gemini\r\n")
		fmt.Printf("## %s\n", message.Subject)
		fmt.Printf("From: %s\n\n", tools.FormatSender(message.Senders[0]))
		fmt.Println(message.MessageBody)
		fmt.Println("=> reply?message-id ✍️ Quick Reply")
		fmt.Println("=> compose?message-id 📝 Draft Reply")
		fmt.Println("=> ?unread ◉ Mark Unread")
		fmt.Println("=> archive 📥 Archive")
		fmt.Println("=> trash 🗑️ Delete")
		fmt.Println("=> ?tag 🏷️ Tag")
		fmt.Println("=> inbox 📪 Inbox")
	} else {
		fmt.Print("30 /\r\n")
	}
}
