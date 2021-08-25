package services

import (
	"strings"

	"github.com/Yapcheekian/rest-go/src/api/config"
	"github.com/Yapcheekian/rest-go/src/api/domain/github"
	"github.com/Yapcheekian/rest-go/src/api/domain/repositories"
	"github.com/Yapcheekian/rest-go/src/api/providers/github_provider"
	"github.com/Yapcheekian/rest-go/src/api/utils/errors"
)

var (
	RepositoryService repoServiceInterface
)

type repoService struct{}

type repoServiceInterface interface {
	CreateRepo(repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
}

func init() {
	RepositoryService = &repoService{}
}

func (s *repoService) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
	input.Name = strings.TrimSpace(input.Name)

	if input.Name == "" {
		return nil, errors.NewBadRequestError("invalid repository name")
	}

	request := github.CreateRepoRequest{
		Name:        input.Name,
		Private:     false,
		Description: input.Description,
	}

	response, err := github_provider.CreateRepo(config.GetGithubAccessToken(), request)

	if err != nil {
		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}

	result := &repositories.CreateRepoResponse{
		Id:    response.Id,
		Owner: response.Owner.Login,
		Name:  response.Name,
	}

	return result, nil
}
