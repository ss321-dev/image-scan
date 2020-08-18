package dockle

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	ImageName    = "docker.io/goodwithtech/dockle:"
	GetLatestUrl = "https://api.github.com/repos/goodwithtech/dockle/releases/latest"
)

type versionCheck struct {
	TagName string `json:"tag_name"`
}

func GetLatestImageName() (string, error) {
	req, err := http.NewRequest(http.MethodGet, GetLatestUrl, nil)
	if err != nil {
		return "", err
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var versionCheck versionCheck
	err = json.Unmarshal(bytes, &versionCheck)
	if err != nil {
		return "", err
	}

	return ImageName + versionCheck.TagName, nil
}
