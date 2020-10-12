# Домашнее задание к занятию «4.1 Kubernetes: основы»

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

**Важно**: для выполнения данного ДЗ вам необходимо установить приложение OpenSSL ([пошаговая инструкция по установке](openssl.md)).

## Stateless Services

### Предварительные требования

В данной задаче предполагается, что вы умеете собирать Docker-образы для своих приложений и публиковать их на GitHub Packages. А кроме того, у вас есть Personal Access Token.

Если вы немного подзабыли, что это всё значит, то обратитесь к [описанию в курсе BGO](https://github.com/netology-code/bgo-homeworks/blob/master/12_docker/packages.md).

### Задача

Необходимо реализовать вычислительный кластер с балансировкой запросов на базе Kubernetes.

Для этого мы с вами создадим Stateless-сервисы, задача которых - выполнять вычисления (а значит, масштабироваться они должны хорошо).

Возьмите за основу код из [вашей любимой задачи с транзакциями](https://github.com/netology-code/bgo-homeworks/tree/master/06_goroutines#%D0%B7%D0%B0%D0%B4%D0%B0%D1%87%D0%B0-2--%D1%81%D1%83%D0%BC%D0%BC%D0%B0-%D0%BF%D0%BE-%D0%BC%D0%B5%D1%81%D1%8F%D1%86%D0%B0%D0%BC).

В чём суть: клиент вам присылает JSON с массивом транзакций, а вы ему отвечаете:
```json
{
  "worker_id": "uuid-here",
  "data": [
    0,
    0,
    10000
  ]
}
```

Где:
1. `worker_id` - ad-hoc идентификатор конкретного сервера, который производил вычисления (мы его сгенерируем программно, поскольку только-только знакомимся с Kubernetes)
2. `data` - это данные по месяцам (индекс соответствует месяцу), если в каком-то месяце нет транзакций, то должен быть 0 (судя по этим данным у нас - март)

Формат запроса вы определяете сами.

Как генерировать `worker_id`:
```go
import (
	"log"
	"net/http"
	"github.com/google/uuid"
)

type Server struct {
	id string
	mux *http.ServeMux
}

func NewServer(mux *http.ServeMux) *Server {
    // т.е. генерируем один раз при создании сервера
	return &Server{id: uuid.New().String(), mux: mux}
}
```

Упакуйте всё в Docker Image и опубликуйте на GitHub Packages (для этого вам достаточно поместить `Dockerfile` в корень вашего проекта и написать следующий Workflow):
```shell script
name: Go

on:
  push:
    branches: [ master ]
  pull_request:
    branches: [ master ]

jobs:

  build:
    name: Build
    runs-on: ubuntu-latest
    steps:
      - name: Push to GitHub Packages
        uses: docker/build-push-action@v1
        with:
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
          registry: docker.pkg.github.com
          repository: netology-code/ago-docker-k8s/replica
          tag_with_ref: true
```

Естественно, `repository` вы должны заменить на свой.

### minikube

#### Установка

<details>
<summary>Инструкция</summary>
    
Установите minikube согласно [руководству по установке](https://minikube.sigs.k8s.io/docs/start/):

Linux:
```shell script
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube_latest_amd64.deb
sudo dpkg -i minikube_latest_amd64.deb
```

Windows: [инсталлятор](https://storage.googleapis.com/minikube/releases/latest/minikube-installer.exe)

MacOS:
```shell script
brew install minikube
```
</details>

#### Запуск

В первом терминале:
```shell script
minikube start # дождитесь запуска
minikube dashboard
```

Во втором терминале:
```shell script
minikube tunnel
```

#### Создание секрета

Kubernetes, так же, как и GitHub, умеет хранить секреты. В качестве секрета будут выступать наши credentials к GitHub Packages.

Команда:
```shell script
kubectl create secret docker-registry github-packages --docker-server=https://docker.pkg.github.com --docker-username=coursar --docker-password=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx --docker-email=coursar@localhost
```

Где:
* `github-packages` - имя секрета
* `--docker-server` - URL реестра (не изменяете)
* `--docker-username` - username на GitHub (**изменяете на свой**)
* `--docker-password` - Access Token на GitHub (**изменяете на свой**)
* `--docker-email` - Email на GitHub (**изменяете на свой**)

Посмотреть, что секрет действительно создался,можно с помощью команды:

```shell script
kubectl get secret github-packages --output=yaml
```

Детальнее про секреты вы можете узнать из [документации]().

**Важно**: мы специально не делаем секреты в виде yaml-файлы, потому что студенты достаточно часто коммитят их прямо в GitHub (что не очень хорошо, поскольку там ваш токен в открытом виде*).

Примечание*: он закодирован base64, но это ни в коем случае не защита.

### Развёртывание

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: replica
spec:
  selector:
    matchLabels:
      app: replica
  replicas: 3
  template:
    metadata:
      labels:
        app: replica
    spec:
      containers:
        - name: replica
          image: docker.pkg.github.com/netology-code/ago-docker-k8s/replica:latest
          ports:
            - containerPort: 9999
      imagePullSecrets:
        - name: github-packages
---
apiVersion: v1
kind: Service
metadata:
  name: replica
  labels:
    app: replica
spec:
  type: LoadBalancer
  ports:
    - port: 9999
  selector:
    app: replica
---
```

Шаблон взят с лекции. Мы разворачиваем 3 реплики, взяв образ с нашего реестра образов (вы можете заменить `replica` на своё имя).

Ключевое: именно вот эта строка отвечает за то, какой секрет мы используем для получения образов:
```yaml
      imagePullSecrets:
        - name: github-packages
```

Для применения:
```shell script
kubectl apply -f <имя_файла или каталога>
```

### Client

Клиент у нас будет исследовательский - нам нужно нагенерировать нагрузку и посмотреть статистику ответов (клиент храните в том же репозитории на GitHub - упаковывать его в Docker образ не нужно).

Что должен делать клиент? Запускать несколько сотен горутин с одним и тем же запросом на сервер, собирать результаты и считать, сколько раз ответил каждый из сервисов (выводите в консоль).

### Результат

В качестве результата пришлите ссылку на ваш GitHub-проект, в котором реализованы описанные выше требования.

### Reference проект

[ago-docker-k8s](https://github.com/netology-code/ago-docker-k8s)
