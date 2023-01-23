package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-elastic-siem/internal/helpers"
	"time"

	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &DetectionRuleResource{}
var _ resource.ResourceWithImportState = &DetectionRuleResource{}

func NewDetectionRuleResource() resource.Resource {
	return &DetectionRuleResource{}
}

// DetectionRuleResource defines the resource implementation.
type DetectionRuleResource struct {
	client *helpers.Client
}

// DetectionRuleResourceModel describes the resource data model.
type DetectionRuleResourceModel struct {
	RuleContent types.String `tfsdk:"rule_content"`
	Id          types.String `tfsdk:"id"`
}

type ThreatItem struct {
	Framework string `json:"framework,omitempty"`
	Tactic    struct {
		ID        string `json:"id,omitempty"`
		Name      string `json:"name,omitempty"`
		Reference string `json:"reference,omitempty"`
	} `json:"tactic,omitempty"`
	Technique []struct {
		ID        string `json:"id,omitempty"`
		Name      string `json:"name,omitempty"`
		Reference string `json:"reference,omitempty"`
	} `json:"technique,omitempty"`
}

type ExecutionHistoryItem struct {
	LastExecution struct {
		Date        time.Time `json:"date,omitempty"`
		Status      string    `json:"status,omitempty"`
		StatusOrder int       `json:"status_order,omitempty"`
		Message     string    `json:"message,omitempty"`
		Metrics     struct {
			TotalSearchDurationMs     int `json:"total_search_duration_ms,omitempty"`
			TotalIndexingDurationMs   int `json:"total_indexing_duration_ms,omitempty"`
			TotalEnrichmentDurationMs int `json:"total_enrichment_duration_ms,omitempty"`
		} `json:"metrics,omitempty"`
	} `json:"last_execution,omitempty"`
}

type RiskScoreMapping struct {
	Field    string `json:"field,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
}

type SeverityMapping struct {
	Field    string `json:"field,omitempty"`
	Value    string `json:"value,omitempty"`
	Operator string `json:"operator,omitempty"`
	Severity string `json:"severity,omitempty"`
}

type ActionItem struct {
	Group  string `json:"group,omitempty"`
	ID     string `json:"id,omitempty"`
	Params struct {
		Body string `json:"body,omitempty"`
	} `json:"params,omitempty"`
	ActionTypeID string `json:"action_type_id,omitempty"`
}

type MetaItem struct {
	From             string `json:"from,omitempty"`
	KibanaSiemAppURL string `json:"kibana_siem_app_url,omitempty"`
}

type DetectionRuleResponse struct {
	Actions             []ActionItem         `json:"actions,omitempty"`
	Author              []string             `json:"author,omitempty"`
	CreatedAt           time.Time            `json:"created_at,omitempty"`
	CreatedBy           string               `json:"created_by,omitempty"`
	Description         string               `json:"description,omitempty"`
	Enabled             bool                 `json:"enabled,omitempty"`
	ExecutionSummary    ExecutionHistoryItem `json:"execution_summary,omitempty"`
	FalsePositives      []interface{}        `json:"false_positives,omitempty"`
	Filters             []interface{}        `json:"filters,omitempty"`
	From                string               `json:"from,omitempty"`
	ID                  string               `json:"id,omitempty"`
	Immutable           bool                 `json:"immutable,omitempty"`
	Interval            string               `json:"interval,omitempty"`
	Language            string               `json:"language,omitempty"`
	MaxSignals          int                  `json:"max_signals,omitempty"`
	Meta                MetaItem             `json:"meta,omitempty"`
	Name                string               `json:"name,omitempty"`
	OutputIndex         string               `json:"output_index,omitempty"`
	Query               string               `json:"query,omitempty"`
	References          []interface{}        `json:"references,omitempty"`
	RelatedIntegrations []interface{}        `json:"related_integrations,omitempty"`
	RequiredFields      []interface{}        `json:"required_fields,omitempty"`
	RiskScore           int                  `json:"risk_score,omitempty"`
	RiskScoreMapping    []RiskScoreMapping   `json:"risk_score_mapping,omitempty"`
	RuleID              string               `json:"rule_id,omitempty"`
	Setup               string               `json:"setup,omitempty"`
	Severity            string               `json:"severity,omitempty"`
	SeverityMapping     []SeverityMapping    `json:"severity_mapping,omitempty"`
	Tags                []string             `json:"tags,omitempty"`
	Threat              []ThreatItem         `json:"threat,omitempty"`
	Throttle            string               `json:"throttle,omitempty"`
	To                  string               `json:"to,omitempty"`
	Type                string               `json:"type,omitempty"`
	UpdatedAt           time.Time            `json:"updated_at,omitempty"`
	UpdatedBy           string               `json:"updated_by,omitempty"`
	Version             int                  `json:"version,omitempty"`
}

type DetectionRule struct {
	Actions             []ActionItem       `json:"actions,omitempty"`
	Author              []string           `json:"author,omitempty"`
	Description         string             `json:"description,omitempty"`
	Enabled             bool               `json:"enabled,omitempty"`
	ExecutionSummary    []interface{}      `json:"execution_summary,omitempty"`
	FalsePositives      []interface{}      `json:"false_positives,omitempty"`
	Filters             []interface{}      `json:"filters,omitempty"`
	From                string             `json:"from,omitempty"`
	ID                  string             `json:"id,omitempty"`
	Immutable           bool               `json:"immutable,omitempty"`
	Interval            string             `json:"interval,omitempty"`
	Language            string             `json:"language,omitempty"`
	MaxSignals          int                `json:"max_signals,omitempty"`
	Name                string             `json:"name,omitempty"`
	OutputIndex         string             `json:"output_index,omitempty"`
	Query               string             `json:"query,omitempty"`
	References          []interface{}      `json:"references,omitempty"`
	RelatedIntegrations []interface{}      `json:"related_integrations,omitempty"`
	RequiredFields      []interface{}      `json:"required_fields,omitempty"`
	RiskScore           int                `json:"risk_score,omitempty"`
	RiskScoreMapping    []RiskScoreMapping `json:"risk_score_mapping,omitempty"`
	RuleID              string             `json:"rule_id,omitempty"`
	Setup               string             `json:"setup,omitempty"`
	Severity            string             `json:"severity,omitempty"`
	SeverityMapping     []SeverityMapping  `json:"severity_mapping,omitempty"`
	Tags                []string           `json:"tags,omitempty"`
	Threat              []ThreatItem       `json:"threat,omitempty"`
	Throttle            string             `json:"throttle,omitempty"`
	To                  string             `json:"to,omitempty"`
	Type                string             `json:"type,omitempty"`
	UpdatedBy           string             `json:"updated_by,omitempty"`
	Version             int                `json:"version,omitempty"`
}

func extractRuleFronJSONStrging(ctx context.Context, yamlString string) (*DetectionRule, error) {
	result := DetectionRule{}
	tflog.Debug(ctx, yamlString)
	err := json.Unmarshal([]byte(yamlString), &result)
	return &result, err
}

func (r *DetectionRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_detection_rule"
}

func (r *DetectionRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Detection rule resource",

		Attributes: map[string]schema.Attribute{
			"rule_content": schema.StringAttribute{
				MarkdownDescription: "The content of the rule (JSON encoded string)",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Rule identifier (in UUID format)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *DetectionRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*helpers.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DetectionRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DetectionRuleResourceModel
	var body *DetectionRule

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Process the rule content
	body, err := extractRuleFronJSONStrging(ctx, data.RuleContent.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parser Error", fmt.Sprintf("Unable to parse file, got error: %s", err))
		return
	}

	// Create the rule through API
	var response DetectionRuleResponse
	if err := r.client.Post("/detection_engine/rules", body, &response); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: \n%s", err))
		return
	}

	// Save id into the Terraform state.
	data.Id = types.StringValue(response.ID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DetectionRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DetectionRuleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the rule through the API
	var response DetectionRuleResponse
	path := fmt.Sprintf("/detection_engine/rules?id=%s", data.Id.ValueString())
	if err := r.client.Get(path, &response); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DetectionRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DetectionRuleResourceModel
	var body *DetectionRule

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Process the rule content
	body, err := extractRuleFronJSONStrging(ctx, data.RuleContent.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parser Error", fmt.Sprintf("Unable to parse file, got error: %s", err))
		return
	}

	body.ID = data.Id.ValueString()

	// Create the rule through API
	var response DetectionRuleResponse
	if err := r.client.Put("/detection_engine/rules", body, &response); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: \n%s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DetectionRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DetectionRuleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the rule through the API
	path := fmt.Sprintf("/detection_engine/rules?id=%s", data.Id.ValueString())
	if err := r.client.Delete(path); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: %s", err))
		return
	}
}

func (r *DetectionRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
