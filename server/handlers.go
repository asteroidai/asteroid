package sentinel

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// respondJSON writes a JSON response with status 200 OK
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	// If the header is not set, set it to StatusOK
	if w.Header().Get("Content-Type") == "" {
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func apiCreateProjectHandler(w http.ResponseWriter, r *http.Request, store ProjectStore) {
	ctx := r.Context()

	log.Printf("received new project registration request")
	var request struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingProject, err := store.GetProjectFromName(ctx, request.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting project: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	if existingProject != nil {
		w.WriteHeader(http.StatusOK)
		respondJSON(w, existingProject.Id.String())
		return
	}

	// Generate a new Project ID
	id := uuid.New()

	// Create the Project struct
	project := Project{
		Id:        id,
		Name:      request.Name,
		CreatedAt: time.Now(),
	}

	// Store the project in the global projects map
	err = store.CreateProject(ctx, project)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to register project, %v", err), http.StatusInternalServerError)
		return
	}

	// Send the response
	w.WriteHeader(http.StatusCreated)
	respondJSON(w, id.String())
}

func apiCreateProjectRunHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID, store RunStore) {
	ctx := r.Context()

	log.Printf("received new run request for project ID: %s", id)

	run := Run{
		Id:        uuid.Nil,
		ProjectId: id,
		CreatedAt: time.Now(),
	}

	runID, err := store.CreateRun(ctx, run)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("created run with ID: %s", run.Id)

	respondJSON(w, runID)
}

func apiGetProjectRunsHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID, store RunStore) {
	ctx := r.Context()

	runs, err := store.GetProjectRuns(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, runs)
}

func apiCreateRunToolHandler(w http.ResponseWriter, r *http.Request, runId uuid.UUID, store ToolStore) {
	ctx := r.Context()

	var t struct {
		Attributes        map[string]interface{} `json:"attributes"`
		Name              string                 `json:"name"`
		Description       string                 `json:"description"`
		IgnoredAttributes []string               `json:"ignoredAttributes"`
	}
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err.Error())
		return
	}

	var existingTool *Tool
	if t.Attributes != nil && t.Name != "" && t.Description != "" && t.IgnoredAttributes != nil {
		found, err := store.GetToolFromValues(ctx, t.Attributes, t.Name, t.Description, t.IgnoredAttributes)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "error trying to locate an existing tool", err.Error())
			return
		}
		if found != nil {
			existingTool = found
		}
	}

	if existingTool != nil {
		respondJSON(w, existingTool)
		return
	}

	toolId, err := store.CreateTool(ctx, runId, t.Attributes, t.Name, t.Description, t.IgnoredAttributes)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error creating tool", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, toolId)
}

func apiGetSupervisorHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID, store SupervisorStore) {
	ctx := r.Context()

	supervisor, err := store.GetSupervisor(ctx, id)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting supervisor", err.Error())
		return
	}

	respondJSON(w, supervisor)
}

func apiCreateToolSupervisorChainsHandler(w http.ResponseWriter, r *http.Request, toolId uuid.UUID, store SupervisorStore) {
	ctx := r.Context()

	var request []ChainRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid JSON format", err.Error())
		return
	}

	// TODO do we want to return the chains here?
	chainIds := make([]uuid.UUID, 0)
	for _, chain := range request {
		chainId, err := store.CreateSupervisorChain(ctx, toolId, chain)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "error creating supervisor chain", err.Error())
			return
		}
		chainIds = append(chainIds, *chainId)
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, chainIds)
}

func apiCreateSupervisorHandler(w http.ResponseWriter, r *http.Request, _ uuid.UUID, store SupervisorStore) {
	ctx := r.Context()

	var request Supervisor
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid JSON format", err.Error())
		return
	}

	// Create new supervisor
	supervisorId, err := store.CreateSupervisor(ctx, request)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error creating supervisor", err.Error())
		return
	}

	request.Id = &supervisorId

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, request)
}

func apiGetToolSupervisorChainsHandler(w http.ResponseWriter, r *http.Request, toolId uuid.UUID, store Store) {
	ctx := r.Context()

	// First check if tool exists
	tool, err := store.GetTool(ctx, toolId)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting tool", err.Error())
		return
	}

	if tool == nil {
		sendErrorResponse(w, http.StatusNotFound, "Tool not found", "")
		return
	}

	chains, err := store.GetSupervisorChains(ctx, toolId)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting tool supervisor chains", err.Error())
		return
	}

	respondJSON(w, chains)
}

func apiGetRunToolsHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID, store Store) {
	ctx := r.Context()

	// First check if run exists
	run, err := store.GetRun(ctx, id)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting run", err.Error())
		return
	}

	if run == nil {
		sendErrorResponse(w, http.StatusNotFound, "Run not found", "")
		return
	}

	tools, err := store.GetRunTools(ctx, id)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting run tools", err.Error())
		return
	}

	respondJSON(w, tools)
}

func apiCreateToolRequestGroupHandler(w http.ResponseWriter, r *http.Request, toolId uuid.UUID, store ToolRequestStore) {
	ctx := r.Context()

	var request ToolRequestGroup
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err.Error())
		return
	}

	trg, err := store.CreateToolRequestGroup(ctx, toolId, request)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error creating tool request group", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, trg.Id)
}

func apiGetRequestGroupHandler(w http.ResponseWriter, r *http.Request, requestGroupId uuid.UUID, store Store) {
	ctx := r.Context()

	requestGroup, err := store.GetRequestGroup(ctx, requestGroupId)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting request group", err.Error())
		return
	}

	respondJSON(w, requestGroup)
}

func apiGetSupervisionResultHandler(w http.ResponseWriter, r *http.Request, supervisionRequestId uuid.UUID, store Store) {
	ctx := r.Context()

	// Check that the supervision request exists
	supervisionRequest, err := store.GetSupervisionRequest(ctx, supervisionRequestId)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting supervision request", err.Error())
		return
	}

	if supervisionRequest == nil {
		sendErrorResponse(w, http.StatusNotFound, "Supervision request not found", "")
		return
	}

	supervisionResult, err := store.GetSupervisionResult(ctx, supervisionRequestId)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting supervision result", err.Error())
		return
	}

	respondJSON(w, supervisionResult)
}

func apiGetRunRequestGroupsHandler(w http.ResponseWriter, r *http.Request, runId uuid.UUID, store Store) {
	ctx := r.Context()

	run, err := store.GetRun(ctx, runId)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting run", err.Error())
		return
	}

	if run == nil {
		sendErrorResponse(w, http.StatusNotFound, "Run not found", "")
		return
	}

	requestGroups, err := store.GetRunRequestGroups(ctx, runId)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting run request groups", err.Error())
		return
	}

	respondJSON(w, requestGroups)
}

func apiGetSupervisorsHandler(w http.ResponseWriter, r *http.Request, projectId uuid.UUID, store Store) {
	ctx := r.Context()

	// First check if project exists
	project, err := store.GetProject(ctx, projectId)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting project", err.Error())
		return
	}

	if project == nil {
		sendErrorResponse(w, http.StatusNotFound, "Project not found", "")
		return
	}

	supervisors, err := store.GetSupervisors(ctx, projectId)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting supervisors", err.Error())
		return
	}

	respondJSON(w, supervisors)
}

func apiGetProjectToolsHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID, store ToolStore) {
	ctx := r.Context()

	tools, err := store.GetProjectTools(ctx, id)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting project tools", err.Error())
		return
	}

	respondJSON(w, tools)
}

func apiGetToolHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID, store ToolStore) {
	ctx := r.Context()

	tool, err := store.GetTool(ctx, id)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting tool", err.Error())
		return
	}

	if tool == nil {
		sendErrorResponse(w, http.StatusNotFound, "Tool not found", "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tool)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error encoding tool", err.Error())
		return
	}
}

func apiGetSupervisionRequestStatusHandler(w http.ResponseWriter, r *http.Request, reviewID uuid.UUID, store Store) {
	ctx := r.Context()
	// Use the reviewID directly
	supervisor, err := store.GetSupervisionRequest(ctx, reviewID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting supervisor", err.Error())
		return
	}

	if supervisor == nil {
		sendErrorResponse(w, http.StatusNotFound, "Supervisor not found", "")
		return
	}

	respondJSON(w, supervisor.Status)
}

func apiCreateSupervisionRequestHandler(
	w http.ResponseWriter,
	r *http.Request,
	requestGroupId uuid.UUID,
	chainId uuid.UUID,
	supervisorId uuid.UUID,
	store Store,
) {
	ctx := r.Context()

	var request SupervisionRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err.Error())
		return
	}

	// Check that the request, chain and supervisor exist
	requestGroup, err := store.GetRequestGroup(ctx, requestGroupId)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting request group", err.Error())
		return
	}

	if requestGroup == nil {
		sendErrorResponse(w, http.StatusNotFound, "Request group not found", "")
		return
	}

	if len(requestGroup.ToolRequests) > 1 {
		sendErrorResponse(w, http.StatusBadRequest, "Request group must contain only one tool request", "")
		return
	}

	chain, err := store.GetSupervisorChain(ctx, chainId)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting supervisor chain", err.Error())
		return
	}

	if chain == nil {
		sendErrorResponse(w, http.StatusNotFound, "Supervisor chain not found", "")
		return
	}

	supervisor, err := store.GetSupervisor(ctx, supervisorId)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting supervisor", err.Error())
		return
	}

	if supervisor == nil {
		sendErrorResponse(w, http.StatusNotFound, "Supervisor not found", "")
		return
	}

	// Check that the supervisor is associated with the tool/request/chain
	found := false
	pos := -1
	for i, chainSupervisor := range chain.Supervisors {
		if chainSupervisor.Id.String() == supervisorId.String() {
			found = true
			pos = i
			break
		}
	}

	if !found {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Supervisor %s not associated with chain %s", supervisorId, chainId), "")
		return
	}

	if pos != request.PositionInChain {
		sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Supervisor %s is not in the correct position in chain %s", supervisorId, chainId), "")
		return
	}

	// Store the supervision in the database
	reviewID, err := store.CreateSupervisionRequest(ctx, request)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error creating supervision request", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, reviewID)
}

func apiCreateSupervisionResultHandler(
	w http.ResponseWriter,
	r *http.Request,
	supervisionRequestId uuid.UUID,
	store Store,
) {
	ctx := r.Context()

	var result SupervisionResult
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err.Error())
		return
	}

	// Check that the group, chain and supervisor, and request exist
	id, err := store.CreateSupervisionResult(ctx, result, supervisionRequestId)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error creating supervision result", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, id)
}

func apiGetHubStatsHandler(w http.ResponseWriter, _ *http.Request, hub *Hub) {
	stats, err := hub.getStats()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting stats", err.Error())
		return
	}

	respondJSON(w, stats)
}

// apiGetProjectsHandler returns all projects
func apiGetProjectsHandler(w http.ResponseWriter, r *http.Request, store ProjectStore) {
	ctx := r.Context()

	projects, err := store.GetProjects(ctx)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting projects", err.Error())
		return
	}

	respondJSON(w, projects)
}

func apiGetProjectHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID, store ProjectStore) {
	ctx := r.Context()

	project, err := store.GetProject(ctx, id)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "error getting project", err.Error())
		return
	}

	if project == nil {
		sendErrorResponse(w, http.StatusNotFound, "Project not found", "")
		return
	}

	respondJSON(w, project)
}

// func apiGetSupervisionRequestHandler(w http.ResponseWriter, r *http.Request, supervisionRequestId uuid.UUID, store Store) {
// 	ctx := r.Context()

// 	supervisionRequest, err := store.GetSupervisionRequest(ctx, supervisionRequestId)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "error getting supervision request", err.Error())
// 		return
// 	}

// 	respondJSON(w, supervisionRequest)
// }

// apiGetExecutionSupervisionsHandler handles the GET /api/executions/{executionId}/supervisions endpoint
// func apiGetExecutionSupervisionsHandler(w http.ResponseWriter, r *http.Request, executionId uuid.UUID, store Store) {
// 	ctx := r.Context()

// 	// First check if execution exists
// 	execution, err := store.GetExecution(ctx, executionId)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "error getting execution", err.Error())
// 		return
// 	}

// 	if execution == nil {
// 		sendErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Execution not found for ID %s", executionId), "")
// 		return
// 	}

// 	supervisions, err := store.GetExecutionSupervisions(ctx, executionId)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "error getting execution supervisions", err.Error())
// 		return
// 	}

// 	// Combine the data into ExecutionSupervisions response
// 	response := &ExecutionSupervisions{
// 		ExecutionId:  executionId,
// 		Supervisions: supervisions,
// 	}

// 	respondJSON(w, response)
// }

// func apiCreateSupervisorChainHandler(w http.ResponseWriter, r *http.Request, runId uuid.UUID, store SupervisorStore) {
// 	ctx := r.Context()

// 	var request CreateSupervisorChainRequest
// 	err := json.NewDecoder(r.Body).Decode(&request)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err.Error())
// 		return
// 	}

// 	supervisorChainId, err := store.CreateSupervisorChain(ctx, runId, request.SupervisorIds)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "error creating supervisor chain", err.Error())
// 		return
// 	}

// 	respondJSON(w, supervisorChainId)
// }

// func apiGetRunToolSupervisorChainsHandler(w http.ResponseWriter, r *http.Request, runId uuid.UUID, toolId uuid.UUID, store SupervisorStore) {
// 	ctx := r.Context()

// 	supervisorChains, err := store.GetRunToolSupervisorChains(ctx, runId, toolId)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "error getting supervisor chains", err.Error())
// 		return
// 	}

// 	respondJSON(w, supervisorChains)
// }

// func apiGetRunExecutionsHandler(w http.ResponseWriter, r *http.Request, runId uuid.UUID, store Store) {
// 	ctx := r.Context()

// 	run, err := store.GetRun(ctx, runId)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "error getting run", err.Error())
// 		return
// 	}

// 	if run == nil {
// 		sendErrorResponse(w, http.StatusNotFound, "Run not found", "")
// 		return
// 	}

// 	executions, err := store.GetRunExecutions(ctx, runId)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "error getting run executions", err.Error())
// 		return
// 	}

// 	respondJSON(w, executions)
// }

// func apiCreateExecutionHandler(w http.ResponseWriter, r *http.Request, runId uuid.UUID, store Store) {
// 	ctx := r.Context()

// 	var request struct {
// 		ToolId uuid.UUID `json:"toolId"`
// 	}

// 	err := json.NewDecoder(r.Body).Decode(&request)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err.Error())
// 		return
// 	}

// 	executionId, err := store.CreateExecution(ctx, runId, request.ToolId)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "error creating execution", err.Error())
// 		return
// 	}

// 	execution, err := store.GetExecution(ctx, executionId)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "error getting execution", err.Error())
// 		return
// 	}

// 	if execution == nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "Something went wrong, execution not found", "")
// 		return
// 	}

// 	respondJSON(w, execution)
// }

// apiGetSupervisionRequestHandler handles the GET /api/supervisor/{id} endpoint
// func apiGetSupervisionRequestHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID, store Store) {
// 	ctx := r.Context()

// 	supervisor, err := store.GetSupervisionRequest(ctx, id)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "error getting supervisor", err.Error())
// 		return
// 	}

// 	if supervisor == nil {
// 		sendErrorResponse(w, http.StatusNotFound, "Supervisor not found", "")
// 		return
// 	}

// 	respondJSON(w, supervisor)
// }

// apiGetSupervisionRequestsHandler handles the GET /api/supervisor endpoint
// func apiGetSupervisionRequestsHandler(w http.ResponseWriter, r *http.Request, store Store) {
// 	ctx := r.Context()

// 	reviews, err := store.GetSupervisionRequests(ctx)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "error getting supervision requests", err.Error())
// 		return
// 	}

// 	respondJSON(w, reviews)
// }

// apiGetSupervisionResultsHandler handles the GET /api/supervisor/{id}/results endpoint
// func apiGetSupervisionResultsHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID, store Store) {
// 	ctx := r.Context()

// 	// First check if supervisor exists
// 	supervisor, err := store.GetSupervisionRequest(ctx, id)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "error getting supervisor", err.Error())
// 		return
// 	}

// 	if supervisor == nil {
// 		sendErrorResponse(w, http.StatusNotFound, "Supervisor not found", "")
// 		return
// 	}

// 	results, err := store.GetSupervisionResults(ctx, id)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "error getting supervision results", err.Error())
// 		return
// 	}

// 	respondJSON(w, results)
// }

// apiGetReviewToolRequestsHandler handles the GET /api/supervisor/{id}/toolrequests endpoint
// func apiGetReviewToolRequestsHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID, store Store) {
// 	ctx := r.Context()

// 	// First check if supervisor exists
// 	supervisor, err := store.GetSupervisionRequest(ctx, id)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "error getting supervisor", err.Error())
// 		return
// 	}

// 	if supervisor == nil {
// 		sendErrorResponse(w, http.StatusNotFound, "Supervisor not found", "")
// 		return
// 	}

// 	results, err := store.GetSupervisionToolRequests(ctx, id)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "error getting supervisor tool requests", err.Error())
// 		return
// 	}

// 	respondJSON(w, results)
// }

// // apiLLMExplanationHandler receives a code snippet and returns an explanation and a danger score by calling an LLM
// func apiLLMExplanationHandler(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	var request struct {
// 		Text string `json:"text"`
// 	}

// 	err := json.NewDecoder(r.Body).Decode(&request)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format", err.Error())
// 		return
// 	}

// 	explanation, score, err := getExplanationFromLLM(ctx, request.Text)
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 		sendErrorResponse(w, http.StatusInternalServerError, "Failed to get explanation from LLM", err.Error())
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	err = json.NewEncoder(w).Encode(map[string]string{"explanation": explanation, "score": score})
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "Failed to get explanation from LLM", err.Error())
// 		return
// 	}
// }

// func apiGetRunsHandler(w http.ResponseWriter, r *http.Request, projectId uuid.UUID, store RunStore) {
// 	ctx := r.Context()

// 	runs, err := store.GetRuns(ctx, projectId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	respondJSON(w, runs)
// }

// func apiGetRunHandler(w http.ResponseWriter, r *http.Request, id uuid.UUID, store RunStore) {
// 	ctx := r.Context()

// 	run, err := store.GetRun(ctx, id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if run == nil {
// 		http.Error(w, "Run not found", http.StatusNotFound)
// 		return
// 	}

// 	respondJSON(w, run)
// }
