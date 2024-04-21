# SSO
go run ./cmd/sso/main.go --config=./config/local.yaml

## Задачи 
[х]. Использовать кластер PostgreSQL на базе Patroni для хранения информации о пользователях
[х]. Использовать кластер Redis Sentinel для сохранения отозванных токенов
[ ]. Мониторинг
[ ]. Сделать возможность отозвать все токены выпущенные для конкретного пользователя  
[ ]. Дописать тесты.
[ ]. Сделать автоматический запуск кода для локальной проверки, используя Docker-compose и bash скрипты
[ ]. Сделать автоматический запуск тестов для локальной проверки, используя Docker-compose и bash скрипты
[ ]. Ревью
[ ]. Сделать отправку сообщения в общую шину данных на базе Apache Kafka при успешной регистрации пользвателя.
[ ]. Сделать Helm Charts и запустить в k8s сервис
[ ]. Найти ansible playbooks для развертывания stateful приложений (Patroni,Redis Sentinel) на ВМ. (или забить и просто все в кубере запустить для тестов...)
[ ]. Сделать CI/CD.
[ ]. Ревью


1. Установить cleanenv
https://pkg.go.dev/github.com/ilyakaznacheev/cleanenv - to get env from file
go get github.com/ilyakaznacheev/cleanenv

2. накатываем миграции
sqlite 
go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations

postgres
alex@black:~/GolandProjects/sso$ go run ./cmd/migrator/postgres  --migrations-path=./migrations
                                 go run ../cmd/migrator/postgres  --migrations-path=./migrations

3. После того как добавили чтение из конфиг файла можем запуститься для проверки
   go run cmd/sso/main.go --config=./config/local.yaml

Для тестирования
4. Регестрируем тестовую app через миграции   
INSERT INTO apps (id, name, secret)
VALUES (1, 'test', 'test-secret')
ON CONFLICT DO NOTHING;

go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./tests/migrations --migrations-table=migrations_test

Example:
https://github.com/open-telemetry/opentelemetry-go-contrib/tree/instrumentation/google.golang.org/grpc/otelgrpc/example/v0.46.1/instrumentation/google.golang.org/grpc/otelgrpc/example


Use only TestRegisterLogin_Login_HappyPath for jaeger

redis instrumentation 
https://redis.uptrace.dev/guide/go-redis-monitoring.html#opentelemetry-instrumentation

postgres pgx What driver name do I use to connect Go sqlx to Postgres using the pgx driver?
https://stackoverflow.com/a/74350022

Запуск тестов:
```
alex@black:~/GolandProjects/sso/tests$ ./run_test.sh
no migrations to apply
ok      command-line-arguments  0.190s
```


// patroni
https://github.com/zalando/patroni/blob/master/docker-compose.yml

//docs for docker
https://github.com/zalando/patroni/blob/master/docker/README.md

// Partitioning
https://habr.com/ru/articles/273933/

https://www.commandprompt.com/education/how-to-use-partition-by-in-postgresql/

# Build and run
go build cmd/sso/main.go
./main --config=./config/local.yaml

docker 
```
docker build -t auth-go . --progress=plain --no-cache
docker run -p 44044:44044 auth-go 
```
```bash
docker rm -f $(docker ps -aq)
```
```bash
sudo netstat -nlp | grep :44044
```

redis sentinel client
https://redis.uptrace.dev/guide/go-redis-sentinel.html#redis-server-client


## Окружение развёртывания программного обеспечения - ДЕМО

#### Запуск
1. Переходим в папку с инфраструктурой и запускаем docker-compose 
``` bash 
docker-compose -f docker-compose.demo.yaml up --force-recreate --build sso
 docker rmi infra-sso
 cd infra
docker-compose -f docker-compose.demo.yaml up
```

2. Из корня проекта накатываем миграции 
```bash
go run ./cmd/migrator/postgres  --migrations-path=./migrations 
```

3. Опционально (тестирование) из папки тестов **tests**
Тестирование необходимо проводить на ПУСТОЙ базе.  
```bash 
cd tests
./run_demo_test.sh 
```

#### Мониторинг
grafana: http://localhost:3000/grafana 
jaeger: http://localhost:16686/jaeger/search


for database tracing
https://pkg.go.dev/github.com/XSAM/otelsql#section-readme

for redis tracing
https://redis.uptrace.dev/guide/go-redis-monitoring.html#opentelemetry-instrumentation

Go pq and Postgres appropriate error handling for constraints
https://stackoverflow.com/questions/34963064/go-pq-and-postgres-appropriate-error-handling-for-constraints