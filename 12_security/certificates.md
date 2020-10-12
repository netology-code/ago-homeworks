# Генерация сертификатов

В этом документе мы научимся генерировать сертификаты, которые затем будем использовать для TLS.

Для этого нужно выбрать криптографический алгоритм, размер ключа (в байтах).

#### Создайте закрытый ключ и сертификат, со сроком действия в 365 дней:

```shell script
sudo openssl req -newkey rsa:2048 -nodes -x509 -days 365 -keyout key.pem -out certificate.pem -subj "/C=RU/ST=Moscow/L=Moscow/O=Dev/OU=Dev/CN=netology.local" -addext "subjectAltName=DNS:netology.local"
```

Где [`req`](https://www.openssl.org/docs/manmaster/man1/openssl-req.html) - подкоманда, отвечающая за генерацию ключей.

Примечание*: если вы на Mac и команда не срабатывает, запустите её в Docker контейнере с Ubuntu (предварительно поставьте туда OpenSSL). Если не получилось, тегайте в слаке coursar.

Откройте файл `certificate.pem` в любом текстовом редакторе (можете просмотреть в терминале с помощью команды `cat certificate.pem`), удостоверьтесь, что он выглядит примерно так:

```text
-----BEGIN CERTIFICATE-----
...
МНОГО БУКВ
...
-----END CERTIFICATE-----
```

Откройте файл `key.pem` (можете просмотреть в терминале с помощью команды `cat key.pem`) в любом текстовом редакторе, удостоверьтесь, что он выглядит примерно так:

```text
-----BEGIN PRIVATE KEY-----
...
МНОГО БУКВ
...
-----END  PRIVATE KEY-----
```
