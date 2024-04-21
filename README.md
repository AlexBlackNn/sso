# SSO
go run ./cmd/sso/main.go --config=./config/local.yaml

## Задачи 
1. [x] Использовать кластер PostgreSQL на базе Patroni для хранения информации о пользователях.
2. [x] Использовать кластер Redis Sentinel для сохранения отозванных токенов.
3. [ ] Мониторинг.
4. [ ] Сделать возможность отозвать все токены выпущенные для конкретного пользователя.
5. [ ] Дописать тесты.
6. [ ] Сделать автоматический запуск кода для локальной проверки, используя Docker-compose и bash скрипты
7. [ ] Сделать автоматический запуск тестов для локальной проверки, используя Docker-compose и bash скрипты
8. [ ] Ревью
9. [ ] Сделать отправку сообщения в общую шину данных на базе Apache Kafka при успешной регистрации пользвателя.
10. [ ] Сделать Helm Charts и запустить в k8s сервис
11. [ ] Найти ansible playbooks для развертывания stateful приложений (Patroni,Redis Sentinel) на ВМ. (или забить и просто все в кубере запустить для тестов...)
12. [ ] Сделать CI/CD.
13. [ ] Ревью

## Окружение развёртывания программного обеспечения - локально

#### Запуск
1. Переходим в папку с инфраструктурой и запускаем docker-compose
``` bash 
cd infra
docker-compose up
```

2. Из корня проекта накатываем миграции
```bash
go run ./cmd/migrator/postgres  --migrations-path=./migrations 
```
Примечание: в случае ошибки, подождать когда все контейнеры запустяться
```
panic: EOF
goroutine 1 [running]:
main.main()
.../sso/cmd/migrator/postgres/main_postgres.go:43 +0x29c
exit status 2
```
3. Запускаем приложение локально
```bash
go run cmd/sso/main.go --config=./config/local.yaml
```

## Окружение развёртывания программного обеспечения - ДЕМО

#### Запуск
1. Переходим в папку с инфраструктурой и запускаем docker-compose 
``` bash 
cd infra
docker-compose -f docker-compose.demo.yaml up --force-recreate --build
```

2. Из корня проекта накатываем миграции
```bash
go run ./cmd/migrator/postgres  --migrations-path=./migrations 
```
Примечание: в случае ошибки, подождать когда все контейнеры запустяться
```
panic: EOF
goroutine 1 [running]:
main.main()
.../sso/cmd/migrator/postgres/main_postgres.go:43 +0x29c
exit status 2
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

#### Мысли об открытых портах в docker-compose
Не вижу смысла закрывать порты на базе, редисе и т.д. в docker-compose.
В проде врятли кто будет использовать docker-compose на 1 машине. Скорее всего 
это будет k8s или еще какой-то оркестратор, а stateful приложения, вероятно, будут
вынесены из кубера (холивар). 

При тестировании через Postman необходимо добавить x-trace-id в metadata и 
менять его перед каждым запросом.
Например:11116f9a6be295d4ef5a6e030ef11110
![postman.png](docs%2Fpostman.png)


Пример трейсинга в jaeger
![jaeger.png](docs%2Fjaeger.png)


Сбор логов в локи и отображение в графане
![loki-grafana.png](docs%2Floki-grafana.png)

