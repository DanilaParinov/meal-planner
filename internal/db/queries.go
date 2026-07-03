package db

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"meal-planner/internal/models"
)

// Repository инкапсулирует все операции с БД
type Repository struct {
	db *sql.DB
}

// NewRepository создает новый репозиторий
func NewRepository(database *Database) *Repository {
	return &Repository{db: database.GetDB()}
}

// ===== Restaurants =====

// GetAllRestaurants возвращает все рестораны
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

// CreateRestaurant добавляет новый ресторан
func (r *Repository) CreateRestaurant(name string) (*models.Restaurant, error) {
	var rest models.Restaurant
	err := r.db.QueryRow(
		`INSERT INTO restaurants (name) VALUES ($1) RETURNING id, name, created_at`,
		name,
	).Scan(&rest.ID, &rest.Name, &rest.CreatedAt)
	return &rest, err
}

// GetRestaurantByID возвращает ресторан по ID
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

// GetMealsByRestaurant возвращает все блюда ресторана
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

// GetMealByID возвращает одно блюдо по ID
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

// CreateMeal добавляет новое блюдо
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

// GetUserByAPIKey возвращает пользователя по API ключу
func (r *Repository) GetUserByAPIKey(apiKey string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		"SELECT id, device_id, api_key, created_at FROM users WHERE api_key = $1",
		apiKey,
	).Scan(&user.ID, &user.DeviceID, &user.APIKey, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser добавляет нового пользователя
func (r *Repository) CreateUser(user *models.User) (string, error) {
	var id string
	err := r.db.QueryRow(
		`INSERT INTO users (device_id, api_key) 
		 VALUES ($1, $2) 
		 RETURNING id`,
		user.DeviceID, user.APIKey,
	).Scan(&id)

	return id, err
}

// ===== Meal Collections =====

// SaveMealCollection сохраняет набор блюд
func (r *Repository) SaveMealCollection(collection *models.MealSet) (string, error) {
	// Преобразуем []Meal в []string (ID'ы)
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

// GetUserCollections возвращает все наборы пользователя
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

		// Загружаем полные данные о блюдах
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
