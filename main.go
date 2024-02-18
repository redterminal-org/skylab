package main

import (
    "fmt"
	"os"
	"strings"
	"git.sr.ht/~satchlj/skylab/tools"
	"git.sr.ht/~satchlj/skylab/config"
)

func main() {
	query := os.Getenv("QUERY_STRING")
	path := os.Getenv("PATH_INFO")
	if os.Getenv("TLS_CLIENT_HASH") == "" {
		if query != "login" {
			fmt.Print("20 text/gemini\r\n")
			fmt.Println("# Welcome to the Skylab Misfin Client")
			fmt.Println("=> ?login Login")
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
	sitecfg, err := config.ExportSiteConfig()
	if err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Couldn't export site-wide config", err)
		return
	}
	skylabPath := sitecfg.SkylabPath
	if username == "" {
		if query == "" {
			fmt.Print("20 text/gemini\r\nDidn't find you anywhere in my list. If you have an account already, then you can log in with a single-use password:\n")
			fmt.Println("=> ?login Login")
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
			fmt.Printf("30 /%s\r\n", skylabPath)
		}
	}
	if query == "login" {
		fmt.Printf("30 /%s\r\n", skylabPath)
		return
	}
	if skylabPath != "" {
		skylabPath = "/" + skylabPath
	}
	if path == skylabPath || path == skylabPath + "/" {
		home()
	} else if strings.HasPrefix(path, skylabPath+"/mail") {
		mail()
	} else if strings.HasPrefix(path, skylabPath+"/archive") {
		archive()
	} else if strings.HasPrefix(path, skylabPath+"/compose") {
		compose()
	} else if strings.HasPrefix(path, skylabPath+"/delete") {
		deleteMail()
	} else if strings.HasPrefix(path, skylabPath+"/drafts") {
		drafts()
	} else if strings.HasPrefix(path, skylabPath+"/inbox") {
		inbox()
	} else if strings.HasPrefix(path, skylabPath+"/reply") {
		reply()
	} else if strings.HasPrefix(path, skylabPath+"/send") {
		sender()
	} else if strings.HasPrefix(path, skylabPath+"/sent") {
		sent()
	} else if strings.HasPrefix(path, skylabPath+"/settings") {
		settings()
	} else if strings.HasPrefix(path, skylabPath+"/tags") {
		tags()
	} else if strings.HasPrefix(path, skylabPath+"/tag") {
		tag()
	} else if strings.HasPrefix(path, skylabPath+"/trash") {
		trash()
	} else if strings.HasPrefix(path, skylabPath+"/unread") {
		unread()
	} else {
		fmt.Print("51 Skylab doesn't know what you're looking for. If you think this is an error, send mail to mail@satch.xyz\r\n")
	}
}
