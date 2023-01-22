package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
	"terraform-provider-elastic-siem/internal/helpers"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ExceptionItemResource{}
var _ resource.ResourceWithImportState = &ExceptionItemResource{}

func NewExceptionItemResource() resource.Resource {
	return &ExceptionItemResource{}
}

// ExceptionItemResource defines the resource implementation.
type ExceptionItemResource struct {
	client *helpers.Client
}

// ExceptionItemResourceModel describes the resource data model.
type ExceptionItemResourceModel struct {
	RuleContent types.String `tfsdk:"exception_item_content"`
	Id          types.String `tfsdk:"id"`
}

type ExceptionItemResponse struct {
	ExceptionItemBase
	HTags        []interface{}               `json:"_tags,omitempty"`
	Comments     []ExceptionCommentsResponse `json:"comments,omitempty"`
	CreatedAt    time.Time                   `json:"created_at,omitempty"`
	CreatedBy    string                      `json:"created_by,omitempty"`
	TieBreakerID string                      `json:"tie_breaker_id,omitempty"`
	UpdatedAt    time.Time                   `json:"updated_a,omitemptyt"`
	UpdatedBy    string                      `json:"updated_by,omitempty"`
}

type ExceptionComments struct {
	Comment string `json:"comment,omitempty"`
}

type ExceptionCommentsResponse struct {
	ExceptionComments
	CreatedAt time.Time `json:"created_at,omitempty"`
	CreatedBy string    `json:"created_by,omitempty"`
	ID        string    `json:"id,omitempty"`
}

type ExceptionItemBase struct {
	Description string `json:"description,omitempty"`
	Entries     []struct {
		Field    string `json:"field,omitempty"`
		Operator string `json:"operator,omitempty"`
		Type     string `json:"type,omitempty"`
		Value    string `json:"value,omitempty"`
	} `json:"entries,omitempty"`
	ID            string   `json:"id,omitempty"`
	ListID        string   `json:"list_id,omitempty"`
	ItemID        string   `json:"item_id,omitempty"`
	Name          string   `json:"name,omitempty"`
	NamespaceType string   `json:"namespace_type,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Type          string   `json:"type,omitempty"`
}

type ExceptionItem struct {
	ExceptionItemBase
	Comments []ExceptionComments `json:"comments,omitempty"`
}

func extractExceptionItemFronJSONStrging(ctx context.Context, yamlString string) (*ExceptionItem, error) {
	result := ExceptionItem{}
	tflog.Debug(ctx, yamlString)
	s, err := strconv.Unquote(yamlString)
	if err != nil {
		tflog.Error(ctx, "Error in extractExceptionItemFronJSONStrging (0)")
		return nil, err
	}
	err = json.Unmarshal([]byte(s), &result)
	if err != nil {
		tflog.Error(ctx, "Error in extractExceptionItemFronJSONStrging (1)")
		return nil, err
	}
	tflog.Debug(ctx, result.Name)
	return &result, nil
}

func (r *ExceptionItemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exception_item"
}

func (r *ExceptionItemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Exception item resource",

		Attributes: map[string]schema.Attribute{
			"exception_item_content": schema.StringAttribute{
				MarkdownDescription: "The content of the exception item (JSON encoded string)",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Exception item identifier (in UUID format)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ExceptionItemResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ExceptionItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ExceptionItemResourceModel
	var body *ExceptionItem

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Process the rule content
	body, err := extractExceptionItemFronJSONStrging(ctx, data.RuleContent.String())
	if err != nil {
		resp.Diagnostics.AddError("Parser Error", fmt.Sprintf("Unable to parse file, got error: %s", err))
		return
	}

	// Create the rule through API
	var response ExceptionItemResponse
	if err := r.client.Post("/exception_lists/items", body, &response); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: \n%s", err))
		return
	}

	// Save id into the Terraform state.
	data.Id = types.StringValue(response.ID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ExceptionItemResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the rule through the API
	var response ExceptionItemResponse
	path := fmt.Sprintf("/exception_lists/items?id=%s", data.Id)
	if err := r.client.Get(path, &response); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ExceptionItemResourceModel
	var body *ExceptionItem

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Process the rule content
	body, err := extractExceptionItemFronJSONStrging(ctx, data.RuleContent.String())
	if err != nil {
		resp.Diagnostics.AddError("Parser Error", fmt.Sprintf("Unable to parse file, got error: %s", err))
		return
	}

	body.ID = data.Id.String()

	// Create the rule through API
	var response ExceptionItemResponse
	if err := r.client.Put("/exception_lists/items", body, &response); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: \n%s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ExceptionItemResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the rule through the API
	path := fmt.Sprintf("/exception_lists/items?id=%s", data.Id)
	if err := r.client.Delete(path); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: %s", err))
		return
	}
}

func (r *ExceptionItemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
