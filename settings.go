package main

import (
	"fmt"
	"git.sr.ht/~satchlj/skylab/config"
	"git.sr.ht/~satchlj/skylab/tools"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func settings() {
	username, _ := tools.Authenticate()
	query := os.Getenv("QUERY_STRING")
	decodedQuery, _ := url.QueryUnescape(query)
	path := os.Getenv("PATH_INFO")
	cfg, err := config.ExportConfig(username)
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
	skylabPath := appPrefix(sitecfg.SkylabPath)
	if skylabPath != "" {
		path = strings.TrimPrefix(path, skylabPath)
	} else if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	switch path {
	case "/settings/quote-reply":
		if query == "true" {
			if err := config.Set(username, "QuoteReply", true); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Printf("30 %s/settings\r\n", skylabPath)
			return
		}
		if query == "false" {
			if err := config.Set(username, "QuoteReply", false); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Printf("30 %s/settings\r\n", skylabPath)
			return
		} else {
			fmt.Print("10 Please enter 'true' or 'false':\r\n")
			return
		}
	case "/settings/misfin-link":
		if query == "true" {
			if err := config.Set(username, "MisfinLink", true); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Printf("30 %s/settings\r\n", skylabPath)
			return
		}
		if query == "false" {
			if err := config.Set(username, "MisfinLink", false); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Printf("30 %s/settings\r\n", skylabPath)
			return
		} else {
			fmt.Print("10 Please enter 'true' or 'false':\r\n")
			return
		}
	case "/settings/fetch-address":
		if !sitecfg.ExposeSettings {
			fmt.Printf("30 %s\r\n", appPath(skylabPath, ""))
			return
		}
		if query != "" {
			if err := config.Set(username, "GeminiFetchAddress", decodedQuery); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Printf("30 %s/settings\r\n", skylabPath)
			return
		} else {
			fmt.Print("10 Enter your gembox fetch address (gemini://example.org:1958/fetch or similar):\r\n")
			return
		}
	case "/settings/header-text":
		if query != "" {
			if err := config.Set(username, "HeaderText", decodedQuery); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Printf("30 %s/settings\r\n", skylabPath)
			return
		} else {
			fmt.Print("10 Write your own homepage header. Default is 'Misfin ✉️'\r\n")
			return
		}
	case "/settings/cert-filename":
		if !sitecfg.ExposeSettings {
			fmt.Printf("30 %s\r\n", appPath(skylabPath, ""))
			return
		}
		if query != "" {
			if err := config.Set(username, "CertFilename", decodedQuery); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Printf("30 %s/settings\r\n", skylabPath)
			return
		} else {
			fmt.Print("10 Ask your server admin what this should be if you're not sure. Use an absolute path to make things easier (/home/misfinuser/skylab/users/yourname/cert.pem or similar):\r\n")
			return
		}
	case "/settings/gmap-parse":
		if !sitecfg.ExposeSettings {
			fmt.Printf("30 %s\r\n", appPath(skylabPath, ""))
			return
		}
		if query == "true" {
			if err := config.Set(username, "GMAP", true); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Printf("30 %s/settings\r\n", skylabPath)
			return
		}
		if query == "false" {
			if err := config.Set(username, "GMAP", false); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Printf("30 %s/settings\r\n", skylabPath)
			return
		} else {
			fmt.Print("10 Please enter 'true' or 'false':\r\n")
			return
		}
	case "/settings/gmap-port":
		if !sitecfg.ExposeSettings {
			fmt.Printf("30 %s\r\n", appPath(skylabPath, ""))
			return
		}
		_, err := strconv.Atoi(query)
		if query != "" && err == nil {
			if err := config.Set(username, "GMAPPort", query); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Printf("30 %s/settings\r\n", skylabPath)
			return
		} else {
			fmt.Print("10 Please enter the port by which to make gemini requests to your GMAP misfin server (default 1958):\r\n")
			return
		}
	case "/settings/gmap-host":
		if !sitecfg.ExposeSettings {
			fmt.Printf("30 %s\r\n", appPath(skylabPath, ""))
			return
		}
		if query != "" {
			if err := config.Set(username, "GMAPHost", strings.TrimSpace(decodedQuery)); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Printf("30 %s/settings\r\n", skylabPath)
			return
		} else {
			fmt.Print("10 Enter the hostname of your GMAP misfin server (example.org):\r\n")
			return
		}
	case "/settings/home-display":
		displayNumber, err := strconv.Atoi(query)
		if query != "" && err == nil {
			if err := config.Set(username, "DisplayNumber", displayNumber); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Printf("30 %s/settings\r\n", skylabPath)
			return
		} else {
			fmt.Print("10 Please enter an integer number of recent messages to display on the home page ('4' recommended):\r\n")
			return
		}
	case "/settings/autofetch":
		if query == "true" {
			if err := config.Set(username, "AutoFetch", true); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Printf("30 %s/settings\r\n", skylabPath)
			return
		}
		if query == "false" {
			if err := config.Set(username, "AutoFetch", false); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Printf("30 %s/settings\r\n", skylabPath)
			return
		} else {
			fmt.Print("10 Please enter 'true' or 'false':\r\n")
			return
		}
	case "/settings/timezone":
		if query != "" {
			if err := config.Set(username, "Timezone", decodedQuery); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Printf("30 %s/settings\r\n", skylabPath)
			return
		} else {
			fmt.Print("10 Please enter your tz database time zone (America/New_York or similar):\r\n")
			return
		}
	case "/settings/password":
		password, err := tools.GenerateOneTimePassword(username)
		if err != nil {
			fmt.Println("20 text/gemini\r\nCouldn't generate a one time password:", err)
			return
		}
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Your one-time password is:", password)
		fmt.Println("This is the last time it will be shown. You can have multiple one-time passwords at a time, but that's probably a bad idea. Each password can be used exactly once.")
		return
	case "/settings/clear-passwords":
		if err := tools.ClearTempPasswords(username); err != nil {
			fmt.Println("20 text/gemini\r\nCouldn't clear your passwords", err)
			return
		}
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Your temporary passwords were cleared.")
		fmt.Printf("=> %s/settings ⚙️ Settings\n", skylabPath)
		fmt.Printf("=> %s/ 🏠 Go Home", skylabPath)
	default:
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("# Edit Settings ⚙️")
		fmt.Println("### Client behavior")
		fmt.Printf("=> %s/settings/autofetch Fetch new messages every page load\n", skylabPath)
		fmt.Printf("=> %s/settings/misfin-link Use misfin: links for quick reply\n", skylabPath)
		if cfg.MisfinLink {
			fmt.Printf("=> %s/settings/quote-reply Quote previous message in replies\n", skylabPath)
		}
		fmt.Printf("=> %s/settings/timezone Set your tz database time zone (currently set to '%s')\n", skylabPath, cfg.Timezone)
		fmt.Println("### Display options")
		fmt.Printf("=> %s/settings/home-display Set max number of messages to display on homepage (0 to disable)\n", skylabPath)
		fmt.Printf("=> %s/settings/header-text Set your homepage header text\n", skylabPath)
		if sitecfg.ExposeSettings {
			fmt.Println("### Misfin account configuration")
			fmt.Printf("=> %s/settings/cert-filename Edit your cert filename (careful with this one! Maybe have your admin do it for you?)\n", skylabPath)
			fmt.Printf("=> %s/settings/fetch-address Edit gembox fetch address\n", skylabPath)
			fmt.Printf("=> %s/settings/gmap-parse Use GMAP Gembox parsing\n", skylabPath)
			fmt.Println("### GMAP settings")
			fmt.Println("You can ignore this section if you're not using multiple clients with the same server")
			fmt.Printf("=> %s/settings/gmap-host Hostname of GMAP misfin server\n", skylabPath)
			fmt.Printf("=> %s/settings/gmap-port Port number of GMAP misfin server (1958 by default)\n", skylabPath)
			// fmt.Printf("=> %s/settings/tag-sync Tell the server when you modify message tags (sync unread status, etc)\n", skylabPath)
		}
		fmt.Println("### Gemini client certificates for Skylab")
		fmt.Printf("=> %s/settings/password Get a temporary password\n(so that you can add a new gemini client certificate to your account)\n", skylabPath)
		fmt.Printf("=> %s/settings/clear-passwords Clear all previously created temporary passwords\n", skylabPath)
		fmt.Println("### Bug Reports and Feature Requests")
		fmt.Println("=> gemini://bbs.geminispace.org/s/Skylab Submit bug reports and feature requests here")
	}
}
