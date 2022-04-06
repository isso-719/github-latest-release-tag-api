package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func GetTag(owner, repo string) (string, error) {
	// GitHub API からリポジトリのデータを取得
	url := "https://api.github.com/repos/" + owner + "/" + repo + "/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// レスポンスを読み込む
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// レスポンスを json 化
	var data GitHubResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	// レスポンスからリリースのタグを取得
	tag := data.TagName

	return tag, nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// API レスポンスようの構造体を定義
	var response Response

	// パラメータから owner/repo を取得
	var owner, repo string
	params := r.URL.Query()
	owner = params.Get("owner")
	repo = params.Get("repo")
	if owner == "" || repo == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "Owner Or Repo Is Empty"
	} else {
		// リリースのタグを取得
		tag, err := GetTag(owner, repo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.Message = "Internal Server Error"
		} else if tag == "" {
			response.Message = "Tag Not Found"
		} else {
			response.Tag = tag
		}
	}

	// レスポンスを json 化
	json, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}

type Response struct {
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

type GitHubResponse struct {
	URL       string `json:"url"`
	AssetsURL string `json:"assets_url"`
	UploadURL string `json:"upload_url"`
	HTMLURL   string `json:"html_url"`
	ID        int    `json:"id"`
	Author    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []struct {
		URL      string      `json:"url"`
		ID       int         `json:"id"`
		NodeID   string      `json:"node_id"`
		Name     string      `json:"name"`
		Label    interface{} `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
}
