<h3 align="center">Тестовое задание на Golang Developer</h3>

### Использованные Технологии

- PostgresSQL
- Chi router
- Docker
- Golang(1.22.0)

<!-- GETTING STARTED -->
## Начало Работы

Чтобы запустить приложение следуйте следующим шагам.

### Установка

1. Клонируйте репозиторий
   ```sh
   git clone https://github.com/Uikola/carCatalog.git
   ```

2. Запустите build докерфайла:
   ```sh
   docker build -t car_catalog -f Dockerfile .
   ```

3. Создайте файл .env и скопируйте туда содержимое .example.env, внеся соответствующие изменения.

4. Запустите docker-compose
   ```sh
   docker-compose up -d
   ```

5. Приложение готово к использованию. Swagger документация располагается по адресу http://localhost:8000/swagger/


<!-- CONTACT -->
## Contact(Если возникли вопросы)

Yuri - [@telegram](https://t.me/uikola) - ugulaev806@yandex.ru

Project Link: [https://github.com/Uikola/carCatalog](https://github.com/Uikola/carCatalog)