package nasa

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"url-collector/internal/collector"
	"url-collector/internal/timeutils"
)

// Collector of nasa data
type Collector struct {
	client        http.Client
	url           string
	apiKey        string
	semaphore     chan struct{}
	requestCaller func(client http.Client, url string) (*http.Response, error)
}

type result struct {
	metadata collector.Metadata
	err      error
}

func requestCaller(client http.Client, url string) (*http.Response, error) {
	return client.Get(url)
}

// New creates new instance of nasa collector
func New(client http.Client, url, apiKey string, concurrentRequests uint) *Collector {
	return &Collector{
		client:        client,
		url:           url,
		apiKey:        apiKey,
		semaphore:     make(chan struct{}, concurrentRequests),
		requestCaller: requestCaller,
	}
}

func (nasa *Collector) Collect(ctx context.Context, from, to time.Time) (set collector.MetadataSet, err error) {
	var (
		dates = timeutils.List(from, to)
		sink  = make(chan result, len(dates))
	)

	for _, date := range dates {
		go nasa.collectOne(ctx, date, sink)
	}

	for range dates {
		res := <-sink
		if res.err != nil {
			return nil, res.err
		}
		set = append(set, res.metadata)
	}

	close(sink)
	return
}

func (nasa *Collector) collectOne(ctx context.Context, date string, sink chan<- result) {
	var (
		rsp *http.Response
		r   Resource
		err error
	)

	nasa.semaphore <- struct{}{}
	defer func() { <-nasa.semaphore }()

	if rsp, err = nasa.requestCaller(nasa.client, nasa.buildURL(date)); err != nil {
		sink <- result{err: err}
		return
	}

	if rsp.StatusCode != http.StatusOK {
		sink <- result{err: errors.New("not a 200 status code")}
		return
	}

	defer rsp.Body.Close()

	if err = json.NewDecoder(rsp.Body).Decode(&r); err != nil {
		sink <- result{err: err}
		return
	}

	sink <- result{metadata: r.ToMetadata()}
}

func (nasa *Collector) buildURL(date string) string {
	return nasa.url + "?api_key=" + nasa.apiKey + "&date=" + date
}
