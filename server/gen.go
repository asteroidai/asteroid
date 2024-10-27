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
	Assigned  ReviewStatusStatus = "assigned"
	Completed ReviewStatusStatus = "completed"
	Pending   ReviewStatusStatus = "pending"
)

// Defines values for SupervisorType.
const (
	All   SupervisorType = "all"
	Code  SupervisorType = "code"
	Human SupervisorType = "human"
	Llm   SupervisorType = "llm"
)

// CreateReviewResult defines model for CreateReviewResult.
type CreateReviewResult struct {
	ReviewResult ReviewResult       `json:"review_result"`
	RunId        openapi_types.UUID `json:"run_id"`
	SupervisorId openapi_types.UUID `json:"supervisor_id"`
	ToolId       openapi_types.UUID `json:"tool_id"`

	// ToolRequest A tool request is a request to use a tool. It must be approved by a supervisor.
	ToolRequest ToolRequest `json:"tool_request"`
}

// HubStats defines model for HubStats.
type HubStats struct {
	AssignedReviews       map[string]int `json:"assigned_reviews"`
	AssignedReviewsCount  int            `json:"assigned_reviews_count"`
	BusyClients           int            `json:"busy_clients"`
	CompletedReviewsCount int            `json:"completed_reviews_count"`
	ConnectedClients      int            `json:"connected_clients"`
	FreeClients           int            `json:"free_clients"`
	PendingReviewsCount   int            `json:"pending_reviews_count"`
	ReviewDistribution    map[string]int `json:"review_distribution"`
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

// ReviewRequest defines model for ReviewRequest.
type ReviewRequest struct {
	Id           *openapi_types.UUID    `json:"id,omitempty"`
	Messages     []LLMMessage           `json:"messages"`
	RunId        openapi_types.UUID     `json:"run_id"`
	Status       *ReviewStatus          `json:"status,omitempty"`
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

// ToolRequest A tool request is a request to use a tool. It must be approved by a supervisor.
type ToolRequest struct {
	Arguments       map[string]interface{} `json:"arguments"`
	Id              *openapi_types.UUID    `json:"id,omitempty"`
	MessageId       *openapi_types.UUID    `json:"message_id,omitempty"`
	ReviewRequestId *openapi_types.UUID    `json:"review_request_id,omitempty"`
	ToolId          openapi_types.UUID     `json:"tool_id"`
}

// GetReviewRequestsParams defines parameters for GetReviewRequests.
type GetReviewRequestsParams struct {
	Type *SupervisorType `form:"type,omitempty" json:"type,omitempty"`
}

// CreateRunToolSupervisorsJSONBody defines parameters for CreateRunToolSupervisors.
type CreateRunToolSupervisorsJSONBody = []openapi_types.UUID

// CreateProjectJSONRequestBody defines body for CreateProject for application/json ContentType.
type CreateProjectJSONRequestBody = ProjectCreate

// CreateReviewRequestJSONRequestBody defines body for CreateReviewRequest for application/json ContentType.
type CreateReviewRequestJSONRequestBody = ReviewRequest

// CreateReviewResultJSONRequestBody defines body for CreateReviewResult for application/json ContentType.
type CreateReviewResultJSONRequestBody = CreateReviewResult

// CreateRunToolSupervisorsJSONRequestBody defines body for CreateRunToolSupervisors for application/json ContentType.
type CreateRunToolSupervisorsJSONRequestBody = CreateRunToolSupervisorsJSONBody

// CreateSupervisorJSONRequestBody defines body for CreateSupervisor for application/json ContentType.
type CreateSupervisorJSONRequestBody = Supervisor

// CreateToolJSONRequestBody defines body for CreateTool for application/json ContentType.
type CreateToolJSONRequestBody = Tool

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
	// List all review requests
	// (GET /api/reviews)
	GetReviewRequests(w http.ResponseWriter, r *http.Request, params GetReviewRequestsParams)
	// Create a review request
	// (POST /api/reviews)
	CreateReviewRequest(w http.ResponseWriter, r *http.Request)
	// Get review request by ID
	// (GET /api/reviews/{reviewId})
	GetReviewRequest(w http.ResponseWriter, r *http.Request, reviewId openapi_types.UUID)
	// Get review results
	// (GET /api/reviews/{reviewId}/results)
	GetReviewResults(w http.ResponseWriter, r *http.Request, reviewId openapi_types.UUID)
	// Create a review result
	// (POST /api/reviews/{reviewId}/results)
	CreateReviewResult(w http.ResponseWriter, r *http.Request, reviewId openapi_types.UUID)
	// Get review status
	// (GET /api/reviews/{reviewId}/status)
	GetReviewStatus(w http.ResponseWriter, r *http.Request, reviewId openapi_types.UUID)
	// Get tool requests for a review
	// (GET /api/reviews/{reviewId}/toolrequests)
	GetReviewToolRequests(w http.ResponseWriter, r *http.Request, reviewId openapi_types.UUID)
	// Get run by ID
	// (GET /api/runs/{runId})
	GetRun(w http.ResponseWriter, r *http.Request, runId openapi_types.UUID)
	// Get tools for a run
	// (GET /api/runs/{runId}/tools)
	GetRunTools(w http.ResponseWriter, r *http.Request, runId openapi_types.UUID)
	// Get the supervisors assigned to a tool
	// (GET /api/runs/{runId}/tools/{toolId}/supervisors)
	GetRunToolSupervisors(w http.ResponseWriter, r *http.Request, runId openapi_types.UUID, toolId openapi_types.UUID)
	// Assign a list of supervisors to a tool for a given run
	// (POST /api/runs/{runId}/tools/{toolId}/supervisors)
	CreateRunToolSupervisors(w http.ResponseWriter, r *http.Request, runId openapi_types.UUID, toolId openapi_types.UUID)
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
	// List all tools
	// (GET /api/tools)
	GetTools(w http.ResponseWriter, r *http.Request)
	// Create a new tool
	// (POST /api/tools)
	CreateTool(w http.ResponseWriter, r *http.Request)
	// Get tool by ID
	// (GET /api/tools/{toolId})
	GetTool(w http.ResponseWriter, r *http.Request, toolId openapi_types.UUID)
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

// GetReviewRequests operation middleware
func (siw *ServerInterfaceWrapper) GetReviewRequests(w http.ResponseWriter, r *http.Request) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetReviewRequestsParams

	// ------------- Optional query parameter "type" -------------

	err = runtime.BindQueryParameter("form", true, false, "type", r.URL.Query(), &params.Type)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "type", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetReviewRequests(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CreateReviewRequest operation middleware
func (siw *ServerInterfaceWrapper) CreateReviewRequest(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateReviewRequest(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetReviewRequest operation middleware
func (siw *ServerInterfaceWrapper) GetReviewRequest(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "reviewId" -------------
	var reviewId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "reviewId", r.PathValue("reviewId"), &reviewId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "reviewId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetReviewRequest(w, r, reviewId)
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

// CreateReviewResult operation middleware
func (siw *ServerInterfaceWrapper) CreateReviewResult(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "reviewId" -------------
	var reviewId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "reviewId", r.PathValue("reviewId"), &reviewId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "reviewId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateReviewResult(w, r, reviewId)
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

// GetRun operation middleware
func (siw *ServerInterfaceWrapper) GetRun(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "runId" -------------
	var runId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "runId", r.PathValue("runId"), &runId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "runId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetRun(w, r, runId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetRunTools operation middleware
func (siw *ServerInterfaceWrapper) GetRunTools(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "runId" -------------
	var runId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "runId", r.PathValue("runId"), &runId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "runId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetRunTools(w, r, runId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetRunToolSupervisors operation middleware
func (siw *ServerInterfaceWrapper) GetRunToolSupervisors(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "runId" -------------
	var runId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "runId", r.PathValue("runId"), &runId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "runId", Err: err})
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
		siw.Handler.GetRunToolSupervisors(w, r, runId, toolId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CreateRunToolSupervisors operation middleware
func (siw *ServerInterfaceWrapper) CreateRunToolSupervisors(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "runId" -------------
	var runId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "runId", r.PathValue("runId"), &runId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: false})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "runId", Err: err})
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
		siw.Handler.CreateRunToolSupervisors(w, r, runId, toolId)
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

// GetTools operation middleware
func (siw *ServerInterfaceWrapper) GetTools(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetTools(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CreateTool operation middleware
func (siw *ServerInterfaceWrapper) CreateTool(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateTool(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetTool operation middleware
func (siw *ServerInterfaceWrapper) GetTool(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "toolId" -------------
	var toolId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "toolId", r.PathValue("toolId"), &toolId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "toolId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetTool(w, r, toolId)
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
	m.HandleFunc("GET "+options.BaseURL+"/api/reviews", wrapper.GetReviewRequests)
	m.HandleFunc("POST "+options.BaseURL+"/api/reviews", wrapper.CreateReviewRequest)
	m.HandleFunc("GET "+options.BaseURL+"/api/reviews/{reviewId}", wrapper.GetReviewRequest)
	m.HandleFunc("GET "+options.BaseURL+"/api/reviews/{reviewId}/results", wrapper.GetReviewResults)
	m.HandleFunc("POST "+options.BaseURL+"/api/reviews/{reviewId}/results", wrapper.CreateReviewResult)
	m.HandleFunc("GET "+options.BaseURL+"/api/reviews/{reviewId}/status", wrapper.GetReviewStatus)
	m.HandleFunc("GET "+options.BaseURL+"/api/reviews/{reviewId}/toolrequests", wrapper.GetReviewToolRequests)
	m.HandleFunc("GET "+options.BaseURL+"/api/runs/{runId}", wrapper.GetRun)
	m.HandleFunc("GET "+options.BaseURL+"/api/runs/{runId}/tools", wrapper.GetRunTools)
	m.HandleFunc("GET "+options.BaseURL+"/api/runs/{runId}/tools/{toolId}/supervisors", wrapper.GetRunToolSupervisors)
	m.HandleFunc("POST "+options.BaseURL+"/api/runs/{runId}/tools/{toolId}/supervisors", wrapper.CreateRunToolSupervisors)
	m.HandleFunc("GET "+options.BaseURL+"/api/stats", wrapper.GetHubStats)
	m.HandleFunc("GET "+options.BaseURL+"/api/supervisors", wrapper.GetSupervisors)
	m.HandleFunc("POST "+options.BaseURL+"/api/supervisors", wrapper.CreateSupervisor)
	m.HandleFunc("GET "+options.BaseURL+"/api/supervisors/{supervisorId}", wrapper.GetSupervisor)
	m.HandleFunc("GET "+options.BaseURL+"/api/tools", wrapper.GetTools)
	m.HandleFunc("POST "+options.BaseURL+"/api/tools", wrapper.CreateTool)
	m.HandleFunc("GET "+options.BaseURL+"/api/tools/{toolId}", wrapper.GetTool)

	return m
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"9F8AAAD//1iIH92yNQAA",
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
