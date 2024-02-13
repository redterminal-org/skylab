package main

import (
	"fmt"
	"strings"
	"os"
	"net/url"
	"git.sr.ht/~satchlj/skylab/tools"
	"git.sr.ht/~satchlj/skylab/config"
)

func main() {
	if os.Getenv("TLS_CLIENT_HASH") == "" {
		fmt.Print("30 /\r\n")
		return
	}
	username, err := tools.Authenticate()
	if err != nil {
		fmt.Print("20 text/gemini\r\nError with somethin, prolly not ur fault :)")
		return
	}
	if username == "" {
		fmt.Print("30 /\r\n")
		return
	}
	query := os.Getenv("QUERY_STRING")
	decodedQuery, _ := url.QueryUnescape(query)
	path := os.Getenv("PATH_INFO")
	if err != nil {
		fmt.Println("Couldn't get the config info:", err)
		return
	}
	sitecfg, err := config.ExportSiteConfig()
	if err != nil {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Couldn't export site-wide config", err)
		return
	}
	if strings.HasPrefix(path, "/settings/fetch-address") {
		if !sitecfg.ExposeSettings {
			fmt.Print("30 /\r\n")
			return
		}
		if query != "" {
			if err := config.Set(username, "GeminiFetchAddress", decodedQuery); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Print("30 /settings\r\n")
			return
		} else {
			fmt.Print("10 Enter your gembox fetch address (gemini://example.org:1958/fetch or similar):\r\n")
			return
		}
	} else if strings.HasPrefix(path, "/settings/cert-filename") {
		if !sitecfg.ExposeSettings {
			fmt.Print("30 /\r\n")
			return
		}
		if query != "" {
			if err := config.Set(username, "CertFilename", decodedQuery); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Print("30 /settings\r\n")
			return
		} else {
			fmt.Print("10 Ask your server admin what this should be if you're not sure. Use an absolute path to make things easier (/home/misfinuser/skylab/users/yourname/cert.pem or similar):\r\n")
			return
		}
	} else if strings.HasPrefix(path, "/settings/gmap") {
		if query == "true" {
			if err := config.Set(username, "GMAP", true); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Print("30 /settings\r\n")
			return
		}
		if query == "false" {
			if err := config.Set(username, "GMAP", false); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Print("30 /settings\r\n")
			return
		} else {
			fmt.Print("10 Please enter 'true' or 'false':\r\n")
			return
		}
	} else if strings.HasPrefix(path, "/settings/autofetch") {
		if query == "true" {
			if err := config.Set(username, "AutoFetch", true); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Print("30 /settings\r\n")
			return
		}
		if query == "false" {
			if err := config.Set(username, "AutoFetch", false); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Print("30 /settings\r\n")
			return
		} else {
			fmt.Print("10 Please enter 'true' or 'false':\r\n")
			return
		}
	} else if strings.HasPrefix(path, "/settings/timezone") {
		if query != "" {
			if err := config.Set(username, "Timezone", decodedQuery); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Print("30 /settings\r\n")
			return
		} else {
			fmt.Print("10 Please enter your tz database time zone (America/New_York or similar):\r\n")
			return
		}
	} else if strings.HasPrefix(path, "/settings/password") {
		password, err := tools.GenerateOneTimePassword(username)
		if err != nil {
			fmt.Println("Couldn't generate a one time password:", err)
			return
		}
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Your one-time password is:", password)
		fmt.Println("This is the last time it will be shown. You can have multiple one-time passwords at a time, but that's probably a bad idea. Each password can be used exactly once.")
		return
	} else {
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("# Edit Settings ⚙️")
		fmt.Println("=> /settings/autofetch Set autofetch true/false")
		fmt.Println("=> /settings/timezone Set your tz database time zone\n")
		fmt.Println("=> /settings/password Get a temporary password\n(so that you can add a new gemini client certificate to your account)")
		if sitecfg.ExposeSettings {
			fmt.Println("=> /settings/cert-filename Edit your cert filename (careful with this one! Maybe have your admin do it for you?)")
			fmt.Println("=> /settings/fetch-address Edit gembox fetch address")
			fmt.Println("=> /settings/gmap Use GMAP Gembox parsing")
		}
		fmt.Println("=> gemini://bbs.geminispace.org/s/Skylab Submit bug reports and feature requests here")
	}
}
