package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type PlaylistItemsResponse struct {
	NextPageToken string `json:"nextPageToken"`
	Items         []struct {
		ContentDetails struct {
			VideoID string `json:"videoId"`
		} `json:"contentDetails"`
	} `json:"items"`
}

type VideoDurationResponse struct {
	Items []struct {
		ContentDetails struct {
			Duration string `json:"duration"`
		} `json:"contentDetails"`
	} `json:"items"`
}

func getPlaylistDetails(playlistID, apiKey string) (int, time.Duration, error) {
	var totalVideos int
	var totalDuration time.Duration

	nextPageToken := ""
	for {
		url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=contentDetails&maxResults=50&pageToken=%s&playlistId=%s&key=%s", nextPageToken, playlistID, apiKey)

		resp, err := http.Get(url)
		if err != nil {
			return 0, 0, err
		}
		defer resp.Body.Close()

		var playlistResp PlaylistItemsResponse
		err = json.NewDecoder(resp.Body).Decode(&playlistResp)
		if err != nil {
			return 0, 0, err
		}

		totalVideos += len(playlistResp.Items)

		var videoIDs []string
		for _, item := range playlistResp.Items {
			videoIDs = append(videoIDs, item.ContentDetails.VideoID)
		}

		videoDuration, err := getVideosDuration(videoIDs, apiKey)
		if err != nil {
			return 0, 0, err
		}
		totalDuration += videoDuration

		if playlistResp.NextPageToken == "" {
			break
		}
		nextPageToken = playlistResp.NextPageToken
	}

	return totalVideos, totalDuration, nil
}

func getVideosDuration(videoIDs []string, apiKey string) (time.Duration, error) {
	var totalDuration time.Duration

	for _, videoID := range videoIDs {
		duration, err := getVideoDuration(videoID, apiKey)
		if err != nil {
			return 0, err
		}
		totalDuration += duration
	}

	return totalDuration, nil
}

func getVideoDuration(videoID, apiKey string) (time.Duration, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=contentDetails&id=%s&key=%s", videoID, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var videoDurationResp VideoDurationResponse
	err = json.NewDecoder(resp.Body).Decode(&videoDurationResp)
	if err != nil {
		return 0, err
	}

	if len(videoDurationResp.Items) == 0 {
		return 0, fmt.Errorf("video not found")
	}

	durationStr := videoDurationResp.Items[0].ContentDetails.Duration
	duration, err := parseDuration(durationStr)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

func parseDuration(durationStr string) (time.Duration, error) {
	durationStr = strings.Replace(durationStr, "PT", "", 1)

	hoursIndex := strings.Index(durationStr, "H")
	minutesIndex := strings.Index(durationStr, "M")
	secondsIndex := strings.Index(durationStr, "S")

	var hours, minutes, seconds int

	if hoursIndex != -1 {
		hoursVal, err := strconv.Atoi(durationStr[:hoursIndex])
		if err != nil {
			return 0, err
		}
		hours = hoursVal
	}

	if minutesIndex != -1 {
		minutesVal, err := strconv.Atoi(durationStr[hoursIndex+1 : minutesIndex])
		if err != nil {
			return 0, err
		}
		minutes = minutesVal
	} else if hoursIndex != -1 {
		minutes = 0
	}

	if secondsIndex != -1 {
		var secondsVal int
		if minutesIndex != -1 {
			secondsVal, _ = strconv.Atoi(durationStr[minutesIndex+1 : secondsIndex])
		} else if hoursIndex != -1 {
			secondsVal, _ = strconv.Atoi(durationStr[hoursIndex+1 : secondsIndex])
		} else {
			secondsVal, _ = strconv.Atoi(durationStr[:secondsIndex])
		}
		seconds = secondsVal
	} else {
		seconds = 0
	}

	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
}

func extractPlaylistID(url string) (string, error) {
	re := regexp.MustCompile(`list=([a-zA-Z0-9_-]+)`)
	match := re.FindStringSubmatch(url)
	if len(match) < 2 {
		return "", fmt.Errorf("invalid playlist URL")
	}
	return match[1], nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set the YOUTUBE_API_KEY environment variable with your YouTube Data API key.")
		return
	}

	var playlistURL string
	fmt.Print("Enter the YouTube playlist URL: ")
	fmt.Scanln(&playlistURL)

	playlistID, err := extractPlaylistID(playlistURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	numVideos, totalDuration, err := getPlaylistDetails(playlistID, apiKey)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	averageLength := totalDuration / time.Duration(numVideos)

	totalDurationDays := int(totalDuration.Hours() / 24)
	totalDurationHours := int(totalDuration.Hours()) % 24
	totalDurationMinutes := int(totalDuration.Minutes()) % 60
	totalDurationSeconds := int(totalDuration.Seconds()) % 60

	fmt.Println("Playlist Information:")
	fmt.Printf("No of videos: %d\n", numVideos)
	fmt.Printf("Average length of video: %s\n", formatDuration(averageLength))
	fmt.Printf("Total length of playlist: %d days, %d hours, %d minutes, %d seconds\n",
		totalDurationDays, totalDurationHours, totalDurationMinutes, totalDurationSeconds)

	playbackSpeeds := []float64{1.25, 1.5, 1.75, 2.0}
	fmt.Println("Duration at Different Playback Speeds:")
	for _, speed := range playbackSpeeds {
		adjustedDuration := totalDuration / time.Duration(speed)
		days := int(adjustedDuration.Hours() / 24)
		adjustedDuration -= time.Duration(days*24) * time.Hour
		hours := int(adjustedDuration.Hours())
		adjustedDuration -= time.Duration(hours) * time.Hour
		minutes := int(adjustedDuration.Minutes())
		adjustedDuration -= time.Duration(minutes) * time.Minute
		seconds := int(adjustedDuration.Seconds())
		fmt.Printf("At %.2fx: %d days, %d hours, %d minutes, %d seconds\n", speed, days, hours, minutes, seconds)
	}
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d hours, %02d mins, %02d secs", h, m, s)
}
