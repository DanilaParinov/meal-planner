const API_BASE = 'http://localhost:8080/api';
let currentApiKey = null;
let lastSolutions = null;

// ===== Initialization =====

document.addEventListener('DOMContentLoaded', () => {
    const savedApiKey = localStorage.getItem('apiKey');
    if (savedApiKey) {
        currentApiKey = savedApiKey;
        showApp(savedApiKey);
    }

    linkSlider('caloriesSlider', 'calories', 'caloriesLabel', v => `${v} ккал`);
    linkSlider('weightSlider',   'weight',   'weightLabel',   v => `${v} г`);
});

function linkSlider(sliderId, inputId, labelId, format) {
    const slider = document.getElementById(sliderId);
    const input  = document.getElementById(inputId);
    const label  = document.getElementById(labelId);
    const update = v => { label.textContent = format(v); };
    slider.addEventListener('input', e => { input.value = e.target.value; update(e.target.value); });
    input.addEventListener('input',  e => { slider.value = e.target.value; update(e.target.value); });
}

// ===== Authentication =====

async function login() {
    const apiKey = document.getElementById('apiKey').value.trim();
    if (!apiKey) { showNotification('Введи API ключ', 'error'); return; }

    try {
        const response = await fetch(`${API_BASE}/restaurants`, {
            headers: { 'X-API-Key': apiKey }
        });
        if (!response.ok) { showNotification('Неверный API ключ', 'error'); return; }
        currentApiKey = apiKey;
        localStorage.setItem('apiKey', apiKey);
        showNotification('Успешно залогинен!', 'success');
        showApp(apiKey);
    } catch {
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
    document.getElementById('userInfo').textContent = `API Key: ${apiKey.substring(0, 10)}...`;
    await loadRestaurants(apiKey);
    await loadHistory(apiKey);
}

// ===== Restaurants =====

async function loadRestaurants(apiKey) {
    try {
        const res  = await fetch(`${API_BASE}/restaurants`, { headers: { 'X-API-Key': apiKey } });
        if (!res.ok) throw new Error();
        const data = await res.json();
        const select = document.getElementById('restaurant');
        select.innerHTML = '<option value="">-- Выбери ресторан --</option>';
        (data.data || []).forEach(r => {
            const opt = document.createElement('option');
            opt.value = r.id; opt.textContent = r.name;
            select.appendChild(opt);
        });
    } catch {
        showNotification('Ошибка загрузки ресторанов', 'error');
    }
}

// ===== Meal Search =====

async function searchMeals() {
    const restaurantID    = document.getElementById('restaurant').value;
    const maxCalories     = parseInt(document.getElementById('calories').value);
    const maxWeight       = parseInt(document.getElementById('weight').value);
    const includeDrinks   = document.getElementById('includeDrinks').checked;

    if (!restaurantID) { showNotification('Выбери ресторан', 'error'); return; }

    document.getElementById('loading').style.display = 'block';
    document.getElementById('results').innerHTML = '';

    try {
        const res = await fetch(`${API_BASE}/suggest`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', 'X-API-Key': currentApiKey },
            body: JSON.stringify({
                restaurant_id:  restaurantID,
                max_calories:   maxCalories,
                max_weight:     maxWeight,
                include_drinks: includeDrinks
            })
        });
        if (!res.ok) throw new Error();
        const data = await res.json();
        lastSolutions = data.data.solutions || [];
        displayResults(data.data, restaurantID);
        showNotification('✓ Результаты найдены', 'success');
    } catch {
        showNotification('Ошибка поиска', 'error');
    } finally {
        document.getElementById('loading').style.display = 'none';
    }
}

function displayResults(data, restaurantID) {
    const resultsDiv  = document.getElementById('results');
    const solutions   = data.solutions || [];
    const maxCals     = data.max_calories || 0;
    const maxWeight   = data.max_weight   || 0;

    if (solutions.length === 0) {
        resultsDiv.innerHTML = '<p class="empty-state">Нет подходящих комбинаций. Попробуй увеличить лимиты.</p>';
        return;
    }

    resultsDiv.innerHTML = '';

    solutions.forEach((solution, index) => {
        const meals      = solution.meals         || [];
        const totalCals  = solution.total_calories || 0;
        const totalW     = solution.total_weight   || 0;
        const score      = solution.score          || 0;

        const calPct     = maxCals   > 0 ? Math.round(totalCals / maxCals   * 100) : 0;
        const wPct       = maxWeight > 0 ? Math.round(totalW    / maxWeight  * 100) : 0;

        const mealsHTML = meals.map(meal => {
            const w = meal.weight_g != null ? `<span class="meal-weight">${meal.weight_g} г</span>` : '';
            return `<div class="meal-item">
                <span class="meal-name">${meal.name}</span>
                <span class="meal-meta">${meal.calories} ккал${w ? ' · ' + meal.weight_g + ' г' : ''}</span>
            </div>`;
        }).join('');

        const weightRow = maxWeight > 0
            ? `<span class="result-stat">${totalW} / ${maxWeight} г (${wPct}%)</span>`
            : '';

        const card = document.createElement('div');
        card.className = 'result-card';
        card.innerHTML = `
            <div class="result-header">
                <div class="result-title">Вариант ${index + 1}</div>
                <div class="result-score">${score.toFixed(0)}%</div>
            </div>
            <div class="result-stats">
                <span class="result-stat">${totalCals} / ${maxCals} ккал (${calPct}%)</span>
                ${weightRow}
            </div>
            <div class="result-meals">${mealsHTML}</div>
            <div class="result-actions">
                <button onclick="saveCollection('${restaurantID}', ${index})" class="btn btn-primary btn-small">
                    💾 Сохранить
                </button>
            </div>`;
        resultsDiv.appendChild(card);
    });
}

// ===== Saving Collections =====

async function saveCollection(restaurantID, solutionIndex) {
    if (!lastSolutions || !lastSolutions[solutionIndex]) {
        showNotification('Ошибка: решение не найдено', 'error'); return;
    }
    const solution    = lastSolutions[solutionIndex];
    const mealIDs     = solution.meals.map(m => m.id);
    const totalCalories = solution.total_calories;

    try {
        const res = await fetch(`${API_BASE}/collections`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', 'X-API-Key': currentApiKey },
            body: JSON.stringify({ restaurant_id: restaurantID, meal_ids: mealIDs, total_calories: totalCalories })
        });
        if (!res.ok) throw new Error();
        showNotification('✓ Набор сохранен!', 'success');
        await loadHistory(currentApiKey);
    } catch {
        showNotification('Ошибка сохранения набора', 'error');
    }
}

// ===== History =====

async function loadHistory(apiKey) {
    try {
        const res  = await fetch(`${API_BASE}/collections`, { headers: { 'X-API-Key': apiKey } });
        if (!res.ok) throw new Error();
        const data = await res.json();
        const historyDiv = document.getElementById('history');
        historyDiv.innerHTML = '';
        const collections = data.data || [];
        if (collections.length === 0) {
            historyDiv.innerHTML = '<p class="empty-state">История пуста</p>'; return;
        }
        collections.forEach(c => {
            const date = new Date(c.created_at).toLocaleDateString('ru-RU');
            const item = document.createElement('div');
            item.className = 'history-item';
            item.innerHTML = `
                <div class="history-date">${date}</div>
                <div class="history-calories">${c.total_calories} ккал</div>
                <div>${c.meals.length} блюд</div>`;
            historyDiv.appendChild(item);
        });
    } catch {
        // fail silently
    }
}

// ===== Utilities =====

function showNotification(message, type = 'info') {
    const n = document.createElement('div');
    n.className = `notification ${type}`;
    n.textContent = message;
    document.body.appendChild(n);
    setTimeout(() => {
        n.style.animation = 'slideOut 0.3s ease';
        setTimeout(() => n.remove(), 300);
    }, 3000);
}
