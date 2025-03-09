package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func SteamIDFromTicket(ticket string) ([]byte, error) {
	if ticket == "" {
		return nil, fmt.Errorf("missing ticket")
	}

	v := make(url.Values)
	v.Set("ticket", ticket)

	resp, err := http.Get(fmt.Sprintf("%s/auth/getid?%s", os.Getenv("API_URL"), v.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to get steamid: %s", err)
	}

	steamid, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read steamid: %s", err)
	}

	return steamid, nil
}
