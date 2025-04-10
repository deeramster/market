Shop API
# Только для генерации файла с данными:
go run cmd/main.go --store-id=store_123 --generate-only
docker-compose run --rm kafka-producer --generate-only --store-id=store_123

# Для генерации и отправки в Kafka:
go run cmd/main.go --store-id=store_123
docker-compose run --rm kafka-producer --store-id=store_123

# С указанием пути к файлу конфигурации:
go run cmd/main.go --config=config.json
docker-compose run --rm kafka-producer --config=/app/config/config.json

Status
docker ps --format "{{.Names}}\t{{.Status}}" | grep kafka-source-1


Просмотр списка запрещенных товаров
docker compose exec banned-products list
Добавление товара в список запрещенных
docker-compose run --rm banned-products add -sku "XYZ-12345" -reason "Запрещенный товар"
docker compose exec banned-products add -sku "XYZ-12345" -reason "Запрещенный товар"
Удаление товара из списка запрещенных
docker compose exec banned-products remove -sku "XYZ-12345"
Обновление причины запрета товара
docker compose exec banned-products update -sku "XYZ-12345" -reason "Новая причина запрета"
Запуск потокового процессора
docker compose exec banned-products stream
