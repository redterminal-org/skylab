package main

import (
    "git.sr.ht/~satchlj/skylab/fetch"
//	"github.com/clseibold/go-gemini"
//    "fmt"
//	"io/ioutil"
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
//    mailText := testGemini(config.GeminiFetchAddress, config.CertFilename)
//    fmt.Println(mailText)
}

// func testGemini(address, certFilename string) string {
//     certFile, err := ioutil.ReadFile(certFilename)
//     if err != nil {
//         panic(err)
//     }
//
//     // Fetch gembox from Gemini server
//     client := gemini.DefaultClient
//     resp, err := client.FetchWithCert(address, certFile, certFile)
//     if err != nil {
//         panic(err)
//     }
//     defer resp.Body.Close()
//
//     // Read gembox data
//     gemboxData, err := ioutil.ReadAll(resp.Body)
//     if err != nil {
//         panic(err)
//     }
//
//     // Parse gembox
//     gbox := string(gemboxData)
//     return gbox
// }
