// Package lifxutil is the utils for interacting with the lifx rest API.
package lifxutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://api.lifx.com/v1/"

// Client holds the Lifx API token.
type Client struct {
	Token string
}

// Light represents a Lifx light (partial fields).
type Light struct {
	ID         string  `json:"id"`
	Label      string  `json:"label"`
	Power      string  `json:"power"`
	Color      string  `json:"color"`
	Brightness float64 `json:"brightness"`
}

// NewClient creates a new Lifx API client.
func NewClient(token string) *Client {
	return &Client{Token: token}
}

// ListLights returns all lights for the account.
func (c *Client) ListLights() ([]Light, error) {
	url := baseURL + "lights/all"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Lifx API error: %s", string(body))
	}
	var lights []Light
	if err := json.NewDecoder(resp.Body).Decode(&lights); err != nil {
		return nil, err
	}
	return lights, nil
}

// TogglePower toggles the power state of a light by selector (e.g., "id:xxxx" or "label:MyLight").
func (c *Client) TogglePower(selector string) error {
	url := baseURL + "lights/" + selector + "/toggle"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)
	if resp.StatusCode != 207 && resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Lifx API error: %s", string(body))
	}
	return nil
}

// SetState sets the state of a light (color, brightness, etc.).
func (c *Client) SetState(selector string, state map[string]interface{}) error {
	url := baseURL + "lights/" + selector + "/state"
	bodyBytes, err := json.Marshal(state)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)
	if resp.StatusCode != 207 && resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("lifx API error: %s", string(body))
	}
	return nil
}

// SetBusy sets the state of the specified light to busy, using the provided color.
func (c *Client) SetBusy(light Light, color string) error {
	if color == "" {
		color = "red saturation:0.5" // fallback default
	}
	state := map[string]interface{}{
		"power": "on",
		"color": color,
	}
	return c.SetState("id:"+light.ID, state)
}

// SetFree sets the state of the specified light to available, using the provided color.
func (c *Client) SetFree(light Light, color string) error {
	if color == "" {
		color = "kelvin:2671" // fallback default
	}
	state := map[string]interface{}{
		"power":      "on",
		"color":      color,
		"brightness": 0.5,
	}
	return c.SetState("id:"+light.ID, state)
}
