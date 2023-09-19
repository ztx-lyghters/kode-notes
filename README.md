# kode-notes

## Запуск

Для запуска требуются `go`, `docker`, `docker-compose`, `GNU Make`
___
Для полной пересборки и запуска:

`make all` или `make all-attach`

Компиляция и сборка контейнеров:

`make build`

Запуск Postgres и сервера:

`make run` или `make run-attach`

Подключение к фоновым контейнерам:
`make attach`

## Взаимодействие с API
Скрипт interact.sh позволит проверить все методы, доступные через API.
___
Список аргументов и примеров:
`./interact.sh -h`

Регистрация нового пользователя:\
`./interact.sh -u johndoe -p qwerty1234 sign_up`

Аутентификация и получение токена авторизации:\
`./interact.sh -u johndoe -p qwerty1234 sign_in`

Зная токен, можно обратиться к двум оставшимся методам:\
`./interact.sh -T "asdfghjkl0987654321" -n foo -d bar new_note`\
`./interact.sh -T "asdfghjkl0987654321" list_notes`
___
Режимы можно комбинировать. При аутентификации, токен будет храниться в переменной скрипта для дальнейшей авторизации:\
`./interact.sh -u johndoe -p qwerty1234 sign_in list_notes`
