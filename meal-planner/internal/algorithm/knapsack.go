package algorithm

import (
	"sort"

	"github.com/yourusername/meal-planner/internal/models"
)

// MealKnapsack решает задачу подбора блюд в пределах калорийного лимита
// Использует динамическое программирование (knapsack problem)
type MealKnapsack struct {
	meals     []models.Meal
	maxCals   int
	solutions []models.MealCombination
}

// NewMealKnapsack создает новый решатель
func NewMealKnapsack(meals []models.Meal, maxCals int) *MealKnapsack {
	return &MealKnapsack{
		meals:     meals,
		maxCals:   maxCals,
		solutions: make([]models.MealCombination, 0),
	}
}

// Solve решает задачу и возвращает все возможные комбинации блюд
func (mk *MealKnapsack) Solve() []models.MealCombination {
	// Для MVP используем рекурсивный подход с отсеиванием
	// (более интуитивный, чем DP table для этой задачи)

	// Сортируем блюда по калорийности (для оптимизации поиска)
	sort.Slice(mk.meals, func(i, j int) bool {
		return mk.meals[i].Calories < mk.meals[j].Calories
	})

	mk.solutions = make([]models.MealCombination, 0)
	mk.backtrack([]models.Meal{}, 0, 0)

	// Сортируем результаты по близости к максимуму (чем ближе, тем лучше)
	sort.Slice(mk.solutions, func(i, j int) bool {
		diffI := mk.maxCals - mk.solutions[i].TotalCalories
		diffJ := mk.maxCals - mk.solutions[j].TotalCalories
		return diffI < diffJ
	})

	return mk.solutions
}

// SolveOptimized решает задачу с ограничением на количество результатов
// limit - максимальное количество комбинаций для возврата
func (mk *MealKnapsack) SolveOptimized(limit int) []models.MealCombination {
	mk.solutions = make([]models.MealCombination, 0)

	// Используем BFS-подобный подход для быстрого поиска лучших результатов
	type state struct {
		meals       []models.Meal
		totalCals   int
		startIdx    int // индекс, с которого начинать поиск
	}

	queue := []state{{meals: []models.Meal{}, totalCals: 0, startIdx: 0}}
	seen := make(map[string]bool) // для дедупликации

	for len(queue) > 0 && len(mk.solutions) < limit {
		current := queue[0]
		queue = queue[1:]

		// Пробуем добавить каждое блюдо, начиная с startIdx
		for i := current.startIdx; i < len(mk.meals); i++ {
			meal := mk.meals[i]
			newTotal := current.totalCals + meal.Calories

			// Если превышаем лимит, пропускаем
			if newTotal > mk.maxCals {
				continue
			}

			newMeals := append([]models.Meal{}, current.meals...)
			newMeals = append(newMeals, meal)

			// Проверяем, не видели ли мы эту комбинацию
			key := mk.makeKey(newMeals)
			if seen[key] {
				continue
			}
			seen[key] = true

			// Добавляем в результаты если калории близки к максимуму
			// или если это первое решение
			combination := models.MealCombination{
				Meals:         newMeals,
				TotalCalories: newTotal,
			}
			mk.solutions = append(mk.solutions, combination)

			// Продолжаем поиск дальше
			if newTotal < mk.maxCals {
				queue = append(queue, state{
					meals:     newMeals,
					totalCals: newTotal,
					startIdx:  i + 1,
				})
			}
		}
	}

	// Сортируем по близости к максимуму
	sort.Slice(mk.solutions, func(i, j int) bool {
		diffI := mk.maxCals - mk.solutions[i].TotalCalories
		diffJ := mk.maxCals - mk.solutions[j].TotalCalories
		return diffI < diffJ
	})

	return mk.solutions
}

// backtrack рекурсивно ищет все комбинации
func (mk *MealKnapsack) backtrack(current []models.Meal, startIdx, totalCals int) {
	// Добавляем текущее состояние в решения
	if len(current) > 0 {
		combination := models.MealCombination{
			Meals:         append([]models.Meal{}, current...),
			TotalCalories: totalCals,
		}
		mk.solutions = append(mk.solutions, combination)
	}

	// Пробуем добавить каждое следующее блюдо
	for i := startIdx; i < len(mk.meals); i++ {
		meal := mk.meals[i]
		newTotal := totalCals + meal.Calories

		// Если превышаем лимит, останавливаем эту ветку
		if newTotal > mk.maxCals {
			continue
		}

		// Рекурсивно ищем дальше
		current = append(current, meal)
		mk.backtrack(current, i+1, newTotal)
		current = current[:len(current)-1] // backtrack
	}
}

// makeKey создает уникальный ключ для комбинации (для дедупликации)
func (mk *MealKnapsack) makeKey(meals []models.Meal) string {
	key := ""
	for _, m := range meals {
		key += m.ID + ","
	}
	return key
}

// GetBestCombinations возвращает top N комбинаций
func (mk *MealKnapsack) GetBestCombinations(n int) []models.MealCombination {
	if len(mk.solutions) == 0 {
		return []models.MealCombination{}
	}

	if n > len(mk.solutions) {
		n = len(mk.solutions)
	}

	return mk.solutions[:n]
}
