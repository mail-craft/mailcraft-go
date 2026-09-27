package mailcraft

// SendEmailParams is a transactional email to send.
type SendEmailParams struct {
	From    string            `json:"from"`
	To      []string          `json:"to"`
	Subject string            `json:"subject"`
	HTML    string            `json:"html,omitempty"`
	Text    string            `json:"text,omitempty"`
	CC      []string          `json:"cc,omitempty"`
	BCC     []string          `json:"bcc,omitempty"`
	ReplyTo string            `json:"reply_to,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Tags    map[string]string `json:"tags,omitempty"`
}

// CreateDomainParams registers a sending domain.
type CreateDomainParams struct {
	Name   string `json:"name"`
	Region string `json:"region,omitempty"`
}

// CreateSenderParams creates a sender identity on a verified domain.
type CreateSenderParams struct {
	DomainID int64  `json:"domain_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	ReplyTo  string `json:"reply_to,omitempty"`
}

// UpsertContactParams creates a contact, or updates it if one exists for this email.
type UpsertContactParams struct {
	Email      string         `json:"email"`
	FirstName  string         `json:"first_name,omitempty"`
	LastName   string         `json:"last_name,omitempty"`
	Status     string         `json:"status,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
}

// CreateListParams creates a list. Type is "static" or "dynamic" (with SegmentID).
type CreateListParams struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	SegmentID   int64  `json:"segment_id,omitempty"`
}

// CreateSegmentParams creates a segment from filter conditions.
type CreateSegmentParams struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Filters     SegmentFilters `json:"filters"`
}

// SegmentFilters matches contacts on all ("and") or any ("or") of its conditions.
type SegmentFilters struct {
	Operator   string             `json:"operator,omitempty"`
	Conditions []SegmentCondition `json:"conditions,omitempty"`
}

// SegmentCondition compares a contact field. Operator is eq, neq, gt, gte, lt, lte or contains.
type SegmentCondition struct {
	Field    string `json:"field"`
	Operator string `json:"operator,omitempty"`
	Value    any    `json:"value,omitempty"`
}

// CreatePropertyParams creates a custom contact property.
// Type is "text", "number", "boolean", "date" or "list".
type CreatePropertyParams struct {
	Key          string `json:"key"`
	Label        string `json:"label"`
	Type         string `json:"type"`
	DefaultValue any    `json:"default_value,omitempty"`
}

// TemplateParams creates or updates a template. Name and Subject are required to create.
type TemplateParams struct {
	Name             string           `json:"name,omitempty"`
	Subject          string           `json:"subject,omitempty"`
	HTMLBody         string           `json:"html_body,omitempty"`
	TextBody         string           `json:"text_body,omitempty"`
	TemplateFolderID int64            `json:"template_folder_id,omitempty"`
	Variables        []map[string]any `json:"variables,omitempty"`
}

// CreateCampaignParams creates a draft campaign for a list or a segment.
type CreateCampaignParams struct {
	Name       string `json:"name"`
	Subject    string `json:"subject"`
	TemplateID int64  `json:"template_id"`
	SenderID   int64  `json:"sender_id"`
	ListID     int64  `json:"list_id,omitempty"`
	SegmentID  int64  `json:"segment_id,omitempty"`
}

// CreateWebhookParams subscribes a URL to events.
type CreateWebhookParams struct {
	URL         string   `json:"url"`
	Events      []string `json:"events"`
	Description string   `json:"description,omitempty"`
}

// AddSuppressionParams stops sending to an address.
type AddSuppressionParams struct {
	Email  string `json:"email"`
	Reason string `json:"reason,omitempty"`
}

// ListParams limits how many items a list call returns.
type ListParams struct {
	Limit int
}

// MetricsParams narrows metrics to a date range (YYYY-MM-DD).
type MetricsParams struct {
	StartDate string
	EndDate   string
}
