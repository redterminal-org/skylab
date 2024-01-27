// fetch.go

package fetch

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
    "bufio"
    "errors"

	"github.com/clseibold/go-gemini"
)

type Message struct {
	Senders     []string
	Recipients  []string
	Timestamps  []time.Time
	Subject     string
	MessageBody string
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

	// Create INBOX directory if it doesn't exist
	inboxDir := "./INBOX/"
	if _, err := os.Stat(inboxDir); os.IsNotExist(err) {
		err := os.Mkdir(inboxDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	// Process gemmails
	for _, mail := range gbox {
		// Check if message is new (not already saved)
		if !isNewMessage(mail) {
			continue
		}

		// Generate filename hash
        // Timestamps[0] certainly doesn't work and is all made up. We need to check the ParseGemBox function to see how it works
		filename := generateFilename(mail.Timestamps[0], isolateAddress(mail.Senders[0]))
        var lines []string
    	for _, line := range mail.MessageBody {
    		lines = append(lines, line.String())
    	}
        content := string(mail.Subject) + "\n" + strings.Join(lines, "\n")

		// Save message to file
		err := ioutil.WriteFile(inboxDir+filename, []byte(content), 0644)
		if err != nil {
			fmt.Println("Error saving message:", err)
		}

		// Print confirmation
		fmt.Println("Message saved:", filename)
	}
    return nil
}

// isNewMessage checks if the message is new (not already saved)
func isNewMessage(mail Message) bool {
	// Check if message file already exists
	filename := generateFilename(mail.Timestamps[0].Value, isolateAddress(mail.Senders[0]))
	_, err := os.Stat("../INBOX/" + filename)
	return os.IsNotExist(err)
}

// generateFilename generates a filename hash using the time and sender address
func generateFilename(timestamp time.Time, sender string) string {
	timeStr := timestamp.Format("2006-01-02T15:04:05")
	hash := md5.Sum([]byte(timeStr + sender))
	return fmt.Sprintf("%x.msg", hash)
}

func ParseGembox(gemboxData string) ([]Message, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(strings.NewReader(gemboxData))
	var messages []Message
	var message Message
	var linesCounter int

	for scanner.Scan() {
		line := scanner.Text()
		message.FullText += line + "\n" // Append line to FullText
		if strings.HasPrefix(line, "<=====") {
			if linesCounter < 3 {
				return nil, errors.New("metadata lines are not properly formatted")
			}
			messages = append(messages, message)
			message = Message{} // Start a new message
			linesCounter = 0
		} else {
			switch linesCounter {
			case 0:
				message.Senders = strings.Split(line, ",")
			case 1:
				message.Recipients = strings.Split(line, ",")
			case 2:
				timestampStrings := strings.Split(line, ",")
                for _, ts := range timestampStrings {
                    timestamp, err := time.Parse(time.RFC3339, ts)
                    if err != nil {
                        return nil, err
                    }
                    message.Timestamps = append(message.Timestamps, timestamp)
                }
			default:
				if strings.HasPrefix(line, "# ") && message.Subject == "" {
					message.Subject = line
				} else {
					message.MessageBody += line + "\n"
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

func isolateAddress(input string) string {
    // Trim leading whitespace
    input = strings.TrimSpace(input)
    // Find first space
    spaceIndex := strings.Index(input, " ")
    if spaceIndex == -1 {
		return input
	}
    return input[:spaceIndex]
}
