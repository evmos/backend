// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package requester

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/tharsis/dashboard-backend/internal/v1/db"
)

type Tree struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	Sha  string `json:"sha"`
	URL  string `json:"url"`
}

type TreeResponse struct {
	Sha       string `json:"sha"`
	URL       string `json:"url"`
	Tree      []Tree `json:"tree"`
	Truncated bool   `json:"truncated"`
}

type Content struct {
	Content string `json:"content"`
	Sha     string `json:"sha"`
}

type File struct {
	Content string
	URL     string
}

func QueryGithubWithCache(url string) (string, error) {
	if val, err := db.RedisGetGithubResponse(url); err == nil {
		return val, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	githubKey := os.Getenv("GITHUB_KEY")

	req.Header.Add("authorization", "token "+githubKey)

	resp, err := Client.Do(req)
	if err != nil {
		if val, err := db.RedisGetHithubFallbackResponse(url); err == nil {
			return val, nil
		}
		return "", err
	}

	if resp.StatusCode != 200 {
		if val, err := db.RedisGetHithubFallbackResponse(url); err == nil {
			return val, nil
		}
		return "", fmt.Errorf("github response status code different from 200: %d", resp.StatusCode)

	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil || len(string(body)) == 0 {
		if val, err := db.RedisGetHithubFallbackResponse(url); err == nil {
			return val, nil
		}
		return "", err
	}

	// NOTE: converting to string and back to bytes is not the best,
	// but to support everything in redis as a string it's worth it
	bodyString := string(body)

	db.RedisSetGithubResponse(url, bodyString)
	db.RedisSetGithubFallbackResponse(url, bodyString)
	return bodyString, nil
}

func getChainTokenRegistryURL() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "production" {
		return "https://api.github.com/repos/evmos/chain-token-registry/git/trees/production?recursive=1"
	}
	return "https://api.github.com/repos/evmos/chain-token-registry/git/trees/main?recursive=1"
}

func GetValidatorDirectory() ([]File, error) {
	ValidatorsDirectoryURL := "https://api.github.com/repos/evmos/validator-directory/git/trees/main?recursive=1"
	return GetJsonsFromFolder(ValidatorsDirectoryURL, "mainnet")
}

func GetERC20TokensDirectory() ([]File, error) {
	ERC20TokensDirectoryURL := getChainTokenRegistryURL()
	return GetJsonsFromFolder(ERC20TokensDirectoryURL, "tokens")
}

func GetNetworkConfig() ([]File, error) {
	url := getChainTokenRegistryURL()
	return GetJsonsFromFolder(url, "chainConfig")
}

func GetJsonsFromFolder(url string, folder string) ([]File, error) {
	res := []File{}
	apiResp, err := QueryGithubWithCache(url)
	if err != nil {
		return []File{}, err
	}

	var m TreeResponse
	err = json.Unmarshal([]byte(apiResp), &m)
	if err != nil {
		return []File{}, err
	}

	for _, t := range m.Tree {
		if t.Mode == "100644" {
			// Is file
			if strings.HasPrefix(t.Path, folder+"/") {
				fileResponse, err := QueryGithubWithCache(t.URL)
				if err != nil {
					return []File{}, err
				}

				var m Content
				err = json.Unmarshal([]byte(fileResponse), &m)
				if err == nil {
					rawDecodedText, err := base64.StdEncoding.DecodeString(m.Content)
					if err != nil {
						return []File{}, err
					}
					res = append(res, File{Content: string(rawDecodedText), URL: t.Path})
				}
			}
		}
	}
	return res, nil
}
