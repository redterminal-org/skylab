package fetch
import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
    "bufio"
    "errors"
	"path/filepath"

	"github.com/clseibold/go-gemini"
    "git.sr.ht/~satchlj/skylab/tools"
	"git.sr.ht/~satchlj/skylab/config"
)

type Message struct {
	Senders     []string
	Timestamps  []time.Time
    FullText    string
	MessageBody	string
	Msgid		string
	Tags		[]string
	Subject		string
	Recipients	[]string
}

// FetchGembox fetches gembox from the given address using the client certificate
func FetchGembox(username string) error {
	// Read certificate file
    cfg, err := config.ExportConfig(username)
	if err != nil {
		return fmt.Errorf("Couldn't get the config info: %w", err)
	}
	certFile, err := ioutil.ReadFile(cfg.CertFilename)
	if err != nil {
		return fmt.Errorf("Couldn't read cert: %w (%s)", err, cfg.CertFilename)
	}

	// Fetch gembox from Gemini server
	client := gemini.DefaultClient
	resp, err := client.FetchWithCert(cfg.GeminiFetchAddress, certFile, certFile)
	if err != nil {
		return fmt.Errorf("Couldn't fetch gembox file: %w", err)
	}
	defer resp.Body.Close()

	// Read gembox data
	gemboxData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Errno 123: %w")
	}

	// Parse gembox
	gbox, err := ParseGembox(string(gemboxData))
    if err != nil {
        return err
    }
	if string(gemboxData) == "" {
		return errors.New("Empty gembox.")
	}

	// Create STORE directory if it doesn't exist
	inboxDir := filepath.Join("..", "users", username, "STORE") + string(filepath.Separator)
	if _, err := os.Stat(inboxDir); os.IsNotExist(err) {
		err := os.Mkdir(inboxDir, 0755)
		if err != nil {
			return fmt.Errorf("Couldn't make inboxDir: %w", err)
		}
	}
	// Process gemmails
	lastFetched, err := time.Parse("2006-01-02 15:04:05", cfg.LastFetched)
	for _, mail := range gbox {
		// Check if message is new (not already saved)
		if !isNewMessage(mail, lastFetched, username) {
			continue
		}
		msgid := tools.GenerateID(mail.Timestamps[0], tools.IsolateAddress(mail.Senders[0]))
		filename := msgid + ".msfn"
		// Save message to file
		err := ioutil.WriteFile(inboxDir+filename, []byte(mail.FullText), 0644)
		if err != nil {
			fmt.Println("Error saving message:", err)
		}
		if err := tools.Tag(msgid, "Unread", username, false); err != nil {
			fmt.Println("Error tagging new messages as unread:", err)
		}
		if err := tools.Tag(msgid, "Inbox", username, false); err != nil {
			fmt.Println("Error adding new message to Inbox:", err)
		}
		// Print confirmation
		// fmt.Println("Message saved:", filename)
	}
	if err := config.Set(username, "LastFetched", time.Now().UTC().Format("2006-01-02 15:04:05")); err != nil {
		return fmt.Errorf("Could not write to config file: %w", err)
	}
    return nil
}

// isNewMessage checks if the message is new (not already saved)
func isNewMessage(mail Message, lastFetched time.Time, username string) bool {
	// Check if message file already exists
	if mail.Timestamps[0].After(lastFetched) {
		return true
	} else if mail.Timestamps[0].Before(lastFetched) {
		return false
	} else {
		filename := tools.GenerateFilename(mail.Timestamps[0], tools.IsolateAddress(mail.Senders[0]))
		_, err := os.Stat(filepath.Join("..", "users", username, "STORE", filename))
		return os.IsNotExist(err)
	}
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
				message.Senders = strings.Split(strings.TrimPrefix(line, "< "), ",")
			case 2:
				timestampStrings := strings.Split(strings.TrimPrefix(line, "@ "), ",")
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
	if linesCounter >= 3 {
		messages = append(messages, message)
	} else if linesCounter > 0 {
		return nil, errors.New("last message is not complete")
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// FetchGembox fetches gembox from the given address using the client certificate
func FetchGemboxGMAP(username string) error {
	// Read config and cert file
    cfg, err := config.ExportConfig(username)
	if err != nil {
		return fmt.Errorf("Couldn't get the config info: %w", err)
	}
	certFile, err := ioutil.ReadFile(cfg.CertFilename)
	if err != nil {
		return fmt.Errorf("Couldn't read cert: %w (%s)", err, cfg.CertFilename)
	}

	// Fetch gembox from Gemini server
	client := gemini.DefaultClient
	resp, err := client.FetchWithCert(cfg.GeminiFetchAddress, certFile, certFile)
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
	gbox, err := ParseGemboxGMAP(string(gemboxData))
    if err != nil {
        return err
    }

	// Create STORE directory if it doesn't exist
	inboxDir := filepath.Join("..", "users", username, "STORE") + string(filepath.Separator)
	if _, err := os.Stat(inboxDir); os.IsNotExist(err) {
		err := os.Mkdir(inboxDir, 0755)
		if err != nil {
			return fmt.Errorf("Couldn't make STORE: %w", err)
		}
	}
	// Process gemmails
	lastFetched, err := time.Parse("2006-01-02 15:04:05", cfg.LastFetched)
	for i, mail := range gbox {
		if len(mail.Timestamps) == 0 {
			return fmt.Errorf("There are no timestamps for message %s, %d of %d", mail.Msgid, i, len(gbox))
		}
		// Check if message is new (not already saved)
		if !isNewMessageGMAP(mail, lastFetched, username) {
//			tools.UntagAll(mail.Msgid, username)
//			for _, tag := range mail.Tags {
//				if err := tools.Tag(mail.Msgid, tag, username); err != nil {
//					fmt.Println("Error tagging new message", err)
//				}
//			}
			continue
		}
		var builder strings.Builder
		builder.WriteString(strings.Join(mail.Senders, ", ") + "\n")
		builder.WriteString(strings.Join(mail.Recipients, ", ") + "\n")
		for i, timestamp := range mail.Timestamps {
			builder.WriteString(timestamp.Format(time.RFC3339))
			if i != len(mail.Timestamps)-1 {
				builder.WriteString(", ")
			} else {
				builder.WriteString("\n")
			}
		}
		if mail.Subject != "" {
			builder.WriteString("# " + mail.Subject + "\n")
		}
		builder.WriteString(strings.TrimSpace(mail.MessageBody))
		fullText := builder.String()
		filename := mail.Msgid + ".msfn"
		// Save message to file
		err := ioutil.WriteFile(inboxDir+filename, []byte(fullText), 0644)
		if err != nil {
			fmt.Println("Error saving message:", err)
		}
		for _, tag := range mail.Tags {
			if err := tools.Tag(mail.Msgid, tag, username, false); err != nil {
				fmt.Println("Error tagging new message", err)
			}
		}
		// Print confirmation
		// fmt.Println("Message saved:", filename)
	}
	if err := config.Set(username, "LastFetched", time.Now().UTC().Format("2006-01-02 15:04:05")); err != nil {
		return fmt.Errorf("Could not write to config file: %w", err)
	}
    return nil
}

// isNewMessage checks if the message is new (not already saved)
func isNewMessageGMAP(mail Message, lastFetched time.Time, username string) bool {
	// Check if message file already exists
	if mail.Timestamps[0].After(lastFetched) {
		return true
	} else if mail.Timestamps[0].Before(lastFetched) {
		return false
	} else {
		filename := mail.Msgid + ".msfn"
		_, err := os.Stat(filepath.Join("..", "users", username, "STORE", filename))
		return os.IsNotExist(err)
	}
}

func ParseGemboxGMAP(gemboxData string) ([]Message, error) {

	scanner := bufio.NewScanner(strings.NewReader(gemboxData))
	var messages []Message
	var message Message
	var reading bool
	var body bool

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[") {
			if reading && !body {
				return nil, errors.New("message not properly formatted")
			}
			if reading {
				messages = append(messages, message)
			}
			reading = true
			message = Message{} // Start a new message
			body = false
			line = strings.TrimPrefix(line, "[")
    		line = strings.TrimSuffix(line, "]")
			message.Msgid = line
		} else {
			if !reading {
				continue
			}
			if !body {
				if strings.HasPrefix(line, "#") {
					body = true;
					if line == "#" || line == "# " {
						continue
					}
					message.Subject = strings.TrimPrefix(line, "# ")
				} else if strings.HasPrefix(line, "Senders:") {
					message.Senders = strings.Split(strings.TrimPrefix(line, "Senders:"), ",")
				} else if strings.HasPrefix(line, "Recipients:") {
					message.Recipients = strings.Split(strings.TrimPrefix(line, "Recipients:"), ",")
				} else if strings.HasPrefix(line, "Timestamps:") {
					timestampStrings := strings.Split(strings.TrimPrefix(line, "Timestamps:"), ",")
					for _, ts := range timestampStrings {
						ts = strings.TrimSpace(ts)
						timestamp, err := time.Parse(time.RFC3339, ts)
						if err != nil {
							return nil, err
						}
						message.Timestamps = append(message.Timestamps, timestamp)
					}
					if len(message.Timestamps) == 0 {
						return nil, errors.New("Something's a bit wrong with parsing the timestamps")
					}
				} else {
					if !strings.Contains(line, ":") {
						return nil, errors.New("Tagline improperly formatted")
					}
					parts := strings.SplitN(line, ":", 2)
					tagName := parts[0]
					rest := strings.TrimSpace(parts[1])
					bits := strings.Split(rest, " ")
//					timestamp, err := time.Parse(time.RFC3339, bits[0])
//					if err != nil {
//						return nil, err
//					}
					if bits[1] == "false" {
						continue
					} else if bits[1] != "true" {
						return nil, errors.New("gembox tag boolean not properly formatted")
					}
					message.Tags = append(message.Tags, tagName)
				}
			} else {
				if strings.HasPrefix(line, "\\[") {
					line = strings.TrimPrefix(line, "\\")
				}
				message.MessageBody += line + "\n"
			}
		}
	}
	if body {
		messages = append(messages, message)
		if len(message.Timestamps) == 0 {
			return nil, errors.New("Something's a bit wrong with the last message")
		}
	} else if reading {
		return nil, errors.New("last message is not complete")
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}
