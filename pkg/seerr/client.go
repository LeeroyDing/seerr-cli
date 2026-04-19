package seerr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

type SearchResult struct {
	ID          int    `json:"id"`
	MediaType   string `json:"mediaType"`
	Title       string `json:"title,omitempty"`
	Name        string `json:"name,omitempty"` // TV shows use Name
	ReleaseDate string `json:"releaseDate,omitempty"`
	FirstAirDate string `json:"firstAirDate,omitempty"`
	Summary     string `json:"overview"`
}

type SearchResponse struct {
	Page         int            `json:"page"`
	TotalPages   int            `json:"totalPages"`
	TotalResults int            `json:"totalResults"`
	Results      []SearchResult `json:"results"`
}

type RequestInfo struct {
	ID        int            `json:"id"`
	Status    int            `json:"status"` // 1 = PENDING, 2 = APPROVED, 3 = DECLINED
	Media     SearchResult   `json:"media"`
	CreatedAt string         `json:"createdAt"`
}

type PageInfo struct {
	Pages    int `json:"pages"`
	Page     int `json:"page"`
	Results  int `json:"results"`
	PageSize int `json:"pageSize"`
}

type RequestListResponse struct {
	PageInfo PageInfo      `json:"pageInfo"`
	Results  []RequestInfo `json:"results"`
}

type UserInfo struct {
	ID           int    `json:"id"`
	DisplayName  string `json:"displayName"`
	Permissions  int    `json:"permissions"`
	RequestCount int    `json:"requestCount"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MediaDetails struct {
	ID           int     `json:"id"`
	Title        string  `json:"title,omitempty"`
	Name         string  `json:"name,omitempty"`
	Overview     string  `json:"overview"`
	ReleaseDate  string  `json:"releaseDate,omitempty"`
	FirstAirDate string  `json:"firstAirDate,omitempty"`
	VoteAverage  float64 `json:"voteAverage"`
	Runtime      int     `json:"runtime,omitempty"`
	Genres       []Genre `json:"genres"`
}

type IssuePayload struct {
	IssueType int    `json:"issueType"`
	Message   string `json:"message"`
	MediaID   int    `json:"mediaId"`
}

type RequestPayload struct {
	MediaID   int    `json:"mediaId"`
	MediaType string `json:"mediaType"`
	Seasons   []int  `json:"seasons,omitempty"`
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTP:    &http.Client{},
	}
}

func (c *Client) getURL(path string) (*url.URL, error) {
	baseURL := strings.TrimSuffix(c.BaseURL, "/")
	return url.Parse(baseURL + path)
}

func (c *Client) Search(query string) (*SearchResponse, error) {
	u, err := c.getURL("/api/v1/search")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("query", query)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "seerr-cli/0.3.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	var searchResp SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResp)
	if err != nil {
		return nil, err
	}

	return &searchResp, nil
}

func (c *Client) ListRequests(take, skip int) (*RequestListResponse, error) {
	u, err := c.getURL("/api/v1/request")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("take", strconv.Itoa(take))
	q.Set("skip", strconv.Itoa(skip))
	q.Set("filter", "all")
	q.Set("sort", "added")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "seerr-cli/0.4.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	var requestResp RequestListResponse
	err = json.NewDecoder(resp.Body).Decode(&requestResp)
	if err != nil {
		return nil, err
	}

	return &requestResp, nil
}

func (c *Client) CancelRequest(requestID int) error {
	u, err := c.getURL(fmt.Sprintf("/api/v1/request/%d", requestID))
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "seerr-cli/0.4.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	return nil
}
func (c *Client) Request(mediaID int, mediaType string) error {
	u, err := c.getURL("/api/v1/request")
	if err != nil {
		return err
	}

	payload := RequestPayload{
		MediaID:   mediaID,
		MediaType: mediaType,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "seerr-cli/0.4.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	return nil
}

func (c *Client) GetMe() (*UserInfo, error) {
	u, err := c.getURL("/api/v1/auth/me")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "seerr-cli/0.5.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	var userResp UserInfo
	err = json.NewDecoder(resp.Body).Decode(&userResp)
	if err != nil {
		return nil, err
	}

	return &userResp, nil
}

func (c *Client) GetMovieDetails(movieID int) (*MediaDetails, error) {
	u, err := c.getURL(fmt.Sprintf("/api/v1/movie/%d", movieID))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "seerr-cli/0.6.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	var details MediaDetails
	err = json.NewDecoder(resp.Body).Decode(&details)
	if err != nil {
		return nil, err
	}

	return &details, nil
}

func (c *Client) GetTVDetails(tvID int) (*MediaDetails, error) {
	u, err := c.getURL(fmt.Sprintf("/api/v1/tv/%d", tvID))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "seerr-cli/0.6.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	var details MediaDetails
	err = json.NewDecoder(resp.Body).Decode(&details)
	if err != nil {
		return nil, err
	}

	return &details, nil
}

func (c *Client) GetTrending() (*SearchResponse, error) {
	u, err := c.getURL("/api/v1/discover/trending")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "seerr-cli/0.7.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	var searchResp SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResp)
	if err != nil {
		return nil, err
	}

	return &searchResp, nil
}

func (c *Client) GetPopularMovies() (*SearchResponse, error) {
	u, err := c.getURL("/api/v1/discover/movies")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "seerr-cli/0.7.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	var searchResp SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResp)
	if err != nil {
		return nil, err
	}

	return &searchResp, nil
}

func (c *Client) GetPopularTV() (*SearchResponse, error) {
	u, err := c.getURL("/api/v1/discover/tv")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "seerr-cli/0.7.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	var searchResp SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResp)
	if err != nil {
		return nil, err
	}

	return &searchResp, nil
}

func (c *Client) CreateIssue(mediaID int, issueType int, message string) error {
	u, err := c.getURL("/api/v1/issue")
	if err != nil {
		return err
	}

	payload := IssuePayload{
		MediaID:   mediaID,
		IssueType: issueType,
		Message:   message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "seerr-cli/0.8.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	return nil
}

func (c *Client) ApproveRequest(requestID int) error {
	u, err := c.getURL(fmt.Sprintf("/api/v1/request/%d/approve", requestID))
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "seerr-cli/0.9.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	return nil
}

func (c *Client) DeclineRequest(requestID int) error {
	u, err := c.getURL(fmt.Sprintf("/api/v1/request/%d/decline", requestID))
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "seerr-cli/0.9.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	return nil
}
