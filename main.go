package main

import (
    "git.sr.ht/~satchlj/skylab/fetch"
    "git.sr.ht/~satchlj/skylab/config"
)

var Version string = "0.0.1" // Haha

func main() {
    // Configuration
    cfg := config.ExportConfig()
    // Fetch gembox from Gemini server
    err := fetch.FetchGembox(cfg.GeminiFetchAddress, cfg.CertFilename)
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
