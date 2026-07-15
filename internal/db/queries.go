package db

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"meal-planner/internal/models"
)

// Repository encapsulates all database operations
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new repository
func NewRepository(database *Database) *Repository {
	return &Repository{db: database.GetDB()}
}

// ===== Restaurants =====

// GetAllRestaurants returns all restaurants
func (r *Repository) GetAllRestaurants() ([]models.Restaurant, error) {
	rows, err := r.db.Query("SELECT id, name, created_at FROM restaurants ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []models.Restaurant
	for rows.Next() {
		var rest models.Restaurant
		if err := rows.Scan(&rest.ID, &rest.Name, &rest.CreatedAt); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, rest)
	}

	return restaurants, rows.Err()
}

// CreateRestaurant adds a new restaurant
func (r *Repository) CreateRestaurant(name string) (*models.Restaurant, error) {
	var rest models.Restaurant
	err := r.db.QueryRow(
		`INSERT INTO restaurants (name) VALUES ($1) RETURNING id, name, created_at`,
		name,
	).Scan(&rest.ID, &rest.Name, &rest.CreatedAt)
	return &rest, err
}

// GetRestaurantByID returns a restaurant by ID
func (r *Repository) GetRestaurantByID(id string) (*models.Restaurant, error) {
	var rest models.Restaurant
	err := r.db.QueryRow(
		"SELECT id, name, created_at FROM restaurants WHERE id = $1",
		id,
	).Scan(&rest.ID, &rest.Name, &rest.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("restaurant not found")
	}
	if err != nil {
		return nil, err
	}

	return &rest, nil
}

// ===== Meals =====

// GetMealsByRestaurant returns all meals for a restaurant
func (r *Repository) GetMealsByRestaurant(restaurantID string) ([]models.Meal, error) {
	rows, err := r.db.Query(
		`SELECT id, restaurant_id, name, calories, COALESCE(description, ''), COALESCE(price, 0), weight_g, is_drink, created_at
		 FROM meals
		 WHERE restaurant_id = $1
		 ORDER BY name`,
		restaurantID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []models.Meal
	for rows.Next() {
		var meal models.Meal
		if err := rows.Scan(
			&meal.ID, &meal.RestaurantID, &meal.Name, &meal.Calories,
			&meal.Description, &meal.Price, &meal.WeightG, &meal.IsDrink, &meal.CreatedAt,
		); err != nil {
			return nil, err
		}
		meals = append(meals, meal)
	}

	return meals, rows.Err()
}

// GetMealByID returns a single meal by ID
func (r *Repository) GetMealByID(id string) (*models.Meal, error) {
	var meal models.Meal
	err := r.db.QueryRow(
		`SELECT id, restaurant_id, name, calories, COALESCE(description, ''), COALESCE(price, 0), weight_g, is_drink, created_at
		 FROM meals WHERE id = $1`,
		id,
	).Scan(
		&meal.ID, &meal.RestaurantID, &meal.Name, &meal.Calories,
		&meal.Description, &meal.Price, &meal.WeightG, &meal.IsDrink, &meal.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("meal not found")
	}
	if err != nil {
		return nil, err
	}

	return &meal, nil
}

// CreateMeal adds a new meal
func (r *Repository) CreateMeal(meal *models.Meal) (string, error) {
	var id string
	err := r.db.QueryRow(
		`INSERT INTO meals (restaurant_id, name, calories, description, price) 
		 VALUES ($1, $2, $3, $4, $5) 
		 RETURNING id`,
		meal.RestaurantID, meal.Name, meal.Calories, meal.Description, meal.Price,
	).Scan(&id)

	return id, err
}

// ===== Users =====

// GetUserByAPIKey returns a user by API key
func (r *Repository) GetUserByAPIKey(apiKey string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		"SELECT id, COALESCE(name, ''), device_id, api_key, is_admin, created_at FROM users WHERE api_key = $1",
		apiKey,
	).Scan(&user.ID, &user.Name, &user.DeviceID, &user.APIKey, &user.IsAdmin, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser adds a new user, filling in the generated ID and creation time
func (r *Repository) CreateUser(user *models.User) (*models.User, error) {
	err := r.db.QueryRow(
		`INSERT INTO users (name, device_id, api_key, is_admin)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at`,
		user.Name, user.DeviceID, user.APIKey, user.IsAdmin,
	).Scan(&user.ID, &user.CreatedAt)

	return user, err
}

// GetAllUsers returns all users, without their API keys
func (r *Repository) GetAllUsers() ([]models.User, error) {
	rows, err := r.db.Query(
		"SELECT id, COALESCE(name, ''), device_id, is_admin, created_at FROM users ORDER BY created_at",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.DeviceID, &user.IsAdmin, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// ===== Meal Collections =====

// SaveMealCollection saves a meal collection
func (r *Repository) SaveMealCollection(collection *models.MealSet) (string, error) {
	// Convert []Meal to []string (IDs)
	mealIDs := make([]string, len(collection.Meals))
	for i, m := range collection.Meals {
		mealIDs[i] = m.ID
	}

	var id string
	err := r.db.QueryRow(
		`INSERT INTO meal_collections (user_id, restaurant_id, total_calories, meal_ids) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id`,
		collection.UserID, collection.RestaurantID, collection.TotalCalories, pq.Array(mealIDs),
	).Scan(&id)

	return id, err
}

// GetUserCollections returns all collections for a user
func (r *Repository) GetUserCollections(userID string) ([]models.MealSet, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, restaurant_id, total_calories, meal_ids, created_at 
		 FROM meal_collections 
		 WHERE user_id = $1 
		 ORDER BY created_at DESC 
		 LIMIT 50`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []models.MealSet
	for rows.Next() {
		var collection models.MealSet
		var mealIDs []string

		if err := rows.Scan(
			&collection.ID, &collection.UserID, &collection.RestaurantID,
			&collection.TotalCalories, pq.Array(&mealIDs), &collection.CreatedAt,
		); err != nil {
			return nil, err
		}

		// Load full meal data
		meals := make([]models.Meal, len(mealIDs))
		for i, mealID := range mealIDs {
			meal, err := r.GetMealByID(mealID)
			if err != nil {
				return nil, err
			}
			meals[i] = *meal
		}
		collection.Meals = meals

		collections = append(collections, collection)
	}

	return collections, rows.Err()
}
