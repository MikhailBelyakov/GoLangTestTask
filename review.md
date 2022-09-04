1. main.go

```
func Migrate() {
    users.AutoMigrate()
    balances.AutoMigrate()
    transactions.AutoMigrate()
}
```
Замалчиваются ошибки в миграциях. Ошибки всё же лучше обрабатывать.+


```
db := common.Init()
```
Замалчивается ошибка подключения к БД. Точнее, она пишется, как `log.Println` внутри функции `Init`, но всё же `Println` - это не тот уровень лога, сервис всё равно попытается стартануть при провале коннекта к БД.

2. dockerfile


```
FROM golang
```
Использовать latest не рекомендуется. Не только в golang, в любом стеке.

```
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/gorilla/mux
```
Тут не понял, зачем ты делаешь отдельно `go get`, когда обе этих зависимости прописаны в `go.mod`?

Аналог `composer install` для голенга звучит как `go mod tidy`, его кстати не увидел.
```
ADD . /go/src
```
При использовании `ADD . ` уместно иметь в репозитории настроенный `.dockerignore`.


3. Makefile

Потерялась `make help`.

4. common/database.go


```
func Init() *gorm.DB {
    // ...
    db.Exec("SET NAMES 'utf8mb4'; SET CHARACTER SET utf8mb4;")
    db.Exec("ALTER TABLE transactions MODIFY operation TEXT CHARACTER SET utf8mb4;")
    // ...
}
```
В GIN есть встроенный механизм миграций, всё лучше через него.

```
func Init() *gorm.DB {
    db, err := gorm.Open(mysql.New(mysql.Config{
        // ...
        DSN: "testuser:testpassword@tcp(db:3306)/testing?parseTime=true&loc=Local&charset=utf8mb4", // data source name
        // ...
    })
}
```
Креды лучше передавать через .env или иными путями, не хардкодить.


4. common/request.go

```
func sortMethods() map[string]bool {
	return map[string]bool{
		"asc":  true,
		"ASC":  true,
		"desc": true,
		"DESC": true,
	}
}

if _, inMap := sortMethods()[orderSplit[1]]; !inMap {
    sortStruct.Sort = defaultSortMethod
}
```

Кмк проще будет как-то так:

```
func sortMethods() map[string]bool {
    return map[string]bool{
        "ASC":  true,
        "DESC": true,
    }
}
if _, inMap := sortMethods()[strings.ToUpper(orderSplit[1])]; !inMap {
    sortStruct.Sort = defaultSortMethod
}
```

Далее. Весь код целиком. Оставил вопросы в комментах. Тут кмк ошибка.

```
func (sortStruct *SortStruct) BindParam(context *gin.Context) {
	orderString := context.Query("order")
	orderSplit := strings.Split(orderString, "_")

	if len(orderSplit) > 1 {
		sortStruct.Order = orderSplit[0]

		// Если не нашли orderSplit[1] в списке разрешённых - ставим дефолтный
		if _, inMap := sortMethods()[orderSplit[1]]; !inMap {
			sortStruct.Sort = defaultSortMethod
		}

		// Независимо от предыдущей проверки назначаем пришедшее значение, даже если его нет в мапе.
		sortStruct.Sort = orderSplit[1]
	} else {
		sortStruct.Order = defaultOrderField
		sortStruct.Sort = defaultSortMethod
	}
}
```

6. balances/models.go, transactions/models.go, users/models.go

```
type Tabler interface {
    TableName() string
}
```

Надо проверить конечно. Но по-моему документация GORM вводит в заблуждение в отношении этого интерфейса.

Сомневаюсь, что они проверяют, объявил ли ты локальный интерфейс Tabler с нужной сигнатурой, и объявляет ли его твоя модель. 

У них в пакете gorm.io/gorm/schema есть свой Tabler, и скорее всего сверка идёт именно с ним через тайпкаст.

Для объявления Tabler'а из пакета gorm.io/gorm/schema достаточно реализовать на модели сам метод TableName, больше ничего. Интерфейс Tabler в локальном пакете скорее всего не нужен.


7. transactions/models.go


```
func FindTransactionsByUser(userModel users.UserModel, condition GetTransactionParamStruct) []TransactionModel {
	var transactionModels []TransactionModel

	db := common.GetDB()
	db.Where(&TransactionModel{
		UserID: userModel.ID,
	}).Limit(condition.Limit).Order(condition.Order + " " + condition.Sort).Offset(condition.Offset).Find(&transactionModels)

	return transactionModels
}
```
Надо проверить, но выглядит как SQL-инъекция через поле condition.Order.

condition.Sort хотя бы проверяется на список разрешённых значений (но там ошибка).

Так же игнорируется возможная ошибка БД.

UPD: нет, нашёл orderFieldAllow(). Но всё равно тонкий лёд.


8. transactions/models.go

Из официальных рекомендаций языка: ID всегда пишем капсом.

https://github.com/golang/go/wiki/CodeReviewComments#initialisms

```
type TransactionModel struct {
	UserID     uint
	TypeID     int
	ReceiverId uint `gorm:"default: null"`
	SenderId   uint `gorm:"default: null"`
}
```

9. Есть возможность отправлять переводы самому себе и дюпать на этом деньги. Учитывается только приход, без соответствующего ему расхода.
```
$ curl -XGET 0:8080/api/balances/1
{"amount":11500}

$ curl -XPOST 0:8080/api/balances/1/sendTo/1  --data 'amount=1000'
{"message":"Отправка прошла успешно"}

$ curl -XPOST 0:8080/api/balances/1/sendTo/1  --data 'amount=1000'
{"message":"Отправка прошла успешно"}

$ curl -XPOST 0:8080/api/balances/1/sendTo/1  --data 'amount=1000'
{"message":"Отправка прошла успешно"}

$ curl -XGET 0:8080/api/balances/1
{"amount":14500}
```

10. ExchangeBetweenUsers

Есть вероятность создать списание у отправителя, но не успеть создать зачисление у получателя.

Ошибка создания списания не проверяется.

Уместно использование транзакции на уровне БД.

11. http.Get

Не рекомендовал бы использовать без указания таймаута на выполнение запроса. Запрос создаётся без контекста, соответственно даже прерывание обработки запроса на уровне http-сервера не отменит отправленный таким образом исходящий запрос и возможна утечка ресурсов.

12. Контекст!

В дополнение к предыдущему пункту - ещё до вызова getCurrency где-то потеряли контекст. 

Не только в этом методе, везде теряем его на уровне контроллера.

Официальная рекомендация языка - передавать context.Context между вызываемыми функциями.

https://github.com/golang/go/wiki/CodeReviewComments#contexts

13. Конкурентность.

Если одномоментно (буквально в одну милисекунду) придёт 100 запросов на получение баланса пользователя, у которого ранее не было баланса - то только 1 из них завершится успешно, остальные 99 скорее всего сфейлят при FirstOrCreate.

Аналогичная проблема: если одномоментно придёт 100 запросов на списание, то все проверки на уровне приложения, которые ты проводишь - не отработают должным образом и пользователь уйдёт в минус.

Посмотри в сторону инструментов для синхронизации. 

Реализация минимум: все модифицирующие БД операции должны перед выполнением пытаться захватить глобальный sync.Mutex на право выполнения.

Реализация посложнее, но производительнее: модифицирующие БД операции также должны захватывать sync.Mutex, но только по затронутым пользователям. sync.Map если что поможет хранить массив user_id => sync.Mutex.

14. Архитектура

В подмодулях в файлах services.go и models.go 
набор объявленных глобальных функций уместно
сгруппировать по классам и инверсировать зависимости
через конструктор.

Условно в таком духе:

```
// transactions/models.go
type TransactionsRepository interface {
    FindByUserID(userID int) ([]Transaction, error)
}

func NewTransactionsRepository(db *gorm.DB) (TransactionsRepository, error) {
    return &transactionsRepositoryImpl{
        db: db,
    }, nil
}

type transactionsRepositoryImpl struct {
    db *gorm.DB
}

func (r *transactionsRepositoryImpl) FindByUserID(userID int) ([]Transaction, error) {
    // ...
}

// transactions/services.go

type TransactionsService interface {
    FindByUserID(userID int) ([]Transaction, error)
}

func NewTransactionsService(repo TransactionsRepository) TransactionsService {
    return &transactionsServiceImpl{
        repo: repo,
    }
}

type transactionsServiceImpl struct {
    repo TransactionsRepository
}

func (s *transactionsServiceImpl) FindByUserID(userID int) ([]Transaction, error) {
    // ...
}

// main.go (с обработкой ошибок, само собой)
db, err := common.Init(/* ... */)

transactionsRepo, err := transactions.NewTransactionsRepository(db)

transactionsService, err := transactions.NewTransactionsService(transactionsRepo)

// и т.д... Пробросить transactionsService в контроллер схожим образом. 

```
