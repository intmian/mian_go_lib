package ai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestProviderOptionsDisableRetriesAndExposeStatus(t *testing.T) {
	var hits atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":{"message":"busy"}}`))
	}))
	defer server.Close()

	provider, err := NewOpenAIProvider(server.URL, "test-token", ProviderSourceTypeOpenAI,
		WithHTTPClient(server.Client()), WithMaxRetries(0))
	if err != nil {
		t.Fatalf("NewOpenAIProvider: %v", err)
	}
	if err := provider.RegisterModel(ModelConfig{ID: "m", Name: "m", Type: ModelTypeText}); err != nil {
		t.Fatalf("RegisterModel: %v", err)
	}
	_, err = provider.Chat(context.Background(), ChatRequest{Model: "m", Messages: []ChatMessage{{Role: ChatRoleUser, Content: "hi"}}})
	if err == nil {
		t.Fatal("expected error")
	}
	if got := hits.Load(); got != 1 {
		t.Fatalf("requests = %d, want 1 (retries disabled)", got)
	}
	if status, ok := HTTPStatusFromError(err); !ok || status != http.StatusServiceUnavailable {
		t.Fatalf("HTTPStatusFromError = %d, %v", status, ok)
	}
}

func TestHTTPStatusFromErrorWithoutResponse(t *testing.T) {
	if _, ok := HTTPStatusFromError(context.DeadlineExceeded); ok {
		t.Fatal("expected no status for transport error")
	}
}
