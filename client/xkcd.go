package client

import (
	"encoding/json"
	"fmt"
	"github.com/redjoker011/go-grab-xkcd/model"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

type ComicNumber int

const (
	// BaseURL of xkcd
	BaseURL string = "https://xkcd.com"
	// DefaultClientTimeout is time to wait before cancelling the request
	DefaultClientTimeout time.Duration = 30 * time.Second
	// LatestComic is the latest comic number according to the xkcd API
	LatestComic ComicNumber = 0
)

type XKCDClient struct {
	client  *http.Client
	baseURL string
}

func NewXKCDClient() *XKCDClient {
	return &XKCDClient{
		client: &http.Client{
			Timeout: DefaultClientTimeout,
		},
		baseURL: BaseURL,
	}
}

func (hc *XKCDClient) SetTimeout(d time.Duration) {
	hc.client.Timeout = d
}

func (hc *XKCDClient) Fetch(n ComicNumber, save bool) (model.Comic, error) {
	resp, err := hc.client.Get(hc.buildURL(n))
	if err != nil {
		return model.Comic{}, err
	}
	defer resp.Body.Close()

	var comicResp model.ComicResponse
	if err := json.NewDecoder(resp.Body).Decode(&comicResp); err != nil {
		return model.Comic{}, err
	}

	if save {
		if err := hc.SaveToDisk(comicResp.Img, "."); err != nil {
			fmt.Println("Failed to save image")
		}
	}
	return comicResp.Comic(), nil
}

func (hc *XKCDClient) SaveToDisk(url, savePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	absSavePath, _ := filepath.Abs(savePath)

	filePath := fmt.Sprintf("%s/%s", absSavePath, path.Base(url))
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return err
	}
	return nil
}

func (hc *XKCDClient) buildURL(n ComicNumber) string {
	var finalURL string
	if n == LatestComic {
		finalURL = fmt.Sprintf("%s/info.0.json", hc.baseURL)
	} else {
		finalURL = fmt.Sprintf("%s/%d/info.0.json", hc.baseURL, n)
	}
	return finalURL
}
