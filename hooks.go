package utils

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"
)

func SendSlackHook(url, channel, username, message string) error {
	r := bytes.NewReader([]byte(
		fmt.Sprintf(`{"channel":"#%s","username":"%s","text":"%s at: %s\n  %s"}`,
			channel, username, username, time.Now().Format(time.UnixDate), message)))

	if _, err := http.Post(url, "application/json", r); err != nil {
		log.Printf("Failed to send slack hook: %s", err)
		return err
	}

	return nil
}
