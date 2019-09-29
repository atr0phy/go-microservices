package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiError_Status(t *testing.T) {
	err := apiError{
		AStatus: 123,
	}

	assert.NotNil(t, err)
	assert.EqualValues(t, 123, err.Status())
}

func TestApiError_Error(t *testing.T) {
	err := apiError{
		AnError: "error",
	}

	assert.NotNil(t, err)
	assert.EqualValues(t, "error", err.Error())

}

func TestApiError_Message(t *testing.T) {
	err := apiError{
		AMessage: "error message",
	}

	assert.NotNil(t, err)
	assert.EqualValues(t, "error message", err.Message())
}

func TestNewApiError(t *testing.T) {
	err := NewApiError(123, "error message")

	assert.NotNil(t, err)
	assert.EqualValues(t, 123, err.Status())
	assert.EqualValues(t, "error message", err.Message())
}

func TestNewApiErrFromBytesWithInvalidJson(t *testing.T) {
	result, err := NewApiErrFromBytes([]byte("not json"))

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid json for creating an api error", err.Error())
}

func TestNewApiErrFromBytesNoError(t *testing.T) {
	result, err := NewApiErrFromBytes([]byte(`{"status": 123,"message": "valid json"}`))

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, 123, result.Status())
	assert.EqualValues(t, "valid json", result.Message())
}

func TestNewInternalServerError(t *testing.T) {
	err := NewInternalServerError("internal server error")

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "internal server error", err.Message())
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("not found")

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
	assert.EqualValues(t, "not found", err.Message())
}

func TestNewBadRequestError(t *testing.T) {
	err := NewBadRequestError("bad request")

	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "bad request", err.Message())
}
