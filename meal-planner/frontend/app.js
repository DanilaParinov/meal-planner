// API Configuration
const API_BASE = 'http://localhost:8080/api';
let currentApiKey = null;
let currentRestaurant = null;

// ===== Инициализация =====

document.addEventListener('DOMContentLoaded', () => {
    // Проверяем, залогинен ли пользователь
    const savedApiKey = localStorage.getItem('apiKey');
    if (savedApiKey) {
        currentApiKey = savedApiKey;
        showApp(savedApiKey);
    }

    // Слушаем изменения слайдера калорий
    const caloriesSlider = document.getElementById('caloriesSlider');
    const caloriesInput = document.getElementById('calories');

    caloriesSlider.addEventListener('input', (e) => {
        caloriesInput.value = e.target.value;
        updateCaloriesLabel();
    });

    caloriesInput.addEventListener('input', (e) => {
        caloriesSlider.value = e.target.value;
        updateCaloriesLabel();
    });
});

function updateCaloriesLabel() {
    const value = document.getElementById('caloriesSlider').value;
    document.getElementById('caloriesLabel').textContent = `${value} ккал`;
}

// ===== Аутентификация =====

async function login() {
    const apiKey = document.getElementById('apiKey').value.trim();

    if (!apiKey) {
        showNotification('Введи API ключ', 'error');
        return;
    }

    try {
        // Проверяем ключ, делая простой запрос
        const response = await fetch(`${API_BASE}/restaurants`, {
            headers: {
                'X-API-Key': apiKey
            }
        });

        if (!response.ok) {
            showNotification('Неверный API ключ', 'error');
            return;
        }

        currentApiKey = apiKey;
        localStorage.setItem('apiKey', apiKey);
        showNotification('Успешно залогинен!', 'success');
        showApp(apiKey);
    } catch (error) {
        console.error('Login error:', error);
        showNotification('Ошибка подключения к серверу', 'error');
    }
}

function logout() {
    currentApiKey = null;
    localStorage.removeItem('apiKey');
    document.getElementById('authSection').style.display = 'flex';
    document.getElementById('appSection').style.display = 'none';
    document.getElementById('apiKey').value = '';
    showNotification('Вы вышли из аккаунта', 'success');
}

async function showApp(apiKey) {
    document.getElementById('authSection').style.display = 'none';
    document.getElementById('appSection').style.display = 'block';

    // Загружаем список ресторанов
    await loadRestaurants(apiKey);

    // Показываем информацию пользователя
    document.getElementById('userInfo').textContent = `API Key: ${apiKey.substring(0, 10)}...`;

    // Загружаем историю
    await loadHistory(apiKey);
}

// ===== Рестораны =====

async function loadRestaurants(apiKey) {
    try {
        const response = await fetch(`${API_BASE}/restaurants`, {
            headers: {
                'X-API-Key': apiKey
            }
        });

        if (!response.ok) throw new Error('Failed to load restaurants');

        const data = await response.json();
        const restaurants = data.data || [];

        const select = document.getElementById('restaurant');
        select.innerHTML = '<option value="">-- Выбери ресторан --</option>';

        restaurants.forEach(restaurant => {
            const option = document.createElement('option');
            option.value = restaurant.id;
            option.textContent = restaurant.name;
            select.appendChild(option);
        });
    } catch (error) {
        console.error('Error loading restaurants:', error);
        showNotification('Ошибка загрузки ресторанов', 'error');
    }
}

// ===== Поиск блюд =====

async function searchMeals() {
    const restaurantID = document.getElementById('restaurant').value;
    const maxCalories = parseInt(document.getElementById('calories').value);

    if (!restaurantID) {
        showNotification('Выбери ресторан', 'error');
        return;
    }

    document.getElementById('loading').style.display = 'block';
    document.getElementById('results').innerHTML = '';

    try {
        const response = await fetch(`${API_BASE}/suggest`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-API-Key': currentApiKey
            },
            body: JSON.stringify({
                restaurant_id: restaurantID,
                max_calories: maxCalories
            })
        });

        if (!response.ok) {
            throw new Error('Failed to search meals');
        }

        const data = await response.json();
        displayResults(data.data, restaurantID);
        showNotification('✓ Результаты найдены', 'success');
    } catch (error) {
        console.error('Search error:', error);
        showNotification('Ошибка поиска', 'error');
    } finally {
        document.getElementById('loading').style.display = 'none';
    }
}

function displayResults(data, restaurantID) {
    const resultsDiv = document.getElementById('results');
    const solutions = data.solutions || [];

    if (solutions.length === 0) {
        resultsDiv.innerHTML = '<p class="empty-state">Нет подходящих комбинаций. Попробуй увеличить лимит калорий.</p>';
        return;
    }

    resultsDiv.innerHTML = '';

    solutions.forEach((solution, index) => {
        const card = document.createElement('div');
        card.className = 'result-card';

        const meals = solution.meals || [];
        const totalCals = solution.total_calories || 0;

        let mealsHTML = meals.map(meal => `
            <div class="meal-item">
                <span class="meal-name">${meal.name}</span>
                <span class="meal-cal">${meal.calories} ккал</span>
            </div>
        `).join('');

        const diffCals = data.max_calories - totalCals;
        const percentFull = ((totalCals / data.max_calories) * 100).toFixed(0);

        card.innerHTML = `
            <div class="result-header">
                <div class="result-title">Вариант ${index + 1}</div>
                <div class="result-calories">${totalCals} / ${data.max_calories} ккал (${percentFull}%)</div>
            </div>
            <div class="result-meals">
                ${mealsHTML}
            </div>
            <div class="result-actions">
                <button onclick="saveCollection('${restaurantID}', ${index})" class="btn btn-primary btn-small">
                    💾 Сохранить
                </button>
            </div>
        `;

        resultsDiv.appendChild(card);
    });
}

// ===== Сохранение наборов =====

let lastSolutions = null;

async function saveCollection(restaurantID, solutionIndex) {
    if (!lastSolutions || !lastSolutions[solutionIndex]) {
        showNotification('Ошибка: решение не найдено', 'error');
        return;
    }

    const solution = lastSolutions[solutionIndex];
    const mealIDs = solution.meals.map(m => m.id);
    const totalCalories = solution.total_calories;

    try {
        const response = await fetch(`${API_BASE}/collections`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-API-Key': currentApiKey
            },
            body: JSON.stringify({
                restaurant_id: restaurantID,
                meal_ids: mealIDs,
                total_calories: totalCalories
            })
        });

        if (!response.ok) {
            throw new Error('Failed to save collection');
        }

        showNotification('✓ Набор сохранен!', 'success');
        await loadHistory(currentApiKey);
    } catch (error) {
        console.error('Save error:', error);
        showNotification('Ошибка сохранения набора', 'error');
    }
}

// Сохраняем последние решения (нужно обновить displayResults)
// TODO: Улучшить это

// ===== История =====

async function loadHistory(apiKey) {
    try {
        const response = await fetch(`${API_BASE}/collections`, {
            headers: {
                'X-API-Key': apiKey
            }
        });

        if (!response.ok) throw new Error('Failed to load history');

        const data = await response.json();
        const collections = data.data || [];

        const historyDiv = document.getElementById('history');
        historyDiv.innerHTML = '';

        if (collections.length === 0) {
            historyDiv.innerHTML = '<p class="empty-state">История пуста</p>';
            return;
        }

        collections.forEach(collection => {
            const date = new Date(collection.created_at).toLocaleDateString('ru-RU');
            const item = document.createElement('div');
            item.className = 'history-item';
            item.innerHTML = `
                <div class="history-date">${date}</div>
                <div class="history-calories">${collection.total_calories} ккал</div>
                <div>${collection.meals.length} блюд</div>
            `;
            historyDiv.appendChild(item);
        });
    } catch (error) {
        console.error('History error:', error);
        // Тихо падаем, если история не загружается
    }
}

// ===== Утилиты =====

function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;

    document.body.appendChild(notification);

    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease';
        setTimeout(() => notification.remove(), 300);
    }, 3000);
}

// Правим displayResults для сохранения решений
const originalDisplayResults = displayResults;
displayResults = function(data, restaurantID) {
    lastSolutions = data.solutions || [];
    originalDisplayResults.call(this, data, restaurantID);
};
