  To-Do List с Telegram 2FA
REST API сервис управления задачами с JWT-аутентификацией и двухфакторной верификацией через Telegram-бота.

  Стек
Go 1.24 · Echo · PostgreSQL · Docker · Docker Compose · Goose · JWT · Telegram Bot API · golangci-lint

  Возможности:

- Регистрация и авторизация пользователей с хешированием паролей (bcrypt)  
- JWT-аутентификация (access + refresh токены)  
- Двухфакторная верификация через Telegram-бота  
- CRUD операции для задач с защитой через JWT middleware  
- Отдельный микросервис Telegram-бота в Docker  
- Миграции базы данных через Goose  
- Middleware с логированием каждого запроса (request_id, метод, путь, статус, время)  
- Graceful shutdown для обоих сервисов  
- Валидация входящих данных  
- Конфигурация через переменные окружения  
- Dockerfile с multi-stage сборкой + Docker Compose  
- pprof для профилирования (опционально)  
- Makefile для удобного запуска  


Как работает авторизация  
1. POST /register        - регистрация  
2. POST /login           - если Telegram не привязан → получаешь link_code  
3. /link code            - отправляешь боту, Telegram привязывается к аккаунту  
4. POST /login           - если Telegram привязан → бот присылает 6-значный код  
5. POST /2fa/verify      - отправляешь код → получаешь JWT токены  
6. Все запросы к /tasks  - Authorization: Bearer access_token 
