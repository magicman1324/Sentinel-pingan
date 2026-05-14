package collector

import (
	"context"

	"github.com/pingan/monitor-agent/internal/model"
)

type Collector interface {
	Name() string
	Collect(ctx context.Context, out *model.Metrics)
}
