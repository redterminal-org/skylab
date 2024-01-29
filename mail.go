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
		if err := tools.Untag(query, "Unread"); err != nil {
			fmt.Print("Error: Could not mark as read")
		}
		fmt.Print("20 text/gemini\r\n")
		fmt.Printf("## %s\n", message.Subject)
		fmt.Printf("From: %s\n\n", tools.FormatSender(message.Senders[0]))
		fmt.Println(message.MessageBody)
		if trashed, _ := tools.HasTag(query, "Trash"); !trashed {
			if unarchived, _ := tools.HasTag(query, "Inbox"); unarchived {
				fmt.Printf("=> reply?%s ✍️ Quick Reply\n", query)
				fmt.Printf("=> compose?%s 📝 Draft Reply\n", query)
				fmt.Printf("=> unread?%s ◉ Mark Unread\n", query)
				fmt.Printf("=> archive?%s 📥 Archive\n", query)
				fmt.Printf("=> trash?%s 🗑️ Delete\n", query)
				fmt.Printf("=> tag?%s 🏷️ Tag\n", query)
				fmt.Println("=> inbox 📪 Go to Inbox")
			} else if archived, _ := tools.HasTag(query, "Archive"); archived {
				fmt.Printf("=> reply?%s ✍️ Quick Reply\n", query)
				fmt.Printf("=> compose?%s 📝 Draft Reply\n", query)
				fmt.Printf("=> unread?%s ◉ Mark Unread\n", query)
				fmt.Printf("=> archive?%s 📤 Move to Inbox\n", query)
				fmt.Printf("=> trash?%s 🗑️ Delete\n", query)
				fmt.Printf("=> tag?%s 🏷️ Tag\n", query)
				fmt.Println("=> inbox 📪 Go to Inbox")
			}

		} else {
			fmt.Printf("=> restore?%s 🛟 Take out of Trash\n", query)
			fmt.Printf("=> delete?%s ❌ Delete Permanently\n", query)
			fmt.Println("=> inbox 📪 Go to Inbox")

		}
	} else {
		fmt.Print("30 /\r\n")
	}
}
