package main

import "strings"

func appPrefix(prefix string) string {
	prefix = strings.Trim(prefix, "/")
	if prefix == "" {
		return ""
	}
	return "/" + prefix
}

func appPath(prefix, suffix string) string {
	prefix = strings.Trim(prefix, "/")
	suffix = strings.Trim(suffix, "/")

	if prefix == "" {
		if suffix == "" {
			return "/"
		}
		return "/" + suffix
	}

	if suffix == "" {
		return "/" + prefix
	}

	return "/" + prefix + "/" + suffix
}
