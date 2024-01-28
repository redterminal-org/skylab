package fetch

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
    "bufio"
    "errors"

	"github.com/clseibold/go-gemini"
    "git.sr.ht/~satchlj/skylab/tools"
)

type Message struct {
	Senders     []string
	Timestamps  []time.Time
    FullText    string
}

// FetchGembox fetches gembox from the given address using the client certificate
func FetchGembox(address, certFilename string) error {
	// Read certificate file
	certFile, err := ioutil.ReadFile(certFilename)
	if err != nil {
		return err
	}

	// Fetch gembox from Gemini server
	client := gemini.DefaultClient
	resp, err := client.FetchWithCert(address, certFile, certFile)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read gembox data
	gemboxData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse gembox
	gbox, err := ParseGembox(string(gemboxData))
    if err != nil {
        return err
    }

	// Create STORE directory if it doesn't exist
	inboxDir := "../STORE/"
	if _, err := os.Stat(inboxDir); os.IsNotExist(err) {
		err := os.Mkdir(inboxDir, 0755)
		if err != nil {
			return(err)
		}
	}

	// Process gemmails
	for _, mail := range gbox {
		// Check if message is new (not already saved)
		if !isNewMessage(mail) {
			continue
		}

		msgid := tools.GenerateID(mail.Timestamps[0], tools.IsolateAddress(mail.Senders[0]))
		filename := msgid + ".msfn"
		// Save message to file
		err := ioutil.WriteFile(inboxDir+filename, []byte(mail.FullText), 0644)
		if err != nil {
			fmt.Println("Error saving message:", err)
		}
		if err := tools.Tag(msgid, "Unread"); err != nil {
			fmt.Println("Error tagging new messages as unread:", err)
		}
		if err := tools.Tag(msgid, "Inbox"); err != nil {
			fmt.Println("Error adding new message to Inbox:", err)
		}
		// Print confirmation
		// fmt.Println("Message saved:", filename)
	}
    return nil
}

// isNewMessage checks if the message is new (not already saved)
func isNewMessage(mail Message) bool {
	// Check if message file already exists
	filename := tools.GenerateFilename(mail.Timestamps[0], tools.IsolateAddress(mail.Senders[0]))
	_, err := os.Stat("../STORE/" + filename)
	if err != nil {
		fmt.Println(err)
	}
	return os.IsNotExist(err)
}

func ParseGembox(gemboxData string) ([]Message, error) {

	scanner := bufio.NewScanner(strings.NewReader(gemboxData))
	var messages []Message
	var message Message
	var linesCounter int

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "<=====") {
			if linesCounter < 3 {
				return nil, errors.New("metadata lines are not properly formatted")
			}
			messages = append(messages, message)
			message = Message{} // Start a new message
			linesCounter = 0
		} else {
			message.FullText += line + "\n" // Append line to FullText
			switch linesCounter {
			case 0:
				message.Senders = strings.Split(line, ",")
			case 2:
				timestampStrings := strings.Split(line, ",")
                for _, ts := range timestampStrings {
                    ts = strings.TrimSpace(ts)
                    timestamp, err := time.Parse(time.RFC3339, ts)
                    if err != nil {
                        return nil, err
                    }
                    message.Timestamps = append(message.Timestamps, timestamp)
                }
			}
			linesCounter++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
