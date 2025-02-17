# API бронирования номеров в отелях

## Объем задач
- Операция по бронированию номера на нужные даты
- Предусмотрена валидация входящий данных
- Документация Swagger api/swagger/swagger.yaml

## Запуск
- скопировать `.env.example` в `.env`
- `make first-init` инициализация зависимостей проекта (docker)
- `make run` запуск приложения на 8080 порту
- `make test` запуск unit тестов
- `make local-test` запуск тестов, включая функциональные
- `make lint` запуск линтера
- `make swag` перегенерерация документации swagger

## Примечания к решению
- :trident: чистая архитектура (handler->service->repository)
- :book: Стандартная схема проекта GO
- :cd: docker compose + Makefile
- :card_file_box: документация API swagger
- :heavy_check_mark: коллекции для Postman (examples/postman)

Первичная инициализация данных при старте приложения
`
"hotel_1": {
    "room_1": { 2025-04-12, 2025-04-13, 2025-04-14 },
    "room_2": { 2025-05-9, 2025-05-10 }
}
`