package nasa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	"url-collector/internal/collector"

	"github.com/stretchr/testify/assert"
)

func errorRequestCallerMock(_ http.Client, _ string) (*http.Response, error) {
	return nil, fmt.Errorf("error")
}

func failRequestCallerMock(_ http.Client, _ string) (*http.Response, error) {
	return &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("error")),
		StatusCode: http.StatusInternalServerError,
	}, nil
}

func successRequestCallerMock(_ http.Client, _ string) (*http.Response, error) {
	r := Resource{
		Title: "title",
		URL:   "https://mock.com/resource/id",
		HDURL: "https://mock.com/resource-hd/id",
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
		StatusCode: http.StatusOK,
	}, nil
}

func TestCollector_Collect(t *testing.T) {
	type fields struct {
		semaphore     chan struct{}
		requestCaller func(client http.Client, url string) (*http.Response, error)
	}
	type args struct {
		from time.Time
		to   time.Time
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantSet   collector.MetadataSet
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "invalid from and to params",
			fields: fields{
				semaphore:     make(chan struct{}, 1),
				requestCaller: errorRequestCallerMock,
			},
			args: func() args {
				now := time.Now()
				return args{
					from: now.AddDate(0, 0, 1),
					to:   now,
				}
			}(),
			wantSet:   nil,
			assertion: assert.NoError,
		},
		{
			name: "valid from and to params but request call error happened",
			fields: fields{
				semaphore:     make(chan struct{}, 1),
				requestCaller: errorRequestCallerMock,
			},
			args: func() args {
				now := time.Now()
				return args{
					from: now,
					to:   now,
				}
			}(),
			wantSet:   nil,
			assertion: assert.Error,
		},
		{
			name: "valid from and to params but request code is not 200",
			fields: fields{
				semaphore:     make(chan struct{}, 1),
				requestCaller: failRequestCallerMock,
			},
			args: func() args {
				now := time.Now()
				return args{
					from: now,
					to:   now,
				}
			}(),
			wantSet:   nil,
			assertion: assert.Error,
		},
		{
			name: "valid from and to params and no errors",
			fields: fields{
				semaphore:     make(chan struct{}, 1),
				requestCaller: successRequestCallerMock,
			},
			args: func() args {
				now := time.Now()
				return args{
					from: now,
					to:   now,
				}
			}(),
			wantSet: collector.MetadataSet{
				{Title: "title", URL: "https://mock.com/resource/id", HDURL: "https://mock.com/resource-hd/id"},
			},
			assertion: assert.NoError,
		},
		{
			name: "valid from and to params and no errors with bigger semaphore",
			fields: fields{
				semaphore:     make(chan struct{}, 3),
				requestCaller: successRequestCallerMock,
			},
			args: func() args {
				now := time.Now()
				return args{
					from: now.AddDate(0, 0, -5),
					to:   now,
				}
			}(),
			wantSet: collector.MetadataSet{
				{Title: "title", URL: "https://mock.com/resource/id", HDURL: "https://mock.com/resource-hd/id"},
				{Title: "title", URL: "https://mock.com/resource/id", HDURL: "https://mock.com/resource-hd/id"},
				{Title: "title", URL: "https://mock.com/resource/id", HDURL: "https://mock.com/resource-hd/id"},
				{Title: "title", URL: "https://mock.com/resource/id", HDURL: "https://mock.com/resource-hd/id"},
				{Title: "title", URL: "https://mock.com/resource/id", HDURL: "https://mock.com/resource-hd/id"},
				{Title: "title", URL: "https://mock.com/resource/id", HDURL: "https://mock.com/resource-hd/id"},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nasa := &Collector{
				semaphore:     tt.fields.semaphore,
				requestCaller: tt.fields.requestCaller,
			}
			gotSet, err := nasa.Collect(context.Background(), tt.args.from, tt.args.to)
			tt.assertion(t, err)
			assert.Equal(t, tt.wantSet, gotSet)
		})
	}
}
