package actions_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shandler/go-expert-temperature/internal/core/search"
	"github.com/shandler/go-expert-temperature/internal/domain"
	"github.com/shandler/go-expert-temperature/internal/dto"
	"github.com/shandler/go-expert-temperature/internal/infra/actions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ZipCodeMock struct {
	mock.Mock
}

func (s *ZipCodeMock) Search(request dto.SearchRequest) dto.SearchResponse {
	args := s.Called(request)
	return args.Get(0).(dto.SearchResponse)
}

var zipCodeMock *ZipCodeMock = &ZipCodeMock{}

func TestZipCodeNew(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		a        domain.Action
		want     string
		response string
		mockRun  func()
	}{
		{
			name:     "success",
			a:        actions.NewZipCode("/zip-code", zipCodeMock),
			want:     "/zip-code",
			response: `{"temp_c":0,"temp_f":32,"temp_k":50}`,
			mockRun: func() {
				zipCodeMock.On("Search", mock.Anything).Return(dto.SearchResponse{
					Status: 200,
					Body: search.WeatherResponse{
						TempC: 0,
						TempF: 32,
						TempK: 50,
					},
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Path(); got != tt.want {
				t.Errorf("zipCodeAction.Path() = %v, want %v", got, tt.want)
			}

			req := httptest.NewRequest("GET", tt.a.Path(), nil)
			w := httptest.NewRecorder()

			tt.mockRun()
			tt.a.Handle(w, req)

			assert.Equal(t, tt.response, strings.TrimSpace(w.Body.String()))

		})
	}
}
