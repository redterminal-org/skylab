package tools

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"
	"time"
    "bufio"
    "path/filepath"
	"io"
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
	return fmt.Sprintf("%x.msfn", hash)
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
            message.Senders = strings.Split(line, ",")
        case 1:
            message.Recipients = strings.Split(line, ",")
        case 2:
            timestampStrings := strings.Split(line, ",")
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

func PrintLinkAndPreview(msgid string) (error) {
    // Get file contents as message string
    filename := filepath.Join("..", "STORE", msgid+".msfn")
    message, err := ParseSkymail(filename)
    if err != nil {
        return err
    }
    linkline, err := GetMessageLinkLine(message, msgid)
    if err != nil {
        return err
    }
    preview, err := GetMailPreview(message)
    if err != nil {
        return err
    }
	fmt.Printf("%s\n%s\n", linkline, preview)
	return nil
}

// Create message link line
func GetMessageLinkLine(message Skymail, msgid string) (string, error) {
	icon := "◎"
	// Check to see if unread
	unread, err := HasTag(msgid, "Unread")
	if err != nil {
		return "", err
	}
	if unread {
		icon = "◉"
	}
	// Create string
	linkString := fmt.Sprintf("=> /mail?%s %s %s - %s", msgid, icon, message.Subject, IsolateAddress(message.Senders[0]))
	return linkString, nil
}

func GetMailPreview(message Skymail) (string, error) {
    // Count out 70 characters and replace newlines with spaces
    var preview strings.Builder
    count := 0
    afterNewline := false

    for _, char := range message.MessageBody {
        if count >= 80 {
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

	return preview.String(), nil
}

func Tag(msgid, tagName string) error {
    // Create the path to the symlink destination
    destDir := filepath.Join("..", "tags", tagName)
    destPath := filepath.Join(destDir, msgid+".msfn")

	// Check if the symlink already exists
    if _, err := os.Lstat(destPath); err == nil {
        // Symlink already exists, do nothing
        return nil
    } else if !os.IsNotExist(err) {
        // Other error occurred, return the error
        return err
    }

    // Create the parent directory for the symlink
    if err := os.MkdirAll(destDir, 0755); err != nil {
        return err
    }

    // Create the symlink
    if err := os.Symlink(filepath.Join("..", "STORE", msgid+".msfn"), destPath); err != nil {
        return err
    }

    fmt.Printf("Tagged message %s as %s\n", msgid, tagName)
    return nil
}

// RemoveSymlink removes the symlink from ../tags/tagName/msgid.msfn
func Untag(msgid, tagName string) error {
    // Create the path to the symlink
    symlinkPath := filepath.Join("..", "tags", tagName, msgid+".msfn")

    // Check if the symlink exists
    _, err := os.Lstat(symlinkPath)
    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("Message is not in %s", tagName)
        }
        return err
    }

    // Remove the symlink
    if err := os.Remove(symlinkPath); err != nil {
        return err
    }

    fmt.Printf("Removed message from %s\n", symlinkPath)

	    // Check if the tagName directory is empty
    dir := filepath.Join("..", "tags", tagName)
    isEmpty, err := isDirEmpty(dir)
    if err != nil {
        return err
    }

    // If the directory is empty, remove it
    if isEmpty {
        if err := os.Remove(dir); err != nil {
            return err
        }
        fmt.Printf("Removed directory: %s\n", dir)
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

func HasTag(msgid, tagName string) (bool, error) {
    // Construct the path to the symlink
    symlinkPath := filepath.Join("..", "tags", tagName, msgid+".msfn")

    // Check if the symlink exists
    _, err := os.Lstat(symlinkPath)
    if err != nil {
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

func WalkTag(tagName string) error {
	dirName := filepath.Join("..", "tags", tagName)
	return filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("filepath.Walk error: %w", err)
		}
		lstat, err := os.Lstat(path)
		if err != nil {
			return fmt.Errorf("os.Lstat error: %w", err)
		} else if lstat.Mode()&os.ModeSymlink != 0 {
			symlinkName := filepath.Base(path)
            msgid := strings.TrimSuffix(symlinkName, ".msfn")
			if err := PrintLinkAndPreview(msgid); err != nil {
				fmt.Println(err)
			}
		}
        return nil
    })
}
