# ТЕСТОВЫЙ ПРОЕТ GO LANG

### Тестовое задание
https://github.com/avito-tech/job-backend-trainee-assignment

### API документация
https://documenter.getpostman.com/view/20010603/VUxKTV9e#6ff4e989-c940-48ee-947f-dec6dcc0c83c <hr />

**Запуск приложения**<p>
$ *make start* <hr />
**Остановка приложения** <p>
$ *make stop*
<hr />


## Тонкости реалзиации 

Для пагинации используется GET параметры - `limit`, `offset`

Для сортировки использует GET параметр `order`, где `_` является разделителем имени колонки сортировки и
метода сортировки. Примеры: "date_asc", "date_desc", "amount_desc"


