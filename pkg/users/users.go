package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kevop-s/n8n-client-go/pkg/client"
)

type Users struct {
	Client *client.Client
}

type N8nUser struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	IsPending bool   `json:"isPending"`
	Role      string `json:"role"`
}

func NewUsers(client *client.Client) *Users {
	return &Users{Client: client}
}

func (u *Users) GetUser(id string) (N8nUser, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/%s", u.Client.HostURL, id), nil)
	if err != nil {
		return N8nUser{}, err
	}
	resp, err := u.Client.GetPaginated(req)

	if err != nil {
		return N8nUser{}, err
	}

	var user N8nUser
	err = json.Unmarshal(resp, &user)

	if err != nil {
		return N8nUser{}, err
	}

	return user, nil
}

func (u *Users) CreateUser(email, role string) (N8nUser, error) {
	payload := strings.NewReader(fmt.Sprintf("[{\"email\": \"%s\", \"role\": \"%s\"}]", email, role))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users", u.Client.HostURL), payload)
	if err != nil {
		return N8nUser{}, err
	}
	_, err = u.Client.DoRequest(req)

	if err != nil {
		return N8nUser{}, err
	}

	return u.GetUser(email)
}

func (u *Users) UpdateUser(email, role string) (N8nUser, error) {
	payload := strings.NewReader(fmt.Sprintf("{ \"newRoleName\": \"%s\"}", role))
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/users/%s/role", u.Client.HostURL, email), payload)
	if err != nil {
		return N8nUser{}, err
	}
	_, err = u.Client.DoRequest(req)

	if err != nil {
		return N8nUser{}, err
	}

	return u.GetUser(email)
}

func (u *Users) DeleteUser(email string) (bool, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%s", u.Client.HostURL, email), nil)
	if err != nil {
		return false, err
	}
	_, err = u.Client.DoRequest(req)

	if err != nil {
		return false, err
	}

	return true, nil
}
