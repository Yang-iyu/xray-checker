package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"xray-checker/checker"
)

func newEmptyProxyChecker() *checker.ProxyChecker {
	return checker.NewProxyChecker(nil, 10000, "", 5, "", "", 5, 1, "status", 0)
}

func TestOnlineSubscriptionHandlerEmptySubscription(t *testing.T) {
	handler := OnlineSubscriptionHandler(newEmptyProxyChecker())
	req := httptest.NewRequest(http.MethodGet, "/online-nodes.txt", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	if resp.Body.String() != "" {
		t.Fatalf("expected an empty response body, got %q", resp.Body.String())
	}
	if contentType := resp.Header().Get("Content-Type"); contentType != "text/plain; charset=utf-8" {
		t.Fatalf("unexpected Content-Type %q", contentType)
	}
	if disposition := resp.Header().Get("Content-Disposition"); !strings.Contains(disposition, "filename=\"online-nodes.txt\"") {
		t.Fatalf("unexpected Content-Disposition %q", disposition)
	}
	if cacheControl := resp.Header().Get("Cache-Control"); cacheControl != "no-store" {
		t.Fatalf("unexpected Cache-Control %q", cacheControl)
	}
}

func TestOnlineSubscriptionHandlerHead(t *testing.T) {
	handler := OnlineSubscriptionHandler(newEmptyProxyChecker())
	req := httptest.NewRequest(http.MethodHead, "/online-nodes.txt", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	if resp.Body.Len() != 0 {
		t.Fatalf("expected HEAD response without a body, got %q", resp.Body.String())
	}
}

func TestOnlineSubscriptionHandlerRejectsUnsupportedMethods(t *testing.T) {
	handler := OnlineSubscriptionHandler(newEmptyProxyChecker())
	req := httptest.NewRequest(http.MethodPost, "/online-nodes.txt", nil)
	resp := httptest.NewRecorder()

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, resp.Code)
	}
	if allow := resp.Header().Get("Allow"); allow != "GET, HEAD" {
		t.Fatalf("unexpected Allow header %q", allow)
	}
}
