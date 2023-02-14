package registry

import (
	"encoding/json"
	"fmt"
	"time"
)

type tagsResponse struct {
	Tags []string `json:"tags"`
}

func (registry *Registry) Tags(repository string) (tags []string, err error) {
	url := registry.url("/v2/%s/tags/list", repository)

	var response tagsResponse
	for {
		registry.Logf("registry.tags url=%s repository=%s", url, repository)
		url, err = registry.getPaginatedJSON(url, &response)
		switch err {
		case ErrNoMorePages:
			tags = append(tags, response.Tags...)
			return tags, nil
		case nil:
			tags = append(tags, response.Tags...)
			continue
		default:
			return nil, err
		}
	}
}

type tagsFullResponse struct {
	Count    int       `json:"count"`
	Next     *string   `json:"next"`
	Previous *string   `json:"previous"`
	Results  []TagFull `json:"results"`
}

type TagFull struct {
	Creator int64          `json:"creator"`
	ID      int64          `json:"id"`
	Images  []TagFullImage `json:"images"`
}

type TagFullImage struct {
	Architecture string    `json:"architecture"`
	Features     string    `json:"features"`
	Variant      *string   `json:"variant"`
	Digest       string    `json:"digest"`
	OS           string    `json:"os"`
	OSFeatures   string    `json:"os_features"`
	OSVersion    *string   `json:"os_version"`
	Size         int64     `json:"size"`
	Status       string    `json:"status"`
	LastPulled   time.Time `json:"last_pulled"`
	LastPushed   time.Time `json:"last_pushed"`
}

func (registry *Registry) TagsFull(repository string) ([]TagFull, error) {
	url := registry.url("/v2/repositories/%s/tags", repository)

	resp, err := registry.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error calling repository tag endpoint: %w", err)
	}
	defer resp.Body.Close()

	var tagsResp tagsFullResponse
	if json.NewDecoder(resp.Body).Decode(&tagsResp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return tagsResp.Results, nil
}
