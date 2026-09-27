# mailcraft-go

Official Go SDK for the [MailCraft](https://mailcraft.host) email API: transactional email, contacts, lists, segments, templates, campaigns, webhooks and more.

Requires Go 1.21+. No dependencies beyond the standard library.

## Install

```bash
go get github.com/mail-craft/mailcraft-go
```

## Quick start

```go
import (
	"context"
	"os"

	"github.com/mail-craft/mailcraft-go"
)

client := mailcraft.NewClient(os.Getenv("MAILCRAFT_API_KEY"))

email, err := client.Emails.Send(context.Background(), &mailcraft.SendEmailParams{
	From:    "hello@yourdomain.com",
	To:      []string{"person@example.com"},
	Subject: "Welcome!",
	HTML:    "<p>Thanks for signing up.</p>",
})
if err != nil {
	return err
}
fmt.Println(email["data"].(map[string]any)["id"])
```

Create an API key under **Settings → API keys** in your MailCraft dashboard. Every method takes a `context.Context`, so you can cancel requests or set deadlines per call.

## Resources

Every MailCraft SDK has the same resources and methods:

| Resource | Methods |
| --- | --- |
| `Emails` | `Send`, `List`, `Get`, `Validate` |
| `Domains` | `Create`, `List`, `Get`, `Verify`, `Delete` |
| `Senders` | `Create`, `List`, `Get`, `Delete` |
| `Contacts` | `Upsert`, `List`, `Get`, `Delete`, `Unsubscribe`, `AddToLists`, `Lists`, `RemoveFromList` |
| `Lists` | `Create`, `List`, `Get`, `Delete` |
| `Segments` | `Create`, `List`, `Get`, `Delete` |
| `Properties` | `Create`, `List`, `Delete` |
| `Templates` | `Create`, `List`, `Get`, `Update`, `Delete` |
| `TemplateFolders` | `Create`, `List`, `Delete` |
| `Campaigns` | `Create`, `List`, `Get`, `Send`, `Delete` |
| `Webhooks` | `Create`, `List`, `Delete` |
| `Suppressions` | `Add`, `List`, `Delete` |
| `Metrics` | `Get`, `Reputation` |

Requests take typed params (`SendEmailParams`, `UpsertContactParams`, ...); empty optional fields are left out of the request. Responses are the API's JSON as a `mailcraft.Response` (`map[string]any`). See the [API reference](https://docs.mailcraft.host/api-reference) for every field.

## Errors

Any non-2xx response returns an `*mailcraft.APIError`:

```go
_, err := client.Domains.Create(ctx, &mailcraft.CreateDomainParams{Name: "acme.com"})

var apiErr *mailcraft.APIError
if errors.As(err, &apiErr) {
	fmt.Println(apiErr.StatusCode) // e.g. 402
	fmt.Println(apiErr.Type)       // e.g. "plan_limit_reached" (business-rule errors)
	fmt.Println(apiErr.Errors)     // field errors for 422 validation failures
}
```

## Options

```go
client := mailcraft.NewClient(apiKey,
	mailcraft.WithBaseURL("https://api.mailcraft.host/v1"),         // override for staging or self-hosting
	mailcraft.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}), // default timeout is 30s
)
```

## Development

```bash
go test ./...
```

The tests run the real client against an `httptest.Server`, so every request goes through the same code a user's would.

## License

MIT
