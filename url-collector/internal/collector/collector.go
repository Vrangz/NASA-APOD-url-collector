package collector

import (
	"context"
	"time"
)

// Collector defines the abstraction of how the data can be collected
type Collector interface {
	Collect(ctx context.Context, from, to time.Time) (MetadataSet, error)
}
