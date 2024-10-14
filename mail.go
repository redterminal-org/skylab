package main

import (
    "fmt"
	"os"
	"strings"
	"path/filepath"
	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
)

func mail() {
	username, _ := tools.Authenticate()
	query := os.Getenv("QUERY_STRING")
	sitecfg, err := config.ExportSiteConfig()
	if err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Couldn't export site-wide config", err)
		return
	}
    cfg, err := config.ExportConfig(username)
	if err != nil {
		fmt.Println("Couldn't get the config info:", err)
		return
	}
	skylabPath := sitecfg.SkylabPath
	if query != "" {
		if strings.Contains(query, "..") {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("That query is not allowed")
			return
		}
		filename := filepath.Join("..", "users", username, "STORE", query+".msfn")
		message, err := tools.ParseSkymail(filename)
		if err != nil {
			fmt.Printf("30 /%s\r\n", skylabPath)
			return
		}
		sent, _ := tools.HasTag(query, "Sent", username)
		draft, _ := tools.HasTag(query, "Drafts", username)
		if err := tools.Untag(query, "Unread", username, true); err != nil {
			fmt.Print("Error: Could not mark as read")
		}
		fmt.Print("20 text/gemini\r\n")
		fmt.Printf("## %s\n", message.Subject)
		if !sent && !draft {
			var sendersBuilder strings.Builder
			for i, str := range message.Senders {
				formattedString := tools.FormatSender(str)
				sendersBuilder.WriteString(formattedString)
				if i != len(message.Senders)-1 {
					sendersBuilder.WriteString(", ")
				}
			}
			senders := sendersBuilder.String()
			fmt.Printf("From: %s\n\n", senders)
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
        var quickReplyLink string
        var quoteText string
        if cfg.QuoteReply {
            quoteText = message.MessageBody
        }
        if cfg.MisfinLink {
            quickReplyLink = tools.BuildMisfinURL(message.Senders, message.Subject, quoteText)
        } else {
            quickReplyLink = "reply/" + query + " ✍️"
        }
		if trashed, _ := tools.HasTag(query, "Trash", username); !trashed {
			if unarchived, _ := tools.HasTag(query, "Inbox", username); unarchived {
				fmt.Printf("=> %s Quick Reply\n", quickReplyLink)
				fmt.Printf("=> compose/new/%s 📝 Draft Reply\n", query)
				fmt.Printf("=> unread?%s ◉ Mark Unread\n", query)
				fmt.Printf("=> archive?%s 📥 Archive\n", query)
				fmt.Printf("=> trash?%s 🗑️ Delete\n", query)
				fmt.Printf("=> tag/%s 🏷️ Manage tags\n", query)
				fmt.Println("=> inbox 📪 Go to Inbox")
				fmt.Printf("=> /%s 🏠 Go Home\n", skylabPath)
				return
			}
			if archived, _ := tools.HasTag(query, "Archive", username); archived {
				fmt.Printf("=> %s Quick Reply\n", quickReplyLink)
				fmt.Printf("=> compose/new/%s 📝 Draft Reply\n", query)
				fmt.Printf("=> unread?%s ◉ Mark Unread\n", query)
				fmt.Printf("=> archive?%s 📤 Move to Inbox\n", query)
				fmt.Printf("=> trash?%s 🗑️ Delete\n", query)
				fmt.Printf("=> tag/%s 🏷️ Manage tags\n", query)
				fmt.Println("=> inbox 📪 Go to Inbox")
				fmt.Printf("=> /%s 🏠 Go Home\n", skylabPath)
				return
			}
			if sent {
                if cfg.MisfinLink {
                    fmt.Printf("=> %s Quick Follow Up\n")
                }
				fmt.Printf("=> compose/new/%s 📝 Follow Up\n", query)
				fmt.Printf("=> unread?%s ◉ Mark Unread\n", query)
				fmt.Printf("=> trash?%s 🗑️ Delete\n", query)
				fmt.Printf("=> tag/%s 🏷️ Manage tags\n", query)
				fmt.Println("=> inbox 📪 Go to Inbox")
				fmt.Printf("=> /%s 🏠 Go Home\n", skylabPath)
				return
			}
			if draft {
				fmt.Printf("=> compose/%s 📝 Edit Draft\n", query)
				fmt.Printf("=> trash?%s 🗑️ Delete\n", query)
				fmt.Printf("=> tag/%s 🏷️ Manage tags\n", query)
				fmt.Printf("=> send?%s ☞ Send\n", query)
				fmt.Printf("=> /%s 🏠 Go Home\n", skylabPath)
				return
			}
		} else {
			fmt.Printf("=> trash?%s 🛟 Take out of Trash\n", query)
			fmt.Printf("=> delete?%s ❌ Delete Permanently\n", query)
			fmt.Println("=> inbox 📪 Go to Inbox")
			fmt.Printf("=> /%s 🏠 Go Home\n", skylabPath)
			return
		}
	} else {
		fmt.Printf("30 /%s\r\n", skylabPath)
	}
}
