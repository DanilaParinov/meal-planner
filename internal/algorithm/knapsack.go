package algorithm

import (
	"sort"

	"meal-planner/internal/models"
)

const maxDishesPerCombo = 3

func mealWeight(m models.Meal) int {
	if m.WeightG == nil {
		return 0
	}
	return *m.WeightG
}

// FindBestCombinations returns the top N combinations, sorted by the sum of
// limit-usage percentages: (totalCals/maxCals + totalWeight/maxWeight) * 100.
// If maxWeight == 0, the weight limit is not applied.
func FindBestCombinations(meals []models.Meal, maxCals, maxWeight int, includeDrinks bool, topN int) []models.MealCombination {
	valid := make([]models.Meal, 0, len(meals))
	for _, m := range meals {
		if m.IsDrink && !includeDrinks {
			continue
		}
		if m.Calories <= maxCals {
			valid = append(valid, m)
		}
	}

	sort.Slice(valid, func(i, j int) bool {
		return valid[i].Calories < valid[j].Calories
	})

	var solutions []models.MealCombination
	backtrack(valid, maxCals, maxWeight, maxDishesPerCombo, 0, 0, 0, []models.Meal{}, &solutions)

	for i := range solutions {
		calPct := float64(solutions[i].TotalCalories) / float64(maxCals)
		var weightPct float64
		if maxWeight > 0 {
			weightPct = float64(solutions[i].TotalWeight) / float64(maxWeight)
		}
		solutions[i].Score = (calPct + weightPct) * 100
	}

	sort.Slice(solutions, func(i, j int) bool {
		return solutions[i].Score > solutions[j].Score
	})

	if topN > len(solutions) {
		topN = len(solutions)
	}
	return solutions[:topN]
}

func backtrack(meals []models.Meal, maxCals, maxWeight, depthLeft, startIdx, totalCals, totalWeight int, current []models.Meal, out *[]models.MealCombination) {
	if len(current) > 0 {
		combo := make([]models.Meal, len(current))
		copy(combo, current)
		*out = append(*out, models.MealCombination{
			Meals:         combo,
			TotalCalories: totalCals,
			TotalWeight:   totalWeight,
		})
	}

	if depthLeft == 0 {
		return
	}

	for i := startIdx; i < len(meals); i++ {
		m := meals[i]
		if totalCals+m.Calories > maxCals {
			break
		}
		w := mealWeight(m)
		if maxWeight > 0 && totalWeight+w > maxWeight {
			continue
		}
		backtrack(meals, maxCals, maxWeight, depthLeft-1, i+1, totalCals+m.Calories, totalWeight+w, append(current, m), out)
	}
}
