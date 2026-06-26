package main

import (
	m "gin_example/middlewares"
	"log"
	"net/http"

	gtonic "github.com/TickLabVN/tonic/adapters/gin"
	"github.com/TickLabVN/tonic/core"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gin-gonic/gin"
)

type GetUserRequest struct {
	ID             string `uri:"id" binding:"required,uuid4"`
	IncludeOrders  bool   `form:"includeOrders"`
	IncludeMetrics bool   `form:"includeMetrics"`
	APIKey         string `header:"x-api-key" binding:"required,min=10"`
}

type User struct {
	ID    string `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type CreateUserRequest struct {
	Name      string `json:"name" binding:"required,min=2,max=60"`
	Email     string `json:"email" binding:"required,email"`
	Role      string `json:"role" binding:"required,oneof=admin manager member"`
	Invite    bool   `json:"invite"`
	RequestID string `header:"x-request-id" binding:"required,uuid4"`
}

type UpdateUserSettingsRequest struct {
	ID             string `uri:"id" binding:"required,uuid4"`
	DryRun         bool   `form:"dryRun"`
	RequestID      string `header:"x-request-id" binding:"required,uuid4"`
	DisplayName    string `json:"displayName" binding:"required,min=2,max=50"`
	Timezone       string `json:"timezone" binding:"required"`
	MarketingOptIn bool   `json:"marketingOptIn"`
}

type ListUserOrdersRequest struct {
	ID       string `uri:"id" binding:"required,uuid4"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Cursor   string `form:"cursor"`
	Region   string `header:"x-region" binding:"required,oneof=apac emea us"`
	Statuses string `form:"statuses"`
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

func getUserByID(c *gin.Context) {
	data := c.MustGet("data").(GetUserRequest)
	segments := []string{"starter", "beta"}
	if data.IncludeMetrics {
		segments = append(segments, "usage-insights")
	}
	c.JSON(http.StatusOK, UserDetailsResponse{
		User: User{
			ID:    data.ID,
			Name:  "John Doe",
			Email: "john.doe@example.com",
		},
		Segments: segments,
	})
}

func createUser(c *gin.Context) {
	data := c.MustGet("data").(CreateUserRequest)
	c.JSON(http.StatusCreated, UserCreatedResponse{
		User: User{
			ID:    "7f4d8a3f-3f3e-4a58-a9f0-1b0e1776b001",
			Name:  data.Name,
			Email: data.Email,
		},
		Invited: data.Invite,
	})
}

func updateUserSettings(c *gin.Context) {
	data := c.MustGet("data").(UpdateUserSettingsRequest)
	c.JSON(http.StatusOK, UserSettingsResponse{
		UserID:         data.ID,
		DisplayName:    data.DisplayName,
		Timezone:       data.Timezone,
		MarketingOptIn: data.MarketingOptIn,
		DryRun:         data.DryRun,
	})
}

func listUserOrders(c *gin.Context) {
	data := c.MustGet("data").(ListUserOrdersRequest)
	c.JSON(http.StatusOK, OrderListResponse{
		UserID:     data.ID,
		NextCursor: "cursor:page:2",
		Orders: []OrderSummary{
			{ID: "ord_1001", Status: "paid", Currency: "USD", Total: 129.5},
			{ID: "ord_1002", Status: "pending", Currency: "USD", Total: 42.0},
		},
	})
}

func main() {
	gin.SetMode(gin.ReleaseMode) // Set Gin to release mode for production
	g := gin.Default()
	if err := g.SetTrustedProxies(nil); err != nil {
		log.Fatalf("configure trusted proxies: %v", err)
	}
	schema := gtonic.New(&docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Version: "1.0.0",
			Title:   "Gin Example API",
			Contact: &docs.ContactObject{
				Name:  "Author",
				URL:   "https://github.com/phucvinh57",
				Email: "npvinh0507@gmail.com",
			},
		},
	})
	api := g.Group("/api/v1")
	ug := api.Group("/users")
	gtonic.For[GetUserRequest, UserDetailsResponse](schema).
		GET(ug, "/:id", m.Bind[GetUserRequest], getUserByID, gtonic.WithOperation(docs.OperationObject{
			Summary:     "Get a user profile",
			Description: "Returns a user profile with optional segments and relationship toggles.",
			Tags:        []string{"Users"},
		}))
	gtonic.For[CreateUserRequest, UserCreatedResponse](schema).
		POST(ug, "/", m.Bind[CreateUserRequest], createUser, gtonic.WithOperation(docs.OperationObject{
			Summary:     "Create a user",
			Description: "Creates a new user and optionally triggers an invite flow.",
			Tags:        []string{"Users"},
		}))
	gtonic.For[UpdateUserSettingsRequest, UserSettingsResponse](schema).
		PATCH(ug, "/:id/settings", m.Bind[UpdateUserSettingsRequest], updateUserSettings, gtonic.WithOperation(docs.OperationObject{
			Summary:     "Update user settings",
			Description: "Demonstrates path, query, header, and JSON body binding in one operation.",
			Tags:        []string{"Users", "Settings"},
		}))
	gtonic.For[ListUserOrdersRequest, OrderListResponse](schema).
		GET(ug, "/:id/orders", m.Bind[ListUserOrdersRequest], listUserOrders, gtonic.WithOperation(docs.OperationObject{
			Summary:     "List user orders",
			Description: "Shows filtered collections with pagination and regional headers.",
			Tags:        []string{"Orders"},
		}))
	// Default renderer is Swagger UI. Pass a core.UI to pick another, e.g.:
	// schema.UIHandle(g, "/docs", core.ReDoc)   // or core.Scalar, core.RapiDoc
	schema.UIHandle(g, "/docs", core.Scalar, core.SwaggerUI)

	g.Run(":1234")
}
