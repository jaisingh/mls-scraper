package mls

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

const (
	CONFIG_FILE = "config.json"
)

type ConfigData struct {
	User      string `json:"user"`
	Pass      string `json:"pass"`
	BaseUrl   string `json:"baseUrl"`
	FileMode  bool   `json:"fileMode"`
	DumpPages bool   `json:"dumpPages"`
}

var Config ConfigData
var Client *http.Client

func Init() error {
	// Load config file
	Config = ConfigData{}
	if err := Config.Load(CONFIG_FILE); err != nil {
		return err
	}

	// Connect to mls server

	if !Config.FileMode {
		cookies, _ := cookiejar.New(nil)
		Client = &http.Client{
			Jar: cookies,
		}
		_, err := webGet(Config.AuthURL())
		if err != nil {
			return err
		}

	}

	return nil
}

func webGet(u string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.10 Safari/537.36")

	return Client.Do(req)
}

func (c *ConfigData) Load(fileName string) error {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.New("Could not load ConfigData file")
	}

	err = json.Unmarshal(file, &c)
	if err != nil {
		return errors.New("Couldn't Unsmarhal ConfigData file")
	}

	return nil
}

func (c ConfigData) AuthURL() string {
	return c.BaseUrl + "?cid=" + c.User + "&pass=" + c.Pass
}
