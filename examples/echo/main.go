package main

import (
	"net/http"

	"echo_example/middlewares"
	"echo_example/utils"

	etonic "github.com/TickLabVN/tonic/adapters/echo"
	"github.com/TickLabVN/tonic/core"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type GetUserRequest struct {
	ID             string `param:"id" validate:"required,uuid4"`
	IncludeOrders  bool   `query:"includeOrders"`
	IncludeMetrics bool   `query:"includeMetrics"`
	ApiKey         string `header:"x-api-key" validate:"required,min=10"`
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
	ID             string `param:"id" validate:"required,uuid4"`
	DryRun         bool   `query:"dryRun"`
	RequestID      string `header:"x-request-id" validate:"required,uuid4"`
	DisplayName    string `json:"displayName" validate:"required,min=2,max=50"`
	Timezone       string `json:"timezone" validate:"required"`
	MarketingOptIn bool   `json:"marketingOptIn"`
}

type ListUserOrdersRequest struct {
	ID       string `param:"id" validate:"required,uuid4"`
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

func getUser(c echo.Context) error {
	data := c.Get("data").(GetUserRequest)
	segments := []string{"starter", "beta"}
	if data.IncludeMetrics {
		segments = append(segments, "usage-insights")
	}
	return c.JSON(http.StatusOK, UserDetailsResponse{
		User: User{
			ID:    data.ID,
			Name:  "John Doe",
			Email: "john.doe@example.com",
		},
		Segments: segments,
	})
}

func createUser(c echo.Context) error {
	data := c.Get("data").(CreateUserRequest)
	return c.JSON(http.StatusCreated, UserCreatedResponse{
		User: User{
			ID:    "7f4d8a3f-3f3e-4a58-a9f0-1b0e1776b001",
			Name:  data.Name,
			Email: data.Email,
		},
		Invited: data.Invite,
	})
}

func updateUserSettings(c echo.Context) error {
	data := c.Get("data").(UpdateUserSettingsRequest)
	return c.JSON(http.StatusOK, UserSettingsResponse{
		UserID:         data.ID,
		DisplayName:    data.DisplayName,
		Timezone:       data.Timezone,
		MarketingOptIn: data.MarketingOptIn,
		DryRun:         data.DryRun,
	})
}

func listUserOrders(c echo.Context) error {
	data := c.Get("data").(ListUserOrdersRequest)
	return c.JSON(http.StatusOK, OrderListResponse{
		UserID:     data.ID,
		NextCursor: "cursor:page:2",
		Orders: []OrderSummary{
			{ID: "ord_1001", Status: "paid", Currency: "USD", Total: 129.5},
			{ID: "ord_1002", Status: "pending", Currency: "USD", Total: 42.0},
		},
	})
}

func main() {
	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	schema := etonic.New(&docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Version: "1.0.0",
			Title:   "Echo Example API",
		},
	})
	api := e.Group("/api/v1")

	etonic.For[GetUserRequest, UserDetailsResponse](schema).AddRoute(
		api.GET("/users/:id", getUser, middlewares.Bind[GetUserRequest]),
		docs.OperationObject{
			Summary:     "Get a user profile",
			Description: "Returns a user profile with optional segments and relationship toggles.",
			Tags:        []string{"Users"},
		},
	)
	etonic.For[CreateUserRequest, UserCreatedResponse](schema).AddRoute(
		api.POST("/users", createUser, middlewares.Bind[CreateUserRequest]),
		docs.OperationObject{
			Summary:     "Create a user",
			Description: "Creates a new user and optionally triggers an invite flow.",
			Tags:        []string{"Users"},
		},
	)
	etonic.For[UpdateUserSettingsRequest, UserSettingsResponse](schema).AddRoute(
		api.PATCH("/users/:id/settings", updateUserSettings, middlewares.Bind[UpdateUserSettingsRequest]),
		docs.OperationObject{
			Summary:     "Update user settings",
			Description: "Demonstrates path, query, header, and JSON body binding in one operation.",
			Tags:        []string{"Users", "Settings"},
		},
	)
	etonic.For[ListUserOrdersRequest, OrderListResponse](schema).AddRoute(
		api.GET("/users/:id/orders", listUserOrders, middlewares.Bind[ListUserOrdersRequest]),
		docs.OperationObject{
			Summary:     "List user orders",
			Description: "Shows filtered collections with pagination and regional headers.",
			Tags:        []string{"Orders"},
		},
	)
	// Default renderer is Swagger UI. Pass a core.UI to pick another, e.g.:
	// schema.UIHandle(e, "/docs", core.Redoc)   // or core.Scalar, core.RapiDoc
	// schema.UIHandle(e, "/docs", core.SwaggerUI)
	schema.UIHandle(e, "/docs", core.Scalar)

	e.Logger.Fatal(e.Start(":1323"))
}
