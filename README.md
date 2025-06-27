# Сервис аналитики товаров с визуализацией данных на фронтенде. #

## Инструкция по запуску:
```bash
git clone github.com/e0m-ru/wb_analytics
cd wb_analytics
go run .
```

## 1. Backend
  ### Фаил конфигурации [config/config.go](./config/config.go)
  ```
APIAddress:      ":8081",
FrontendAddress: ":8080",
APIBaseURL:      "http://localhost:8081",
WBURL:           "https://search.wb.ru/exactmatch/ru/common/v4/search?",
Pages:           100, // количество страниц парсинга
DBPath:          "./data/products.db",
TimeOut:         time.Millisecond * 100, // чтобы не нервировать wb
  ```

  ### Парсер [internal/parser](./internal/parser/parser.go):
Парсер данных товаров с сайта Wildberries.
Открытого API WB видимо не существует, поэтому пользуемся поиском https://search.wb.ru/exactmatch/ru/common/v4/search?
    
Парсим в цикле до __config.Load().Pages=100__ постранично. Ну такой вот поиск. Если страница не вернёт товаров return __error=nil__. 

Страницы сейчас парсятся упорядоченно по популярности __?sort=popular__. Стоит сделать выбор. Не стал это всё выносить на фронтенд пока. Займусь - если будет желание.
    

  ### Хранилище [internal/storage](./internal/storage/storage.go):
Данные сохраняются в базу. Сначала хранил все запросы в отдельные таблицы, но потом решил очищать после каждого запроса. 

  ### API-эндпоинт: [api/](./api/)
Решил не делать запросы к API каждый раз при изменении фильтров. А загрузить один раз всё и фильтровать на фронтенде. Хотя API-эндпоинт работает по Т.З. с фильтрацией. [README.md](./api/README.md)

## 2. Frontend
ну там есть над чем поработать...
