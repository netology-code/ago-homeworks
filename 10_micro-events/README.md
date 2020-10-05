# Домашнее задание к занятию «3.3 Микросервисы: Event-Driven Communication»

Все задачи этого занятия можно делать в **одном репозитории**.

В качестве результата пришлите ссылки на ваши GitHub-проекты через личный кабинет студента на сайте [netology.ru](https://netology.ru).

**Важно**: если у вас что-то не получилось, то оформляйте Issue [по установленным правилам](../report-requirements.md).

## Как сдавать задачи

1. Создайте на вашем компьютере проект, в котором разместите Go-модули (на основании исходников с лекции).
1. Инициализируйте в нём пустой Git-репозиторий.
1. Добавьте в него готовый файл [.gitignore](../.gitignore).
1. Сделайте необходимые коммиты.
1. Создайте публичный репозиторий на GitHub и свяжите свой локальный репозиторий с удалённым.
1. Сделайте пуш (удостоверьтесь, что ваш код появился на GitHub).
1. Ссылку на ваш проект отправьте в личном кабинете на сайте [netology.ru](https://netology.ru).
1. Задачи, отмеченные как необязательные, можно не сдавать, это не повлияет на получение зачета (в этом ДЗ все задачи являются обязательными).

## Merge

### Задача

На лекции мы интегрировали Kafka и Event Driven communications в наши микросервисы, которые используют REST-интерфейсы. А у вас с прошлого ДЗ есть "свои" продвинутые микросервисы с gRPC коммуникациями.

Задача достаточно простая: вам нужно из нашего проекта с REST "перетащить" реализацию в свой проект.

**Важно**: внешнее API Backend сервиса, через которое к нему обращаются клиенты, менять не нужно.

Ваш сервис по-прежнему должен запускаться с помощью Docker Compose и работать.

В качестве результата пришлите ссылку на ваш GitHub-проект, в котором реализованы описанные выше требования.

## Push notifications

### Задача

Задача достаточно жизненная: мы хотим сделать новый сервис, который подписывается на события успешного и не успешного завершения платежей и посылает Push-уведомление на мобильное устройство пользователя.

Как это происходит в реальной жизни: вы в обычной БД храните связку между ID пользователя и Push-токеном. Push-токен присылают вам сами мобильные приложения и выглядит этот токен просто как "длинная строка": `bk3RNwTe3H0:CI2k_HHwgIpoDKCIZvvDMExUdFQ3P1...`.

Но поскольку у нас с вами нет этих самых мобильных приложений, которые могут это всё делать*, мы с вами проэмулируем отправку обычным выводом в лог вида:
```text
send notification to <token> with message: <message>
```

<details>
<summary>Примечание*: об используемых в промышленном коде инструментах</summary>

Для отправки обычно используется [Firebase Cloud Messaging (FCM)](https://firebase.google.com/docs/cloud-messaging) и есть отдельно [SDK для Go](https://firebase.google.com/docs/admin/setup/#go).

Если вы хотите поработать с живыми приложениями и отправлять настоящие Push'и на моб.устройства Android, тегайте в слаке @coursar, мы предоставим вам соответствующую информацию и инструкции. Но будьте готовы к тому, что задача реализации всего проекта займёт больше часа (вполне возможно, что несколько дней) и вам потребуется аккаунт Google Play разработчика (Google требует при регистрации единоразовую оплату в 25$).

</details>

Соответственно, вам нужно решить две задачи:
1. Как вы через `backend` сервис предоставите "наружу" API для сохранения Push-токена (оно не обязательно должно быть синхронным)
1. Как вы будете получать события об успешных/неуспешных платежах

Поскольку вы - главный архитектор, вам и решать, как организовать коммуникации, сколько сервисов для этого нужно и какие зависимости и API им требуется.

Единственное требование: как минимум одна из двух описанных выше задач должна быть решена с использованием Kafka.