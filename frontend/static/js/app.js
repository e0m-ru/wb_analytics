document.addEventListener('DOMContentLoaded', function () {
    const form = document.getElementById('searchForm');
    const submitBtn = document.getElementById('submitBtn');
    const loadingIndicator = document.getElementById('loadingIndicator');

    form.addEventListener('submit', async function (e) {
        e.preventDefault();
        const query = document.getElementById('searchQuery').value.trim();

        if (!query) return;

        // Показываем индикатор загрузки
        submitBtn.disabled = true;
        loadingIndicator.classList.remove('d-none');

        try {
            // Отправляем запрос на сервер для запуска парсера
            const response = await fetch('/api/parse', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ query })
            });

            if (!response.ok) {
                throw new Error('Ошибка сервера');
            }

            // Перенаправляем на страницу аналитики
            window.location.href = `/analytics.html?query=${encodeURIComponent(query)}`;

        } catch (error) {
            console.error('Error:', error);
            alert('Произошла ошибка при парсинге данных', error);
        } finally {
            submitBtn.disabled = false;
            loadingIndicator.classList.add('d-none');
        }
    });
});