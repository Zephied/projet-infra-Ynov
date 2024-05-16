package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type SpotifyResponse struct {
	Tracks struct {
		Items []struct {
			PreviewURL string `json:"preview_url"`
		} `json:"items"`
	} `json:"tracks"`
}

type Data struct {
	Preview string `json:"preview_url"`
}

type Song struct {
	Name string
}

func main() {
	// S E T U P  H T T P  R O U T E  H A N D L E R S
	// O P E N A I
	http.HandleFunc("/chat", handleChat) // set router

	// S P O T I F Y
	http.HandleFunc("/spotify/token", handleSpotifyToken)     // Spotify token generation
	http.HandleFunc("/spotify/random-song", handleRandomSong) // Fetch random song from Spotify playlist
	http.HandleFunc("/spotify/getPreviewURL", handlePlayer)   // Fetch preview URL for a song

	// S T A R T  S E R V E R
	fmt.Println("Server is running on http://localhost:8080") // print message
	log.Fatal(http.ListenAndServe(":8080", nil))              // set listen port
}

// S P O T I F Y  H A N D L E R S
// generate Spotify access token
func handleSpotifyToken(w http.ResponseWriter, r *http.Request) {
	accessToken, err := generateAccessToken()
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Access Token: %s", accessToken)
}

func handlePlayer(w http.ResponseWriter, r *http.Request) {
	// Get the preview URL
	previewURL, err := getPreviewURL(songs.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the response
	resp := Data{Preview: previewURL}

	// Convert the response to JSON
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

var songs Song

// fetch a random song
func handleRandomSong(w http.ResponseWriter, r *http.Request) {

	log.Println("Received request for random song")

	accessToken, err := generateAccessToken()
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	log.Println("Access token generated, fetching random song...")
	song, err := getRandomSongFromPlaylist(accessToken)
	if err != nil {
		log.Printf("Error fetching song: %v", err)
		http.Error(w, "Failed to fetch song", http.StatusInternalServerError)
		return
	}

	log.Println("Song fetched successfully:", song)
	songs.Name = song
	fmt.Fprintf(w, "Random song: %s", song)
}

// F E T C H   P L A Y L I S T  A N D  R A N D O M  S O N G
func getFeaturedPlaylistID(accessToken string) (string, error) {
	url := "https://api.spotify.com/v1/browse/featured-playlists" // Correct endpoint for public data
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		Playlists struct {
			Items []struct {
				ID string `json:"id"`
			} `json:"items"`
		} `json:"playlists"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if len(data.Playlists.Items) == 0 {
		return "", fmt.Errorf("no playlists found")
	}

	rand.Seed(time.Now().UnixNano())
	selectedPlaylist := data.Playlists.Items[rand.Intn(len(data.Playlists.Items))]
	return selectedPlaylist.ID, nil
}

// F E T C H  R A N D O M  S O N G  F R O M  P L A Y L I S T
func getRandomSongFromPlaylist(accessToken string) (string, error) {
	// Get a random featured playlist ID
	playlistID, err := getFeaturedPlaylistID(accessToken)
	if err != nil {
		return "", err
	}

	// Fetch tracks from the selected playlist
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		Items []struct {
			Track struct {
				Name string `json:"name"`
			} `json:"track"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if len(data.Items) == 0 {
		return "", fmt.Errorf("no songs found in the playlist")
	}

	selectedSong := data.Items[rand.Intn(len(data.Items))].Track.Name
	return selectedSong, nil
}

// C H A T G P T  R E S P O N S E  H A N D L E R
func handleChat(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY") // Load API key from environment variable
	if apiKey == "" {
		log.Fatal("API key not found in the environment variables")
	}
	client := openai.NewClient(apiKey) // Create a new client

	// Decode incoming request
	var req struct {
		Prompt string `json:"prompt"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send request to OpenAI
	completion, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo", // or "gpt-4"
		Messages: []openai.ChatCompletionMessage{ // Correct type for messages
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: req.Prompt},
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response back to Streamlit
	response := struct {
		Answer string `json:"answer"`
	}{
		Answer: completion.Choices[0].Message.Content, // Adjust this depending on the actual response structure
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

// S P O T I F Y
// A P I  T O K E N  G E N E R A T I O N
func generateAccessToken() (string, error) {
	// environment variable for the client credentials
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	authHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+authHeader)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return "", err
	}
	return tokenResp.AccessToken, nil
}
func getPreviewURL(song string) (string, error) {
	accessToken, err := generateAccessToken()
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return "", err
	}
	url := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track&limit=1", url.QueryEscape(song))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", err
	}

	var data SpotifyResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("Error unmarshalling response body: %v", err)
		return "", err
	}

	if len(data.Tracks.Items) > 0 {
		return data.Tracks.Items[0].PreviewURL, nil
	}

	return "", fmt.Errorf("no preview URL found for song: %s", song)
}
