Сервис регистрирует персон в шине rabbitmq при получении POST запроса с данными о персоне

Параметры:
 -addr string
        web api address (default ":7878")
 -amqp string
        amqp url as amqp://user:psw@host
 -queue string
        bus queue in default exchange (default "person-reg")

 АПИ:

 POST /person
 firstname=anya&lastname=kozyreva

 Пример:

 curl -X POST -d "firstname=anya&lastname=kozyreva" http://127.0.0.1:7878/person


