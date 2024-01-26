// fetch.go

package fetch

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/clseibold/go-gemini"
	"gitlab.com/clseibold/misfin-server/gembox"
)

// FetchGembox fetches gembox from the given address using the client certificate
func FetchGembox(address, certFilename string) ([]byte, error) {
	// Read certificate file
	certFile, err := ioutil.ReadFile(certFilename)
	if err != nil {
		return nil, err
	}

	// Fetch gembox from Gemini server
	client := gemini.DefaultClient
	resp, err := client.FetchWithCert(address, certFile, certFile)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read gembox data
	gemboxData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse gembox
	gbox := gembox.ParseGemBox(gemboxData)

	// Create INBOX directory if it doesn't exist
	inboxDir := "../INBOX/"
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
		filename := generateFilename(mail.Timestamps[0].Value, mail.Sender)

		// Save message to file
		err := ioutil.WriteFile(inboxDir+filename, []byte(mail.Body), 0644)
		if err != nil {
			fmt.Println("Error saving message:", err)
			continue
		}

		// Print confirmation
		fmt.Println("Message saved:", filename)
	}
}

// isNewMessage checks if the message is new (not already saved)
func isNewMessage(mail gembox.GemMail) bool {
	// Check if message file already exists
	filename := generateFilename(mail.Timestamps[0].Value, mail.Sender)
	_, err := os.Stat("../INBOX/" + filename)
	return os.IsNotExist(err)
}

// generateFilename generates a filename hash using the time and sender address
func generateFilename(timestamp time.Time, sender string) string {
	timeStr := timestamp.Format("2006-01-02T15:04:05")
	hash := md5.Sum([]byte(timeStr + sender))
	return fmt.Sprintf("%x.msg", hash)
}
