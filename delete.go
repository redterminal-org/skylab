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
	if query != "" {
		filename := filepath.Join("..", "users", username, "STORE", query+".msfn")
		_, err := os.Stat(filename)
		if err != nil {
			fmt.Println("30 /\r\n")
			return
		}
		if err := tools.UntagAll(query, username); err != nil {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("Error: Could not remove from somewhere")
			fmt.Println(err)
			fmt.Println("=> inbox 📪 Inbox")
			return
		} else {
			filename := filepath.Join("..", "users", username, "STORE", query+".msfn")
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("File was already gone.")
				fmt.Println("=> inbox 📪 Inbox")
				return
			}

			err := os.Remove(filename)
			if err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Could not remove file.")
				fmt.Println("=> inbox 📪 Inbox")
				return
			}

			cfg, err := config.ExportConfig(username)
			if err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Error getting config for Tag:", err)
				fmt.Println("=> inbox 📪 Inbox")
				return
			}
			if cfg.TagSync {
				client := gemini.DefaultClient
				certFile, err := ioutil.ReadFile(cfg.CertFilename)
				if err != nil {
					fmt.Print("20 text/gemini\r\n")
					fmt.Printf("Couldn't read cert: %w (%s)\n", err, cfg.CertFilename)
					fmt.Println("=> inbox 📪 Inbox")
					return
				}
				if cfg.GMAPHost == "" {
					fmt.Print("20 text/gemini\r\n")
					fmt.Println("Your GMAP host configuration is wrong. Go fix it:")
					fmt.Println("=> settings ⚙️ Settings")
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
				if err != nil || resp.Status != 20 {
					fmt.Print("20 text/gemini\r\n")
					fmt.Printf("Couldn't delete message: %w (%s)", err, deletePath)
					fmt.Println("=> inbox 📪 Inbox")
					return
				}
				defer resp.Body.Close()
			}
			fmt.Print("30 /trash\r\n")
		}
	} else {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("# Empty Graveyard 🪦")
		fmt.Println("There is no trace left of your deleted messages. I'd return home if I were you:")
		fmt.Println("=> / 🏠 Home")
	}
}
