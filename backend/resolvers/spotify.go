package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type SpotifyResolver struct {
	client       *http.Client
	clientID     string
	clientSecret string
	token        string
	tokenExpiry  time.Time
	mu           sync.Mutex
}

func NewSpotifyResolver(clientID, clientSecret string) *SpotifyResolver {
	if clientID == "" || clientSecret == "" {
		return nil
	}
	return &SpotifyResolver{
		client:       SafeHttpClient(5 * time.Second), // Slightly longer for potential auth
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (r *SpotifyResolver) Name() string {
	return "spotify"
}

func (r *SpotifyResolver) CanHandle(u *url.URL) bool {
	return strings.ToLower(u.Host) == "open.spotify.com"
}

func (r *SpotifyResolver) getToken(ctx context.Context) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Add 5 min buffer to expiry
	if r.token != "" && time.Now().Add(5*time.Minute).Before(r.tokenExpiry) {
		return r.token, nil
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, "POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(r.clientID, r.clientSecret)

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get spotify token, status: %d", resp.StatusCode)
	}

	var res struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"` // seconds
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	r.token = res.AccessToken
	r.tokenExpiry = time.Now().Add(time.Duration(res.ExpiresIn) * time.Second)

	return r.token, nil
}

func (r *SpotifyResolver) Resolve(ctx context.Context, u *url.URL) (*Result, error) {
	// Path is like /track/ID or /album/ID
	parts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid spotify path")
	}

	entityType := parts[0]
	entityID := parts[1]

	token, err := r.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("spotify auth error: %v", err)
	}

	apiURL := fmt.Sprintf("https://api.spotify.com/v1/%ss/%s", entityType, entityID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify api error, status: %d", resp.StatusCode)
	}

	var data struct {
		Name    string `json:"name"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Owner struct { // for playlists
			DisplayName string `json:"display_name"`
		} `json:"owner"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	title := data.Name
	if len(data.Artists) > 0 {
		var artists []string
		for _, a := range data.Artists {
			artists = append(artists, a.Name)
		}
		title = fmt.Sprintf("%s - %s", strings.Join(artists, ", "), title)
	} else if data.Owner.DisplayName != "" {
		title = fmt.Sprintf("%s - Playlist by %s", title, data.Owner.DisplayName)
	}

	return &Result{
		Title:    title,
		Platform: "Spotify",
	}, nil
}
