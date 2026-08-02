package metrics

import "testing"

func TestCachedCollector(t *testing.T) {
	mock := &MockCollector{}

	cache := &CachedCollector{
		collector: mock,
	}

	_, err := cache.Collect()
	if err != nil {
		t.Fatalf("Failed to collect metrics: %v", err)
	}

	_, err = cache.Collect()
	if err != nil {
		t.Fatalf("Failed to collect metrics: %v", err)
	}

	if mock.Count != 1 {
		t.Fatalf("Expected mock count to be 1, got %d", mock.Count)
	}
}
