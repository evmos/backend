// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package requester

import (
	"io"
	"os"
	"strings"
)

func MakeAirtableGetRequest(path string) (string, error) {
	airtableKey := os.Getenv("AIRTABLE_KEY")

	key := "&api_key=" + airtableKey
	url := "https://api.airtable.com/v0/appv4nSZwDWdKTA8t" + path + key

	var sb strings.Builder
	sb.WriteString(url)

	resp, err := Client.Get(sb.String())
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 404 {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil || len(string(body)) == 0 {
		return "", err
	}

	return string(body), nil
}
