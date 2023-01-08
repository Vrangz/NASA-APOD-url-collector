package nasa

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"url-collector/internal/collector"
	"url-collector/internal/timeutils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type successNasaMock struct{}

func (successNasaMock) Collect(ctx context.Context, from, to time.Time) (set collector.MetadataSet, err error) {
	for i := range timeutils.List(from, to) {
		m := collector.Metadata{
			Title: fmt.Sprintf("%d", i),
			URL:   fmt.Sprintf("https://mock.nasa.com/picture/%d", i),
			HDURL: fmt.Sprintf("https://mock.nasa.com/picture-hd/%d", i),
		}
		set = append(set, m)
	}
	return

}

type errorNasaMock struct{}

func (errorNasaMock) Collect(ctx context.Context, from, to time.Time) (set collector.MetadataSet, err error) {
	return nil, fmt.Errorf("error")
}

func toDayFormat(myTime time.Time) string {
	return myTime.Format(timeutils.DayFormat)
}

func setupRouter(mock collector.Collector) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	ctrl := &Controller{nasa: mock}
	r.GET("/test-pictures", ctrl.GetPictures)
	return r
}

func TestController_GetPictures(t *testing.T) {
	tests := []struct {
		name        string
		nasa        collector.Collector
		queryParams string
		wantCode    int
		wantBody    string
	}{
		{
			name:        "'to' in wrong format",
			queryParams: "?to=1673184713",
			wantCode:    http.StatusBadRequest,
			wantBody:    "{\"error\":\"couldn't parse query params properly. They must be in 2006-01-02 format\"}",
		},
		{
			name:        "'from' in wrong format",
			queryParams: "?from=2006-01-02T15:04:05Z07:00",
			wantCode:    http.StatusBadRequest,
			wantBody:    "{\"error\":\"couldn't parse query params properly. They must be in 2006-01-02 format\"}",
		},
		{
			name:        "'to' defined only",
			nasa:        &successNasaMock{},
			queryParams: "?to=" + toDayFormat(time.Now()),
			wantCode:    http.StatusBadRequest,
			wantBody:    "{\"error\":\"validation failed: 'from' must be defined if 'to' is set\"}",
		},
		{
			name:     "failed to collect",
			nasa:     &errorNasaMock{},
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"error\":\"failed to collect urls\"}",
		},
		{
			name: "'to' and 'from' after now",
			queryParams: func() string {
				now := time.Now()
				from := now.AddDate(0, 0, 1)
				to := now.AddDate(0, 0, 1)
				return fmt.Sprintf("?from=%s&to=%s", toDayFormat(from), toDayFormat(to))
			}(),
			wantCode: http.StatusBadRequest,
			wantBody: "{\"error\":\"validation failed: the given time cannot be future\"}",
		},
		{
			name: "'from' after 'to'",
			nasa: &successNasaMock{},
			queryParams: func() string {
				now := time.Now()
				from := now
				to := now.AddDate(0, 0, -1)
				return fmt.Sprintf("?from=%s&to=%s", toDayFormat(from), toDayFormat(to))
			}(),
			wantCode: http.StatusBadRequest,
			wantBody: "{\"error\":\"validation failed: 'from' must be before 'to' param\"}",
		},
		{
			name:        "empty query params",
			nasa:        &successNasaMock{},
			queryParams: "",
			wantCode:    http.StatusOK,
			wantBody:    "{\"urls\":[\"https://mock.nasa.com/picture/0\"]}",
		},
		{
			name:        "'from' defined only",
			nasa:        &successNasaMock{},
			queryParams: "?from=" + toDayFormat(time.Now()),
			wantCode:    http.StatusOK,
			wantBody:    "{\"urls\":[\"https://mock.nasa.com/picture/0\"]}",
		},
		{
			name: "'from' before 'to'",
			nasa: &successNasaMock{},
			queryParams: func() string {
				now := time.Now()
				from := now.AddDate(0, 0, -1)
				to := now
				return fmt.Sprintf("?from=%s&to=%s", toDayFormat(from), toDayFormat(to))
			}(),
			wantCode: http.StatusOK,
			wantBody: "{\"urls\":[\"https://mock.nasa.com/picture/0\",\"https://mock.nasa.com/picture/1\"]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.nasa)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/test-pictures"+tt.queryParams, nil)
			require.Nil(t, err)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Result().StatusCode)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}
