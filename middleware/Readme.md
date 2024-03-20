# Echo middleware EN

## Task

You have 1 hour to solve the task. If you are lucky enough and have time left, please try to optimise your solution.

**Build a middleware using echo framework**

First of all, you should create a handler which sends how many days left until 1 Jan 2025 and response with HTTP 200 OK
status code.

Secondly, build a middleware, which checks HTTP header User-Role presents and contains admin and prints red button user
detected to the console (using default log package or any 3rd party) if so.

## Test

```
curl --location --request GET '127.0.0.1:8080/status' \
--header 'User-Role: admin'
```

## References

1. [Author repository](https://github.com/spatecon/echo-middleware-assessment.git)
2. [YouTube channel of author](https://www.youtube.com/watch?v=Lsh3ylmXdJ8&t=2037s)
3. [Repository of web framework](https://github.com/labstack/echo)
4. [Standard Go Project Layout repository](https://github.com/golang-standards/project-layout?tab=readme-ov-file)

# Echo middleware RU

## Task

На решение задачи отводится 1 час. Если Вам повезло и у Вас есть достаточно времени, пожалуйста попробуйте
оптимизировать свое решение.

**Создайте связующее программное обеспечение/обёртку/middleware с использованием фреймворка echo**

Для начала Вам следует создать обработчик, который отправляет количество оставшихся дней до 1 Января 2025 и отвечает
на запрос HTTP кодом `200 OK`.

Во-вторых, создайте middleware, который проверяет HTTP заголовок `User-Role`. Если в нём содержится строка `admin`,
тогда выведите в консоль строку "red button user detected" (используйте стандартный пакет `log` или любую другую).