package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-elastic-siem/internal/helpers"
	"terraform-provider-elastic-siem/internal/provider/transferobjects"
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
	ExceptionContent types.String `tfsdk:"exception_item_content"`
	ListIdOverride   types.String `tfsdk:"list_id_override"`
	Id               types.String `tfsdk:"id"`
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
			"list_id_override": schema.StringAttribute{
				MarkdownDescription: "The list ID that should be used for the item (overrides id in exception_item_content)",
				Optional:            true,
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
	var body transferobjects.ExceptionItem

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Process the rule content
	err := helpers.ObjectFronJSON(data.ExceptionContent.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError("Parser Error", fmt.Sprintf("Unable to parse file, got error: %s", err))
		return
	}

	if !data.ListIdOverride.IsNull() {
		body.ListID = data.ListIdOverride.ValueString()
	}

	// Create the rule through API
	var response transferobjects.ExceptionItemResponse
	if err := r.client.Post("/exception_lists/items", body, &response, []string{}); err != nil {
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
	var response transferobjects.ExceptionItemResponse
	path := fmt.Sprintf("/exception_lists/items?id=%s", data.Id.ValueString())
	if err := r.client.Get(path, &response); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ExceptionItemResourceModel
	var body *transferobjects.ExceptionItem

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Process the rule content
	err := helpers.ObjectFronJSON(data.ExceptionContent.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError("Parser Error", fmt.Sprintf("Unable to parse file, got error: %s", err))
		return
	}

	body.ID = data.Id.ValueString()

	// Create the rule through API
	var response transferobjects.ExceptionItemResponse
	if err := r.client.Put("/exception_lists/items", body, &response, []string{}); err != nil {
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
	path := fmt.Sprintf("/exception_lists/items?id=%s", data.Id.ValueString())
	if err := r.client.Delete(path); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: %s", err))
		return
	}
}

func (r *ExceptionItemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
