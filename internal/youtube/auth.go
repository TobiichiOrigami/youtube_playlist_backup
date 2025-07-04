package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

const tokenFile = ".youtube-token.json"

func getToken(config *oauth2.Config) (*oauth2.Token, error) {
	usr, _ := user.Current()
	tokenPath := filepath.Join(usr.HomeDir, tokenFile)
	
	if token, err := readToken(tokenPath); err == nil {
		return token, nil
	}

	// 使用 OOB (Out-of-Band) 流程
	config.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("請前往以下網址授權: \n%v\n", authURL)
	fmt.Print("輸入授權碼: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("讀取代碼失敗: %v", err)
	}

	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, fmt.Errorf("取得token失敗: %v", err)
	}
	saveToken(tokenPath, token)
	return token, nil
}

func readToken(file string) (*oauth2.Token, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	
	tok := &oauth2.Token{}
	if err := json.Unmarshal(data, tok); err != nil {
		return nil, err
	}
	return tok, nil
}

func saveToken(file string, token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0600)
}

func NewClient(clientID, clientSecret string) (*youtube.Service, error) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{youtube.YoutubeReadonlyScope},
	}

	token, err := getToken(config)
	if err != nil {
		return nil, err
	}

	client := config.Client(context.Background(), token)
	service, err := youtube.New(client)
	if err != nil {
		return nil, fmt.Errorf("建立YouTube服務失敗: %v", err)
	}

	return service, nil
}