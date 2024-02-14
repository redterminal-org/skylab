package main

import (
    "fmt"
	"os"
	"strings"
	"git.sr.ht/~satchlj/skylab/tools"
)

func main() {
	query := os.Getenv("QUERY_STRING")
	path := os.Getenv("PATH_INFO")
	if os.Getenv("TLS_CLIENT_HASH") == "" {
		if query != "login" {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("# Welcome to the Skylab Misfin Client")
			fmt.Println("=> /?login Login")
		} else {
			fmt.Print("60 Please use a client certificate\r\n")
		}
		return
	}
	username, err := tools.Authenticate()
	if err != nil {
		fmt.Print("20 text/gemini\r\nError with somethin, prolly the server admin's fault :)")
		return
	}
	if username == "" {
		if query == "" {
			fmt.Print("20 text/gemini\r\nDidn't find you anywhere in my list. If you have an account already, then you can log in with a single-use password:\n")
			fmt.Println("=> /?login Login")
			return
		} else if query == "login" {
			fmt.Print("11 Enter your single-use password:\r\n")
			return
		} else {
			if err := tools.CheckTempPassword(query); err != nil {
				fmt.Print("20 text/gemini\r\n")
				fmt.Println("Something went wrong trying to add your cert:", err)
				return
			}
			fmt.Print("30 /\r\n")
		}
	}
	if query == "login" {
		fmt.Print("30 /\r\n")
		return
	}
	if path == "/" || path == "" {
		home()
	} else if strings.HasPrefix(path, "/mail") {
		mail()
	} else if strings.HasPrefix(path, "/archive") {
		archive()
	} else if strings.HasPrefix(path, "/compose") {
		compose()
	} else if strings.HasPrefix(path, "/delete") {
		deleteMail()
	} else if strings.HasPrefix(path, "/drafts") {
		drafts()
	} else if strings.HasPrefix(path, "/inbox") {
		inbox()
	} else if strings.HasPrefix(path, "/reply") {
		reply()
	} else if strings.HasPrefix(path, "/send") {
		sender()
	} else if strings.HasPrefix(path, "/sent") {
		sent()
	} else if strings.HasPrefix(path, "/settings") {
		settings()
	} else if strings.HasPrefix(path, "/tag") {
		tag()
	} else if strings.HasPrefix(path, "/tags") {
		tags()
	} else if strings.HasPrefix(path, "/trash") {
		trash()
	} else if strings.HasPrefix(path, "/unread") {
		unread()
	} else {
		fmt.Print("51 Skylab doesn't know what you're looking for. If you think this is an error, send mail to mail@satch.xyz\r\n")
	}
}
