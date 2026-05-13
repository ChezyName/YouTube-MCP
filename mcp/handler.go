package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/ChezyName/YouTube-MCP/youtube"
	"github.com/gorilla/mux"
)

// Calls an existing handler internally and returns the response body as a string
func callHandler(handler http.HandlerFunc, method, path string, vars map[string]string, queryParams map[string]string) (string, error) {
	req := httptest.NewRequest(method, path, nil)

	// Add query params
	if len(queryParams) > 0 {
		q := req.URL.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Add mux vars (path params like {id})
	if len(vars) > 0 {
		req = mux.SetURLVars(req, vars)
	}

	rr := httptest.NewRecorder()
	handler(rr, req)

	res := rr.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("handler error %d: %s", res.StatusCode, string(body))
	}

	return string(body), nil
}

func HandleMCP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, nil, -32700, "Parse error")
		return
	}

	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		respondError(w, nil, -32700, "Parse error")
		return
	}

	switch req.Method {

	case "initialize":
		respond(w, req.ID, map[string]any{
			"protocolVersion": "2025-03-26",
			"serverInfo": map[string]any{
				"name":    "youtube-mcp",
				"version": "1.0.0",
			},
			"capabilities": map[string]any{
				"tools": map[string]any{},
			},
		})

	case "notifications/initialized":
		w.WriteHeader(http.StatusNoContent)

	case "tools/list":
		respond(w, req.ID, map[string]any{
			"tools": GetTools(),
		})

	case "tools/call":
		var params struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			respondError(w, req.ID, -32600, "Invalid params")
			return
		}

		result, err := dispatchTool(params.Name, params.Arguments)
		if err != nil {
			respondError(w, req.ID, -32603, err.Error())
			return
		}

		respond(w, req.ID, ToolResult{
			Content: []ToolContent{{Type: "text", Text: result}},
		})

	default:
		respondError(w, req.ID, -32601, "Method not found")
	}
}

func dispatchTool(name string, args map[string]any) (string, error) {
	getString := func(key string) string {
		if v, ok := args[key]; ok {
			return fmt.Sprintf("%v", v)
		}
		return ""
	}

	switch name {
	case "list_videos":
		return callHandler(youtube.ListVideos, "GET", "/videos", nil, nil)

	case "get_video":
		id := getString("id")
		if id == "" {
			return "", fmt.Errorf("missing required argument: id")
		}
		return callHandler(youtube.GetVideo, "GET", "/videos/"+id,
			map[string]string{"id": id}, nil)

	case "get_video_analytics":
		id := getString("id")
		if id == "" {
			return "", fmt.Errorf("missing required argument: id")
		}
		queryParams := map[string]string{}
		if r := getString("range"); r != "" {
			queryParams["range"] = r
		}
		return callHandler(youtube.GetAnalyticsForVideo, "GET", "/videos/"+id+"/analytics",
			map[string]string{"id": id}, queryParams)

	case "get_video_comments":
		id := getString("id")
		if id == "" {
			return "", fmt.Errorf("missing required argument: id")
		}
		queryParams := map[string]string{}
		if l := getString("limit"); l != "" {
			queryParams["limit"] = l
		}
		return callHandler(youtube.GetVideoComments, "GET", "/videos/"+id+"/comments",
			map[string]string{"id": id}, queryParams)

	case "get_channel":
		return callHandler(youtube.GetChannel, "GET", "/channel", nil, nil)

	case "get_channel_analytics":
		queryParams := map[string]string{}
		if r := getString("range"); r != "" {
			queryParams["range"] = r
		}
		return callHandler(youtube.GetChannelAnalytics, "GET", "/channel/analytics",
			nil, queryParams)

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}
