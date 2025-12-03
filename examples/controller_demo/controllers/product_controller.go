package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/web"
)

// Product represents a product entity.
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// ProductService provides product operations.
type ProductService interface {
	GetProduct(id int) (*Product, error)
	ListProducts() ([]*Product, error)
	CreateProduct(product *Product) error
}

// productService is the implementation.
type productService struct {
	products map[int]*Product
}

// NewProductService creates a new ProductService.
func NewProductService() ProductService {
	return &productService{
		products: map[int]*Product{
			1: {ID: 1, Name: "Laptop", Description: "High-performance laptop", Price: 999.99},
			2: {ID: 2, Name: "Mouse", Description: "Wireless mouse", Price: 29.99},
		},
	}
}

func (s *productService) GetProduct(id int) (*Product, error) {
	product, ok := s.products[id]
	if !ok {
		return nil, nil
	}
	return product, nil
}

func (s *productService) ListProducts() ([]*Product, error) {
	products := make([]*Product, 0, len(s.products))
	for _, p := range s.products {
		products = append(products, p)
	}
	return products, nil
}

func (s *productService) CreateProduct(product *Product) error {
	s.products[product.ID] = product
	return nil
}

// ============================================
// ProductController - Controller Pattern
// ============================================

// ProductController handles product-related HTTP requests.
type ProductController struct {
	app            *web.WebApplication
	productService ProductService
}

// NewProductController creates a new ProductController.
func NewProductController(app *web.WebApplication) *ProductController {
	// Resolve service from DI container
	productService := di.GetRequiredService[ProductService](app.Services)

	return &ProductController{
		app:            app,
		productService: productService,
	}
}

// MapRoutes registers all routes for this controller.
// Implements web.IController interface.
func (ctrl *ProductController) MapRoutes(app *web.WebApplication) {
	products := app.MapGroup("/api/products")
	products.WithTags("Products")

	// GET /api/products
	products.MapGet("", ctrl.GetAll).
		WithSummary("Get all products").
		WithDescription("Returns a list of all products")

	// GET /api/products/{id}
	products.MapGet("/{id}", ctrl.GetByID).
		WithSummary("Get product by ID").
		WithDescription("Returns a single product by its ID")

	// POST /api/products
	products.MapPost("", ctrl.Create).
		WithSummary("Create a new product").
		WithDescription("Creates a new product")
}

// ============================================
// Action Methods
// ============================================

// GetAll handles GET /api/products
func (ctrl *ProductController) GetAll(c *gin.Context) {
	products, err := ctrl.productService.ListProducts()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, products)
}

// GetByID handles GET /api/products/{id}
func (ctrl *ProductController) GetByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	product, err := ctrl.productService.GetProduct(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if product == nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(200, product)
}

// Create handles POST /api/products
func (ctrl *ProductController) Create(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.productService.CreateProduct(&product); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, product)
}

