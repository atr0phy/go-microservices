package services

import (
	"net/http"
	"sync"

	"github.com/atr0phy/go-microservices/internal/api/config"
	"github.com/atr0phy/go-microservices/internal/api/domain/github"
	"github.com/atr0phy/go-microservices/internal/api/domain/repositories"
	"github.com/atr0phy/go-microservices/internal/api/log"
	"github.com/atr0phy/go-microservices/internal/api/providers/github_provider"
	"github.com/atr0phy/go-microservices/internal/api/utils/errors"
)

type reposService struct{}

type reposServiceInterface interface {
	CreateRepo(request repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
	CreateRepos(request []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError)
}

var (
	RepositoryService reposServiceInterface
)

func init() {
	RepositoryService = &reposService{}
}

func (s *reposService) CreateRepo(request repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	req := github.CreateRepoRequest{
		Name:        request.Name,
		Description: request.Description,
		Private:     false,
	}
	log.Info("about to send request to external api",
		log.Field("status", "pending"))

	response, err := github_provider.CreateRepo(config.GetGithubAccessToken(), req)
	if err != nil {
		log.Error("about to send request to external api",
			err,
			log.Field("status", "error"))
		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}
	log.Info("about to send request to external api",
		log.Field("status", "success"))

	result := repositories.CreateRepoResponse{
		Id:    response.Id,
		Name:  response.Name,
		Owner: response.Owner.Login,
	}

	return &result, nil

}

func (s *reposService) CreateRepos(request []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError) {
	input := make(chan repositories.CreateRepositoriesResult)
	output := make(chan repositories.CreateReposResponse)
	defer close(output)

	var wg sync.WaitGroup

	go s.handleRepoResults(&wg, input, output)

	for _, current := range request {
		wg.Add(1)
		go s.createRepoConcurrent(current, input)
	}

	wg.Wait()
	close(input)

	result := <-output

	successCreations := 0
	for _, current := range result.Results {
		if current.Response != nil {
			successCreations++
		}
	}

	if successCreations == 0 {
		result.StatusCode = result.Results[0].Error.Status()
	} else if successCreations == len(request) {
		result.StatusCode = http.StatusCreated
	} else {
		result.StatusCode = http.StatusPartialContent
	}

	return result, nil
}

func (s *reposService) handleRepoResults(wg *sync.WaitGroup, input chan repositories.CreateRepositoriesResult, output chan repositories.CreateReposResponse) {
	var results repositories.CreateReposResponse

	for incomingEvent := range input {
		repoResult := repositories.CreateRepositoriesResult{
			Response: incomingEvent.Response,
			Error:    incomingEvent.Error,
		}
		results.Results = append(results.Results, repoResult)
		wg.Done()
	}

	output <- results

}

func (s *reposService) createRepoConcurrent(request repositories.CreateRepoRequest, output chan repositories.CreateRepositoriesResult) {
	if err := request.Validate(); err != nil {
		output <- repositories.CreateRepositoriesResult{Error: err}
		return
	}

	result, err := s.CreateRepo(request)
	if err != nil {
		output <- repositories.CreateRepositoriesResult{Error: err}
		return
	}

	output <- repositories.CreateRepositoriesResult{Response: result}
}
