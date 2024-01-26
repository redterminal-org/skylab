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
	// Fetch gembox from Gemini server
	client := gemini.DefaultClient
	resp, err := client.FetchWithCert(address, certFilename, certFilename)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read gembox data
	gemboxData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return gemboxData, nil
}

// isNewMessage checks if the message is new (not already saved)
func IsNewMessage(mail gembox.GemMail) bool {
	// Check if message file already exists
	filename := generateFilename(mail.Timestamps[0].Value, mail.Sender)
	_, err := os.Stat("../INBOX/" + filename)
	return os.IsNotExist(err)
}
