package main

import (
    "fmt"
	"os"
	"strings"
	"path/filepath"
//	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
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
	if query != "" {
		filename := filepath.Join("..", "users", username, "STORE", query+".msfn")
		message, err := tools.ParseSkymail(filename)
		if err != nil {
			fmt.Println(err)
			fmt.Println("30 /\r\n")
			return
		}
		sent, _ := tools.HasTag(query, "Sent", username)
		draft, _ := tools.HasTag(query, "Drafts", username)
		if err := tools.Untag(query, "Unread", username); err != nil {
			fmt.Print("Error: Could not mark as read")
		}
		fmt.Print("20 text/gemini\r\n")
		fmt.Printf("## %s\n", message.Subject)
		if !sent && !draft {
			fmt.Printf("From: %s\n\n", tools.FormatSender(message.Senders[0]))
		} else {
			var recipientsBuilder strings.Builder
			for i, str := range message.Recipients {
				formattedString := tools.FormatSender(str)
				recipientsBuilder.WriteString(formattedString)
				if i != len(message.Recipients)-1 {
					recipientsBuilder.WriteString(", ")
				}
			}
			recipients := recipientsBuilder.String()
			fmt.Printf("To: %s\n\n", recipients)
		}
		fmt.Println(message.MessageBody)
		if trashed, _ := tools.HasTag(query, "Trash", username); !trashed {
			if unarchived, _ := tools.HasTag(query, "Inbox", username); unarchived {
				fmt.Printf("=> reply/%s ✍️ Quick Reply\n", query)
				fmt.Printf("=> compose/new/%s 📝 Draft Reply\n", query)
				fmt.Printf("=> unread?%s ◉ Mark Unread\n", query)
				fmt.Printf("=> archive?%s 📥 Archive\n", query)
				fmt.Printf("=> trash?%s 🗑️ Delete\n", query)
				fmt.Printf("=> tag/%s 🏷️ Manage tags\n", query)
				fmt.Println("=> inbox 📪 Go to Inbox")
				return
			}
			if archived, _ := tools.HasTag(query, "Archive", username); archived {
				fmt.Printf("=> reply/%s ✍️ Quick Reply\n", query)
				fmt.Printf("=> compose/new/%s 📝 Draft Reply\n", query)
				fmt.Printf("=> unread?%s ◉ Mark Unread\n", query)
				fmt.Printf("=> archive?%s 📤 Move to Inbox\n", query)
				fmt.Printf("=> trash?%s 🗑️ Delete\n", query)
				fmt.Printf("=> tag/%s 🏷️ Manage tags\n", query)
				fmt.Println("=> inbox 📪 Go to Inbox")
				return
			}
			if sent {
				fmt.Printf("=> compose/new/%s 📝 Follow Up\n", query)
				fmt.Printf("=> unread?%s ◉ Mark Unread\n", query)
				fmt.Printf("=> trash?%s 🗑️ Delete\n", query)
				fmt.Printf("=> tag/%s 🏷️ Manage tags\n", query)
				return
			}
			if draft {
				fmt.Printf("=> compose/%s 📝 Edit Draft\n", query)
				fmt.Printf("=> trash?%s 🗑️ Delete\n", query)
				fmt.Printf("=> tag/%s 🏷️ Manage tags\n", query)
				fmt.Printf("=> send?%s ☞ Send\n", query)
				return
			}
		} else {
			fmt.Printf("=> trash?%s 🛟 Take out of Trash\n", query)
			fmt.Printf("=> delete?%s ❌ Delete Permanently\n", query)
			fmt.Println("=> inbox 📪 Go to Inbox")
			return
		}
	} else {
		fmt.Print("30 /\r\n")
	}
}
