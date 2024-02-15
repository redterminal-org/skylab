package tools

import (
	"crypto/md5"
	"math/rand"
	"fmt"
	"os"
	"strings"
	"time"
    "bufio"
    "path/filepath"
	"io"
	"io/ioutil"
	"sort"

	"github.com/clseibold/go-gemini"
	"gopkg.in/yaml.v2"
	"git.sr.ht/~satchlj/skylab/config"
)

type Skymail struct {
	Senders     []string
	Recipients  []string
	Timestamps  []time.Time
	Subject     string
	MessageBody string
}

func GenerateFilename(timestamp time.Time, sender string) string {
	hash := GenerateID(timestamp, sender)
	return fmt.Sprintf("%s.msfn", hash)
}

func GenerateID(timestamp time.Time, sender string) string {
	timeStr := timestamp.Format("2006-01-02T15:04:05")
	hash := md5.Sum([]byte(timeStr + sender))
    return fmt.Sprintf("%x", hash)
}

func ParseSkymail(skymailFilename string) (Skymail, error) {
    fileContent, err := os.ReadFile(skymailFilename)
    if err != nil {
        return Skymail{}, err
    }
    // Convert byte slice to string
    mailText := string(fileContent)

	scanner := bufio.NewScanner(strings.NewReader(mailText))
	var message Skymail
	var linesCounter int

    for scanner.Scan() {
        line := scanner.Text()
        switch linesCounter {
        case 0:
			senders := strings.Split(strings.TrimPrefix(line, "< "), ",")
            for _, sender := range senders {
                sender = strings.TrimSpace(sender)
                message.Senders = append(message.Senders, sender)
            }
        case 1:
			recipients := strings.Split(strings.TrimPrefix(line, ": "), ",")
            for _, recipient := range recipients {
                recipient = strings.TrimSpace(recipient)
                message.Recipients = append(message.Recipients, recipient)
            }
        case 2:
            timestampStrings := strings.Split(strings.TrimPrefix(line, "@ "), ",")
            for _, ts := range timestampStrings {
                ts = strings.TrimSpace(ts)
                timestamp, err := time.Parse(time.RFC3339, ts)
                if err != nil {
                    return Skymail{}, err
                }
                message.Timestamps = append(message.Timestamps, timestamp)
            }
        default:
            if strings.HasPrefix(line, "# ") && message.Subject == "" {
                message.Subject = strings.TrimPrefix(line, "# ")
            } else {
                message.MessageBody += line + "\n"
            }
        }
        linesCounter++
    }

	if err := scanner.Err(); err != nil {
		return Skymail{}, err
	}
	if message.Subject == "" {
		message.Subject = "No Subject"
	}

	return message, nil
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

func PrintLinkAndPreview(msgid, username string, space bool) (string, error) {
    // Get file contents as message string
    filename := filepath.Join("..", "users", username, "STORE", msgid+".msfn")
    message, err := ParseSkymail(filename)
    if err != nil {
        return "", err
    }
    linkline, err := GetMessageLinkLine(message, msgid, username)
    if err != nil {
        return "", err
    }
    preview, err := GetMailPreview(message)
    if err != nil {
        return "", err
    }
    // Construct the formatted text
	var formattedText string
	if space {
		formattedText = fmt.Sprintf("%s\n%s\n", linkline, preview)
	} else {
		formattedText = fmt.Sprintf("%s\n%s", linkline, preview)
	}
	return formattedText, nil
}

// Create message link line
func GetMessageLinkLine(message Skymail, msgid, username string) (string, error) {
	icon := "◎"
	subject := "No Subject"
	// Check to see if unread
	unread, err := HasTag(msgid, "Unread", username)
	if err != nil {
		return "", err
	}
	if unread {
		icon = "◉"
	}
	if message.Subject != "" {
		subject = message.Subject
	}
	maxLength := 50
    if len(subject) > maxLength {
        // If the input is larger than maxLength, truncate it and add "..."
        subject = subject[:maxLength-4] + "... "
    }
	cfg, err := config.ExportConfig(username)
	if err != nil {
		return "", fmt.Errorf("Error getting config for timezone:", err)
	}
	location, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		return "", fmt.Errorf("Error loading timezone:", err)
	}
	ts := message.Timestamps[0].In(location)
	// Create string
	linkString := fmt.Sprintf("=> /mail?%s %s %s — %s: %s", msgid, icon, ts.Format("2006-01-02 03:04 PM"), IsolateAddress(message.Senders[0]), subject)
	return linkString, nil
}

func GetMailPreview(message Skymail) (string, error) {
    // Count out 70 characters and replace newlines with spaces
    var preview strings.Builder
    count := 0
    afterNewline := false

    for _, char := range message.MessageBody {
        if count >= 91 {
            break
        }
        if char == '\n' || char == '\r' {
            if !afterNewline && preview.Len() > 0 && preview.String()[preview.Len()-1] != ' ' {
			    preview.WriteByte(' ')
		    	count++
			}
			afterNewline = true
		} else {
			preview.WriteRune(char)
			count++
			afterNewline = false
		}
	}

	previewString := preview.String()
	if len(previewString) >= 3 && previewString[:3] == "=> " {
		previewString = previewString[3:]
	}
	if len(previewString) >= 2 && previewString[:2] == "# " {
		previewString = previewString[2:]
	}
	if len(previewString) >= 3 && previewString[:3] == "## " {
		previewString = previewString[3:]
	}
	if len(previewString) >= 4 && previewString[:4] == "### " {
		previewString = previewString[4:]
	}
	if len(previewString) >= 3 && previewString[:3] == "```" {
		previewString = "[preformatted block]"
	}
	maxLength := 90
    if len(previewString) > maxLength {
        // If the input is larger than maxLength, truncate it and add "..."
		if previewString[maxLength-3] == ' ' && previewString[maxLength-4] != ' ' {
			previewString = previewString[:maxLength-4] + "..."
		} else {
			previewString = previewString[:maxLength-3] + "..."
		}
	}
	return previewString, nil
}

func Tag(msgid, tagName, username string, sync bool) error {
    // Create the path to the symlink destination
	destDir := filepath.Join("..", "users", username, "tags", tagName)
    destPath := filepath.Join(destDir, msgid+".msfn")

	// Check if the symlink already exists
    if _, err := os.Lstat(destPath); err == nil {
        // Symlink already exists, do nothing
        return nil
    } else if !os.IsNotExist(err) {
        // Other error occurred, return the error
        return err
    }
	cfg, err := config.ExportConfig(username)
	if err != nil {
		return fmt.Errorf("Error getting config for Tag:", err)
	}
	if cfg.TagSync && sync {
		client := gemini.DefaultClient
		certFile, err := ioutil.ReadFile(cfg.CertFilename)
		if err != nil {
			return fmt.Errorf("Couldn't read cert: %w (%s)", err, cfg.CertFilename)
		}
		if cfg.GMAPHost == "" {
			return fmt.Errorf("Your GMAP host configuration is wrong. Go fix it:\n=> settings ⚙️ Settings")
		}
		var port string
		if cfg.GMAPPort == "" {
			port = "1958"
		} else {
			port = cfg.GMAPPort
		}
		tagPath := "gemini://" + cfg.GMAPHost + ":" + port + "/tag/" + tagName + "?" + msgid
		// Fetch gembox from Gemini server
		resp, err := client.FetchWithCert(tagPath, certFile, certFile)
		if err != nil || resp.Status != 20 {
			return fmt.Errorf("Couldn't tag the message: %w (%s)", err, tagPath)
		}
		defer resp.Body.Close()
	}

    // Create the parent directory for the symlink
    if err := os.MkdirAll(destDir, 0755); err != nil {
        return err
    }

    // Create the symlink
    if err := os.Symlink(filepath.Join("..", "users", username, "STORE", msgid+".msfn"), destPath); err != nil {
        return err
    }
    return nil
}

// RemoveSymlink removes the symlink from ../tags/tagName/msgid.msfn
func Untag(msgid, tagName, username string, sync bool) error {
    // Create the path to the symlink
    symlinkPath := filepath.Join("..", "users", username, "tags", tagName, msgid+".msfn")

    // Check if the symlink exists
    _, err := os.Lstat(symlinkPath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil
        }
        return err
    }

	cfg, err := config.ExportConfig(username)
	if err != nil {
		return fmt.Errorf("Error getting config for Tag:", err)
	}
	if cfg.TagSync && sync {
		client := gemini.DefaultClient
		certFile, err := ioutil.ReadFile(cfg.CertFilename)
		if err != nil {
			return fmt.Errorf("Couldn't read cert: %w (%s)", err, cfg.CertFilename)
		}
		if cfg.GMAPHost == "" {
			return fmt.Errorf("Your GMAP host configuration is wrong. Go fix it:\n=> settings ⚙️ Settings")
		}
		var port string
		if cfg.GMAPPort == "" {
			port = "1958"
		} else {
			port = cfg.GMAPPort
		}
		tagPath := "gemini://" + cfg.GMAPHost + ":" + port + "/untag/" + tagName + "?" + msgid
		time.Sleep(50 * time.Millisecond)
		// Fetch gembox from Gemini server
		resp, err := client.FetchWithCert(tagPath, certFile, certFile)
		if err != nil || resp.Status != 20 {
			return fmt.Errorf("Couldn't untag the message: %w (%s)", err, tagPath)
		}
		defer resp.Body.Close()
	}
    // Remove the symlink
    if err := os.Remove(symlinkPath); err != nil {
        return err
    }

	    // Check if the tagName directory is empty
    dir := filepath.Join("..", "users", username, "tags", tagName)
    isEmpty, err := isDirEmpty(dir)
    if err != nil {
        return err
    }

    // If the directory is empty, remove it
    if isEmpty {
        if err := os.Remove(dir); err != nil {
            return err
        }
    }

    return nil
}

func isDirEmpty(name string) (bool, error) {
    f, err := os.Open(name)
    if err != nil {
        return false, err
    }
    defer f.Close()

    _, err = f.Readdirnames(1) // Try to read one file
    if err == nil {
        // If no error, directory is not empty
        return false, nil
    }
    if err == io.EOF {
        // If EOF error, directory is empty
        return true, nil
    }
    // Otherwise, an error occurred
    return false, err
}

func HasTag(msgid, tagName, username string) (bool, error) {
    // Construct the path to the symlink
    symlinkPath := filepath.Join("..", "users", username, "tags", tagName, msgid+".msfn")

    // Check if the symlink exists
    if _, err := os.Lstat(symlinkPath); err != nil {
        if os.IsNotExist(err) {
            // Symlink does not exist
            return false, nil
        }
        // Other error occurred
        return false, err
    }

    // Check if the symlink points to a valid target
    target, err := os.Readlink(symlinkPath)
    if err != nil {
        return false, err
    }

    // If the symlink points to a target, it exists
    return target != "", nil
}

func WalkTag(tagName, username string, space bool, n ...int) error {
	dirName := filepath.Join("..", "users", username, "tags", tagName)
	var msgs []string
	err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		lstat, err := os.Lstat(path)
		if err != nil {
			return fmt.Errorf("os.Lstat error: %w", err)
		} else if lstat.Mode()&os.ModeSymlink != 0 {
			symlinkName := filepath.Base(path)
			msgid := strings.TrimSuffix(symlinkName, ".msfn")
			hide, err := HasTag(msgid, "Trash", username)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			if tagName == "Trash" {
				hide = false
			}
			if !hide {
				output, err := PrintLinkAndPreview(msgid, username, space)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				msgs = append(msgs, output)
			}
		}
		return nil
	})
	if err != nil {
        return err
    }
	sortByTime := func(i, j int) bool {
        timeFormat := "2006-01-02 03:04 PM"
        strI := msgs[i][46:65] // Extract timestamp from string i
        strJ := msgs[j][46:65] // Extract timestamp from string j
        timeI, err := time.Parse(timeFormat, strI)
		if err != nil {
			fmt.Println(err)
		}
        timeJ, err := time.Parse(timeFormat, strJ)
		if err != nil {
			fmt.Println(err)
		}
        return timeI.After(timeJ)
    }
	// Sort the outputs slice
    sort.Slice(msgs, sortByTime)
	// Limit the number of printed things
	var limit int
	if len(n) > 0 && n[0] < len(msgs) {
		limit = n[0]
	} else {
		limit = len(msgs)
	}

	// Print the sorted outputs
	for _, output := range msgs[:limit] {
        fmt.Println(output)
    }

    return nil
}

func FormatSender(sender string) string {
    parts := strings.SplitN(sender, " ", 2)
    if len(parts) == 2 {
        return parts[1] + " (" + parts[0] + ")"
    }
    // If the input string doesn't contain both email and name, return it as is
    return sender
}

// CountTags counts the number of symlinks in the directory for a given tag.
// It returns the count as a string in the format "(n)" or an empty string if the directory is empty.
func CountTags(tagName, username string) (string, error) {
    dirName := filepath.Join("..", "users", username, "tags", tagName)
    var count int

    err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }
		symlinkName := filepath.Base(path)
        msgid := strings.TrimSuffix(symlinkName, ".msfn")
		hide, err := HasTag(msgid, "Trash", username)
		if err != nil {
            return fmt.Errorf("HasTag(Trash) error: %w", err)
        }
		if tagName == "Trash" {
			hide = false
		}
        if info.Mode()&os.ModeSymlink != 0 && !hide {
            count++
        }
        return nil
    })
    if err != nil {
        return "", err
    }
    if count == 0 {
        return "", nil
    }
    return fmt.Sprintf(" (%d)", count), nil
}
func UntagAll(msgID, username string) error {
    tagsDir := filepath.Join("..", "users", username, "tags")
    err := filepath.Walk(tagsDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }
        if info != nil && info.IsDir() && path != tagsDir {
            tag := filepath.Base(path)
            tagDir := filepath.Join(tagsDir, tag)
            if _, err := os.Stat(tagDir); os.IsNotExist(err) {
                fmt.Printf("Skipping non-existent directory: %s\n", tagDir) // Debug print
                return nil
            }
            err := filepath.Walk(tagDir, func(tagPath string, tagInfo os.FileInfo, err error) error {
                if err != nil {
                    return fmt.Errorf("filepath.Walk error 2: %w", err)
                }
                if tagInfo.Mode()&os.ModeSymlink != 0 {
                    symlinkName := filepath.Base(tagPath)
                    if strings.HasPrefix(symlinkName, msgID) {
                        if err := os.Remove(tagPath); err != nil {
                            return fmt.Errorf("failed to remove symlink %s: %w", tagPath, err)
                        }
                    }
                }
                return nil
            })
            if err != nil {
                return err
            }
        }
        return nil
    })
    if err != nil {
        return err
    }
    return nil
}

func ListTags(username string) ([]string, error) {
    tagsDir := filepath.Join("..", "users", username, "tags")
	var tagNames []string
	err := filepath.Walk(tagsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != tagsDir {
			tagName := filepath.Base(path)
			reservedTags := []string{"Inbox", "Archive", "Trash", "Sent", "Unread", "Drafts"}
			var reserved bool
			for _, tag := range reservedTags {
				if tagName == tag {
					reserved = true
					break
				}
			}
			if !reserved {
				tagNames = append(tagNames, tagName)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tagNames, nil
}

type UserHashes map[string]string

func Authenticate() (string, error) {
	filename := "../users.yaml"
	hash := os.Getenv("TLS_CLIENT_HASH")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Could not read users file.")
		return "", err
	}

	var users UserHashes
	if err := yaml.Unmarshal(yamlFile, &users); err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	username, exists := users[hash]
	if !exists {
		return "", nil
	}
	return username, nil
}

func AddUserHash(username string) error {
	filename := "../users.yaml"
	hash := os.Getenv("TLS_CLIENT_HASH")
	yamlFile, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        return err
    }
	defer yamlFile.Close()

    // Prepare the new line to append to the file
    newLine := fmt.Sprintf("%s: %s\n", hash, username)

    // Append the new line to the end of the file
    _, err = yamlFile.WriteString(newLine)
    if err != nil {
        return err
    }

    return nil
}

func CheckTempPassword(password string) error {
	filename := "../passwords.yaml"
	passHash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Could not read users file.")
		return err
	}
	var temppasses UserHashes
	if err := yaml.Unmarshal(yamlFile, &temppasses); err != nil {
		fmt.Println("Error:", err)
		return err
	}
	username, exists := temppasses[passHash]
	if !exists {
		return fmt.Errorf("invalid temporary password")
	}
	if err := AddUserHash(username); err != nil {
		return err
	}
	delete(temppasses, passHash)
	updatedYAML, err := yaml.Marshal(&temppasses)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	if err := ioutil.WriteFile(filename, updatedYAML, 0644); err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	return nil
}

// GenerateOneTimePassword generates a one-time use password, writes its hash to a YAML file along with the username,
// and returns the generated password.
func GenerateOneTimePassword(username string) (string, error) {
	// Generate a random password
	password := generateRandomPassword()

	// Hash the password
	passHash := fmt.Sprintf("%x", md5.Sum([]byte(password)))

	// Read existing passwords from YAML file
	filename := "../passwords.yaml"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
        // If the file doesn't exist, create it with empty content
        if err := ioutil.WriteFile(filename, []byte{}, 0644); err != nil {
            return "", fmt.Errorf("error creating passwords file: %w", err)
        }
    }
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read passwords file: %w", err)
	}
	// Unmarshal existing passwords
	userHashes := make(map[string]string)
	if err := yaml.Unmarshal(yamlFile, &userHashes); err != nil {
		return "", fmt.Errorf("failed to unmarshal passwords data: %w", err)
	}

	// Add new password hash to userHashes
	userHashes[passHash] = username

	// Marshal userHashes to YAML
	updatedYAML, err := yaml.Marshal(&userHashes)
	if err != nil {
		return "", fmt.Errorf("failed to marshal passwords data: %w", err)
	}

	// Write updated YAML data back to the file
	if err := ioutil.WriteFile(filename, updatedYAML, 0644); err != nil {
		return "", fmt.Errorf("failed to write passwords file: %w", err)
	}

	return password, nil
}

// generateRandomPassword generates a random password.
func generateRandomPassword() string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func contains(slice []string, msgid string) bool {
    for _, item := range slice {
        if strings.TrimSpace(item) == strings.TrimSpace(msgid) {
            return true
        }
    }
    return false
}
