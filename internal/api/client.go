package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL   string
	Token     string
	HTTP      *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) SetToken(token string) {
	c.Token = token
}

func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	url := c.BaseURL + path

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Error != "" {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error)
		}
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *Client) Get(path string) ([]byte, error) {
	return c.doRequest(http.MethodGet, path, nil)
}

func (c *Client) Post(path string, body interface{}) ([]byte, error) {
	return c.doRequest(http.MethodPost, path, body)
}

func (c *Client) Register(name, email, password string) (*AuthResponse, error) {
	body := map[string]string{
		"name":     name,
		"email":    email,
		"password": password,
	}
	data, err := c.Post("/auth/register", body)
	if err != nil {
		return nil, err
	}
	var resp AuthResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &resp, nil
}

func (c *Client) Login(email, password string) (*AuthResponse, error) {
	body := map[string]string{
		"email":    email,
		"password": password,
	}
	data, err := c.Post("/auth/authenticate", body)
	if err != nil {
		return nil, err
	}
	var resp AuthResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetTeams(group string) (*TeamsResponse, error) {
	path := "/get/teams"
	if group != "" {
		path += "?group=" + group
	}
	data, err := c.Get(path)
	if err != nil {
		return nil, err
	}
	var resp TeamsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetTeamByID(id string) (*Team, error) {
	data, err := c.Get("/get/team/" + id)
	if err != nil {
		return nil, err
	}
	var resp TeamResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &resp.Team, nil
}

func (c *Client) GetTeamByName(name string) (*Team, error) {
	data, err := c.Get("/get/team?name=" + name)
	if err != nil {
		return nil, err
	}
	var resp TeamResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &resp.Team, nil
}

func (c *Client) GetGroups() (*GroupsResponse, error) {
	data, err := c.Get("/get/groups")
	if err != nil {
		return nil, err
	}
	var resp GroupsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetGroup(name string) (*GroupResponse, error) {
	data, err := c.Get("/get/group?name=" + name)
	if err != nil {
		return nil, err
	}
	var resp GroupResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetGames() (*GamesResponse, error) {
	data, err := c.Get("/get/games")
	if err != nil {
		return nil, err
	}
	var resp GamesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetGameByID(id string) (*Game, error) {
	data, err := c.Get("/get/game/" + id)
	if err != nil {
		return nil, err
	}
	var resp GameResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &resp.Game, nil
}

func (c *Client) GetStadiums() (*StadiumsResponse, error) {
	data, err := c.Get("/get/stadiums")
	if err != nil {
		return nil, err
	}
	var resp StadiumsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &resp, nil
}

func (c *Client) GetStadiumByID(id string) (*Stadium, error) {
	data, err := c.Get("/get/stadium/" + id)
	if err != nil {
		return nil, err
	}
	var resp StadiumResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &resp.Stadium, nil
}

func (c *Client) Health() (*HealthResponse, error) {
	data, err := c.Get("/health")
	if err != nil {
		return nil, err
	}
	var resp HealthResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &resp, nil
}
