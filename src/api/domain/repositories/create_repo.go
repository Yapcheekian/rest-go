package repositories

import (
	"strings"

	"github.com/Yapcheekian/rest-go/src/api/utils/errors"
)

type CreateRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *CreateRepoRequest) Validate() errors.ApiError {
	r.Name = strings.TrimSpace(r.Name)

	if r.Name == "" {
		return errors.NewBadRequestError("invalid repository name")
	}

	return nil
}

type CreateRepoResponse struct {
	Id    int64  `json:"id"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

type CreateReposResponse struct {
	StatusCode int                      `json:"status"`
	Results    []CreateRepositoryResult `json"results"`
}

type CreateRepositoryResult struct {
	Response *CreateRepoResponse `json:"repo"`
	Error    errors.ApiError     `json:"error"`
}
