package models

import "time"

// Restaurant представляет ресторан
type Restaurant struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Meal представляет блюдо в меню ресторана
type Meal struct {
	ID           string    `json:"id" db:"id"`
	RestaurantID string    `json:"restaurant_id" db:"restaurant_id"`
	Name         string    `json:"name" db:"name"`
	Calories     int       `json:"calories" db:"calories"`
	Description  string    `json:"description,omitempty" db:"description"`
	Price        float64   `json:"price,omitempty" db:"price"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// User представляет тестового пользователя
type User struct {
	ID        string    `json:"id" db:"id"`
	DeviceID  string    `json:"device_id" db:"device_id"`
	APIKey    string    `json:"api_key" db:"api_key"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// MealSet представляет набор блюд (результат подбора)
type MealSet struct {
	ID            string    `json:"id" db:"id"`
	UserID        string    `json:"user_id" db:"user_id"`
	RestaurantID  string    `json:"restaurant_id" db:"restaurant_id"`
	TotalCalories int       `json:"total_calories" db:"total_calories"`
	Meals         []Meal    `json:"meals"`         // JSON array in DB
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// MealCombination используется для внутренних расчетов
type MealCombination struct {
	Meals         []Meal
	TotalCalories int
}

// API Request/Response structs

// FindMealsRequest параметры для поиска блюд
type FindMealsRequest struct {
	MaxCalories  int    `json:"max_calories" binding:"required,min=100,max=5000"`
	RestaurantID string `json:"restaurant_id" binding:"required"`
}

// SaveMealSetRequest сохранение подобранного набора
type SaveMealSetRequest struct {
	RestaurantID  string   `json:"restaurant_id" binding:"required"`
	MealIDs       []string `json:"meal_ids" binding:"required"`
	TotalCalories int      `json:"total_calories" binding:"required"`
}

// CreateRestaurantRequest создание ресторана
type CreateRestaurantRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

// CreateMealRequest добавление блюда
type CreateMealRequest struct {
	RestaurantID string  `json:"restaurant_id" binding:"required"`
	Name         string  `json:"name" binding:"required,min=1,max=255"`
	Calories     int     `json:"calories" binding:"required,min=1,max=10000"`
	Price        float64 `json:"price" binding:"min=0"`
	Description  string  `json:"description"`
}

// ErrorResponse общее сообщение об ошибке
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// SuccessResponse общее сообщение успеха
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
