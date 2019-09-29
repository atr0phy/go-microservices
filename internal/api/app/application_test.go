package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartApp(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/marco", nil)
	router.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusOK, w.Code)
	assert.EqualValues(t, "polo", w.Body.String())

}
