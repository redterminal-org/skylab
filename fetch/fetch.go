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

	"github.com/clseibold/go-gemini"
    "gitlab.com/clseibold/misfin-server/gemmail"
)

// GemBox struct defined with necessary metadata pieces
type GemBox struct {
	Identifier any // For your use
	Mails      []gemmail.GemMail
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
	gbox := ParseGemBox(address, string(gemboxData))

	// Create INBOX directory if it doesn't exist
	inboxDir := "./INBOX/"
	if _, err := os.Stat(inboxDir); os.IsNotExist(err) {
		err := os.Mkdir(inboxDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	// Process gemmails
	for _, mail := range gbox.Mails {
		// Check if message is new (not already saved)
		if !isNewMessage(mail) {
			continue
		}

		// Generate filename hash
        // Timestamps[0] certainly doesn't work and is all made up. We need to check the ParseGemBox function to see how it works
		filename := generateFilename(mail.Timestamps[0].Value, mail.Senders[0].Address)
        var lines []string
    	for _, line := range mail.GemtextLines {
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
func isNewMessage(mail gemmail.GemMail) bool {
	// Check if message file already exists
	filename := generateFilename(mail.Timestamps[0].Value, mail.Senders[0].Address)
	_, err := os.Stat("../INBOX/" + filename)
	return os.IsNotExist(err)
}

// generateFilename generates a filename hash using the time and sender address
func generateFilename(timestamp time.Time, sender string) string {
	timeStr := timestamp.Format("2006-01-02T15:04:05")
	hash := md5.Sum([]byte(timeStr + sender))
	return fmt.Sprintf("%x.msg", hash)
}

func ParseGemBox(identifier any, gembox string) GemBox {
	var result GemBox
	result.Mails = make([]gemmail.GemMail, 0)
	result.Identifier = identifier

	scanner := bufio.NewScanner(strings.NewReader(gembox))
	// TODO: scanner.Err()

	var currentGemMail strings.Builder
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "<=====") {
			// Append current GemMail to list
			result.Mails = append(result.Mails, gemmail.ParseGemMail(currentGemMail.String()))

			// Start a new GemMail
			currentGemMail.Reset()
		} else {
			currentGemMail.WriteString(line)
			currentGemMail.WriteByte('\n')
		}
	}

	// Append last GemMail, if not empty
	if strings.TrimSpace(currentGemMail.String()) != "" {
		parsedGemMail := gemmail.ParseGemMail(currentGemMail.String())
		result.Mails = append(result.Mails, parsedGemMail)
	}

	return result
}
