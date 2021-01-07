package httplib

import "fmt"

func PerformHTTPRequestUrl(host, url string, args map[string]string) string {
	var returnStr = fmt.Sprintf("%s/%s", host, url)

	if len(args) > 0 {
		returnStr += `?`
		for i := range args {
			returnStr += fmt.Sprintf("%s=%s&", i, args[i])
		}
	}

	return returnStr
}
