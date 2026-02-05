package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FaintLocket424/rc-timing-api/internal/api"
	"github.com/FaintLocket424/rc-timing-api/internal/service"
	"github.com/FaintLocket424/rc-timing-api/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestEventHandler_GetMeta(t *testing.T) {
	tests := []struct {
		Name           string
		Request        *http.Request
		ExpectedStatus int
	}{
		{
			"name",
			testutil.GenerateValidRequest("/api/v1/event/meta", "https://example.com", "fake"),
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			store := service.NewStore()
			r := api.SetupRouter(store)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, tt.Request)

			assert.Equal(t, tt.ExpectedStatus, w.Code)
		})
	}
}
