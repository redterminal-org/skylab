package main

import (
	"fmt"
	"strings"
	"os"
	"io/ioutil"
	"path/filepath"
	"github.com/clseibold/go-gemini"
	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
)

func deleteMail() {
	username, _ := tools.Authenticate()
	query := os.Getenv("QUERY_STRING")
	if strings.Contains(query, "..") {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("That query is not allowed")
		return
	}
	sitecfg, err := config.ExportSiteConfig()
	if err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Couldn't export site-wide config", err)
		return
	}
	skylabPath := sitecfg.SkylabPath
	if query != "" {
		filename := filepath.Join("..", "users", username, "STORE", query+".msfn")
		_, err := os.Stat(filename)
		if err != nil {
			fmt.Printf("30 /%s\r\n", skylabPath)
			return
		}
		if err := tools.UntagAll(query, username); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("Error: Could not remove from somewhere")
			fmt.Println(err)
			fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
			return
		} else {
			filename := filepath.Join("..", "users", username, "STORE", query+".msfn")
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("File was already gone.")
				fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
				return
			}

			err := os.Remove(filename)
			if err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Could not remove file.")
				fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
				return
			}

			cfg, err := config.ExportConfig(username)
			if err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Error getting config for Tag:", err)
				fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
				return
			}
			if cfg.TagSync {
				client := gemini.DefaultClient
				certFile, err := ioutil.ReadFile(cfg.CertFilename)
				if err != nil {
					fmt.Print("20 text/gemini\r\n")
					fmt.Printf("Couldn't read cert: %w (%s)\n", err, cfg.CertFilename)
					fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
					return
				}
				if cfg.GMAPHost == "" {
					fmt.Print("20 text/gemini\r\n")
					fmt.Println("Your GMAP host configuration is wrong. Go fix it:")
					fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
					return
				}
				var port string
				if cfg.GMAPPort == "" {
					port = "1958"
				} else {
					port = cfg.GMAPPort
				}
				deletePath := "gemini://" + cfg.GMAPHost + ":" + port + "/delete?" + query
				resp, err := client.FetchWithCert(deletePath, certFile, certFile)
				if err != nil || (resp.Status != 20 && resp.Status != 51) {
					fmt.Print("20 text/gemini\r\n")
					fmt.Printf("Couldn't delete message: %w (%s)", err, deletePath)
					fmt.Printf("=> %s/inbox 📪 Inbox\n", skylabPath)
					return
				}
				defer resp.Body.Close()
			}
			fmt.Printf("30 /%s/trash\r\n", skylabPath)
		}
	} else {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("# Empty Graveyard 🪦")
		fmt.Println("There is no trace left of your deleted messages. I'd return home if I were you:")
		fmt.Printf("=> /%s 🏠 Home\n", skylabPath)
	}
}
