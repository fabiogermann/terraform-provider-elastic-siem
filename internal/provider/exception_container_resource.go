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
var _ resource.Resource = &ExceptionContainerResource{}
var _ resource.ResourceWithImportState = &ExceptionContainerResource{}

func NewExceptionContainerResource() resource.Resource {
	return &ExceptionContainerResource{}
}

// ExceptionContainerResource defines the resource implementation.
type ExceptionContainerResource struct {
	client *helpers.Client
}

// ExceptionContainerResourceModel describes the resource data model.
type ExceptionContainerResourceModel struct {
	RuleContent types.String `tfsdk:"exception_item_content"`
	Id          types.String `tfsdk:"id"`
}

type ExceptionContainerResponse struct {
	Tags          []interface{} `json:"_tags,omitempty"`
	CreatedAt     time.Time     `json:"created_at,omitempty"`
	CreatedBy     string        `json:"created_by,omitempty"`
	Description   string        `json:"description,omitempty"`
	ID            string        `json:"id,omitempty"`
	ListID        string        `json:"list_id,omitempty"`
	Name          string        `json:"name,omitempty"`
	NamespaceType string        `json:"namespace_type,omitempty"`
	Tags0         []string      `json:"tags,omitempty"`
	TieBreakerID  string        `json:"tie_breaker_id,omitempty"`
	Type          string        `json:"type,omitempty"`
	UpdatedAt     time.Time     `json:"updated_at,omitempty"`
	UpdatedBy     string        `json:"updated_by,omitempty"`
}

type ExceptionContainer struct {
	ID            string   `json:"id,omitempty"`
	Description   string   `json:"description,omitempty"`
	Name          string   `json:"name,omitempty"`
	ListID        string   `json:"list_id,omitempty"`
	Type          string   `json:"typ,omitemptye"`
	NamespaceType string   `json:"namespace_type,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

func extractExceptionContainerFronJSONStrging(ctx context.Context, yamlString string) (*ExceptionContainer, error) {
	result := ExceptionContainer{}
	tflog.Debug(ctx, yamlString)
	s, err := strconv.Unquote(yamlString)
	if err != nil {
		tflog.Error(ctx, "Error in extractExceptionContainerFronJSONStrging (0)")
		return nil, err
	}
	err = json.Unmarshal([]byte(s), &result)
	if err != nil {
		tflog.Error(ctx, "Error in extractExceptionContainerFronJSONStrging (1)")
		return nil, err
	}
	tflog.Debug(ctx, result.Name)
	return &result, nil
}

func (r *ExceptionContainerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exception_item"
}

func (r *ExceptionContainerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

func (r *ExceptionContainerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ExceptionContainerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ExceptionContainerResourceModel
	var body *ExceptionContainer

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Process the rule content
	body, err := extractExceptionContainerFronJSONStrging(ctx, data.RuleContent.String())
	if err != nil {
		resp.Diagnostics.AddError("Parser Error", fmt.Sprintf("Unable to parse file, got error: %s", err))
		return
	}

	// Create the rule through API
	var response ExceptionContainerResponse
	if err := r.client.Post("/exception_lists", body, &response); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: \n%s", err))
		return
	}

	// Save id into the Terraform state.
	data.Id = types.StringValue(response.ID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionContainerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ExceptionContainerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the rule through the API
	var response ExceptionContainerResponse
	path := fmt.Sprintf("/exception_lists?id=%s", data.Id)
	if err := r.client.Get(path, &response); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionContainerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ExceptionContainerResourceModel
	var body *ExceptionContainer

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Process the rule content
	body, err := extractExceptionContainerFronJSONStrging(ctx, data.RuleContent.String())
	if err != nil {
		resp.Diagnostics.AddError("Parser Error", fmt.Sprintf("Unable to parse file, got error: %s", err))
		return
	}

	body.ID = data.Id.String()

	// Create the rule through API
	var response ExceptionContainerResponse
	if err := r.client.Put("/exception_lists", body, &response); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: \n%s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionContainerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ExceptionContainerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the rule through the API
	path := fmt.Sprintf("/exception_lists?id=%s", data.Id)
	if err := r.client.Delete(path); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: %s", err))
		return
	}
}

func (r *ExceptionContainerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
