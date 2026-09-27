package mailcraft

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type recorded struct {
	method string
	path   string
	query  string
	header http.Header
	body   map[string]any
}

func newTestClient(t *testing.T, status int, response string) (*Client, *recorded) {
	t.Helper()

	rec := &recorded{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.path = r.URL.Path
		rec.query = r.URL.RawQuery
		rec.header = r.Header.Clone()
		raw, _ := io.ReadAll(r.Body)
		if len(raw) > 0 {
			_ = json.Unmarshal(raw, &rec.body)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(server.Close)

	return NewClient("mc_test_key", WithBaseURL(server.URL+"/v1/")), rec
}

func TestSendEmail(t *testing.T) {
	client, rec := newTestClient(t, 202, `{"data":{"id":"em_1","status":"queued"}}`)

	res, err := client.Emails.Send(context.Background(), &SendEmailParams{
		From:    "hello@yourdomain.com",
		To:      []string{"person@example.com"},
		Subject: "Welcome!",
		HTML:    "<p>Hi</p>",
		Tags:    map[string]string{"flow": "welcome"},
	})
	if err != nil {
		t.Fatal(err)
	}

	if rec.method != "POST" || rec.path != "/v1/emails" {
		t.Fatalf("got %s %s", rec.method, rec.path)
	}
	if got := rec.header.Get("Authorization"); got != "Bearer mc_test_key" {
		t.Fatalf("authorization header %q", got)
	}
	if got := rec.header.Get("User-Agent"); got != "mailcraft-go/"+Version {
		t.Fatalf("user agent %q", got)
	}
	if rec.body["from"] != "hello@yourdomain.com" || rec.body["subject"] != "Welcome!" {
		t.Fatalf("body %v", rec.body)
	}
	if _, ok := rec.body["text"]; ok {
		t.Fatal("empty optional fields should be omitted")
	}
	if res["data"].(map[string]any)["id"] != "em_1" {
		t.Fatalf("response %v", res)
	}
}

func TestListSendsLimit(t *testing.T) {
	client, rec := newTestClient(t, 200, `{"data":[]}`)

	if _, err := client.Contacts.List(context.Background(), &ListParams{Limit: 5}); err != nil {
		t.Fatal(err)
	}
	if rec.path != "/v1/contacts" || rec.query != "limit=5" {
		t.Fatalf("got %s?%s", rec.path, rec.query)
	}
}

func TestValidateEscapesEmail(t *testing.T) {
	client, rec := newTestClient(t, 200, `{"email":"a+b@example.com","valid":true}`)

	if _, err := client.Emails.Validate(context.Background(), "a+b@example.com"); err != nil {
		t.Fatal(err)
	}
	if rec.path != "/v1/emails/validate" || rec.query != "email=a%2Bb%40example.com" {
		t.Fatalf("got %s?%s", rec.path, rec.query)
	}
}

func TestTemplateUpdateUsesPatch(t *testing.T) {
	client, rec := newTestClient(t, 200, `{"data":{"id":3}}`)

	if _, err := client.Templates.Update(context.Background(), 3, &TemplateParams{Subject: "New"}); err != nil {
		t.Fatal(err)
	}
	if rec.method != "PATCH" || rec.path != "/v1/templates/3" {
		t.Fatalf("got %s %s", rec.method, rec.path)
	}
	if len(rec.body) != 1 || rec.body["subject"] != "New" {
		t.Fatalf("body %v", rec.body)
	}
}

func TestContactListMembership(t *testing.T) {
	client, rec := newTestClient(t, 204, ``)
	ctx := context.Background()

	if err := client.Contacts.AddToLists(ctx, "c_1", []int64{1, 2}); err != nil {
		t.Fatal(err)
	}
	if rec.path != "/v1/contacts/c_1/lists" || len(rec.body["list_ids"].([]any)) != 2 {
		t.Fatalf("got %s %v", rec.path, rec.body)
	}

	if err := client.Contacts.RemoveFromList(ctx, "c_1", 2); err != nil {
		t.Fatal(err)
	}
	if rec.method != "DELETE" || rec.path != "/v1/contacts/c_1/lists/2" {
		t.Fatalf("got %s %s", rec.method, rec.path)
	}
}

func TestDeleteReturnsNoContent(t *testing.T) {
	client, rec := newTestClient(t, 204, ``)

	if err := client.Domains.Delete(context.Background(), 9); err != nil {
		t.Fatal(err)
	}
	if rec.method != "DELETE" || rec.path != "/v1/domains/9" {
		t.Fatalf("got %s %s", rec.method, rec.path)
	}
}

func TestMetricsAndReputation(t *testing.T) {
	client, rec := newTestClient(t, 200, `{"data":[]}`)
	ctx := context.Background()

	if _, err := client.Metrics.Get(ctx, &MetricsParams{StartDate: "2026-01-01", EndDate: "2026-01-31"}); err != nil {
		t.Fatal(err)
	}
	if rec.path != "/v1/metrics" || rec.query != "end_date=2026-01-31&start_date=2026-01-01" {
		t.Fatalf("got %s?%s", rec.path, rec.query)
	}

	if _, err := client.Metrics.Reputation(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.path != "/v1/reputation" {
		t.Fatalf("got %s", rec.path)
	}
}

func TestTemplateFoldersUseHyphenatedPath(t *testing.T) {
	client, rec := newTestClient(t, 201, `{"data":{"id":1}}`)

	if _, err := client.TemplateFolders.Create(context.Background(), "Onboarding"); err != nil {
		t.Fatal(err)
	}
	if rec.path != "/v1/template-folders" || rec.body["name"] != "Onboarding" {
		t.Fatalf("got %s %v", rec.path, rec.body)
	}
}

func TestBusinessRuleError(t *testing.T) {
	client, _ := newTestClient(t, 402, `{"error":{"type":"plan_limit_reached","message":"Monthly limit reached."}}`)

	_, err := client.Campaigns.Send(context.Background(), 4)

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %v", err)
	}
	if apiErr.StatusCode != 402 || apiErr.Type != "plan_limit_reached" || apiErr.Message != "Monthly limit reached." {
		t.Fatalf("got %+v", apiErr)
	}
}

func TestValidationError(t *testing.T) {
	client, _ := newTestClient(t, 422, `{"message":"The email field is required.","errors":{"email":["The email field is required."]}}`)

	_, err := client.Contacts.Upsert(context.Background(), &UpsertContactParams{})

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %v", err)
	}
	if apiErr.StatusCode != 422 || apiErr.Type != "" || apiErr.Errors["email"][0] != "The email field is required." {
		t.Fatalf("got %+v", apiErr)
	}
}

func TestUnparseableErrorFallsBackToStatusText(t *testing.T) {
	client, _ := newTestClient(t, 502, `<html>Bad gateway</html>`)

	_, err := client.Lists.List(context.Background())

	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 502 || apiErr.Message != "Bad Gateway" {
		t.Fatalf("got %v", err)
	}
}

func TestNewClientRequiresAPIKey(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic for an empty API key")
		}
	}()

	NewClient("")
}
