# GoNews - APIGateWay

Сервис шлюз перенаправляет запросы к различным сервисам системы.

Входные данные указаны в файле конфигураций.

## Методы проверки

Проверка производилась с помощью _Postman_:

- GET `/news?s={query}&page={num}&limit={count}`, где query - поисковой запрос, num - номер страницы (по-умолчанию 1),
  count - количество новостей на страницу. Возвращает все статьи с пагинацией, соответствующие параметрам.
- GET `/news/id/{id}`, id - id - идентификатор ObjectID новостной статьи, берется из БД MongoDB. Возвращает статью с
  переданным ID и дерево комментариев к новости.
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

## References

- [Сервис новостного агрегатора](https://github.com/MoJIoToK/learning_go/tree/master/GoNews)
- [Сервис комментариев](https://github.com/MoJIoToK/learning_go/tree/master/Comments)
- [Сервис цензурирования комментариев](https://github.com/MoJIoToK/learning_go/tree/master/Cenzor)
