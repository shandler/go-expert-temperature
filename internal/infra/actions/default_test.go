package actions_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shandler/go-expert-temperature/internal/domain"
	"github.com/shandler/go-expert-temperature/internal/infra/actions"
)

func TestNewDefault(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		a        domain.Action
		want     string
		response string
	}{
		{
			name:     "success",
			a:        actions.NewDefault("/"),
			want:     "/",
			response: "access: /zip-code?zipCode=07987110",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Path(); got != tt.want {
				t.Errorf("defaultAction.Path() = %v, want %v", got, tt.want)
			}

			req := httptest.NewRequest(http.MethodGet, tt.a.Path(), nil)
			w := httptest.NewRecorder()

			tt.a.Handle(w, req)

			if w.Body.String() != tt.response {
				t.Errorf("defaultAction.Handle() = %v, want %v", w.Body.String(), tt.response)
			}

		})
	}

}
