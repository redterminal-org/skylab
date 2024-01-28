package tools

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"
	"time"
    "bufio"
    "errors"
)

type Skymail struct {
	Senders     []string
	Recipients  []string
	Timestamps  []time.Time
	Subject     string
	MessageBody string
    FullText    string
}

func GenerateFilename(timestamp time.Time, sender string) string {
	hash := GenerateID(timestamp, sender)
	return fmt.Sprintf("%x.msfn", hash)
}

func GenerateID(timestamp time.Time, sender string) string {
	timeStr := timestamp.Format("2006-01-02T15:04:05")
	hash := md5.Sum([]byte(timeStr + sender))
    return fmt.Sprintf("%x", hash)
}

func ParseSkymail(skymailFilename string) ([]Skymail, error) {
    fileContent, err := os.ReadFile(skymailFilename)
    if err != nil {
        return nil, err
    }
    // Convert byte slice to string
    mailText := string(fileContent)

	scanner := bufio.NewScanner(strings.NewReader(mailText))
	var messages []Skymail
	var message Skymail
	var linesCounter int

	for scanner.Scan() {
		line := scanner.Text()
		message.FullText += line + "\n" // Append line to FullText
		if strings.HasPrefix(line, "<=====") {
			if linesCounter < 3 {
				return nil, errors.New("metadata lines are not properly formatted")
			}
			messages = append(messages, message)
			message = Skymail{} // Start a new message
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
                    ts = strings.TrimSpace(ts)
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

func IsolateAddress(input string) string {
    // Trim leading whitespace
    input = strings.TrimSpace(input)
    // Find first space
    spaceIndex := strings.Index(input, " ")
    if spaceIndex == -1 {
		return input
	}
    return input[:spaceIndex]
}
