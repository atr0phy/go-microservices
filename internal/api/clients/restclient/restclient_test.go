package restclient

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartMockups(t *testing.T) {
	assert.False(t, enabledMocks)
	StartMockups()
	assert.True(t, enabledMocks)
}

func TestStopMockups(t *testing.T) {
	enabledMocks = true
	StopMockups()
	assert.False(t, enabledMocks)
}

func TestFlushMockups(t *testing.T) {
	FlushMockups()

	assert.EqualValues(t, make(map[string]*Mock), mocks)
}

func TestAddMockup(t *testing.T) {
	FlushMockups()

	AddMockup(Mock{
		HttpMethod: http.MethodGet,
		Url:        "http://localhost/mocktest",
	})

	assert.EqualValues(t,
		&Mock{HttpMethod: http.MethodGet, Url: "http://localhost/mocktest"},
		mocks[fmt.Sprintf("%s_%s", http.MethodGet, "http://localhost/mocktest")])
}

func TestPostMockUrlError(t *testing.T) {
	FlushMockups()
	StartMockups()

	response, err := Post("http://localhost/mocktest", nil, nil)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "no mockup found for give request", err.Error())
}

func TestPostMockUrlFound(t *testing.T) {
	FlushMockups()
	StartMockups()

	AddMockup(Mock{
		HttpMethod: http.MethodPost,
		Url:        "http://localhost/mocktest",
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`mock response`)),
		},
	})
	response, err := Post("http://localhost/mocktest", nil, nil)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.EqualValues(t, http.StatusOK, response.StatusCode)

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	assert.Nil(t, err)
	assert.NotNil(t, body)
	assert.EqualValues(t, "mock response", string(body))
}

func TestPostInvalidJson(t *testing.T) {
	FlushMockups()
	StopMockups()

	response, err := Post("http://localhost/mocktest", math.NaN(), nil)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "json: unsupported value: NaN", err.Error())
}

func TestPostInvalidRequest(t *testing.T) {
	FlushMockups()
	StopMockups()

	response, err := Post("\t", nil, nil)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "parse \t: net/url: invalid control character in URL", err.Error())
}

func TestPostInvalidUrl(t *testing.T) {
	FlushMockups()
	StopMockups()

	response, err := Post("localhost/mocktest", nil, nil)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Post localhost/mocktest: unsupported protocol scheme \"\"", err.Error())
}

func TestPostNoError(t *testing.T) {
	mux := http.NewServeMux()
	mockServer := httptest.NewServer(mux)
	mux.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("success"))
			return
		}
		writer.WriteHeader(http.StatusNotImplemented)
	})
	response, err := Post(mockServer.URL, nil, nil)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.EqualValues(t, http.StatusOK, response.StatusCode)

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	assert.Nil(t, err)
	assert.NotNil(t, body)
	assert.EqualValues(t, "success", string(body))
}
