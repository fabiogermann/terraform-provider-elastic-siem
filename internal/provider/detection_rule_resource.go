package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
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

type DetectionRuleResponse struct {
	CreatedAt           time.Time     `json:"created_at,omitempty"`
	UpdatedAt           time.Time     `json:"updated_at,omitempty"`
	CreatedBy           string        `json:"created_by,omitempty"`
	Description         string        `json:"description,omitempty"`
	Enabled             bool          `json:"enabled,omitempty"`
	FalsePositives      []interface{} `json:"false_positives,omitempty"`
	From                string        `json:"from,omitempty"`
	ID                  string        `json:"id,omitempty"`
	Immutable           bool          `json:"immutable,omitempty"`
	Interval            string        `json:"interval,omitempty"`
	RuleID              string        `json:"rule_id,omitempty"`
	OutputIndex         string        `json:"output_index,omitempty"`
	MaxSignals          int           `json:"max_signals,omitempty"`
	RiskScore           int           `json:"risk_score,omitempty"`
	Name                string        `json:"name,omitempty"`
	References          []interface{} `json:"references,omitempty"`
	Severity            string        `json:"severity,omitempty"`
	UpdatedBy           string        `json:"updated_by,omitempty"`
	Tags                []string      `json:"tags,omitempty"`
	To                  string        `json:"to,omitempty"`
	Type                string        `json:"type,omitempty"`
	Threat              []interface{} `json:"threat,omitempty"`
	Version             int           `json:"version,omitempty"`
	Actions             []interface{} `json:"actions,omitempty"`
	Filters             []interface{} `json:"filters,omitempty"`
	Throttle            string        `json:"throttle,omitempty"`
	Query               string        `json:"query,omitempty"`
	Language            string        `json:"language,omitempty"`
	RelatedIntegrations []interface{} `json:"related_integrations,omitempty"`
	RequiredFields      []interface{} `json:"required_fields,omitempty"`
	ExecutionSummary    []interface{} `json:"execution_summary,omitempty"`
	Setup               string        `json:"setup,omitempty"`
}

type DetectionRule struct {
	Description         string        `json:"description,omitempty"`
	Enabled             bool          `json:"enabled,omitempty"`
	FalsePositives      []interface{} `json:"false_positives,omitempty"`
	From                string        `json:"from,omitempty"`
	ID                  string        `json:"id,omitempty"`
	Immutable           bool          `json:"immutable,omitempty"`
	Interval            string        `json:"interval,omitempty"`
	RuleID              string        `json:"rule_id,omitempty"`
	OutputIndex         string        `json:"output_index,omitempty"`
	MaxSignals          int           `json:"max_signals,omitempty"`
	RiskScore           int           `json:"risk_score,omitempty"`
	Name                string        `json:"name,omitempty"`
	References          []interface{} `json:"references,omitempty"`
	Severity            string        `json:"severity,omitempty"`
	UpdatedBy           string        `json:"updated_by,omitempty"`
	Tags                []string      `json:"tags,omitempty"`
	To                  string        `json:"to,omitempty"`
	Type                string        `json:"type,omitempty"`
	Threat              []interface{} `json:"threat,omitempty"`
	Version             int           `json:"version,omitempty"`
	Actions             []interface{} `json:"actions,omitempty"`
	Filters             []interface{} `json:"filters,omitempty"`
	Throttle            string        `json:"throttle,omitempty"`
	Query               string        `json:"query,omitempty"`
	Language            string        `json:"language,omitempty"`
	RelatedIntegrations []interface{} `json:"related_integrations,omitempty"`
	RequiredFields      []interface{} `json:"required_fields,omitempty"`
	ExecutionSummary    []interface{} `json:"execution_summary,omitempty"`
	Setup               string        `json:"setup,omitempty"`
}

func extractRuleFronYAMLStrgin(ctx context.Context, yamlString string) (*DetectionRule, error) {
	result := DetectionRule{}
	tflog.Debug(ctx, yamlString)
	s, err := strconv.Unquote(yamlString)
	if err != nil {
		tflog.Error(ctx, "Error in extractRuleFronYAMLStrgin (0)")
		return nil, err
	}
	err = json.Unmarshal([]byte(s), &result)
	if err != nil {
		tflog.Error(ctx, "Error in extractRuleFronYAMLStrgin (1)")
		return nil, err
	}
	tflog.Debug(ctx, result.Name)
	return &result, nil
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
				MarkdownDescription: "The content of the rule",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Rule identifier",
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
	body, err := extractRuleFronYAMLStrgin(ctx, data.RuleContent.String())
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
	path := fmt.Sprintf("/detection_engine/rules?id=%s", data.Id)
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
	body, err := extractRuleFronYAMLStrgin(ctx, data.RuleContent.String())
	if err != nil {
		resp.Diagnostics.AddError("Parser Error", fmt.Sprintf("Unable to parse file, got error: %s", err))
		return
	}

	body.ID = data.Id.String()

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
	path := fmt.Sprintf("/detection_engine/rules?id=%s", data.Id)
	if err := r.client.Delete(path); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: %s", err))
		return
	}
}

func (r *DetectionRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
