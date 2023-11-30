package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"githubEvents/shared"
	"githubEvents/shared/models"
	"log"
	"net/http"
)

type GithubClient struct {
	http.Client
}

func (x GithubClient) GetEvents(url string, page, perPage int) ([]models.Event, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	values := req.URL.Query()
	values.Add("page", fmt.Sprintf("%v", page))
	values.Add("per_page", fmt.Sprintf("%v", perPage))
	req.URL.RawQuery = values.Encode()

	resp, err := x.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= shared.LowestBadHttpStatusCode {
		return nil, errors.New(fmt.Sprintf("Got HTTP error: %v", resp.Status))
	}
	var events []models.Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	return events, nil
}

func (x GithubClient) GetRepo(url string) (models.ApiRepo, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return models.ApiRepo{}, err
	}

	resp, err := x.Do(req)
	if err != nil {
		return models.ApiRepo{}, err
	}
	if resp.StatusCode >= shared.LowestBadHttpStatusCode {
		return models.ApiRepo{}, errors.New(fmt.Sprintf("Got HTTP error: %v", resp.Status))
	}
	var repo models.ApiRepo
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return models.ApiRepo{}, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	return repo, nil
}

func NewGithubClient() GithubClient {
	return GithubClient{http.Client{}}
}
