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
    docker compose down

    # Запускаем еще раз 
    docker compose up
```

### Examples:
`user/save`
```bash
    curl --location --request POST 'http://localhost:8080/user/save'
    {
      "status":"OK",
      "id":16
    }
```

`user/delete`
```bash
    curl --location --request DELETE 'http://localhost:8080/user/delete' \
    --header 'Content-Type: application/json' \
    --data '{
              "id": 9
            }'
    {
      "status":"OK",
      "id":9
    }
```

`user/segments`

```bash 
    curl --location --request GET 'http://localhost:8080/user/segments' \
    --header 'Content-Type: application/json' \
    --data '{
        "id": 1
        }'
    {
      "status":"OK",
      "segments":{
        "UserId":1,
        "Segments":[
            {
              "ID":3,
              "Name":"test1"
            },
            {
              "ID":4,
              "Name":"test2"
            },
            {
              "ID":5,
              "Name":"test3"
            },
            {
              "ID":6,"Name":"test4"
            }
        ]
      }
    }
```

`segment/save` 

```bash
    curl --location 'http://localhost:8080/segment/save' \
    --header 'Content-Type: application/json' \
    --data '{
        "Name": "test10"
    }'
    {
      "status":"OK",
      "id":14,
      "name":"test10"
    }    
```

`segment/delete` 

```bash
    curl --location --request DELETE 'http://localhost:8080/segment/delete' \
    --header 'Content-Type: application/json' \
    --data '{
        "Name": "test9"
        }'
    {
      "status":"OK",
      "name":"test9"
    }

    
```

`segment/addToUser`

```bash
    curl --location 'http://localhost:8080/segment/addToUser' \
    --header 'Content-Type: application/json' \
    --data '{
        "SegmentsToSave": [ "test1", "test2", "test3", "test4", "test5" ],
        "SegmentsToDelete": ["qwerty", "qwerty1", "test5"],
        "UserID": 1   
    }'
    {
      "status":"OK",
      "UserId ":1,
      "NotAddedSegments":["test1","test2","test3","test4","test5"],
      "DeletedSegments ":["qwerty","qwerty1","test5"]
    }
```
