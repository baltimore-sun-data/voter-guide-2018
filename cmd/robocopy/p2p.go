package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type P2PClient struct{ Token string }

func NewP2PClient() P2PClient {
	return P2PClient{
		Token: os.Getenv("P2P_ACCESS_TOKEN"),
	}
}

func (pc P2PClient) Update(slug, content string) error {
	const baseURL = "https://content-api.p2p.tribuneinteractive.com/content_items/%s.json?preserve_embedded_tags=false"
	url := fmt.Sprintf(baseURL, slug)

	type (
		Body struct {
			Body string `json:"body"`
		}
		Wrapper struct {
			Content Body `json:"content_item"`
		}
	)

	var body bytes.Buffer
	enc := json.NewEncoder(&body)
	_ = enc.Encode(Wrapper{Body{content}})

	req, _ := http.NewRequest("PUT", url, &body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+pc.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("P2P connection error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response code from P2P: %s", resp.Status)
	}
	return nil
}
