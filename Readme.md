# Infra

## Описание системы
Проект представляет собой систему из двух Kafka-кластеров (source и target) с настройкой репликации топиков между ними с помощью MirrorMaker 2.0. В систему также включены:
- Kafka Connect для управления коннекторами
- Schema Registry для работы со схемами данных
- Kafka UI для визуализации кластеров
- Мониторинг через Prometheus, Grafana и Alertmanager
- Интеграция с Hadoop и Spark для обработки данных

## Основные компоненты
-docker-compose.yml - описание всей инфраструктуры
-generate-certs.sh - генерация SSL-сертификатов
-setup_acls.sh - настройка прав доступа
-run-infra.sh - основной скрипт запуска
-deploy_kafka_connect.sh - управление коннекторами MirrorMaker
-mm2-*.json - конфигурации коннекторов MirrorMaker 2.0

## Последовательность развертывания
1. **Запуск инфраструктуры**
   ```bash
   chmod +x run-infra.sh
   ./run-infra.sh
   ```
2. **Генерация сертификатов** (выполняется скриптом на первом шаге)
   ```bash
   chmod +x generate-certs.sh
   ./generate-certs.sh
   ```
3. **Настройка прав доступа (ACL)** (выполняется скриптом на первом шаге)
   ```bash
   chmod +x setup_acls.sh
   ./setup_acls.sh
   ```

4. **Создание топика filtered_products в source-кластере**
   ```bash
   docker-compose exec kafka-source-0 kafka-topics.sh --create \
     --bootstrap-server kafka-source-0:9092 \
     --topic filtered_products \
     --partitions 1 \
     --replication-factor 1 \
     --command-config /bitnami/kafka/config/client-ssl.properties
   ```
5. **Развертывание MirrorMaker 2.0**
Поддерживаются следующие операции
-deploy
  первоначальная настройка
-undeploy
  откат настройки
-redeploy
  модификация (наприме, при добавлении нового топика для зеркалирования)
-status
  статус настройки

   ```bash
   chmod +x deploy_kafka_connect.sh
   ./deploy_kafka_connect.sh deploy source-0
   ```

## Проверка работы системы

1. **Проверка статуса коннекторов**
   ```bash
   ./deploy_kafka_connect.sh status source-0
   ```
2. **Доступ к интерфейсам**
-Kafka UI: http://localhost:8080
-Grafana: http://localhost:3000 (логин: admin, пароль: grafana)
-Prometheus: http://localhost:9090

## Доступы ##
1. **client-api (CN=client-api.kafka.ssl)**
Права на топик filtered_products (source):
-Read
-Describe
Права на топик client_search (source):
-Create
-Alter
-Write
-Read
-Describe
Права на топик analytics (target):
-Read
-Describe
2. **shop-api (CN=shop-api.kafka.ssl)**
Права на топик products (source):
-Create
-Alter
-Write
-Read
-Describe
Права на топик filtered_products (target):
-Create
-Alter
-Write
-Read
-Describe
Права на топик analytics (target):
-Create
-Alter
-Write
-Read
-Describe

## Замечания ##
### Перед развертыванием MirrorMaker убедитесь, что:
    -Топик filtered_products существует в source-кластере
    -Оба кластера (source и target) работают
      В композе настроен healthcheck в колонке Status должно быть примерно так Up 7 hours (healthy)
    -Настроены все необходимые ACL
### Для отладки можно использовать команды:
   ```bash
   # Просмотр логов Kafka Connect
   docker-compose logs kafka-connect-0
   # Проверка списка топиков в target-кластере
   docker-compose exec kafka-target-0 kafka-topics.sh --list \
     --bootstrap-server kafka-target-0:9092 \
     --command-config /bitnami/kafka/config/client-ssl.properties
   #Проверка ACL для SOURCE кластера:
   docker-compose exec kafka-source-0 kafka-acls.sh \
     --bootstrap-server kafka-source-0:9092 --list \
     --command-config /bitnami/kafka/config/client-ssl.properties
   #Проверка ACL для TARGET кластера:
   docker-compose exec kafka-target-0 kafka-acls.sh \
     --bootstrap-server kafka-target-0:9092 --list \
     --command-config /bitnami/kafka/config/client-ssl.properties
   ```



# Shop API

## Основные компоненты

### Генерация данных о товарах
- Модуль `generator` создает файл с данными о товарах в формате JSON.
- Поддерживает настраиваемые параметры (количество товаров, ID магазина).

### Чтение и отправка данных в Kafka
- Модуль `producer` обрабатывает взаимодействие с Kafka.
- Использует Schema Registry для сериализации данных.
- Обеспечивает SSL-соединение с Kafka (настройки указаны в `docker-compose`).

### Конфигурация
- Модуль `config` поддерживает загрузку конфигурации из файла или использование значений по умолчанию.
- Настройки Kafka, пути к файлам и другие параметры.

## Как это работает
1. Приложение генерирует файл с данными о товарах согласно указанной структуре JSON.
2. Затем читает данные из файла по одному товару, используя потоковую обработку.
3. Каждый товар сериализуется и отправляется в топик Kafka с использованием схемы из Schema Registry.

## Примеры команд
```bash
# Только для генерации файла с данными:
go run cmd/main.go --store-id=store_123 --generate-only
docker-compose run --rm kafka-producer --generate-only --store-id=store_123

# Для генерации и отправки в Kafka:
go run cmd/main.go --store-id=store_123
docker-compose run --rm kafka-producer --store-id=store_123

# С указанием пути к файлу конфигурации:
go run cmd/main.go --config=config.json
docker-compose run --rm kafka-producer --config=/app/config/config.json
```



# Client API

## Основные компоненты:
### Поиск товара по SKU:

- Подключается к Kafka-source и читает из топика filtered_products
- Ищет товар с указанным SKU
- Отображает полную информацию о найденном товаре
- Сохраняет запрос в файл и отправляет в топик client_search

### Получение персонализированных рекомендаций:

- Считывает топ часов по созданию объявлений из топика analytics в kafka-target
- Находит товары из топика filtered_products, созданные в эти часы
- Отображает рекомендованные товары
- Также сохраняет запрос в файл и отправляет в топик client_search

cmd/main.go - точка входа, обработка CLI команд
config - загрузка и хранение конфигурации
models - структуры данных
kafka - работа с Kafka (создание consumer и producer)
storage - сохранение данных о запросах
products - логика работы с товарами и рекомендациями

## Примеры команд
```bash
# Поиск по SKU
go run cmd/main.go search --sku "XYZ-12345"

# Получение рекомендаций
go run cmd/main.go recommend
```


# Banned Products System

## Основные компоненты:

### Управление списком запрещенных товаров:

```bash
# Добавление:
go run cmd/main.go add -sku "XYZ-12345" -reason "Запрещенный товар"
# Удаление:
go run cmd/main.go remove -sku "XYZ-12345"
# Обновление:
go run cmd/main.go update -sku "XYZ-12345" -reason "Новая причина запрета"
# Просмотр:
go run cmd/main.go list
```

### Потоковая обработка данных:

```bash
#Запуск:
go run cmd/main.go stream
Читает данные из топика Kafka
Проверяет каждый товар на наличие в списке запрещенных
Передает только разрешенные товары в выходной топик
```

Spark APP
# Аналитический Пайплайн для Обработки Данных о Товарах

Этот проект представляет собой аналитический пайплайн для обработки данных о товарах с использованием Kafka, HDFS и Spark. Пайплайн обеспечивает получение данных из Kafka-топика, их сохранение в HDFS и последующую аналитическую обработку с отправкой результатов обратно в Kafka.

## Архитектура пайплайна

```
Kafka Topic (filtered_products) → Consumer → HDFS → Spark Analytics → Kafka Topic (analytics)
```

### Компоненты системы:

1. **Kafka Consumer** (Go)
   - Читает сообщения из топика `filtered_products`
   - Валидирует через Schema Registry
   - Сохраняет в HDFS

2. **HDFS Writer** (Go)
   - Записывает данные в HDFS в формате JSONL
   - Файлы организованы по датам (`/data/filtered_products/YYYY-MM-DD.jsonl`)

3. **Spark Analytics** (Go)
   - Читает данные из HDFS
   - Анализирует распределение объявлений по часам
   - Находит пиковые часы

4. **Kafka Producer** (Go)
   - Отправляет результаты анализа в топик `analytics`
   - Использует Schema Registry для валидации выходных сообщений


## Примеры данных

### Входное сообщение (топик filtered_products)

```json
{
  "product_id": "12345",
  "name": "Умные часы XYZ",
  "description": "Умные часы с функцией мониторинга здоровья, GPS и уведомлениями.",
  "price": {
    "amount": 4999.99,
    "currency": "RUB"
  },
  "category": "Электроника",
  "brand": "XYZ",
  "stock": {
    "available": 150,
    "reserved": 20
  },
  "sku": "XYZ-12345",
  "tags": ["умные часы", "гаджеты", "технологии"],
  "images": [
    {
      "url": "https://example.com/images/product1.jpg",
      "alt": "Умные часы XYZ - вид спереди"
    },
    {
      "url": "https://example.com/images/product1_side.jpg",
      "alt": "Умные часы XYZ - вид сбоку"
    }
  ],
  "specifications": {
    "weight": "50g",
    "dimensions": "42mm x 36mm x 10mm",
    "battery_life": "24 hours",
    "water_resistance": "IP68"
  },
  "created_at": "2023-10-01T12:00:00Z",
  "updated_at": "2023-10-10T15:30:00Z",
  "index": "products",
  "store_id": "store_001"
}
```

### Выходное сообщение (топик analytics)

```json
{
  "hour": 12,
  "count": 148
}
```

## Запуск приложения

Приложение можно запустить в нескольких режимах:

```bash
# Запуск всего пайплайна
go run main.go

# Запуск только аналитики
go run main.go --analytics

# Запуск только consumer
go run main.go --consumer
```

## Требования

- Docker [Container Runtime]
- Docker Compose
- Go 1.24.1
- Kafka с поддержкой SSL
- Schema Registry
- HDFS
- Spark
