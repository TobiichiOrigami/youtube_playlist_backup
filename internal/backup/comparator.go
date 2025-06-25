package backup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("讀取備份目錄失敗: %w", err)
	}

	var filePaths []string
	for _, file := range files {
		if !file.IsDir() {
			filePaths = append(filePaths, filepath.Join(dir, file.Name()))
		}
	}
	return filePaths, nil
}

func loadBackup(path string) ([]Video, error) {
	data, err := ioutil.ReadFile(path)
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
	currentMap := make(map[string]bool)
	for _, v := range current {
		currentMap[v.ID] = true
	}

	var missing []Video
	for _, v := range previous {
		if !currentMap[v.ID] {
			missing = append(missing, v)
		}
	}
	return missing
}