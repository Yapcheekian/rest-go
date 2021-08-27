package services

import (
	"net/http"
	"sync"

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
	CreateRepo(request repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
	CreateRepos(requests []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError)
}

func init() {
	RepositoryService = &repoService{}
}

func (s *repoService) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
	if err := input.Validate(); err != nil {
		return nil, err
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

func (s *repoService) CreateRepos(requests []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError) {
	input := make(chan repositories.CreateRepositoryResult)
	output := make(chan repositories.CreateReposResponse)
	defer close(output)

	var wg sync.WaitGroup

	go s.handleRepoResults(&wg, input, output)

	for _, request := range requests {
		wg.Add(1)
		go s.createRepoConcurrent(request, input)
	}

	wg.Wait()
	close(input)

	result := <-output

	successCreations := 0
	for _, result := range result.Results {
		if result.Response != nil {
			successCreations++
		}
	}

	if successCreations == 0 {
		result.StatusCode = result.Results[0].Error.GetStatus()
	} else if successCreations == len(requests) {
		result.StatusCode = http.StatusCreated
	} else {
		result.StatusCode = http.StatusPartialContent
	}

	return result, nil
}

func (s *repoService) handleRepoResults(wg *sync.WaitGroup, input chan repositories.CreateRepositoryResult, output chan repositories.CreateReposResponse) {
	var results repositories.CreateReposResponse

	for result := range input {
		results.Results = append(results.Results, result)

		wg.Done()
	}

	output <- results
}

func (s *repoService) createRepoConcurrent(input repositories.CreateRepoRequest, output chan repositories.CreateRepositoryResult) {
	if err := input.Validate(); err != nil {
		output <- repositories.CreateRepositoryResult{
			Error: err,
		}
		return
	}

	response, err := s.CreateRepo(input)

	if err != nil {
		output <- repositories.CreateRepositoryResult{
			Error: err,
		}
		return
	}

	output <- repositories.CreateRepositoryResult{
		Response: response,
	}
}
