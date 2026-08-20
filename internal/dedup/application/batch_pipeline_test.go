package application

import (
	"testing"

	"github.com/acme/signalforge/internal/dedup/domain"
	"github.com/acme/signalforge/internal/dedup/infrastructure"
)

func TestBatchPipelineDoesNotMutateCallerSlice(t *testing.T) {
	input := []domain.BatchItem{{Key: "a", Labels: []string{"one"}}, {Key: "", Labels: []string{"skip"}}}
	normalized := NormalizeBatch(input)
	normalized[0].Labels[0] = "changed"
	if input[0].Labels[0] != "one" || input[1].Key != "" {
		t.Fatalf("normalization mutated caller data: %#v", input)
	}
	stored := infrastructure.StoreBatch(normalized)
	stored[0].Labels[0] = "stored-change"
	if normalized[0].Labels[0] != "changed" || AggregateBatch(normalized)["a"] != 1 {
		t.Fatalf("stored batch must be isolated: %#v", normalized)
	}
}
