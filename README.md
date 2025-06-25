# YouTube 播放清單備份工具

此工具可備份指定 YouTube 播放清單的影片資訊，並比對兩次備份間的差異，找出被刪除的影片。

> ⚠️ **授權方式警告**
> 本專案使用 OAuth 2.0 的 OOB (Out-of-Band) 授權流程，此方式已被 Google 標記為[已淘汰](https://developers.google.com/identity/protocols/oauth2/resources/oob-migration)。
> 雖然目前仍可運作，但請注意：
> 1. Google 可能隨時停止支援此授權方式
> 2. 未來需要遷移至 PKCE 或其他更安全的授權流程
> 3. 生產環境應用應考慮使用更現代的授權方式

## 功能
- 備份播放清單中的影片標題和 ID
- 比對兩次備份差異，找出缺失影片
- 支援 OAuth 2.0 授權流程 (使用 OOB 方式)

## 使用前準備
此工具需要在本地環境運行，但必須先透過 Google Cloud Console 取得 API 憑證：

1. **建立 Google Cloud 專案**
   前往 [Google Cloud Console](https://console.cloud.google.com/) 建立新專案

2. **啟用 YouTube Data API v3**
   在 API 與服務 > 資料庫中搜尋並啟用 "YouTube Data API v3"

3. **建立 OAuth 2.0 憑證**
   在「憑證」頁面建立 OAuth 2.0 用戶端 ID，選擇應用程式類型為「桌面應用」

4. **設定** (選擇以下任一方式)：
   
   **選項一：使用 .env 檔案 (推薦)**
   複製範例檔案並編輯：
   ```bash
   cp .env.example .env
   # 編輯 .env 檔案填入您的憑證
   ```
   程式會自動載入 .env 檔案，無需手動設定環境變數

   **選項二：使用設定檔**
   編輯 [`configs/config.yaml`](configs/config.yaml) 填入您的憑證：
   ```yaml
   google:
     client_id: "您的用戶端ID"
     client_secret: "您的用戶端密鑰"
   playlist_id: "播放清單ID"
   backup_dir: "backups"  # 可選
   ```

> **注意**：環境變數優先於設定檔

## 使用方式
```bash
# 1. 建立並設定 .env 檔案
cp .env.example .env
# 編輯 .env 檔案填入您的憑證

# 2. 執行備份
go run cmd/youtube-backup/main.go --backup

# 3. 比對最近兩次備份
go run cmd/youtube-backup/main.go --compare
```

> 提示：程式會自動載入相同目錄下的 .env 檔案

## 目錄結構
```
├── cmd/                  # 主程式入口
├── configs/              # 設定檔
├── internal/
│   ├── backup/           # 備份核心邏輯
│   └── youtube/          # YouTube API 整合
├── backups/              # 備份檔案存放處 (自動建立)
└── README.md             # 使用說明
```

## 授權流程
首次執行時會啟動 OAuth 2.0 授權流程：
1. 終端機顯示授權網址
2. 登入 Google 帳號並授予權限
3. 複製授權碼貼回終端機
4. 存取權杖將保存在 `$HOME/.youtube-token.json`