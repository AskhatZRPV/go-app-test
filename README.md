# go-app-test

## Тестовое задание.

СУБД: **postgres** (github.com/lib/pq)
Роутер: **julienschmidt/httprouter**
Валидация: **go-playground/validator/v10**
Логгер: **logrus**

Для миграций используется [migrate]([https://link-url-here.org](https://github.com/golang-migrate/migrate))
Порт 8080.

Чтобы добавить запись в базу, надо сделать запрос на `/visit` (валидация присутствует)
Запрос должен выглядить вот так
```
{
    "user_id": "93c08399-f4e8-4111-9e82-462eef39304c",
    "doctor_id": "aea75c86-2692-4a84-8632-30ca55a18d5b",
    "slot": "2022-06-08 15:54:00"
}
```

Написано не очень хорошо, но мог написать и лучше.
