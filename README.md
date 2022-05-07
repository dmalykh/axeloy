# UNDER DEVELOPMENT

Axeloy is a is a plugin driven tool that helps send a messages for your customers by a variety of communication channels.

### Main concepts
Subscribe source **profile** to send messages by destinations profiles by **ways**.
Axeloy receives **message** by source **way**, determinate by source **profile** destination profiles. And sent message for destination profiles by ways.

 - **message** contains:
   - source — way received message
   - payload of message, i.e. content
   - destinations — ways to send payload, if doesn't set then Axeloy determinate automatically
 - **profile**  


Axeloy provides several key features:
 - subscription by tag
 - manage drivers in runtime  
 - ways released as golang plugins
 - destinations for messages determinates automatically

### How to make your own driver
Sometimes you may need another way to send messages or interact with Axeloy. 
Driver have to implement [driver.Driver](github.com/dmalykh/axeloy/axeloy/way/driver) interface. Axeloy use drivers as Go plugins.   

#### How driver works?
There is three additional interfaces for drivers: listener, sender and subscriber. 
- Listeners used for receiving messages.
- Senders
- subscriber ?

When Axeloy starts it loads all drivers specified in config. While driver loads method `Init(ctx context.Context, loader ConfigLoader) error` called with `type ConfigLoader func(v interface{}) error` function. If you specified any attributes in config of your plugin, you can get it just calling ConfigLoader func in Init method
```hcl
driver "superdemo" "ways/graphql/graphql.so" {
  port = "998"
}
```
  
And get this config in plugin:
```go

type GraphQl struct {
    config config
}

type config struct {
    port string
}

func (q *GraphQl) Init(ctx context.Context, loader driver.ConfigLoader) error {
	if err := loader(&q.config); err != nil {
		return err
	}
	... do something ...
	return nil
}
```


clean architecture
 

 - **ways** receives events and sends messages
 - ****


How to subscribe?


### Ways
Ways made as go plugins.

helps send and control messages for customers by events.

Axeloy provides several key features:





Commandline interface description




axeloy new [migration|driver] name
axeloy run --config=[path]
axeloy migrate --


channel могут быть входящие, могут быть исходящие.
route маршрут по которому сообщение будет идти дальше

Сервис в channel принимает сообщение, 
у роутера узнаёт куда его дальше, передаёт

#Channel

Принимает/отправляет данные.
Принимающий получает какую-то херню и преобразует в message.
Отправляющий получает message и преобразует в какую-то херню.

IN  | OUT
--- | ---
REST | REST
GRPC |  GRPC
GraphQL |GraphQL
MQ |MQ
Kafka |Kafka
e-mail | e-mail
sms |  sms
call | call
telegram | telegram
 — | push
  —| ws


юзер? Или кальянная? Должен быть как-то потк?

интерфейс доступа к данным — а это тоже по идее тупо out канал:
-   последние N для 

service  -  channel? конвеер  стандартный по идее тут
##Message
Просто сообщение.
Может быть направленое и ненаправленое (с передачей назначения сразу).

messageId —  UUID
source { — Источник, поля должно быть мочь добавлять. Может типа профили? Похер по какому профилю сообщение придёт.
    profile_id — указываем по какому профилю будет обработка данных
    body — ну  как бы то, что по профилю, тип, id, мож ещё что
}
destination {
    channel ?
}

##Preprocessor

Для всякой херни типа составленния текста сообщения и т.п.

##Router
(profile A) subscribe (profile B in chan X, Y, H)

in — source, out — [destination]

Кейс  | Что | Внутри
--- | --- | ---
Я подписался на обновления кальянной.  | При добавлении комментария в кальянную, должно смарштуризировать на меня | В роут будет добавлено place 35790 -> 
Я хочу получать sms
При добавлении комментария  приходит e-mail | А где и как указывается/модифицируется текст письма?

ProfileParser

##Subscription
Для подписки используется вызов в API:
{
  "source": Profile
  "destination": Profile
  "channels": ["sms", "ahues"]
}



##User Cases
#### я хочу получать уведомления о добавлении нового комментария к кальянной
При добавлении комментария к кальянной должно сработать событие on_add_comment. 
В рамках этого события должен быть запрос в модуль уведомлений о том, что place 35790 событие add_comment + (тело комментария? — говно, придётся на каждый месадж писать обработчик).
В результате должно придти internal? Нахера в систему уведомление отправлять, если юзера его может взять сам? Но по вебсокету должно ж придти.)
Тогда два роута — один в вебсокет, второй внутри

place 35790 add_comment 

{
    source: {
        profile_id: 99
        body: {
            type:  place
            id: 35790
            source: add_comment??  а тут ли? И куда текст сам сообщения?
        }
    }
    message: {
        profile_id: 22
        body {
            event: add_comment
            sender: 2
            text: ara_pidara
        }        
    }
}

Регистрация роута 


#### я должен получить код в SMS при вводе сотрудником кальянной моего номера карты hookahid
#### при отправке сообщения на стену кальянной я хочу, чтобы оно продублировалось на страницу hhokahadvisor в facebook
1. Чекаем изменения контента. Например, внесли изменения в цены в кальянной, Новая акция, событие или промокод
2. Новый комментарий или ответ
3. Новый запрос на управление кальянной.
4. Добавлена новая кальянная.
