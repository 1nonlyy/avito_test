# Avito Test Backend (Весенняя стажировка 2025)

## 🚀 Описание
Backend-сервис для сотрудников ПВЗ, реализующий учёт товаров, приёмок и пунктов выдачи заказов.

## 🔧 Функциональность

- Регистрация и авторизация пользователей (`/register`, `/login`, `/dummyLogin`)
- Создание ПВЗ (только для модераторов)
- Управление приёмками и товарами:
  - начало приёмки
  - добавление/удаление товаров (LIFO)
  - закрытие приёмки
- Получение информации по ПВЗ с фильтрацией и пагинацией

## 📦 Стек
- Go 1.24.2
- PostgreSQL 15
- Gin
- Docker & Docker Compose (OrbStack)
- SQLite для unit-тестов

## ▶️ Запуск

```bash
# Клонировать репозиторий
git clone https://github.com/1nonlyy/avito_test.git
cd avito_test

# Собрать и запустить
docker-compose up --build

# Инициализировать БД
docker exec -i $(docker ps -qf "name=avito_test-db-1") \
  psql -U postgres -d avito_test < init.sql
