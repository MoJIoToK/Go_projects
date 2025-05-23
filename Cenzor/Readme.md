# GoNews - Cenzor

Сервис цензурирования комментариев к новостным статьям на наличие недопустимых выражений.

Входные данные указаны в файле конфигураций. В этот же файл записываются недопустимые выражения.

## Методы проверки

Проверка производилась с помощью _Postman_:

- POST `/`, проверка содержимого в теле запроса на наличие слов из списка в конфигурационном файле. Возвращает код 200,
  либо 400. В теле JSON формата должно присутствовать поле `content`, в котором содержится текст комментария.

## References

- [Сервис шлюз](https://github.com/MoJIoToK/Go_projects/tree/main/APIGateWay)
- [Сервис новостного агрегатора](https://github.com/MoJIoToK/Go_projects/tree/main/GoNews)
- [Сервис комментариев](https://github.com/MoJIoToK/Go_projects/tree/main/Comments)