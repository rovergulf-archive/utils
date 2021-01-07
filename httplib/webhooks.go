package httplib

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func SendSlackHook(url, channel, username, message string) error {
	data := fmt.Sprintf(`{"channel": "#%s", "username":"%s", "text":"[%s] %s"}`, channel, username, time.Now().Format(time.RFC850), message)

	r := bytes.NewReader([]byte(data))

	_, err := http.Post(url, "application/json", r)
	if err != nil {
		return err
	}

	return nil
}
