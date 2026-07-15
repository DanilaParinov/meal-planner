package models

import "time"

// Restaurant represents a restaurant
type Restaurant struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Meal represents a dish on a restaurant's menu
type Meal struct {
	ID           string    `json:"id" db:"id"`
	RestaurantID string    `json:"restaurant_id" db:"restaurant_id"`
	Name         string    `json:"name" db:"name"`
	Calories     int       `json:"calories" db:"calories"`
	Description  string    `json:"description,omitempty" db:"description"`
	Price        float64   `json:"price,omitempty" db:"price"`
	WeightG      *int      `json:"weight_g,omitempty" db:"weight_g"`
	IsDrink      bool      `json:"is_drink" db:"is_drink"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// User represents a user of the application
type User struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name,omitempty" db:"name"`
	DeviceID  string    `json:"device_id" db:"device_id"`
	APIKey    string    `json:"api_key,omitempty" db:"api_key"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// MealSet represents a set of meals (a suggestion result)
type MealSet struct {
	ID            string    `json:"id" db:"id"`
	UserID        string    `json:"user_id" db:"user_id"`
	RestaurantID  string    `json:"restaurant_id" db:"restaurant_id"`
	TotalCalories int       `json:"total_calories" db:"total_calories"`
	Meals         []Meal    `json:"meals"`         // JSON array in DB
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// MealCombination is used for internal calculations
type MealCombination struct {
	Meals         []Meal  `json:"meals"`
	TotalCalories int     `json:"total_calories"`
	TotalWeight   int     `json:"total_weight"`
	Score         float64 `json:"score"`
}

// API Request/Response structs

// FindMealsRequest holds the parameters for finding meals
type FindMealsRequest struct {
	MaxCalories   int    `json:"max_calories" binding:"required,min=100,max=5000"`
	MaxWeight     int    `json:"max_weight"`
	RestaurantID  string `json:"restaurant_id" binding:"required"`
	IncludeDrinks bool   `json:"include_drinks"`
}

// SaveMealSetRequest saves a suggested meal set
type SaveMealSetRequest struct {
	RestaurantID  string   `json:"restaurant_id" binding:"required"`
	MealIDs       []string `json:"meal_ids" binding:"required"`
	TotalCalories int      `json:"total_calories" binding:"required"`
}

// CreateRestaurantRequest creates a restaurant
type CreateRestaurantRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

// CreateMealRequest adds a meal
type CreateMealRequest struct {
	RestaurantID string  `json:"restaurant_id" binding:"required"`
	Name         string  `json:"name" binding:"required,min=1,max=255"`
	Calories     int     `json:"calories" binding:"required,min=1,max=10000"`
	Price        float64 `json:"price" binding:"min=0"`
	Description  string  `json:"description"`
}

// CreateUserRequest creates a new user (admin only)
type CreateUserRequest struct {
	Name    string `json:"name" binding:"required,min=1,max=255"`
	IsAdmin bool   `json:"is_admin"`
}

// ErrorResponse is a generic error message
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// SuccessResponse is a generic success message
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
