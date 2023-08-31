# Сервис динамического сегментирования пользователей

### Задача:

Требуется реализовать сервис, хранящий пользователя и сегменты, в которых он состоит (создание, изменение, удаление сегментов, а также добавление и удаление пользователей в сегмент)
Полное условие в [CONDITION.md](CONDITION.md)

### Стек используемых технологий:
- Go v1.20
- Роутер                - [go-chi/chi](https://github.com/go-chi/chi) 
- PostgreSQL v13.3          
- Драйвер БД            - [lib/pq](https://github.com/lib/pq)
- Валидатор пакетов     - [go-playground/validator](https://github.com/go-playground/validator)
- Логгер                - [slog](https://pkg.go.dev/golang.org/x/exp/slog)   

### Handlers:
- `user/save`          - Создание нового пользователя
- `user/delete`        - Удаление пользователя 
- `user/segments`      - Получение сегментов пользователя 
- `segment/save`       - Создание нового сегмента
- `segment/delete`     - Удаление сегмента 
- `segment/addToUser`  - Добавление пользователя в сегмент

### Запуск

```bash
    # Запустить среду для создания БД и создания таблиц
    docker compose up
    docker compose domwn

    # Запускаем еще раз 
    docker compose up
```

### Examples:
`user/save`

![usersave](/Media/usersave.jpg) 

`user/delete`

![userdelete](/Media/userdelete.jpg)

`user/segments`

![usersegments](/Media/usersegments.jpg)

`segment/save` 

![segmentsave](/Media/segmentsave.jpg)

`segment/delete` 

![segmentdelete](/Media/segmentdelete.jpg)

`segment/addToUser`

![segmentAddUser](/Media/segmentAddToUser.jpg)