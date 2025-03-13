# Микросервис для регистрации и авторизации пользователей

## Установка

- Создайте переменные среды
    > cp .env.example .env

- Создайте yaml файл 
    > cd config
    - Для локального запуска
        > cp example/local.example.yaml local.yaml
    - Для запуска в Docker
        > cp example/prod.example.yaml prod.yaml

- Запуск(локально)

  > go run cmd/app/main.go
  
  > go run cmd/migrator/main.go

## Эндпоинты

- - **POST /api/auth/v1/sign-up -** Регистрация
- - **POST /api/auth/v1/sign-in -** Авторизация
- - **POST /api/auth/v1/check-user -** Проверяет пользователя

