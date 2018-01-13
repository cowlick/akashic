package npm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type DistTags struct {
	Latest string `json:"latest"`
}

func GetDistTags(pkg string) (*DistTags, error) {

	resp, err := http.Get(fmt.Sprintf("https://registry.npmjs.org/-/package/%s/dist-tags", url.PathEscape(pkg)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tags DistTags
	err = json.Unmarshal(jsonData, &tags)
	return &tags, err
}
