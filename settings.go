package main
import (
	"fmt"
	"strings"
	"os"
	"net/url"
	"strconv"
	"git.sr.ht/~satchlj/skylab/tools"
	"git.sr.ht/~satchlj/skylab/config"
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
	switch path {
	case "/settings/fetch-address":
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
	case "/settings/header-text":
		if query != "" {
			if err := config.Set(username, "HeaderText", decodedQuery); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Print("30 /settings\r\n")
			return
		} else {
			fmt.Print("10 Write your own homepage header. Default is 'Misfin ✉️'\r\n")
			return
		}
	case "/settings/cert-filename":
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
	case "/settings/gmap-parse":
		if !sitecfg.ExposeSettings {
			fmt.Print("30 /\r\n")
			return
		}
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
	case "/settings/gmap-port":
		if !sitecfg.ExposeSettings {
			fmt.Print("30 /\r\n")
			return
		}
		_, err := strconv.Atoi(query)
		if query != "" && err == nil {
			if err := config.Set(username, "GMAPPort", query); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Println("30 /settings\r\n")
			return
		} else {
			fmt.Print("10 Please enter the port by which to make gemini requests to your GMAP misfin server (default 1958):\r\n")
			return
		}
	case "/settings/gmap-host":
		if !sitecfg.ExposeSettings {
			fmt.Print("30 /\r\n")
			return
		}
		if query != "" {
			if err := config.Set(username, "GMAPHost", strings.TrimSpace(decodedQuery)); err != nil {
				fmt.Println("Couldn't set the config value:", err)
				return
			}
			fmt.Print("30 /settings\r\n")
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
			fmt.Println("30 /settings\r\n")
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
	case "/settings/timezone":
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
	case "/settings/password":
		password, err := tools.GenerateOneTimePassword(username)
		if err != nil {
			fmt.Println("Couldn't generate a one time password:", err)
			return
		}
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("Your one-time password is:", password)
		fmt.Println("This is the last time it will be shown. You can have multiple one-time passwords at a time, but that's probably a bad idea. Each password can be used exactly once.")
		return
	default:
		fmt.Print("20 text/gemini\r\n")
		fmt.Println("# Edit Settings ⚙️")
		fmt.Println("### Client behavior")
		fmt.Println("=> /settings/autofetch Set autofetch true/false")
		fmt.Printf("=> /settings/timezone Set your tz database time zone (currently set to '%s')\n", cfg.Timezone)
		fmt.Println("### Display options")
		fmt.Println("=> /settings/home-display Set max number of messages to display on homepage (0 to disable)")
		fmt.Println("=> /settings/header-text Set your homepage header text")
		if sitecfg.ExposeSettings {
			fmt.Println("### Misfin account configuration")
			fmt.Println("=> /settings/cert-filename Edit your cert filename (careful with this one! Maybe have your admin do it for you?)")
			fmt.Println("=> /settings/fetch-address Edit gembox fetch address")
			fmt.Println("=> /settings/gmap-parse Use GMAP Gembox parsing")
			fmt.Println("### GMAP settings")
			fmt.Println("You can ignore this section if you're not using multiple clients with the same server")
			fmt.Println("=> /settings/gmap-host Hostname of GMAP misfin server")
			fmt.Println("=> /settings/gmap-port Port number of GMAP misfin server (1958 by default)")
			fmt.Println("=> /settings/tag-sync Tell the server when you modify message tags (sync unread status, etc)")
		}
		fmt.Println("### Gemini client certificates for Skylab")
		fmt.Println("=> /settings/password Get a temporary password\n(so that you can add a new gemini client certificate to your account)")
		fmt.Println("### Bug Reports and Feature Requests")
		fmt.Println("=> gemini://bbs.geminispace.org/s/Skylab Submit bug reports and feature requests here")
	}
}
