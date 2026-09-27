package mailcraft

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

func limitQuery(params *ListParams) url.Values {
	if params == nil || params.Limit <= 0 {
		return nil
	}
	return url.Values{"limit": {strconv.Itoa(params.Limit)}}
}

// EmailsService sends and looks up transactional email.
type EmailsService struct{ client *Client }

func (s *EmailsService) Send(ctx context.Context, params *SendEmailParams) (Response, error) {
	return s.client.post(ctx, "/emails", params)
}

func (s *EmailsService) List(ctx context.Context, params *ListParams) (Response, error) {
	return s.client.get(ctx, "/emails", limitQuery(params))
}

func (s *EmailsService) Get(ctx context.Context, id string) (Response, error) {
	return s.client.get(ctx, "/emails/"+url.PathEscape(id), nil)
}

// Validate checks an address's format, MX records and disposable domain, without sending.
func (s *EmailsService) Validate(ctx context.Context, email string) (Response, error) {
	return s.client.get(ctx, "/emails/validate", url.Values{"email": {email}})
}

// DomainsService manages sending domains.
type DomainsService struct{ client *Client }

func (s *DomainsService) Create(ctx context.Context, params *CreateDomainParams) (Response, error) {
	return s.client.post(ctx, "/domains", params)
}

func (s *DomainsService) List(ctx context.Context) (Response, error) {
	return s.client.get(ctx, "/domains", nil)
}

func (s *DomainsService) Get(ctx context.Context, id int64) (Response, error) {
	return s.client.get(ctx, fmt.Sprintf("/domains/%d", id), nil)
}

func (s *DomainsService) Verify(ctx context.Context, id int64) (Response, error) {
	return s.client.post(ctx, fmt.Sprintf("/domains/%d/verify", id), nil)
}

func (s *DomainsService) Delete(ctx context.Context, id int64) error {
	return s.client.delete(ctx, fmt.Sprintf("/domains/%d", id))
}

// SendersService manages sender identities.
type SendersService struct{ client *Client }

func (s *SendersService) Create(ctx context.Context, params *CreateSenderParams) (Response, error) {
	return s.client.post(ctx, "/senders", params)
}

func (s *SendersService) List(ctx context.Context) (Response, error) {
	return s.client.get(ctx, "/senders", nil)
}

func (s *SendersService) Get(ctx context.Context, id int64) (Response, error) {
	return s.client.get(ctx, fmt.Sprintf("/senders/%d", id), nil)
}

func (s *SendersService) Delete(ctx context.Context, id int64) error {
	return s.client.delete(ctx, fmt.Sprintf("/senders/%d", id))
}

// ContactsService manages contacts and their list memberships.
type ContactsService struct{ client *Client }

func (s *ContactsService) Upsert(ctx context.Context, params *UpsertContactParams) (Response, error) {
	return s.client.post(ctx, "/contacts", params)
}

func (s *ContactsService) List(ctx context.Context, params *ListParams) (Response, error) {
	return s.client.get(ctx, "/contacts", limitQuery(params))
}

func (s *ContactsService) Get(ctx context.Context, id string) (Response, error) {
	return s.client.get(ctx, "/contacts/"+url.PathEscape(id), nil)
}

func (s *ContactsService) Delete(ctx context.Context, id string) error {
	return s.client.delete(ctx, "/contacts/"+url.PathEscape(id))
}

func (s *ContactsService) Unsubscribe(ctx context.Context, id string) (Response, error) {
	return s.client.post(ctx, "/contacts/"+url.PathEscape(id)+"/unsubscribe", nil)
}

func (s *ContactsService) AddToLists(ctx context.Context, id string, listIDs []int64) error {
	_, err := s.client.post(ctx, "/contacts/"+url.PathEscape(id)+"/lists", map[string]any{"list_ids": listIDs})
	return err
}

func (s *ContactsService) Lists(ctx context.Context, id string) (Response, error) {
	return s.client.get(ctx, "/contacts/"+url.PathEscape(id)+"/lists", nil)
}

func (s *ContactsService) RemoveFromList(ctx context.Context, id string, listID int64) error {
	return s.client.delete(ctx, fmt.Sprintf("/contacts/%s/lists/%d", url.PathEscape(id), listID))
}

// ListsService manages mailing lists.
type ListsService struct{ client *Client }

func (s *ListsService) Create(ctx context.Context, params *CreateListParams) (Response, error) {
	return s.client.post(ctx, "/lists", params)
}

func (s *ListsService) List(ctx context.Context) (Response, error) {
	return s.client.get(ctx, "/lists", nil)
}

func (s *ListsService) Get(ctx context.Context, id int64) (Response, error) {
	return s.client.get(ctx, fmt.Sprintf("/lists/%d", id), nil)
}

func (s *ListsService) Delete(ctx context.Context, id int64) error {
	return s.client.delete(ctx, fmt.Sprintf("/lists/%d", id))
}

// SegmentsService manages segments.
type SegmentsService struct{ client *Client }

func (s *SegmentsService) Create(ctx context.Context, params *CreateSegmentParams) (Response, error) {
	return s.client.post(ctx, "/segments", params)
}

func (s *SegmentsService) List(ctx context.Context) (Response, error) {
	return s.client.get(ctx, "/segments", nil)
}

func (s *SegmentsService) Get(ctx context.Context, id int64) (Response, error) {
	return s.client.get(ctx, fmt.Sprintf("/segments/%d", id), nil)
}

func (s *SegmentsService) Delete(ctx context.Context, id int64) error {
	return s.client.delete(ctx, fmt.Sprintf("/segments/%d", id))
}

// PropertiesService manages custom contact properties.
type PropertiesService struct{ client *Client }

func (s *PropertiesService) Create(ctx context.Context, params *CreatePropertyParams) (Response, error) {
	return s.client.post(ctx, "/properties", params)
}

func (s *PropertiesService) List(ctx context.Context) (Response, error) {
	return s.client.get(ctx, "/properties", nil)
}

func (s *PropertiesService) Delete(ctx context.Context, id int64) error {
	return s.client.delete(ctx, fmt.Sprintf("/properties/%d", id))
}

// TemplatesService manages versioned templates.
type TemplatesService struct{ client *Client }

func (s *TemplatesService) Create(ctx context.Context, params *TemplateParams) (Response, error) {
	return s.client.post(ctx, "/templates", params)
}

func (s *TemplatesService) List(ctx context.Context) (Response, error) {
	return s.client.get(ctx, "/templates", nil)
}

func (s *TemplatesService) Get(ctx context.Context, id int64) (Response, error) {
	return s.client.get(ctx, fmt.Sprintf("/templates/%d", id), nil)
}

// Update changes a template; each save becomes a new version. Unset fields are left as they are.
func (s *TemplatesService) Update(ctx context.Context, id int64, params *TemplateParams) (Response, error) {
	return s.client.patch(ctx, fmt.Sprintf("/templates/%d", id), params)
}

func (s *TemplatesService) Delete(ctx context.Context, id int64) error {
	return s.client.delete(ctx, fmt.Sprintf("/templates/%d", id))
}

// TemplateFoldersService organises templates into folders.
type TemplateFoldersService struct{ client *Client }

func (s *TemplateFoldersService) Create(ctx context.Context, name string) (Response, error) {
	return s.client.post(ctx, "/template-folders", map[string]string{"name": name})
}

func (s *TemplateFoldersService) List(ctx context.Context) (Response, error) {
	return s.client.get(ctx, "/template-folders", nil)
}

func (s *TemplateFoldersService) Delete(ctx context.Context, id int64) error {
	return s.client.delete(ctx, fmt.Sprintf("/template-folders/%d", id))
}

// CampaignsService manages marketing campaigns.
type CampaignsService struct{ client *Client }

func (s *CampaignsService) Create(ctx context.Context, params *CreateCampaignParams) (Response, error) {
	return s.client.post(ctx, "/campaigns", params)
}

func (s *CampaignsService) List(ctx context.Context) (Response, error) {
	return s.client.get(ctx, "/campaigns", nil)
}

func (s *CampaignsService) Get(ctx context.Context, id int64) (Response, error) {
	return s.client.get(ctx, fmt.Sprintf("/campaigns/%d", id), nil)
}

func (s *CampaignsService) Send(ctx context.Context, id int64) (Response, error) {
	return s.client.post(ctx, fmt.Sprintf("/campaigns/%d/send", id), nil)
}

func (s *CampaignsService) Delete(ctx context.Context, id int64) error {
	return s.client.delete(ctx, fmt.Sprintf("/campaigns/%d", id))
}

// WebhooksService manages webhook endpoints.
type WebhooksService struct{ client *Client }

func (s *WebhooksService) Create(ctx context.Context, params *CreateWebhookParams) (Response, error) {
	return s.client.post(ctx, "/webhooks", params)
}

func (s *WebhooksService) List(ctx context.Context) (Response, error) {
	return s.client.get(ctx, "/webhooks", nil)
}

func (s *WebhooksService) Delete(ctx context.Context, id int64) error {
	return s.client.delete(ctx, fmt.Sprintf("/webhooks/%d", id))
}

// SuppressionsService manages the suppression list.
type SuppressionsService struct{ client *Client }

func (s *SuppressionsService) Add(ctx context.Context, params *AddSuppressionParams) (Response, error) {
	return s.client.post(ctx, "/suppressions", params)
}

func (s *SuppressionsService) List(ctx context.Context, params *ListParams) (Response, error) {
	return s.client.get(ctx, "/suppressions", limitQuery(params))
}

func (s *SuppressionsService) Delete(ctx context.Context, id int64) error {
	return s.client.delete(ctx, fmt.Sprintf("/suppressions/%d", id))
}

// MetricsService reads sending metrics and reputation.
type MetricsService struct{ client *Client }

func (s *MetricsService) Get(ctx context.Context, params *MetricsParams) (Response, error) {
	query := url.Values{}
	if params != nil {
		if params.StartDate != "" {
			query.Set("start_date", params.StartDate)
		}
		if params.EndDate != "" {
			query.Set("end_date", params.EndDate)
		}
	}
	return s.client.get(ctx, "/metrics", query)
}

func (s *MetricsService) Reputation(ctx context.Context) (Response, error) {
	return s.client.get(ctx, "/reputation", nil)
}
