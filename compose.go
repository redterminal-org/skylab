package main

import (
	"fmt"
	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/send"
	"git.sr.ht/~satchlj/skylab/tools"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func compose() {
	username, _ := tools.Authenticate()
	query := os.Getenv("QUERY_STRING")
	decodedQuery, _ := url.QueryUnescape(query)
	path := os.Getenv("PATH_INFO")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if strings.Contains(path, "..") {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("That path is not allowed")
		return
	}
	sitecfg, err := config.ExportSiteConfig()
	if err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Couldn't export site-wide config", err)
		return
	}
	skylabPath := appPrefix(sitecfg.SkylabPath)
	if skylabPath != "" {
		path = strings.TrimPrefix(path, skylabPath)
	}
	cfg, err := config.ExportConfig(username)
	if err != nil {
		fmt.Println("Couldn't get the config info")
		return
	}
	if strings.HasPrefix(path, "/compose/new") {
		ref := strings.TrimPrefix(path, "/compose/new/")
		address, err := send.GetAddressFromCert(cfg.CertFilename)
		message := tools.Skymail{
			Senders:    []string{address},
			Timestamps: []time.Time{time.Now().UTC()},
		}
		msgid := tools.GenerateID(message.Timestamps[0], tools.IsolateAddress(message.Senders[0]))
		refFilename := filepath.Join("..", "users", username, "STORE", ref+".msfn")
		refMessage, err := tools.ParseSkymail(refFilename)
		if err == nil {
			if sent, _ := tools.HasTag(ref, "Sent", username); sent {
				message.Recipients = refMessage.Recipients
			} else {
				message.Recipients = []string{refMessage.Senders[0]}
			}
			message.Subject = refMessage.Subject
		}
		filename := filepath.Join("..", "users", username, "STORE", msgid+".msfn")
		if err := os.WriteFile(filename, []byte(send.SkymailToString(message)), 0644); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println(err)
			fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
			return
		}
		if err := tools.Tag(msgid, "Drafts", username, false); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("Error tagging as draft:", err)
			fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
			return
		}
		fmt.Printf("30 %s/compose/%s\r\n", skylabPath, msgid)
		return
	} else if strings.HasPrefix(path, "/compose/edit-subject") {
		msgid := strings.TrimPrefix(path, "/compose/edit-subject/")
		filename := filepath.Join("..", "users", username, "STORE", msgid+".msfn")
		draft, err := tools.HasTag(msgid, "Drafts", username)
		if err != nil || draft != true {
			fmt.Printf("30 %s\r\n", appPath(skylabPath, ""))
			return
		}
		if query != "" {
			content, err := os.ReadFile(filename)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}
			lines := strings.Split(string(content), "\n")
			if len(lines) >= 4 {
				lines[3] = "# " + decodedQuery
			}
			newMessage := strings.Join(lines, "\n")
			err = os.WriteFile(filename, []byte(newMessage), 0644)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
			fmt.Printf("30 %s/compose/%s\r\n", skylabPath, msgid)
		} else {
			fmt.Print("10 Enter message subject:\r\n")
		}
		return
	} else if strings.HasPrefix(path, "/compose/edit-recipients") {
		msgid := strings.TrimPrefix(path, "/compose/edit-recipients/")
		filename := filepath.Join("..", "users", username, "STORE", msgid+".msfn")
		draft, err := tools.HasTag(msgid, "Drafts", username)
		if err != nil || draft != true {
			fmt.Printf("30 %s\r\n", appPath(skylabPath, ""))
			return
		}
		if query != "" {
			content, err := os.ReadFile(filename)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}
			lines := strings.Split(string(content), "\n")
			if len(lines) >= 4 {
				lines[1] = decodedQuery
			}
			newMessage := strings.Join(lines, "\n")
			err = os.WriteFile(filename, []byte(newMessage), 0644)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
			fmt.Printf("30 %s/compose/%s\r\n", skylabPath, msgid)
		} else {
			fmt.Print("10 Enter message recipients as address and name separated by a single space, with multiple recipients separated by commas:\r\n")
		}
		return
	} else if strings.HasPrefix(path, "/compose/add-text") {
		// Insert text at the appropriate place
		composePath := strings.TrimPrefix(path, "/compose/add-text/")
		n, msgid := send.ParseComposePath(composePath)
		filename := filepath.Join("..", "users", username, "STORE", msgid+".msfn")
		draft, err := tools.HasTag(msgid, "Drafts", username)
		if err != nil || draft != true {
			fmt.Printf("30 %s\r\n", appPath(skylabPath, ""))
			return
		}
		if query != "" {
			content, err := os.ReadFile(filename)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}
			lines := strings.Split(string(content), "\n")
			inserted := false
			if len(lines) == 5 && lines[4] == "" {
				lines[4] = fmt.Sprintf("  \n%s\n  \n", decodedQuery)
				inserted = true
			}
			counter := 0
			for i := 4; i < len(lines); i++ {
				line := lines[i]
				if strings.TrimSpace(line) == "" && len(line) == 2 {
					counter++
					if counter == n {
						lines[i] = fmt.Sprintf("  \n%s\n  ", decodedQuery)
						inserted = true
						break
					}
				}
			}
			if !inserted {
				lines = append(lines, fmt.Sprintf("  \n%s\n  \n", decodedQuery))
			}
			newMessage := strings.Join(lines, "\n")
			err = os.WriteFile(filename, []byte(newMessage), 0644)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
			fmt.Printf("30 %s/compose/%s\r\n", skylabPath, msgid)
		} else {
			fmt.Print("10 Enter new text:\r\n")
		}
		return
	} else if strings.HasPrefix(path, "/compose/edit-text") {
		composePath := strings.TrimPrefix(path, "/compose/edit-text/")
		n, msgid := send.ParseComposePath(composePath)
		filename := filepath.Join("..", "users", username, "STORE", msgid+".msfn")
		draft, err := tools.HasTag(msgid, "Drafts", username)
		if err != nil || draft != true {
			fmt.Printf("30 %s/\r\n", skylabPath)
			return
		}
		if query != "" {
			content, err := os.ReadFile(filename)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}
			blocks := strings.Split(string(content), "\n  \n")
			if decodedQuery != "X" {
				blocks[n] = decodedQuery
			} else {
				blocks = append(blocks[:n], blocks[n+1:]...)
				if len(blocks) == 2 {
					blocks = blocks[:1]
				}
			}
			newMessage := strings.Join(blocks, "\n  \n")
			err = os.WriteFile(filename, []byte(newMessage), 0644)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
			fmt.Printf("30 %s/compose/%s\r\n", skylabPath, msgid)
		} else {
			fmt.Print("10 Enter replacement text ('X' to delete block):\r\n")
		}
		return
	} else {
		msgid := strings.TrimPrefix(path, "/compose/")
		filename := filepath.Join("..", "users", username, "STORE", msgid+".msfn")
		draft, err := tools.HasTag(msgid, "Drafts", username)
		if err != nil || draft != true {
			fmt.Printf("30 %s/\r\n", skylabPath)
			return
		}
		message, _ := tools.ParseSkymail(filename)
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("# Edit Draft 📝")
		if message.Subject != "" {
			fmt.Printf("## %s\n", message.Subject)
		} else {
			fmt.Println("## No Subject")
		}
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

		fmt.Printf("=> edit-subject/%s  ✏️ Edit Subject\n", msgid)
		fmt.Printf("=> edit-recipients/%s  ✏️ Edit Recipients\n", msgid)
		if message.MessageBody == "" {
			fmt.Printf("=> add-text/1/%s  ✏️ Edit Message Text\n", msgid)
			fmt.Printf("=> %s/mail?%s 👁️ Preview and Send\n", skylabPath, msgid)
			return
		}
		fmt.Printf("=> add-text/1/%s 📃 Insert text\n", msgid)
		interblockText := "=> edit-text/%d/%s  ✏️ Edit text\n=> add-text/%d/%s 📃 Insert text\n"
		lines := strings.Split(message.MessageBody, "\n")
		counter := 0
		for i, line := range lines {
			if line == "  " {
				if counter >= 1 {
					lines[i] = fmt.Sprintf(interblockText, counter, msgid, counter+1, msgid)
				} else {
					lines[i] = ""
				}
				counter++
			}
		}
		result := strings.Join(lines, "\n")
		fmt.Println(result)
		fmt.Printf("=> %s/mail?%s 👁️ Preview and Send\n", skylabPath, msgid)
	}
}
