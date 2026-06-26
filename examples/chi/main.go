package main

import (
	"encoding/json"
	"log"
	"net/http"

	"chi_example/middlewares"

	ctonic "github.com/TickLabVN/tonic/adapters/chi"
	"github.com/TickLabVN/tonic/core"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/go-chi/chi/v5"
)

type GetUserRequest struct {
	ID             string `path:"id" validate:"required,uuid4"`
	IncludeOrders  bool   `query:"includeOrders"`
	IncludeMetrics bool   `query:"includeMetrics"`
	APIKey         string `header:"x-api-key" validate:"required,min=10"`
}

type User struct {
	ID    string `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type CreateUserRequest struct {
	Name      string `json:"name" validate:"required,min=2,max=60"`
	Email     string `json:"email" validate:"required,email"`
	Role      string `json:"role" validate:"required,oneof=admin manager member"`
	Invite    bool   `json:"invite"`
	RequestID string `header:"x-request-id" validate:"required,uuid4"`
}

type UpdateUserSettingsRequest struct {
	ID             string `path:"id" validate:"required,uuid4"`
	DryRun         bool   `query:"dryRun"`
	RequestID      string `header:"x-request-id" validate:"required,uuid4"`
	DisplayName    string `json:"displayName" validate:"required,min=2,max=50"`
	Timezone       string `json:"timezone" validate:"required"`
	MarketingOptIn bool   `json:"marketingOptIn"`
}

type ListUserOrdersRequest struct {
	ID       string `path:"id" validate:"required,uuid4"`
	Limit    int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Cursor   string `query:"cursor"`
	Region   string `header:"x-region" validate:"required,oneof=apac emea us"`
	Statuses string `query:"statuses"`
}

type UserDetailsResponse struct {
	User
	Segments []string `json:"segments"`
}

type UserCreatedResponse struct {
	User
	Invited bool `json:"invited"`
}

type UserSettingsResponse struct {
	UserID         string `json:"userId"`
	DisplayName    string `json:"displayName"`
	Timezone       string `json:"timezone"`
	MarketingOptIn bool   `json:"marketingOptIn"`
	DryRun         bool   `json:"dryRun"`
}

type OrderSummary struct {
	ID       string  `json:"id"`
	Status   string  `json:"status"`
	Currency string  `json:"currency"`
	Total    float64 `json:"total"`
}

type OrderListResponse struct {
	UserID     string         `json:"userId"`
	NextCursor string         `json:"nextCursor"`
	Orders     []OrderSummary `json:"orders"`
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	data := middlewares.Data[GetUserRequest](r)
	segments := []string{"starter", "beta"}
	if data.IncludeMetrics {
		segments = append(segments, "usage-insights")
	}
	writeJSON(w, http.StatusOK, UserDetailsResponse{
		User: User{
			ID:    data.ID,
			Name:  "John Doe",
			Email: "john.doe@example.com",
		},
		Segments: segments,
	})
}

func createUser(w http.ResponseWriter, r *http.Request) {
	data := middlewares.Data[CreateUserRequest](r)
	writeJSON(w, http.StatusCreated, UserCreatedResponse{
		User: User{
			ID:    "7f4d8a3f-3f3e-4a58-a9f0-1b0e1776b001",
			Name:  data.Name,
			Email: data.Email,
		},
		Invited: data.Invite,
	})
}

func updateUserSettings(w http.ResponseWriter, r *http.Request) {
	data := middlewares.Data[UpdateUserSettingsRequest](r)
	writeJSON(w, http.StatusOK, UserSettingsResponse{
		UserID:         data.ID,
		DisplayName:    data.DisplayName,
		Timezone:       data.Timezone,
		MarketingOptIn: data.MarketingOptIn,
		DryRun:         data.DryRun,
	})
}

func listUserOrders(w http.ResponseWriter, r *http.Request) {
	data := middlewares.Data[ListUserOrdersRequest](r)
	writeJSON(w, http.StatusOK, OrderListResponse{
		UserID:     data.ID,
		NextCursor: "cursor:page:2",
		Orders: []OrderSummary{
			{ID: "ord_1001", Status: "paid", Currency: "USD", Total: 129.5},
			{ID: "ord_1002", Status: "pending", Currency: "USD", Total: 42.0},
		},
	})
}

func main() {
	r := chi.NewRouter()

	schema := ctonic.New(&docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Version: "1.0.0",
			Title:   "Chi Example API",
		},
	})

	api := schema.Wrap(r)
	api.Route("/api/v1/users", func(users chi.Router) {
		ctonic.For[GetUserRequest, UserDetailsResponse](schema).
			GET(users, "/{id}", middlewares.Bind[GetUserRequest], getUserByID, ctonic.WithOperation(docs.OperationObject{
				Summary:     "Get a user profile",
				Description: "Returns a user profile with optional segments and relationship toggles.",
				Tags:        []string{"Users"},
			}))
		ctonic.For[CreateUserRequest, UserCreatedResponse](schema).
			POST(users, "/", middlewares.Bind[CreateUserRequest], createUser, ctonic.WithOperation(docs.OperationObject{
				Summary:     "Create a user",
				Description: "Creates a new user and optionally triggers an invite flow.",
				Tags:        []string{"Users"},
			}))
		ctonic.For[UpdateUserSettingsRequest, UserSettingsResponse](schema).
			PATCH(users, "/{id}/settings", middlewares.Bind[UpdateUserSettingsRequest], updateUserSettings, ctonic.WithOperation(docs.OperationObject{
				Summary:     "Update user settings",
				Description: "Demonstrates path, query, header, and JSON body binding in one operation.",
				Tags:        []string{"Users", "Settings"},
			}))
		ctonic.For[ListUserOrdersRequest, OrderListResponse](schema).
			GET(users, "/{id}/orders", middlewares.Bind[ListUserOrdersRequest], listUserOrders, ctonic.WithOperation(docs.OperationObject{
				Summary:     "List user orders",
				Description: "Shows filtered collections with pagination and regional headers.",
				Tags:        []string{"Orders"},
			}))
	})
	// Default renderer is Swagger UI. Pass a core.UI to pick another, e.g.:
	// schema.UIHandle(r, "/docs", core.ReDoc) // or core.SwaggerUI, core.RapiDoc
	schema.UIHandle(r, "/docs", core.SwaggerUI)

	log.Fatal(http.ListenAndServe(":1523", r))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
