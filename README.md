# Yandex Lyceum Golang Final Project
**Если есть вопросы [пишите](https://t.me/voronkovmaksim) https://t.me/voronkovmaksim**


## Как это работает 
[Наглядная блок схема](https://excalidraw.com/#json=B6NkXR9ZB0sjSl1JqCuyE,gsL3U_U_fLt_gMiQAuWLJw)

У меня все реализовано в консоле, без графического интерфейса. Есть оркестратор и неограниченное количество агентов, для хранения информации используются csv. Перезапускать проект нельзя. Мониторинг агентов есть, можно узнать на каких портах работают агенты, крайнее время, когда с ними была связь, сколько операций они выполняют сейчас и их максимальное количество операций. Связь оркестратора и агентов реальзованна на GRPC. Так же есть регистрация пользователей.
Оркестратор запускается на порту 8080, он распределяет все задачи между агентами. Сначла нужна зарегестрироваться и получить jwt токен. Далее все операции должны происходить с его использованием. С помощью POST запрса передаем оркестратору json файл, с выражением, которое нужно посчитать и с длительностью выполнения операций. С помощью get запроса можно узнать статус выражения, которое было отправленно оркестратору, или узнать состояние агентов. После того как оркестратор получил выражение запускается горутина где он сначало проверяет правильность ввода выражения, затем переводит в постфиксную форму для простоты подсчета и начинает считать. Для подсчета чего либо оркестратор формирует json с операндами, операцией и длительностью выполнения действий, и отправляет его в так называемы long list, в нем выражения ждут пока их распределят на свободного агента, т.к. Свободного агента может не быть для этого и сделан long list. Выражение приходит агенту он его проверяет и запускает горутину с вычислением, после вычисления он записывает ответ в answer.csv, откуда его позже получает оркестратор и продолжает вычисления, после того, как он выполнил вычисления и получил финальный ответ значение записывается в файл  expression.csv, в нем также хранится основная информация про это выражение. В то время как выражения вычисляется пользователь может наблюдать в реальном времени за состоянием агентов, их занятостью. 

## Мои комментарии по критериям 

1. Регистрация реализованна и все работает в контексте конкретного пользователя. - 20 баллов
2. Нет, все информация хранится в csv, перезагружать систему нельзя. 
3. Общение оркестратора и агентов реализованно на GRPC. - 10 баллов
4. Нет
5. Нет 


## Как это запустить
**Видео иструкция**   
1. Копируем репозиторий 
```commandline 
git clone https://github.com/Zyvexa/Yandex-Lyceum-Golang-Final-Project
```  
2. Открываем в этой папке 4 cmd (2 для пользователей и 2 для оркестратора и агента), для оркестратора нужно перейти в папку _Main_, для агента в  _Agent_.
3. Устананавливаем все для JWT и GRPC (как в уроках)(https://lms.yandex.ru/courses/1051/groups/9299/lessons/6683/tasks/51780 и https://lms.yandex.ru/courses/1051/groups/9299/lessons/6690/tasks/51796)
4. В cmd для агента и оркестратора
```commandline 
set GO111MODULE=on 
``` 
**Теперь все запущенно и готово к работе**


## Тестовые сценарии
**Все должно быть запущено как в предыдущем пункте**   
### Первый тест
Файлы:   
login.json : 
```json 
{
    "login": "1",
    "password": "1"
}   
```
expression.json : 
```json 
{
    "expression" : "2 + 2"
} 
```
time.json : 
```json 
{
    "timeAddition": 10,
    "timeSubtraction": 10,
    "timeMultiplication": 10,
    "timeDivision": 10,
    "timeExponentiation": 10
}   
``` 
Запускаем команды  
**В папке Main**
```commandline 
go run main.go
```
**В папке Agent**
```commandline 
go run agent.go
```
**В окне пользователя**  
Зарегестрируемся и войдем
```commandline 
curl -X POST -d @login.json http://localhost:8080/register
```
```commandline 
curl -X POST -d @login.json http://localhost:8080/login
```
Тут вернется jwt токен которы нужно будет скопировать и истпользовать   
```commandline 
curl -X POST -H "Authorization: Bearer Ваш_токен" -d @time.json http://localhost:8080/time
```
Отправим выражение
```commandline 
curl -X POST -H "Authorization: Bearer Ваш_токен" -d @expression.json http://localhost:8080/expression 
```
Ждем 10 секунд 
```commandline 
curl -H "Authorization: Bearer Ваш_токен" http://localhost:8080/expressions
```
Вернем информауию про выражение

### Второй тест
**Файлы:**   
login.json : 
```json 
{
    "login": "1",
    "password": "1"
}   
```
expression.json : 
```json 
{
    "expression" : "2 + 2"
} 
```

Запускаем команды  
**В папке Main**
```commandline 
go run main.go
```
**В папке Agent**
```commandline 
go run agent.go
```
**В первом окне пользователя**  
Зарегестрируемся и войдем
```commandline 
curl -X POST -d @login.json http://localhost:8080/register
```
```commandline 
curl -X POST -d @login.json http://localhost:8080/login
```
Тут вернется jwt токен которы нужно будет скопировать и истпользовать   
Отправим выражение
```commandline 
curl -X POST -H "Authorization: Bearer Ваш_токен" -d @expression.json http://localhost:8080/expression 
```
Ждем 10 секунд 
```commandline 
curl -H "Authorization: Bearer Ваш_токен" http://localhost:8080/expressions
```
**Файлы:**  
login.json : 
```json 
{
    "login": "2",
    "password": "1"
}   
```
expression.json : 
```json 
{
    "expression" : "2 + 3"
} 
```
**Во втором окне пользователя**  
Зарегестрируемся и войдем
```commandline 
curl -X POST -d @login.json http://localhost:8080/register
```
```commandline 
curl -X POST -d @login.json http://localhost:8080/login
```
Тут вернется jwt токен которы нужно будет скопировать и истпользовать   
Отправим выражение
```commandline 
curl -X POST -H "Authorization: Bearer Ваш_токен" -d @expression.json http://localhost:8080/expression 
```
Ждем 10 секунд 
```commandline 
curl -H "Authorization: Bearer Ваш_токен" http://localhost:8080/expressions
```
Вернет информауию про выражения только этого пользователя 
### Третий тест
**Файлы:**   
login.json : 
```json 
{
    "login": "1",
    "password": "1"
}   
```

Запускаем команды  
**В папке Main**
```commandline 
go run main.go
```
**В папке Agent**
```commandline 
go run agent.go
```
**В первом окне пользователя**  
```commandline 
curl -X POST -d @login.json http://localhost:8080/register
```
**Во втором окне пользователя**  
```commandline 
curl -X POST -d @login.json http://localhost:8080/register
```
Выдаст ошибку т.к. такой пользователь уже зарегестрирован  
**Файлы:**   
login.json : 
```json 
{
    "login": "1",
    "password": "2"
}   
```
**В первом окне пользователя**  
```commandline 
curl -X POST -d @login.json http://localhost:8080/login
```
Выдат ошибку т.к. пароль неправильный 
## Как этим  пользоваться 
Есть неесколько основных команд, которые описаны ниже. Перед началом каждому пользователю нужно зарегестрироваться и войти. Для отправки данных у меня везде используется json. Почти с каждым запросом должен идти JWT токен. Пример должен быть записан через пробел, т.е. ~~"2+2"~~, **"2 + 2"**, **"( 12 + 2 ) ^ 2"**, ~~"(12 + 2) ^ 2"~~.
### Основные команды
#### Регистрация 
Сurl запрос:
```commandline 
curl -X POST -d @login.json http://localhost:8080/register
```
login.json : 
```json 
{
    "login": "Ваш логин",
    "password": "Ваш пароль"
}   
```
#### Вход 
Сurl запрос:
```commandline 
curl -X POST -d @login.json http://localhost:8080/login
```
login.json : 
```json 
{
    "login": "Ваш логин",
    "password": "Ваш пароль"
}   
```
**Возвращает JWT токен**
#### Отправить выражение 
Сurl запрос:
```commandline 
curl -X POST -H "Authorization: Bearer Ваш_токен" -d @expression.json http://localhost:8080/expression 
```
пример с токеном(ваш токен будет другим):
```commandline 
curl -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM3NjMyNjMsImlhdCI6MTcxMzcwMzI2MywibG9naW4iOiIxIiwibmJmIjoxNzEzNzAzMjYzLCJwYXNzd29yZCI6IjEifQ.desjokdTHKf1NXzrXyzF7iXzRPHyjpXVQFPcIY_QZnk" -d @expression.json http://localhost:8080/expression 
```
expression.json : 
```json 
{
    "expression" : "2 + 2"
} 
```
Выражение должно быть через пробелы, например ```"2 + 2 * 2"```
#### Отправить время 
Сurl запрос:
```commandline 
curl -X POST -H "Authorization: Bearer Ваш_токен" -d @time.json http://localhost:8080/time
```
time.json : 
```json 
{
    "timeAddition": 60,
    "timeSubtraction": 60,
    "timeMultiplication": 60,
    "timeDivision": 60,
    "timeExponentiation": 60
}   
``` 
"timeAddition" - сумма  
"timeSubtraction" - разность   
"timeMultiplication" - умножение   
"timeDivision" - деление   
"timeExponentiation" - возведение в степень  
#### Получить информацию про выражения 
Сurl запрос:
```commandline 
curl -H "Authorization: Bearer Ваш_токен" http://localhost:8080/expressions
```
**Возвращает список выражений в формате: expression - выражение, id, time in - время начала расчетов, time out - время конца, answer - ответ, error - ошибка**   
#### Получить информацию про агентов 
Сurl запрос:
```commandline 
curl http://localhost:8080/agents
```
**Возвращает список выражений в формате: port - порт, last_time - последнее время активности, free - свободных мест, total - всего мест**