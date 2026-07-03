package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"meal-planner/internal/algorithm"
	"meal-planner/internal/db"
	"meal-planner/internal/models"
)

// Handler encapsulates all HTTP request handlers
type Handler struct {
	repo *db.Repository
}

// NewHandler creates a new handler
func NewHandler(repo *db.Repository) *Handler {
	return &Handler{repo: repo}
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Public endpoints
	router.GET("/health", h.Health)

	// Protected endpoints (require authentication)
	authorized := router.Group("")
	authorized.Use(AuthMiddleware(h.repo))
	{
		// Restaurants
		authorized.GET("/api/restaurants", h.GetRestaurants)

		// Meal suggestion
		authorized.POST("/api/suggest", h.SuggestMeals)

		// History
		authorized.GET("/api/collections", h.GetCollections)
		authorized.POST("/api/collections", h.SaveCollection)

		// Admin: manual entry of restaurants and meals
		authorized.POST("/api/admin/restaurants", h.AdminCreateRestaurant)
		authorized.POST("/api/admin/meals", h.AdminCreateMeal)
	}
}

// ===== Public Endpoints =====

// Health checks whether the server is running
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "Meal Planner API is running",
	})
}

// ===== Restaurants =====

// GetRestaurants returns a list of all restaurants
// GET /api/restaurants
func (h *Handler) GetRestaurants(c *gin.Context) {
	restaurants, err := h.repo.GetAllRestaurants()
	if err != nil {
		ErrorResponseJSON(c, http.StatusInternalServerError, "database_error", "Failed to fetch restaurants")
		return
	}

	SuccessResponse(c, http.StatusOK, restaurants, "Restaurants fetched successfully")
}

// ===== Meal Suggestion =====

// SuggestMeals finds a set of meals matching the calorie limit
// POST /api/suggest
// Body: {"restaurant_id": "...", "max_calories": 1500}
func (h *Handler) SuggestMeals(c *gin.Context) {
	var req models.FindMealsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorHandler(c, err)
		return
	}

	// Check whether the restaurant exists
	restaurant, err := h.repo.GetRestaurantByID(req.RestaurantID)
	if err != nil {
		ErrorResponseJSON(c, http.StatusNotFound, "restaurant_not_found", "Restaurant does not exist")
		return
	}

	// Fetch the restaurant's meals
	meals, err := h.repo.GetMealsByRestaurant(req.RestaurantID)
	if err != nil {
		ErrorResponseJSON(c, http.StatusInternalServerError, "database_error", "Failed to fetch meals")
		return
	}

	if len(meals) == 0 {
		SuccessResponse(c, http.StatusOK, []models.MealCombination{}, "No meals available for this restaurant")
		return
	}

	// Solve the meal selection problem
	solutions := algorithm.FindBestCombinations(meals, req.MaxCalories, req.MaxWeight, req.IncludeDrinks, 20)

	// Build the response with restaurant info
	response := gin.H{
		"restaurant":      restaurant,
		"max_calories":    req.MaxCalories,
		"max_weight":      req.MaxWeight,
		"solutions_count": len(solutions),
		"solutions":       solutions,
	}

	SuccessResponse(c, http.StatusOK, response, "Meal combinations found")
}

// ===== Collections =====

// GetCollections returns the user's saved collection history
// GET /api/collections
func (h *Handler) GetCollections(c *gin.Context) {
	user := GetUserFromContext(c)
	if user == nil {
		ErrorResponseJSON(c, http.StatusUnauthorized, "unauthorized", "User not found in context")
		return
	}

	collections, err := h.repo.GetUserCollections(user.ID)
	if err != nil {
		ErrorResponseJSON(c, http.StatusInternalServerError, "database_error", "Failed to fetch collections")
		return
	}

	SuccessResponse(c, http.StatusOK, collections, "Collections fetched successfully")
}

// SaveCollection saves the selected set of meals
// POST /api/collections
// Body: {"restaurant_id": "...", "meal_ids": ["...", "..."], "total_calories": 1250}
func (h *Handler) SaveCollection(c *gin.Context) {
	user := GetUserFromContext(c)
	if user == nil {
		ErrorResponseJSON(c, http.StatusUnauthorized, "unauthorized", "User not found in context")
		return
	}

	var req models.SaveMealSetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorHandler(c, err)
		return
	}

	// Check whether the restaurant exists
	_, err := h.repo.GetRestaurantByID(req.RestaurantID)
	if err != nil {
		ErrorResponseJSON(c, http.StatusNotFound, "restaurant_not_found", "Restaurant does not exist")
		return
	}

	// Load meals and calculate calories
	var meals []models.Meal
	var totalCalories int

	for _, mealID := range req.MealIDs {
		meal, err := h.repo.GetMealByID(mealID)
		if err != nil {
			ErrorResponseJSON(c, http.StatusNotFound, "meal_not_found", "One or more meals not found")
			return
		}

		// Check that the meal belongs to the correct restaurant
		if meal.RestaurantID != req.RestaurantID {
			ErrorResponseJSON(c, http.StatusBadRequest, "meal_mismatch", "Meal does not belong to the restaurant")
			return
		}

		meals = append(meals, *meal)
		totalCalories += meal.Calories
	}

	// Save the collection
	collection := &models.MealSet{
		UserID:        user.ID,
		RestaurantID:  req.RestaurantID,
		Meals:         meals,
		TotalCalories: totalCalories,
	}

	id, err := h.repo.SaveMealCollection(collection)
	if err != nil {
		ErrorResponseJSON(c, http.StatusInternalServerError, "database_error", "Failed to save collection")
		return
	}

	response := gin.H{
		"id": id,
		"collection": collection,
	}

	SuccessResponse(c, http.StatusCreated, response, "Collection saved successfully")
}

// ===== Admin =====

// AdminCreateRestaurant creates a new restaurant
// POST /api/admin/restaurants
// Body: {"name": "..."}
func (h *Handler) AdminCreateRestaurant(c *gin.Context) {
	var req models.CreateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorHandler(c, err)
		return
	}

	restaurant, err := h.repo.CreateRestaurant(req.Name)
	if err != nil {
		ErrorResponseJSON(c, http.StatusInternalServerError, "database_error", "Failed to create restaurant")
		return
	}

	SuccessResponse(c, http.StatusCreated, restaurant, "Restaurant created successfully")
}

// AdminCreateMeal adds a meal to a restaurant
// POST /api/admin/meals
// Body: {"restaurant_id": "...", "name": "...", "calories": 350, "price": 450.0, "description": "..."}
func (h *Handler) AdminCreateMeal(c *gin.Context) {
	var req models.CreateMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorHandler(c, err)
		return
	}

	_, err := h.repo.GetRestaurantByID(req.RestaurantID)
	if err != nil {
		ErrorResponseJSON(c, http.StatusNotFound, "restaurant_not_found", "Restaurant does not exist")
		return
	}

	meal := &models.Meal{
		RestaurantID: req.RestaurantID,
		Name:         req.Name,
		Calories:     req.Calories,
		Price:        req.Price,
		Description:  req.Description,
	}

	id, err := h.repo.CreateMeal(meal)
	if err != nil {
		ErrorResponseJSON(c, http.StatusInternalServerError, "database_error", "Failed to create meal")
		return
	}

	meal.ID = id
	SuccessResponse(c, http.StatusCreated, meal, "Meal created successfully")
}
