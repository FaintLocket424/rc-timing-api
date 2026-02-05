package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/FaintLocket424/rc-timing-api/internal/api"
	"github.com/FaintLocket424/rc-timing-api/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestEventHandler_GetMeta(t *testing.T) {
	store := service.NewStore()
	r := api.SetupRouter(store)

	q := url.Values{}
	q.Add("url", "http://test-club.com/timing")
	q.Add("software", "fake")

	basePath := "/api/v1/event/meta"
	fullPath := fmt.Sprintf("%s?%s", basePath, q.Encode())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fullPath, nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "num_competitors")
}
