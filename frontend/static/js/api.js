document.addEventListener('DOMContentLoaded', function () {
    // Текущие параметры сортировки
    let currentSort = {
        field: 'price',
        direction: 'desc'
    };

    let allProducts = []; // все загруженные продукты

    // Элементы DOM
    const productsTable = document.getElementById('productsTable');
    const totalCount = document.getElementById('totalCount');
    const lastUpdate = document.getElementById('lastUpdate');
    const minPriceInput = document.getElementById('minPrice');
    const maxPriceInput = document.getElementById('maxPrice');
    const priceRange = document.getElementById('priceRange');
    const minRatingSelect = document.getElementById('minRating');
    const minFeedbacksInput = document.getElementById('minFeedbacks');
    const sortableHeaders = document.querySelectorAll('.sortable');

    // Инициализация
    initEventListeners();
    loadProducts();

    function initEventListeners() {
        // Фильтры
        minPriceInput.addEventListener('input', debounce(updateProducts, 300));
        maxPriceInput.addEventListener('input', debounce(updateProducts, 300));
        priceRange.addEventListener('input', updatePriceInputs);
        minRatingSelect.addEventListener('change', updateProducts);
        minFeedbacksInput.addEventListener('input', debounce(updateProducts, 300));

        // Сортировка
        sortableHeaders.forEach(header => {
            header.addEventListener('click', function () {
                const field = this.dataset.sort;

                if (currentSort.field === field) {
                    currentSort.direction = currentSort.direction === 'asc' ? 'desc' : 'asc';
                } else {
                    currentSort.field = field;
                    currentSort.direction = 'asc';
                }

                updateSortIcons();
                updateProducts();
            });
        });
    }

    function updatePriceInputs() {
        const value = parseInt(priceRange.value);
        minPriceInput.value = Math.floor(value * 0.3);
        maxPriceInput.value = Math.floor(value * 0.7);
        updateProducts();
    }

    function updateSortIcons() {
        sortableHeaders.forEach(header => {
            const icon = header.querySelector('.sort-icon');
            if (header.dataset.sort === currentSort.field) {
                icon.textContent = currentSort.direction === 'asc' ? '↑' : '↓';
            } else {
                icon.textContent = '↕';
            }
        });
    }

    function loadProducts() {
        fetch('/api/products')
            .then(response => response.json())
            .then(data => {
                allProducts = data.products;
                updateProducts();
                totalCount.textContent = `Всего товаров: ${allProducts.length}`;
                lastUpdate.textContent = `Последнее обновление: ${new Date().toLocaleTimeString()}`;
            })
            .catch(error => {
                console.error('Error:', error);
                productsTable.innerHTML = `<tr><td colspan="5" class="text-center text-danger">Ошибка загрузки данных</td></tr>`;
            });
    }

    function updateProducts() {
        // Фильтрация
        let filteredProducts = allProducts.filter(product => {
            const meetsPrice = (!minPriceInput.value || product.price >= parseFloat(minPriceInput.value)) &&
                (!maxPriceInput.value || product.price <= parseFloat(maxPriceInput.value));
            const meetsRating = !minRatingSelect.value || product.rating >= parseFloat(minRatingSelect.value);
            const meetsFeedbacks = !minFeedbacksInput.value || product.feedbacks >= parseInt(minFeedbacksInput.value);

            return meetsPrice && meetsRating && meetsFeedbacks;
        });

        // Сортировка
        filteredProducts.sort((a, b) => {
            const field = currentSort.field;
            const direction = currentSort.direction === 'asc' ? 1 : -1;

            // Для числовых полей
            if (['price', 'sale_price', 'rating', 'feedbacks'].includes(field)) {
                return (a[field] - b[field]) * direction;
            }

            // Для строковых полей (например, name)
            return a[field].localeCompare(b[field]) * direction;
        });

        renderProducts(filteredProducts);
        renderProducts(filteredProducts);
        updateCharts(filteredProducts);
        totalCount.textContent = `Показано товаров: ${filteredProducts.length} из ${allProducts.length}`;
    }

    function renderProducts(products) {
        if (products.length === 0) {
            productsTable.innerHTML = '<tr><td colspan="5" class="text-center">Товары не найдены</td></tr>';
            return;
        }

        let html = '';
        products.forEach(product => {
            const ratingClass = getRatingClass(product.rating);

            html += `
                <tr>
                    <td>${product.name}</td>
                    <td>${product.price.toLocaleString('ru-RU')} ₽</td>
                    <td class="text-success fw-bold">${product.sale_price.toLocaleString('ru-RU')} ₽</td>
                    <td class="${ratingClass}">${product.rating.toFixed(1)}</td>
                    <td>${product.feedbacks.toLocaleString('ru-RU')}</td>
                </tr>
            `;
        });

        productsTable.innerHTML = html;
    }
    let priceChart = null;
    let discountRatingChart = null;

    function updateCharts(products) {
        updatePriceHistogram(products);
        updateDiscountRatingChart(products);
    }

    function updatePriceHistogram(products) {
        const ctx = document.getElementById('priceHistogram').getContext('2d');

        // Рассчитываем динамические диапазоны
        const prices = products.map(p => p.price);
        const minPrice = Math.min(...prices);
        const maxPrice = Math.max(...prices);
        const rangeSize = (maxPrice - minPrice) / 10;

        // Создаем 10 динамических диапазонов
        const priceRanges = [];
        for (let i = 0; i < 10; i++) {
            const rangeMin = minPrice + (i * rangeSize);
            const rangeMax = minPrice + ((i + 1) * rangeSize);
            priceRanges.push({
                min: rangeMin,
                max: rangeMax,
                label: `${Math.round(rangeMin / 1000)}k-${Math.round(rangeMax / 1000)}k`
            });
        }

        // Считаем количество товаров в каждом диапазоне
        const data = priceRanges.map(range => {
            return products.filter(p => p.price >= range.min && p.price < range.max).length;
        });

        const labels = priceRanges.map(range => range.label);

        if (priceChart) {
            priceChart.data.labels = labels;
            priceChart.data.datasets[0].data = data;
            priceChart.update();
        } else {
            priceChart = new Chart(ctx, {
                type: 'bar',
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'Количество товаров',
                        data: data,
                        backgroundColor: 'rgba(54, 162, 235, 0.7)',
                        borderColor: 'rgba(54, 162, 235, 1)',
                        borderWidth: 1
                    }]
                },
                options: {
                    responsive: true,
                    scales: {
                        y: {
                            beginAtZero: true,
                            title: {
                                display: true,
                                text: 'Количество товаров'
                            }
                        },
                        x: {
                            title: {
                                display: true,
                                text: 'Диапазон цен (руб)'
                            }
                        }
                    },
                    plugins: {
                        tooltip: {
                            callbacks: {
                                label: function (context) {
                                    const range = priceRanges[context.dataIndex];
                                    return [
                                        `Диапазон: ${Math.round(range.min)}-${Math.round(range.max)} руб`,
                                        `Товаров: ${context.raw}`
                                    ];
                                }
                            }
                        }
                    }
                }
            });
        }
    }

    function updateDiscountRatingChart(products) {
        const ctx = document.getElementById('discountRatingChart').getContext('2d');

        // Рассчитываем размер скидки
        const discounts = products.map(p => {
            return ((p.price - p.sale_price) / p.price * 100).toFixed(0);
        });

        const ratings = products.map(p => p.rating);

        if (discountRatingChart) {
            discountRatingChart.data.labels = ratings;
            discountRatingChart.data.datasets[0].data = discounts;
            discountRatingChart.update();
        } else {
            discountRatingChart = new Chart(ctx, {
                type: 'line',
                data: {
                    labels: ratings,
                    datasets: [{
                        label: 'Размер скидки (%)',
                        data: discounts,
                        backgroundColor: 'rgba(75, 192, 192, 0.2)',
                        borderColor: 'rgba(75, 192, 192, 1)',
                        borderWidth: 2,
                        tension: 0.1,
                        pointRadius: 4
                    }]
                },
                options: {
                    responsive: true,
                    scales: {
                        y: {
                            beginAtZero: true,
                            title: {
                                display: true,
                                text: 'Размер скидки (%)'
                            }
                        },
                        x: {
                            title: {
                                display: true,
                                text: 'Рейтинг товара'
                            }
                        }
                    }
                }
            });
        }
    }
    function getRatingClass(rating) {
        if (rating >= 4.5) return 'rating-high';
        if (rating >= 3.5) return 'rating-medium';
        return 'rating-low';
    }

    // Утилиты
    function debounce(func, timeout = 300) {
        let timer;
        return (...args) => {
            clearTimeout(timer);
            timer = setTimeout(() => { func.apply(this, args); }, timeout);
        };
    }
});