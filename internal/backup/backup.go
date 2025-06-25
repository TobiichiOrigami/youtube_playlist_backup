package backup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/api/youtube/v3"
)

type Video struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func RunBackup(service *youtube.Service, playlistID, backupDir string) error {
	playlistItems, err := fetchPlaylistItems(service, playlistID)
	if err != nil {
		return err
	}

	videos := make([]Video, 0, len(playlistItems))
	for _, item := range playlistItems {
		videos = append(videos, Video{
			ID:    item.ContentDetails.VideoId,
			Title: item.Snippet.Title,
		})
	}

	if err := saveBackup(backupDir, videos); err != nil {
		return err
	}

	log.Printf("成功備份 %d 部影片", len(videos))
	return nil
}

func fetchPlaylistItems(service *youtube.Service, playlistID string) ([]*youtube.PlaylistItem, error) {
	var allItems []*youtube.PlaylistItem
	pageToken := ""

	for {
		call := service.PlaylistItems.List([]string{"snippet,contentDetails"}).
			PlaylistId(playlistID).
			MaxResults(50).
			PageToken(pageToken)

		response, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("取得播放清單失敗: %w", err)
		}

		allItems = append(allItems, response.Items...)
		pageToken = response.NextPageToken
		if pageToken == "" {
			break
		}
	}

	return allItems, nil
}

func saveBackup(backupDir string, videos []Video) error {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("建立備份目錄失敗: %w", err)
	}

	filename := fmt.Sprintf("backup_%s.json", time.Now().Format("20060102_150405"))
	path := filepath.Join(backupDir, filename)

	file, err := json.MarshalIndent(videos, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化備份資料失敗: %w", err)
	}

	if err := ioutil.WriteFile(path, file, 0644); err != nil {
		return fmt.Errorf("寫入備份檔案失敗: %w", err)
	}
	return nil
}