package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

var (
	videoData []byte
	mutex     sync.Mutex
)

// Upload video (POST)
func uploadVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, _, err := r.FormFile("video")
		if err != nil {
			http.Error(w, "Failed to read video", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to process video", http.StatusInternalServerError)
			return
		}

		mutex.Lock()
		videoData = data
		mutex.Unlock()
	}
}

// Fetch video (GET)
func getVideo(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	if len(videoData) == 0 {
		mutex.Unlock()
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Write(videoData)
	mutex.Unlock()
}

func main() {
	http.HandleFunc("/vshare1", uploadVideo) // Upload video data
	http.HandleFunc("/vshare2", getVideo)    // Fetch video data

	fmt.Println("Server running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
