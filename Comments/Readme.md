# GoNews - Comments

Сервис комментариев к новостному агрегатору.

Входные данные указаны в файле конфигураций. 

## Запуск приложения

Для запуска приложения необходимо развернуть сервер с БД mongoDB и подключиться к серверу.

БД для сервиса использовалась из сервиса GoNews. См [запуск БД](https://github.com/MoJIoToK/Go_projects/tree/main/GoNews)

## Методы проверки

Проверка производилась с помощью _Postman_:

- POST `/comments/new`, метод создает новый комментарий для новости. Тело запроса в формате JSON:

```JSON
{
  "ParentID": "{ParentID}",
  "NewsID": "{NewsID}",
  "Content": "{Content}"
} 
```

Обязательные поля - NewsID и Content.
NewsID - ID новости в формате ObjectID из MongoDB. Content - содержимое комментария. ParentID - необходим при условии
наличия родительского комментария.

- GET `/comments/{id}`, метод возвращает дерево комментариев новости по переданному ID новости. id является ObjectID из
  MongoDB.

## References

- [Сервис шлюз](https://github.com/MoJIoToK/Go_projects/tree/main/APIGateWay)
- [Сервис новостного агрегатора](https://github.com/MoJIoToK/Go_projects/tree/main/GoNews)
- [Сервис цензурирования комментариев](https://github.com/MoJIoToK/Go_projects/tree/main/Cenzor)