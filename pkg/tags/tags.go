package tags

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kevop-s/n8n-client-go/pkg/client"
)

type Tags struct {
	Client *client.Client
}

type N8nTag struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func NewTags(client *client.Client) *Tags {
	return &Tags{Client: client}
}

// GetTag retrieves a tag by its ID
func (u *Tags) GetTag(id string) (N8nTag, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tags/%s", u.Client.HostURL, id), nil)
	if err != nil {
		return N8nTag{}, err
	}
	resp, err := u.Client.GetPaginated(req)

	if err != nil {
		return N8nTag{}, err
	}

	var tag N8nTag
	err = json.Unmarshal(resp, &tag)

	if err != nil {
		return N8nTag{}, err
	}

	return tag, nil
}

// CreateTag creates a new tag
func (u *Tags) CreateTag(name string) (N8nTag, error) {
	payload := strings.NewReader(fmt.Sprintf("[{\"name\": \"%s\"}]", name))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/tags", u.Client.HostURL), payload)
	if err != nil {
		return N8nTag{}, err
	}
	resp, err := u.Client.DoRequest(req)

	if err != nil {
		return N8nTag{}, err
	}
	var tag N8nTag
	err = json.Unmarshal(resp, &tag)

	if err != nil {
		return N8nTag{}, err
	}

	return tag, nil
}

// UpdateTag updates an existing tag
func (u *Tags) UpdateTag(id, name string) (N8nTag, error) {
	payload := strings.NewReader(fmt.Sprintf("{ \"name\": \"%s\"}", name))
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/tags/%s", u.Client.HostURL, id), payload)
	if err != nil {
		return N8nTag{}, err
	}
	resp, err := u.Client.DoRequest(req)

	if err != nil {
		return N8nTag{}, err
	}
	var tag N8nTag
	err = json.Unmarshal(resp, &tag)

	if err != nil {
		return N8nTag{}, err
	}

	return tag, nil
}

// DeleteTag deletes a tag by its ID
func (u *Tags) DeleteTag(id string) (bool, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/tags/%s", u.Client.HostURL, id), nil)
	if err != nil {
		return false, err
	}
	_, err = u.Client.DoRequest(req)

	if err != nil {
		return false, err
	}

	return true, nil
}
