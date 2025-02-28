# Сервис лидерборда в реальном времени

Это серверное приложение для системы лидерборда в реальном времени, написанное на Go с использованием PostgreSQL, Redis и REST API. Пользователи могут регистрироваться, отправлять очки, просматривать рейтинги и получать отчёты о лучших игроках, с обновлениями в реальном времени через Redis Pub/Sub.

## Возможности
- Регистрация пользователей и аутентификация с помощью JWT.
- Отправка очков с обновлением лидерборда в реальном времени (перезапись текущего счёта).
- Глобальный лидерборд с пагинацией.
- Получение ранга пользователя.
- Отчёт о лучших игроках за период (день, неделя, месяц).
- Обновления в реальном времени через Redis Pub/Sub с топ-5 лидерборда.

## Технологии
- Go: Язык программирования для сервера.
- PostgreSQL: Постоянное хранилище для пользователей и истории очков.
- Redis: Хранилище в памяти для лидерборда (Sorted Sets) и сообщений Pub/Sub.
- Docker: Контейнеризация для удобного запуска.

## Требования
- Установленные Docker и Docker Compose.
- Опционально: Go, PostgreSQL и Redis для локального запуска без Docker.

## Установка и запуск

### С использованием Docker
1. Склонируйте репозиторий:
   git clone <url-репозитория>
   cd leaderboard
2. Запустите сервисы:
   docker-compose up --build
3. Сервер будет доступен на http://localhost:8080.

### Без Docker
1. Установите зависимости:
   - Go: Скачать с https://golang.org/dl/
   - PostgreSQL: Скачать с https://www.postgresql.org/download/
   - Redis: Скачать для Windows с https://github.com/tporadowski/redis/releases
2. Настройте PostgreSQL:
   psql -U postgres
   CREATE DATABASE leaderboard;
   CREATE USER user WITH PASSWORD 'password';
   GRANT ALL PRIVILEGES ON DATABASE leaderboard TO user;
3. Запустите Redis:
   cd <папка-redis>
   redis-server.exe
4. Запустите приложение:
   cd leaderboard
   go mod tidy
   go run main.go

## Эндпоинты API

Все эндпоинты, кроме /register и /login, требуют JWT-токен в заголовке Authorization: Bearer <token>.

### 1. POST /register
Регистрация нового пользователя.

- Запрос:
  curl -X POST -H "Content-Type: application/json" -d "{\"username\":\"player1\",\"password\":\"pass1\"}" http://localhost:8080/register
- Ответ (201 Created):
  {"message":"Пользователь зарегистрирован"}
- Ошибка (409 Conflict):
  {"error":"Пользователь уже существует или ошибка","code":409}

### 2. POST /login
Вход для получения JWT-токена.

- Запрос:
  curl -X POST -H "Content-Type: application/json" -d "{\"username\":\"player1\",\"password\":\"pass1\"}" http://localhost:8080/login
- Ответ (200 OK):
  {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}
- Ошибка (401 Unauthorized):
  {"error":"Неверные учетные данные","code":401}

### 3. POST /submit-score
Отправка очков для пользователя (перезаписывает текущий счёт в Redis, требует JWT).

- Запрос:
  curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer <token>" -d "{\"game\":\"game1\",\"score\":500}" http://localhost:8080/submit-score
- Ответ (200 OK):
  {"message":"Очки успешно обновлены"}
- Ошибка (401 Unauthorized):
  {"error":"Требуется токен авторизации","code":401}
- Ошибка (400 Bad Request):
  {"error":"Неверный запрос","code":400}

### 4. GET /leaderboard
Получение глобального лидерборда с пагинацией (требует JWT).

- Запрос:
  curl -H "Authorization: Bearer <token>" "http://localhost:8080/leaderboard?offset=0&limit=5"
- Ответ (200 OK):
  [
      {"username":"player10","score":1000},
      {"username":"player9","score":900},
      {"username":"player8","score":800},
      {"username":"player7","score":700},
      {"username":"player6","score":600}
  ]
- Ошибка (400 Bad Request):
  {"error":"Неверный параметр limit","code":400}
- Ошибка (401 Unauthorized):
  {"error":"Требуется токен авторизации","code":401}

### 5. GET /rank
Получение ранга пользователя в лидерборде (требует JWT).

- Запрос:
  curl -H "Authorization: Bearer <token>" "http://localhost:8080/rank?username=player1"
- Ответ (200 OK):
  {"rank":1}
- Ошибка (404 Not Found):
  {"error":"Пользователь не найден в лидерборде","code":404}
- Ошибка (400 Bad Request):
  {"error":"Не указан username","code":400}

### 6. GET /top-players
Получение топа игроков за период (день, неделя, месяц; требует JWT).

- Запрос:
  curl -H "Authorization: Bearer <token>" "http://localhost:8080/top-players?period=week"
- Ответ (200 OK):
  [
      {"username":"player1","total_score":600},
      {"username":"player2","total_score":200}
  ]
- Ошибка (400 Bad Request):
  {"error":"Неверный период. Используйте: day, week, month","code":400}
- Ошибка (401 Unauthorized):
  {"error":"Требуется токен авторизации","code":401}

## Обновления в реальном времени
- Сервис использует Redis Pub/Sub для отправки уведомлений при обновлении очков.
- Пример сообщения:
  app-1  | Обновление лидерборда: player1 updated score to 500, top-5: [{"username":"player1","score":500},...]
- Подписчик в фоновом режиме логирует эти обновления.

## Примечания
- Redis хранит последний счёт каждого пользователя (перезаписывается через ZAdd).
- PostgreSQL сохраняет полную историю отправок очков, используемую для /top-players.
- Все эндпоинты, кроме /register и /login, требуют JWT-аутентификацию.

## Итоги
Проект успешно реализует систему лидерборда с аутентификацией через JWT, отправкой очков, пагинированным лидербордом, рангами и отчётами о лучших игроках за период. Обновления в реальном времени реализованы через Redis Pub/Sub, что позволяет подписчикам получать топ-5 после каждого изменения счёта. Тестирование производительности показало, что система обрабатывает 100 запросов на отправку очков за 1.364 секунды (примерно 73 запроса в секунду), демонстрируя хорошую скорость для базовой реализации.