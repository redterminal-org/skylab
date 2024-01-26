package main

import (
    "git.sr.ht/~satchlj/skylab/fetch"
)

var Version string = "0.0.1" // Haha

type Config struct {
	GeminiFetchAddress string
	CertFilename       string
	// Add other configuration fields as needed
}

func main() {
	// Configuration
	config := Config{
        GeminiFetchAddress: "gemini://satch.xyz:1958/fetch",
		CertFilename:       "cert.pem",
		// Add other configuration fields as needed
	}

	// Fetch gembox from Gemini server
    err := fetch.FetchGembox(config.GeminiFetchAddress, config.CertFilename)
	    if err != nil {
		panic(err)
	}
}
