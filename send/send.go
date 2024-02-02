package send

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	"crypto/x509"
	"encoding/pem"
	"strconv"

    "git.sr.ht/~satchlj/skylab/tools"
	"gitlab.com/clseibold/misfin-server/misfin_client"
	"gitlab.com/clseibold/misfin-server/gemmail"
)

func SendSkymail(skymail tools.Skymail, certFilename string) error {
	recipients := skymail.Recipients
	certFile, err := ioutil.ReadFile(certFilename)
	if err != nil {
		return err
	}
	message := "# " + skymail.Subject + "\n" + skymail.MessageBody
	c := misfin_client.DefaultClient
	for _, recipient := range recipients {
		address := tools.IsolateAddress(recipient)
		var redirectAttempt = 0
		var misfinB = false
		for {
			mail := gemmail.CreateGemMailFromBody(message)
			msgToSend := ""
			if misfinB {
				msgToSend = mail.String_MisfinB()
			} else {
				msgToSend = mail.String()
			}
			resp, err := c.SendWithCert("misfin://"+address, certFile, certFile, msgToSend) // TODO: Make sure this strips all '\r'
			if err != nil {
				return fmt.Errorf("Sending gemmail to %s failed: %w\n", address, err)
			} else if (resp.Status == 59 || resp.Status == 40) && !misfinB { // Invalid request, try with Misfin(B) format
				misfinB = true
				c = misfin_client.DefaultClient_MisfinB
				redirectAttempt = 0
				continue
			} else if redirectAttempt < 5 && (resp.Status == misfin_client.StatusRedirect || resp.Status == misfin_client.StatusRedirectPermanent) {
				address = strings.TrimPrefix(strings.TrimSpace(resp.Meta), "misfin://")
				redirectAttempt += 1
				continue
			} else if resp.Status == 20 {
				// should check against stored fingerprint: address data
				return nil
			} else {
				return fmt.Errorf("Qwackers! It didn't work: %d %s\n", resp.Status, resp.Meta)
			}
			break
		}
	}
	return nil
}

func GetAddressFromCert(filename string) (string, error) {
	pemData, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "CERTIFICATE" {
		return "", fmt.Errorf("no certificate found in PEM data")
	}
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", err
	}
	domains := certificate.DNSNames
	var mailbox string
	for _, n := range certificate.Subject.Names {
		if n.Type.String() == "0.9.2342.19200300.100.1.1" {
			mailbox = n.Value.(string)
		}
	}
	mailAddress := mailbox + "@" + domains[0]
	return mailAddress, nil
}

func SkymailToString(mail tools.Skymail) string {
	// Combine Senders into a single string separated by commas
	senderStr := strings.Join(mail.Senders, ", ")

	// Combine Recipients into a single string separated by commas
	recipientStr := strings.Join(mail.Recipients, ", ")

	// Format Timestamps into ISO-8601 format and combine into a single string separated by commas
	var timestampStrs []string
	for _, ts := range mail.Timestamps {
		timestampStrs = append(timestampStrs, ts.Format(time.RFC3339))
	}
	timestampsStr := strings.Join(timestampStrs, ", ")

	// Combine Subject with "# " prepended on its own line
	subjectStr := "# " + mail.Subject

	// Combine all parts into a single string
	result := fmt.Sprintf("%s\n%s\n%s\n%s\n%s", senderStr, recipientStr, timestampsStr, subjectStr, mail.MessageBody)
	return result
}

func ParseComposePath(input string) (int, string) {
	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return 0, ""
	}
	integerPart, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, ""
	}
	alphanumericString := parts[1]
	return integerPart, alphanumericString
}
