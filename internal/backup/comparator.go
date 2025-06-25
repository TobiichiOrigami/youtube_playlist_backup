package backup

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func CompareBackups(backupDir string) error {
	files, err := getBackupFiles(backupDir)
	if err != nil {
		return err
	}

	if len(files) < 2 {
		return fmt.Errorf("需要至少兩個備份檔案進行比對")
	}

	// 按時間排序 (檔名包含時間戳)
	sort.Strings(files)
	latest := files[len(files)-1]
	previous := files[len(files)-2]

	current, err := loadBackup(latest)
	if err != nil {
		return err
	}

	previousBackup, err := loadBackup(previous)
	if err != nil {
		return err
	}

	missing := findMissingVideos(previousBackup, current)
	if len(missing) == 0 {
		log.Println("沒有發現缺失影片")
		return nil
	}

	log.Println("以下影片已被移除:")
	for _, video := range missing {
		log.Printf("ID: %s, 標題: %s", video.ID, video.Title)
	}
	return nil
}

func getBackupFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("讀取備份目錄失敗: %w", err)
	}

	var filePaths []string
	for _, entry := range entries {
		if !entry.IsDir() {
			filePaths = append(filePaths, filepath.Join(dir, entry.Name()))
		}
	}
	return filePaths, nil
}

func loadBackup(path string) ([]Video, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("讀取備份檔案失敗: %w", err)
	}

	var videos []Video
	if err := json.Unmarshal(data, &videos); err != nil {
		return nil, fmt.Errorf("解析備份檔案失敗: %w", err)
	}
	return videos, nil
}

func findMissingVideos(previous, current []Video) []Video {
	currentMap := make(map[string]string)
	for _, v := range current {
		currentMap[v.ID] = v.Title
	}

	var missing []Video
	for _, v := range previous {
		title, exists := currentMap[v.ID]
		if !exists || title == "Deleted video" {
			missing = append(missing, v)
		}
	}
	return missing
}