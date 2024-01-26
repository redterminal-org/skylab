package main

import (
    "fmt"
    "os"
    "errors"
    "net/url"
    "path/filename"
    "strings"
    "time"

    "git.sr.ht/~satchlj/skylab/fetch"
)

var Version string = "0.0.1" // Haha

func main() {
	// Configuration
	config := fetch.Config{
        GeminiFetchAddress: "gemini://satch.xyz:1958/fetch",
		CertFilename:       "cert.pem",
		// Add other configuration fields as needed
	}

	// Fetch gembox from Gemini server
	gemboxData, err := fetch.FetchGembox(config.GeminiFetchAddress, config.CertFilename)
	if err != nil {
		panic(err)
	}
