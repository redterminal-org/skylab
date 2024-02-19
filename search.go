package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
//    "time"
	"net/url"
	"sort"

	"git.sr.ht/~satchlj/skylab/tools"
	"git.sr.ht/~satchlj/skylab/send"
)

type SearchOptions struct {
	From    []string
	To      []string
	Subject []string
	Body    []string
	Tags	[]string
	All		[]string
	Not		*SearchOptions
}

func search() {
	query := os.Getenv("QUERY_STRING")
	query, err := url.QueryUnescape(query)
	query = strings.ToLower(query)
	if err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println(err)
		return
	}
	if query == "" {
		fmt.Print("10 Search mail\r\n")
		return
	}
    searchOptions := parseQuery(query)
    msgs, err := searchMessages(searchOptions)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Print("20 text/gemini\r\n")
	fmt.Println("=> ? 🔎 Search again")
	fmt.Println("=> ./ 🏠 Home\n")
	for _, output := range msgs {
        fmt.Println(output.Text)
    }
	fmt.Println("\n=> gemini://satch.xyz/skylab/search.gmi How to search Skylab")
}

func searchMessages(options SearchOptions) ([]tools.Preview, error) {
    var matches []string
	username, _ := tools.Authenticate()
    directory := filepath.Join("..", "users", username, "STORE")

    err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        message, err := tools.ParseSkymail(path)
        if err != nil {
            return err
        }

        if match(message, options) {
			symlinkName := filepath.Base(path)
			msgid := strings.TrimSuffix(symlinkName, ".msfn")
            matches = append(matches, msgid)
        }

        return nil
    })
	if err != nil {
		return nil, err
	}
	var msgs []tools.Preview
	for _, msgid := range matches {
		hide, err := tools.HasTag(msgid, "Trash", username)
		if err != nil {
			return nil, err
		}
		if !hide {
			output, timestamp, err := tools.PrintLinkAndPreview(msgid, username, true)
			if err != nil {
				return nil, err
			}
			msgs = append(msgs, tools.Preview{output, timestamp})
		}
	}
	sortByTime := func(i, j int) bool {
		timeI := msgs[i].Timestamp
		timeJ := msgs[j].Timestamp
		return timeI.After(timeJ)
	}
	// Sort the outputs slice
	sort.Slice(msgs, sortByTime)

	return msgs, nil
}

func match(message tools.Skymail, options SearchOptions) bool {
    // Check if any 'Not' options match the message
    if options.Not != nil && match(message, *options.Not) {
        return false
    }

    // Check if 'From' option matches any sender
    if len(options.From) > 0 {
		matchedFrom := false
        for _, sender := range message.Senders {
            for _, from := range options.From {
                if strings.Contains(strings.ToLower(sender), strings.ToLower(from)) {
                    matchedFrom = true
                    break
                }
            }
        }
        if !matchedFrom {
            return false
        }
    }

    // Check if 'To' option matches any recipient
    if len(options.To) > 0 {
        matchedTo := false
        for _, recipient := range message.Recipients {
            for _, to := range options.To {
                if strings.Contains(strings.ToLower(recipient), strings.ToLower(to)) {
                    matchedTo = true
                    break
                }
            }
        }
        if !matchedTo {
            return false
        }
    }

    // Check if 'Subject' option matches the message subject
    if len(options.Subject) > 0 {
        matched := false
        for _, subject := range options.Subject {
            if strings.Contains(strings.ToLower(message.Subject), strings.ToLower(subject)) {
                matched = true
                break
            }
        }
        if !matched {
            return false
        }
    }

    // Check if 'Body' option matches the message body
    if len(options.Body) > 0 {
        matched := false
        for _, body := range options.Body {
            if strings.Contains(strings.ToLower(message.MessageBody), strings.ToLower(body)) {
                matched = true
                break
            }
        }
        if !matched {
            return false
        }
    }
    // Add similar checks for other fields like 'Tags' and 'All' if needed
	if len(options.All) > 0 {
		wholeMessage := send.SkymailToString(message)
		matched := false
		for _, term := range options.All {
			if strings.Contains(strings.ToLower(wholeMessage), strings.ToLower(term)) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
    // If all checks passed, return true
    return true
}

func parseQuery(query string) SearchOptions {
	options := SearchOptions{}

	terms := strings.Fields(query)
	var notQuery string
	for _, term := range terms {
		if strings.HasPrefix(term, "from:") {
			options.From = append(options.From, strings.TrimPrefix(term, "from:"))
		} else if strings.HasPrefix(term, "to:") {
			options.To = append(options.To, strings.TrimPrefix(term, "to:"))
		} else if strings.HasPrefix(term, "subject:") {
			options.Subject = append(options.Subject, strings.TrimPrefix(term, "subject:"))
		} else if strings.HasPrefix(term, "body:") {
			options.Body = append(options.Body, strings.TrimPrefix(term, "body:"))
		} else if strings.HasPrefix(term, "tag:") {
			options.Tags = append(options.Tags, strings.TrimPrefix(term, "tag:"))
		} else if term[0] == '!' {
			notQuery += " " + term[1:]
		} else {
			options.All = append(options.All, strings.TrimSpace(term))
		}
	}
	if notQuery != "" {
		notOptions := parseQuery(notQuery)
		options.Not = &notOptions
	}
	return options
}
