package suggest

import (
	"fmt"

	"github.com/alldroll/suggest/pkg/metric"
)

// SearchConfig is a config for NGramIndex Suggest method
type SearchConfig struct {
	query      string
	topK       int
	metric     metric.Metric
	similarity float64
}

// NewSearchConfig returns new instance of SearchConfig
func NewSearchConfig(query string, topK int, metric metric.Metric, similarity float64) (*SearchConfig, error) {
	if topK < 0 {
		return nil, fmt.Errorf("topK is invalid") //TODO fixme
	}

	if similarity <= 0 || similarity > 1 {
		return nil, fmt.Errorf("similarity shouble be in (0.0, 1.0]")
	}

	return &SearchConfig{
		query:      query,
		topK:       topK,
		metric:     metric,
		similarity: similarity,
	}, nil
}
