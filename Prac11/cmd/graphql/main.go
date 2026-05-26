package main

import (
	"Prak_11/graph/model"
	"Prak_11/internal/store"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"Prak_11/graph"
)

// graphqlRequest is the JSON body of a GraphQL request.
type graphqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

func main() {
	//s := store.New()
	//srv := handler.NewDefaultServer(
	//	graph.NewExecutableSchema(graph.Config{
	//		Resolvers: &graph.Resolver{Store: s},
	//	}),
	//)
	//
	////http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	//http.Handle("/", srv)
	//
	//log.Println("GraphQL server started on http://localhost:8080")
	//log.Fatal(http.ListenAndServe(":8080", nil))

	s := store.New()
	r := &graph.Resolver{Store: s}

	mux := http.NewServeMux()

	// GraphQL endpoint + minimal Playground
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(playgroundHTML))
			return
		}

		if req.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var gqlReq graphqlRequest
		if err := json.NewDecoder(req.Body).Decode(&gqlReq); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		result, err := executeQuery(req.Context(), r, gqlReq.Query, gqlReq.Variables)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"errors": []map[string]string{{"message": err.Error()}},
			})
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]any{"data": result})
	})

	log.Println("GraphQL server started on http://localhost:8080")
	log.Println("Open http://localhost:8080 for Playground")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

// executeQuery is a minimal hand-written GraphQL dispatcher.
// In a real project, use the gqlgen-generated handler instead.
func executeQuery(ctx context.Context, r *graph.Resolver, query string, vars map[string]any) (map[string]any, error) {
	q := strings.TrimSpace(query)

	// tasks query
	if strings.Contains(q, "tasks {") && !strings.Contains(q, "task(") {
		tasks, err := r.Query().Tasks(ctx)
		if err != nil {
			return nil, err
		}
		return map[string]any{"tasks": tasks}, nil
	}

	// task(id) query
	if strings.HasPrefix(q, "query") && strings.Contains(q, "task(") {
		id := ""
		if v, ok := vars["id"]; ok {
			id, _ = v.(string)
		}
		t, err := r.Query().Task(ctx, id)
		if err != nil {
			return nil, err
		}
		return map[string]any{"task": t}, nil
	}

	// createTask mutation
	if strings.Contains(q, "createTask") {
		input := model.CreateTaskInput{}
		if v, ok := vars["input"].(map[string]any); ok {
			if t, ok := v["title"].(string); ok {
				input.Title = t
			}
			if d, ok := v["description"].(string); ok {
				input.Description = &d
			}
		}
		t, err := r.Mutation().CreateTask(ctx, input)
		if err != nil {
			return nil, err
		}
		return map[string]any{"createTask": t}, nil
	}

	// updateTask mutation
	if strings.Contains(q, "updateTask") {
		id := ""
		if v, ok := vars["id"]; ok {
			id, _ = v.(string)
		}
		input := model.UpdateTaskInput{}
		if v, ok := vars["input"].(map[string]any); ok {
			if t, ok := v["title"].(string); ok {
				input.Title = &t
			}
			if d, ok := v["description"].(string); ok {
				input.Description = &d
			}
			if done, ok := v["done"].(bool); ok {
				input.Done = &done
			}
		}
		t, err := r.Mutation().UpdateTask(ctx, id, input)
		if err != nil {
			return nil, err
		}
		return map[string]any{"updateTask": t}, nil
	}

	// deleteTask mutation
	if strings.Contains(q, "deleteTask") {
		id := ""
		if v, ok := vars["id"]; ok {
			id, _ = v.(string)
		}
		ok, err := r.Mutation().DeleteTask(ctx, id)
		if err != nil {
			return nil, err
		}
		return map[string]any{"deleteTask": ok}, nil
	}

	return nil, nil
}

// playgroundHTML is a minimal GraphQL Playground HTML page.
const playgroundHTML = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>GraphQL Playground</title>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/3.0.9/graphiql.min.css">
</head>
<body style="margin:0">
  <div id="graphiql" style="height:100vh"></div>
  <script crossorigin src="https://cdnjs.cloudflare.com/ajax/libs/react/18.2.0/umd/react.production.min.js"></script>
  <script crossorigin src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/18.2.0/umd/react-dom.production.min.js"></script>
  <script crossorigin src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/3.0.9/graphiql.min.js"></script>
  <script>
    const fetcher = GraphiQL.createFetcher({ url: 'http://localhost:8080/' });
    ReactDOM.createRoot(document.getElementById('graphiql')).render(
      React.createElement(GraphiQL, { fetcher })
    );
  </script>
</body>
</html>`
