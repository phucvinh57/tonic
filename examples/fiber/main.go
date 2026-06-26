package main

import (
	"log"
	"net/http"

	"fiber_example/middlewares"

	ftonic "github.com/TickLabVN/tonic/adapters/fiber"
	"github.com/TickLabVN/tonic/core"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gofiber/fiber/v3"
)

type GetUserRequest struct {
	ID             string `uri:"id" validate:"required,uuid4"`
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
	ID             string `uri:"id" validate:"required,uuid4"`
	DryRun         bool   `query:"dryRun"`
	RequestID      string `header:"x-request-id" validate:"required,uuid4"`
	DisplayName    string `json:"displayName" validate:"required,min=2,max=50"`
	Timezone       string `json:"timezone" validate:"required"`
	MarketingOptIn bool   `json:"marketingOptIn"`
}

type ListUserOrdersRequest struct {
	ID       string `uri:"id" validate:"required,uuid4"`
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

func getUserByID(c fiber.Ctx) error {
	data := c.Locals("data").(GetUserRequest)
	segments := []string{"starter", "beta"}
	if data.IncludeMetrics {
		segments = append(segments, "usage-insights")
	}
	return c.Status(http.StatusOK).JSON(UserDetailsResponse{
		User: User{
			ID:    data.ID,
			Name:  "John Doe",
			Email: "john.doe@example.com",
		},
		Segments: segments,
	})
}

func createUser(c fiber.Ctx) error {
	data := c.Locals("data").(CreateUserRequest)
	return c.Status(http.StatusCreated).JSON(UserCreatedResponse{
		User: User{
			ID:    "7f4d8a3f-3f3e-4a58-a9f0-1b0e1776b001",
			Name:  data.Name,
			Email: data.Email,
		},
		Invited: data.Invite,
	})
}

func updateUserSettings(c fiber.Ctx) error {
	data := c.Locals("data").(UpdateUserSettingsRequest)
	return c.Status(http.StatusOK).JSON(UserSettingsResponse{
		UserID:         data.ID,
		DisplayName:    data.DisplayName,
		Timezone:       data.Timezone,
		MarketingOptIn: data.MarketingOptIn,
		DryRun:         data.DryRun,
	})
}

func listUserOrders(c fiber.Ctx) error {
	data := c.Locals("data").(ListUserOrdersRequest)
	return c.Status(http.StatusOK).JSON(OrderListResponse{
		UserID:     data.ID,
		NextCursor: "cursor:page:2",
		Orders: []OrderSummary{
			{ID: "ord_1001", Status: "paid", Currency: "USD", Total: 129.5},
			{ID: "ord_1002", Status: "pending", Currency: "USD", Total: 42.0},
		},
	})
}

func main() {
	app := fiber.New()

	schema := ftonic.New(&docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Version: "1.0.0",
			Title:   "Fiber Example API",
		},
	})
	api := app.Group("/api/v1")
	users := api.Group("/users")

	ftonic.For[GetUserRequest, UserDetailsResponse](schema).
		GET(users, "/:id", middlewares.Bind[GetUserRequest], getUserByID, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Get a user profile",
			Description: "Returns a user profile with optional segments and relationship toggles.",
			Tags:        []string{"Users"},
		}))
	ftonic.For[CreateUserRequest, UserCreatedResponse](schema).
		POST(users, "/", middlewares.Bind[CreateUserRequest], createUser, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Create a user",
			Description: "Creates a new user and optionally triggers an invite flow.",
			Tags:        []string{"Users"},
		}))
	ftonic.For[UpdateUserSettingsRequest, UserSettingsResponse](schema).
		PATCH(users, "/:id/settings", middlewares.Bind[UpdateUserSettingsRequest], updateUserSettings, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Update user settings",
			Description: "Demonstrates path, query, header, and JSON body binding in one operation.",
			Tags:        []string{"Users", "Settings"},
		}))
	ftonic.For[ListUserOrdersRequest, OrderListResponse](schema).
		GET(users, "/:id/orders", middlewares.Bind[ListUserOrdersRequest], listUserOrders, ftonic.WithOperation(docs.OperationObject{
			Summary:     "List user orders",
			Description: "Shows filtered collections with pagination and regional headers.",
			Tags:        []string{"Orders"},
		}))
	// Default renderer is Swagger UI. Pass a core.UI to pick another, e.g.:
	// schema.UIHandle(app, "/docs", core.ReDoc) // or core.SwaggerUI, core.RapiDoc
	schema.UIHandle(app, "/docs", core.Scalar)

	if err := app.Listen(":1423"); err != nil {
		log.Fatalf("start fiber server: %v", err)
	}
}
