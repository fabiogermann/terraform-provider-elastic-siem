package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"terraform-provider-elastic-siem/internal/helpers"
	"terraform-provider-elastic-siem/internal/provider/transferobjects"
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
	Description   types.String `tfsdk:"description"`
	Name          types.String `tfsdk:"name"`
	ListId        types.String `tfsdk:"list_id"`
	Type          types.String `tfsdk:"type"`
	NamespaceType types.String `tfsdk:"namespace_type"`
	Tags          types.List   `tfsdk:"tags"`
	Id            types.String `tfsdk:"id"`
}

func (r *ExceptionContainerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exception_container"
}

func (r *ExceptionContainerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Exception container resource",

		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the exception container",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the exception container",
				Required:            true,
			},
			"list_id": schema.StringAttribute{
				MarkdownDescription: "The list id of the exception container (referenced in rule and item)",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the exception container",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("detection", "endpoint")},
			},
			"namespace_type": schema.StringAttribute{
				MarkdownDescription: "The namespace type of the exception container",
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("single", "agnostic")},
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "The tags of the exception container",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Exception container identifier (in UUID format)",
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
	var body *transferobjects.ExceptionContainer

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body = &transferobjects.ExceptionContainer{
		ListID:        data.ListId.ValueString(),
		Name:          data.Name.ValueString(),
		NamespaceType: data.NamespaceType.ValueString(),
		Description:   data.Description.ValueString(),
		Type:          data.Type.ValueString(),
		Tags:          []string{},
	}
	for _, element := range data.Tags.Elements() {
		s, err := strconv.Unquote(element.String())
		if err != nil {
			body.Tags = append(body.Tags, element.String())
		}
		body.Tags = append(body.Tags, s)
	}

	// Create the rule through API
	var response transferobjects.ExceptionContainerResponse
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
	var response transferobjects.ExceptionContainerResponse
	apiPath := fmt.Sprintf("/exception_lists?id=%s", data.Id.ValueString())
	if err := r.client.Get(apiPath, &response); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExceptionContainerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ExceptionContainerResourceModel
	var body *transferobjects.ExceptionContainer

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body = &transferobjects.ExceptionContainer{
		ListID:        data.ListId.ValueString(),
		Name:          data.Name.ValueString(),
		NamespaceType: data.NamespaceType.ValueString(),
		Description:   data.Description.ValueString(),
		Type:          data.Type.ValueString(),
		Tags:          []string{},
		ID:            data.Id.ValueString(),
	}
	for _, element := range data.Tags.Elements() {
		s, err := strconv.Unquote(element.String())
		if err != nil {
			body.Tags = append(body.Tags, element.String())
		}
		body.Tags = append(body.Tags, s)
	}

	// Create the rule through API
	var response transferobjects.ExceptionContainerResponse
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
	apiPath := fmt.Sprintf("/exception_lists?id=%s", data.Id.ValueString())
	if err := r.client.Delete(apiPath); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: %s", err))
		return
	}
}

func (r *ExceptionContainerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
