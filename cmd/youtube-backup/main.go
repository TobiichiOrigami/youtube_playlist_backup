package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	
	"github.com/joho/godotenv"
	"github.com/kevin/golang/youtube_playlist_backup/internal/backup"
	"github.com/kevin/golang/youtube_playlist_backup/internal/youtube"
	"gopkg.in/yaml.v3"
)

// Config 結構對應 config.yaml 設定檔
type Config struct {
	Google struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
	} `yaml:"google"`
	PlaylistID string `yaml:"playlist_id"`
	BackupDir  string `yaml:"backup_dir"`
}

func loadConfig() (*Config, error) {
	// 自動載入 .env 檔案
	_ = godotenv.Load()

	// 從環境變數讀取
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	playlistID := os.Getenv("YOUTUBE_PLAYLIST_ID")
	backupDir := os.Getenv("BACKUP_DIR")

	cfg := &Config{}
	cfg.Google.ClientID = clientID
	cfg.Google.ClientSecret = clientSecret
	cfg.PlaylistID = playlistID
	if backupDir != "" {
		cfg.BackupDir = backupDir
	} else {
		cfg.BackupDir = "backups"
	}

	// 如果必要變數缺失，嘗試從設定檔讀取
	if clientID == "" || clientSecret == "" || playlistID == "" {
		cfgPath := filepath.Join("configs", "config.yaml")
		if _, err := os.Stat(cfgPath); err == nil {
			data, err := os.ReadFile(cfgPath)
			if err != nil {
				return nil, fmt.Errorf("讀取設定檔失敗: %w", err)
			}

			var fileCfg Config
			if err := yaml.Unmarshal(data, &fileCfg); err != nil {
				return nil, fmt.Errorf("解析設定檔失敗: %w", err)
			}
			
			// 只填充缺失的值
			if clientID == "" {
				cfg.Google.ClientID = fileCfg.Google.ClientID
			}
			if clientSecret == "" {
				cfg.Google.ClientSecret = fileCfg.Google.ClientSecret
			}
			if playlistID == "" {
				cfg.PlaylistID = fileCfg.PlaylistID
			}
			if backupDir == "" && fileCfg.BackupDir != "" {
				cfg.BackupDir = fileCfg.BackupDir
			}
		}
	}
	
	// 再次檢查必要變數
	if cfg.Google.ClientID == "" || cfg.Google.ClientSecret == "" || cfg.PlaylistID == "" {
		return nil, fmt.Errorf("缺少必要設定: GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, YOUTUBE_PLAYLIST_ID")
	}
	
	return cfg, nil
}

func main() {
	backupCmd := flag.Bool("backup", false, "執行備份")
	compareCmd := flag.Bool("compare", false, "比對備份")
	flag.Parse()

	// 載入設定
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("載入設定失敗: %v", err)
	}

	client, err := youtube.NewClient(cfg.Google.ClientID, cfg.Google.ClientSecret)
	if err != nil {
		log.Fatalf("建立YouTube客戶端失敗: %v", err)
	}

	switch {
	case *backupCmd:
		if err := backup.RunBackup(client, cfg.PlaylistID, cfg.BackupDir); err != nil {
			log.Fatalf("備份失敗: %v", err)
		}
	case *compareCmd:
		if err := backup.CompareBackups(cfg.BackupDir); err != nil {
			log.Fatalf("比對失敗: %v", err)
		}
	default:
		fmt.Println("請指定指令: --backup 或 --compare")
	}
}