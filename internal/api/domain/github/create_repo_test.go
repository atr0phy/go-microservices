package github

import (
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestCreateRepoRequestsAsJson(t *testing.T) {
	request := CreateRepoRequest{
		Name:        "go-test",
		Description: "test repo",
		Homepage:    "https://github.com/",
		Private:     true,
		HasIssues:   true,
		HasProjects: true,
		HasWiki:     true,
	}

	bytes, err := json.Marshal(request)

	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var target CreateRepoRequest

	err = json.Unmarshal(bytes, &target)
	assert.Nil(t, err)

	assert.EqualValues(t, target.Name, request.Name)
	assert.EqualValues(t, target.HasIssues, request.HasIssues)
}
