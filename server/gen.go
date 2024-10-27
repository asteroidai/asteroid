//go:build go1.22

// Package sentinel provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package sentinel

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Defines values for LLMMessageRole.
const (
	Assistant LLMMessageRole = "assistant"
	System    LLMMessageRole = "system"
	User      LLMMessageRole = "user"
)

// Defines values for ReviewResultDecision.
const (
	Approve   ReviewResultDecision = "approve"
	Escalate  ReviewResultDecision = "escalate"
	Modify    ReviewResultDecision = "modify"
	Reject    ReviewResultDecision = "reject"
	Terminate ReviewResultDecision = "terminate"
)

// Defines values for ReviewStatusStatus.
const (
	Completed ReviewStatusStatus = "completed"
	Pending   ReviewStatusStatus = "pending"
	Timeout   ReviewStatusStatus = "timeout"
)

// Defines values for SupervisorType.
const (
	All   SupervisorType = "all"
	Code  SupervisorType = "code"
	Human SupervisorType = "human"
	Llm   SupervisorType = "llm"
)

// HubStats defines model for HubStats.
type HubStats struct {
	AssignedReviews    map[string]int `json:"assigned_reviews"`
	BusyClients        int            `json:"busy_clients"`
	CompletedReviews   int            `json:"completed_reviews"`
	ConnectedClients   int            `json:"connected_clients"`
	FreeClients        int            `json:"free_clients"`
	QueuedReviews      int            `json:"queued_reviews"`
	ReviewDistribution map[string]int `json:"review_distribution"`
	StoredReviews      int            `json:"stored_reviews"`
}

// LLMMessage defines model for LLMMessage.
type LLMMessage struct {
	Content string              `json:"content"`
	Id      *openapi_types.UUID `json:"id,omitempty"`
	Role    LLMMessageRole      `json:"role"`
}

// LLMMessageRole defines model for LLMMessage.Role.
type LLMMessageRole string

// Project defines model for Project.
type Project struct {
	CreatedAt time.Time          `json:"created_at"`
	Id        openapi_types.UUID `json:"id"`
	Name      string             `json:"name"`
}

// ProjectCreate defines model for ProjectCreate.
type ProjectCreate struct {
	Name string `json:"name"`
}

// Review defines model for Review.
type Review struct {
	Id        openapi_types.UUID     `json:"id"`
	RunId     openapi_types.UUID     `json:"run_id"`
	Status    *ReviewStatus          `json:"status,omitempty"`
	TaskState map[string]interface{} `json:"task_state"`
}

// ReviewRequest defines model for ReviewRequest.
type ReviewRequest struct {
	Id           *openapi_types.UUID    `json:"id,omitempty"`
	Messages     []LLMMessage           `json:"messages"`
	RunId        openapi_types.UUID     `json:"run_id"`
	TaskState    map[string]interface{} `json:"task_state"`
	ToolRequests []ToolRequest          `json:"tool_requests"`
}

// ReviewResult defines model for ReviewResult.
type ReviewResult struct {
	CreatedAt       time.Time            `json:"created_at"`
	Decision        ReviewResultDecision `json:"decision"`
	Id              openapi_types.UUID   `json:"id"`
	Reasoning       string               `json:"reasoning"`
	ReviewRequestId openapi_types.UUID   `json:"review_request_id"`

	// Toolrequest A tool request is a request to use a tool. It must be approved by a supervisor.
	Toolrequest *ToolRequest `json:"toolrequest,omitempty"`
}

// ReviewResultDecision defines model for ReviewResult.Decision.
type ReviewResultDecision string

// ReviewStatus defines model for ReviewStatus.
type ReviewStatus struct {
	CreatedAt time.Time          `json:"created_at"`
	Id        openapi_types.UUID `json:"id"`
	Status    ReviewStatusStatus `json:"status"`
}

// ReviewStatusStatus defines model for ReviewStatus.Status.
type ReviewStatusStatus string

// Run defines model for Run.
type Run struct {
	CreatedAt time.Time          `json:"created_at"`
	Id        openapi_types.UUID `json:"id"`
	ProjectId openapi_types.UUID `json:"project_id"`
}

// Supervisor defines model for Supervisor.
type Supervisor struct {
	Code        *string             `json:"code,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	Description string              `json:"description"`
	Id          *openapi_types.UUID `json:"id,omitempty"`
	Type        SupervisorType      `json:"type"`
}

// SupervisorAssignment defines model for SupervisorAssignment.
type SupervisorAssignment struct {
	SupervisorId openapi_types.UUID `json:"supervisor_id"`
}

// SupervisorType defines model for SupervisorType.
type SupervisorType string

// Tool defines model for Tool.
type Tool struct {
	Attributes  *map[string]interface{} `json:"attributes,omitempty"`
	CreatedAt   *time.Time              `json:"created_at,omitempty"`
	Description string                  `json:"description"`
	Id          openapi_types.UUID      `json:"id"`
	Name        string                  `json:"name"`
}

// ToolCreate defines model for ToolCreate.
type ToolCreate struct {
	Attributes  *map[string]interface{} `json:"attributes,omitempty"`
	Description string                  `json:"description"`
	Name        string                  `json:"name"`
}

// ToolRequest A tool request is a request to use a tool. It must be approved by a supervisor.
type ToolRequest struct {
	Arguments       map[string]interface{} `json:"arguments"`
	Id              *openapi_types.UUID    `json:"id,omitempty"`
	MessageId       *openapi_types.UUID    `json:"message_id,omitempty"`
	ReviewRequestId *openapi_types.UUID    `json:"review_request_id,omitempty"`
	ToolId          openapi_types.UUID     `json:"tool_id"`
}

// GetReviewsParams defines parameters for GetReviews.
type GetReviewsParams struct {
	Type *SupervisorType `form:"type,omitempty" json:"type,omitempty"`
}

// CreateProjectJSONRequestBody defines body for CreateProject for application/json ContentType.
type CreateProjectJSONRequestBody = ProjectCreate

// CreateRunToolJSONRequestBody defines body for CreateRunTool for application/json ContentType.
type CreateRunToolJSONRequestBody = ToolCreate

// AssignSupervisorToToolJSONRequestBody defines body for AssignSupervisorToTool for application/json ContentType.
type AssignSupervisorToToolJSONRequestBody = SupervisorAssignment

// CreateReviewJSONRequestBody defines body for CreateReview for application/json ContentType.
type CreateReviewJSONRequestBody = ReviewRequest

// CreateSupervisorJSONRequestBody defines body for CreateSupervisor for application/json ContentType.
type CreateSupervisorJSONRequestBody = Supervisor

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get the OpenAPI schema
	// (GET /api/openapi.yaml)
	GetOpenAPI(w http.ResponseWriter, r *http.Request)
	// List all projects
	// (GET /api/projects)
	GetProjects(w http.ResponseWriter, r *http.Request)
	// Create a new project
	// (POST /api/projects)
	CreateProject(w http.ResponseWriter, r *http.Request)
	// Get project by ID
	// (GET /api/projects/{projectId})
	GetProject(w http.ResponseWriter, r *http.Request, projectId openapi_types.UUID)
	// Get runs for a project
	// (GET /api/projects/{projectId}/runs)
	GetProjectRuns(w http.ResponseWriter, r *http.Request, projectId openapi_types.UUID)
	// Create a new run for a project
	// (POST /api/projects/{projectId}/runs)
	CreateRun(w http.ResponseWriter, r *http.Request, projectId openapi_types.UUID)
	// Get run by ID
	// (GET /api/projects/{projectId}/runs/{runId})
	GetRun(w http.ResponseWriter, r *http.Request, projectId openapi_types.UUID, runId openapi_types.UUID)
	// List all tools
	// (GET /api/projects/{projectId}/runs/{runId}/tools)
	GetRunTools(w http.ResponseWriter, r *http.Request, projectId openapi_types.UUID, runId openapi_types.UUID)
	// Create a new tool
	// (POST /api/projects/{projectId}/runs/{runId}/tools)
	CreateRunTool(w http.ResponseWriter, r *http.Request, projectId openapi_types.UUID, runId openapi_types.UUID)
	// Get tool by ID
	// (GET /api/projects/{projectId}/tools/{toolId})
	GetTool(w http.ResponseWriter, r *http.Request, projectId openapi_types.UUID, toolId openapi_types.UUID)
	// Get supervisors assigned to a tool
	// (GET /api/projects/{projectId}/tools/{toolId}/supervisors)
	GetToolSupervisors(w http.ResponseWriter, r *http.Request, projectId openapi_types.UUID, toolId openapi_types.UUID)
	// Assign supervisor to tool
	// (POST /api/projects/{projectId}/tools/{toolId}/supervisors)
	AssignSupervisorToTool(w http.ResponseWriter, r *http.Request, projectId openapi_types.UUID, toolId openapi_types.UUID)
	// List all reviews
	// (GET /api/reviews)
	GetReviews(w http.ResponseWriter, r *http.Request, params GetReviewsParams)
	// Create a review request
	// (POST /api/reviews)
	CreateReview(w http.ResponseWriter, r *http.Request)
	// Get review by ID
	// (GET /api/reviews/{reviewId})
	GetReview(w http.ResponseWriter, r *http.Request, reviewId openapi_types.UUID)
	// Get review results
	// (GET /api/reviews/{reviewId}/results)
	GetReviewResults(w http.ResponseWriter, r *http.Request, reviewId openapi_types.UUID)
	// Get review status
	// (GET /api/reviews/{reviewId}/status)
	GetReviewStatus(w http.ResponseWriter, r *http.Request, reviewId openapi_types.UUID)
	// Get tool requests for a review
	// (GET /api/reviews/{reviewId}/toolrequests)
	GetReviewToolRequests(w http.ResponseWriter, r *http.Request, reviewId openapi_types.UUID)
	// Get hub stats
	// (GET /api/stats)
	GetHubStats(w http.ResponseWriter, r *http.Request)
	// List all supervisors
	// (GET /api/supervisors)
	GetSupervisors(w http.ResponseWriter, r *http.Request)
	// Create a new supervisor
	// (POST /api/supervisors)
	CreateSupervisor(w http.ResponseWriter, r *http.Request)
	// Get supervisor by ID
	// (GET /api/supervisors/{supervisorId})
	GetSupervisor(w http.ResponseWriter, r *http.Request, supervisorId openapi_types.UUID)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetOpenAPI operation middleware
func (siw *ServerInterfaceWrapper) GetOpenAPI(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetOpenAPI(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetProjects operation middleware
func (siw *ServerInterfaceWrapper) GetProjects(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetProjects(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CreateProject operation middleware
func (siw *ServerInterfaceWrapper) CreateProject(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateProject(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetProject operation middleware
func (siw *ServerInterfaceWrapper) GetProject(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "projectId" -------------
	var projectId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "projectId", r.PathValue("projectId"), &projectId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "projectId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetProject(w, r, projectId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetProjectRuns operation middleware
func (siw *ServerInterfaceWrapper) GetProjectRuns(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "projectId" -------------
	var projectId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "projectId", r.PathValue("projectId"), &projectId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "projectId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetProjectRuns(w, r, projectId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CreateRun operation middleware
func (siw *ServerInterfaceWrapper) CreateRun(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "projectId" -------------
	var projectId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "projectId", r.PathValue("projectId"), &projectId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "projectId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateRun(w, r, projectId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetRun operation middleware
func (siw *ServerInterfaceWrapper) GetRun(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "projectId" -------------
	var projectId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "projectId", r.PathValue("projectId"), &projectId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "projectId", Err: err})
		return
	}

	// ------------- Path parameter "runId" -------------
	var runId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "runId", r.PathValue("runId"), &runId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "runId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetRun(w, r, projectId, runId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetRunTools operation middleware
func (siw *ServerInterfaceWrapper) GetRunTools(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "projectId" -------------
	var projectId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "projectId", r.PathValue("projectId"), &projectId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "projectId", Err: err})
		return
	}

	// ------------- Path parameter "runId" -------------
	var runId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "runId", r.PathValue("runId"), &runId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "runId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetRunTools(w, r, projectId, runId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CreateRunTool operation middleware
func (siw *ServerInterfaceWrapper) CreateRunTool(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "projectId" -------------
	var projectId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "projectId", r.PathValue("projectId"), &projectId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "projectId", Err: err})
		return
	}

	// ------------- Path parameter "runId" -------------
	var runId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "runId", r.PathValue("runId"), &runId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: false})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "runId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateRunTool(w, r, projectId, runId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetTool operation middleware
func (siw *ServerInterfaceWrapper) GetTool(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "projectId" -------------
	var projectId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "projectId", r.PathValue("projectId"), &projectId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "projectId", Err: err})
		return
	}

	// ------------- Path parameter "toolId" -------------
	var toolId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "toolId", r.PathValue("toolId"), &toolId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "toolId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetTool(w, r, projectId, toolId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetToolSupervisors operation middleware
func (siw *ServerInterfaceWrapper) GetToolSupervisors(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "projectId" -------------
	var projectId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "projectId", r.PathValue("projectId"), &projectId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "projectId", Err: err})
		return
	}

	// ------------- Path parameter "toolId" -------------
	var toolId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "toolId", r.PathValue("toolId"), &toolId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "toolId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetToolSupervisors(w, r, projectId, toolId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// AssignSupervisorToTool operation middleware
func (siw *ServerInterfaceWrapper) AssignSupervisorToTool(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "projectId" -------------
	var projectId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "projectId", r.PathValue("projectId"), &projectId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "projectId", Err: err})
		return
	}

	// ------------- Path parameter "toolId" -------------
	var toolId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "toolId", r.PathValue("toolId"), &toolId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "toolId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.AssignSupervisorToTool(w, r, projectId, toolId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetReviews operation middleware
func (siw *ServerInterfaceWrapper) GetReviews(w http.ResponseWriter, r *http.Request) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetReviewsParams

	// ------------- Optional query parameter "type" -------------

	err = runtime.BindQueryParameter("form", true, false, "type", r.URL.Query(), &params.Type)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "type", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetReviews(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CreateReview operation middleware
func (siw *ServerInterfaceWrapper) CreateReview(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateReview(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetReview operation middleware
func (siw *ServerInterfaceWrapper) GetReview(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "reviewId" -------------
	var reviewId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "reviewId", r.PathValue("reviewId"), &reviewId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "reviewId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetReview(w, r, reviewId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetReviewResults operation middleware
func (siw *ServerInterfaceWrapper) GetReviewResults(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "reviewId" -------------
	var reviewId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "reviewId", r.PathValue("reviewId"), &reviewId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "reviewId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetReviewResults(w, r, reviewId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetReviewStatus operation middleware
func (siw *ServerInterfaceWrapper) GetReviewStatus(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "reviewId" -------------
	var reviewId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "reviewId", r.PathValue("reviewId"), &reviewId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "reviewId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetReviewStatus(w, r, reviewId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetReviewToolRequests operation middleware
func (siw *ServerInterfaceWrapper) GetReviewToolRequests(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "reviewId" -------------
	var reviewId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "reviewId", r.PathValue("reviewId"), &reviewId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "reviewId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetReviewToolRequests(w, r, reviewId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetHubStats operation middleware
func (siw *ServerInterfaceWrapper) GetHubStats(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHubStats(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetSupervisors operation middleware
func (siw *ServerInterfaceWrapper) GetSupervisors(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetSupervisors(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CreateSupervisor operation middleware
func (siw *ServerInterfaceWrapper) CreateSupervisor(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateSupervisor(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetSupervisor operation middleware
func (siw *ServerInterfaceWrapper) GetSupervisor(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "supervisorId" -------------
	var supervisorId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "supervisorId", r.PathValue("supervisorId"), &supervisorId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "supervisorId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetSupervisor(w, r, supervisorId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

// ServeMux is an abstraction of http.ServeMux.
type ServeMux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("GET "+options.BaseURL+"/api/openapi.yaml", wrapper.GetOpenAPI)
	m.HandleFunc("GET "+options.BaseURL+"/api/projects", wrapper.GetProjects)
	m.HandleFunc("POST "+options.BaseURL+"/api/projects", wrapper.CreateProject)
	m.HandleFunc("GET "+options.BaseURL+"/api/projects/{projectId}", wrapper.GetProject)
	m.HandleFunc("GET "+options.BaseURL+"/api/projects/{projectId}/runs", wrapper.GetProjectRuns)
	m.HandleFunc("POST "+options.BaseURL+"/api/projects/{projectId}/runs", wrapper.CreateRun)
	m.HandleFunc("GET "+options.BaseURL+"/api/projects/{projectId}/runs/{runId}", wrapper.GetRun)
	m.HandleFunc("GET "+options.BaseURL+"/api/projects/{projectId}/runs/{runId}/tools", wrapper.GetRunTools)
	m.HandleFunc("POST "+options.BaseURL+"/api/projects/{projectId}/runs/{runId}/tools", wrapper.CreateRunTool)
	m.HandleFunc("GET "+options.BaseURL+"/api/projects/{projectId}/tools/{toolId}", wrapper.GetTool)
	m.HandleFunc("GET "+options.BaseURL+"/api/projects/{projectId}/tools/{toolId}/supervisors", wrapper.GetToolSupervisors)
	m.HandleFunc("POST "+options.BaseURL+"/api/projects/{projectId}/tools/{toolId}/supervisors", wrapper.AssignSupervisorToTool)
	m.HandleFunc("GET "+options.BaseURL+"/api/reviews", wrapper.GetReviews)
	m.HandleFunc("POST "+options.BaseURL+"/api/reviews", wrapper.CreateReview)
	m.HandleFunc("GET "+options.BaseURL+"/api/reviews/{reviewId}", wrapper.GetReview)
	m.HandleFunc("GET "+options.BaseURL+"/api/reviews/{reviewId}/results", wrapper.GetReviewResults)
	m.HandleFunc("GET "+options.BaseURL+"/api/reviews/{reviewId}/status", wrapper.GetReviewStatus)
	m.HandleFunc("GET "+options.BaseURL+"/api/reviews/{reviewId}/toolrequests", wrapper.GetReviewToolRequests)
	m.HandleFunc("GET "+options.BaseURL+"/api/stats", wrapper.GetHubStats)
	m.HandleFunc("GET "+options.BaseURL+"/api/supervisors", wrapper.GetSupervisors)
	m.HandleFunc("POST "+options.BaseURL+"/api/supervisors", wrapper.CreateSupervisor)
	m.HandleFunc("GET "+options.BaseURL+"/api/supervisors/{supervisorId}", wrapper.GetSupervisor)

	return m
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
