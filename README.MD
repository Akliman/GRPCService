# Тестовое задание

Текст задания https://github.com/nakukop/go


### Переменные
Перед запуском в файлах
cmd/localStorage/main.go 

cmd/memcached/main.go

Можно изменить локальные переменные

```
const (
memcacheAdr = "localhost:11211" //Адрес memcached сервера
GRPCPort    = "8080"            //Порт на котором поднимается GRPC сервер
)
```




### Запуск сервиса
* Сервис с хранилищем внутри памяти 
```
go run cmd/localStorage/main.go 
```
* Сервис с memcached хранилищем
```
go run cmd/memcached/main.go 
```

### Контракт
proto контракт располагается в api/grpc/protos/kvstorage.proto

### Тесты
Тесты лежат в папке с бизнесс логикой
/core