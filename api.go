package speedtest

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)


func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}


func FetchSettings() (Settings, error) {
	body, err := Fetch("http://www.speedtest.net/speedtest-servers.php")
	if err != nil {
		return Settings{}, err
	}
	settings := Settings{}
	err = xml.Unmarshal(body, &settings)
	return settings, err
}


func FetchConfig() (Config, error) {
	body, err := Fetch("http://www.speedtest.net/speedtest-config.php")
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	err = xml.Unmarshal(body, &config)
	return config, err
}
