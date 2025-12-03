package controllers

import (
	"strconv"

	"controller_api_demo/services"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/web"
)

// OrderController handles order-related HTTP requests.
type OrderController struct {
	orderService services.OrderService
}

// NewOrderController creates a new OrderController.
func NewOrderController(orderService services.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

// MapRoutes registers all routes for this controller.
// Implements web.IController interface.
func (ctrl *OrderController) MapRoutes(app *web.WebApplication) {
	orders := app.MapGroup("/api/orders")
	orders.WithTags("Orders")

	// GET /api/orders/{id}
	orders.MapGet("/{id}", ctrl.GetOrder).
		WithSummary("Get order by ID").
		WithDescription("Returns a single order by its ID")

	// GET /api/orders/user/{userId}
	orders.MapGet("/user/{userId}", ctrl.GetOrdersByUser).
		WithSummary("Get orders by user ID").
		WithDescription("Returns all orders for a specific user")

	// POST /api/orders
	orders.MapPost("", ctrl.CreateOrder).
		WithSummary("Create a new order").
		WithDescription("Creates a new order in the system")
}

// GetOrder handles GET /api/orders/{id}
func (ctrl *OrderController) GetOrder(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := ctrl.orderService.GetOrder(idInt)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, order)
}

// GetOrdersByUser handles GET /api/orders/user/{userId}
func (ctrl *OrderController) GetOrdersByUser(c *gin.Context) {
	userID := c.Param("userId")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	orders, err := ctrl.orderService.GetOrdersByUser(userIDInt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, orders)
}

// CreateOrder handles POST /api/orders
func (ctrl *OrderController) CreateOrder(c *gin.Context) {
	var order services.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.orderService.CreateOrder(&order); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, order)
}

