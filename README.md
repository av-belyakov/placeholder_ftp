# README

Приложение 'placeholder_ftp' выполняет скачивание с FTP сервера-источника файлы сетевого трафика,  преобразовывает их в текстовый формат по аналогии с TShark и загружает файлы в текстовом виде на FTP сервера-назначения, в нашем случае это FTP сервер ГЦМ.
Выполняет действия в следующем порядке:

1. Принимает команду в виде JSON объекта, структура JSON объекта подробно описана ниже, через брокер сообщения NATS.
2. Забирает по указанному пути с локального ftp сервера, pcap файл, с фрагментом сетевого трафика.
3. Преобразует полученный pcap файл в текстовый формат, результат преобразования похож на файл полученный в результате
   обработки файла pcap программой TShark.
4. Отправляет полученный текстовый файл на ftp сервер-агрегатор по заданному в полученном JSON пути.
5. Отправляет JSON объект с ответом инициатору запроса.

## Конфигурационные настройки

Конфигурационные параметры для сервиса могут быть заданы как через конфигурационный файл так и методом установки переменных окружения. Однако, все пароли и
ключевые токены, используемые для авторизации, задаются ТОЛЬКО через переменные окружения.

#### Типы конфигурационных файлов:

- config.yaml общий конфигурационный файл
- config_dev.yaml конфигурационный файл используемый для тестов при разработке
- config_prod.yaml конфигурационный файл применяемый в продуктовом режиме

В конфигурационных файлах config_dev.yaml и config_prod.yaml есть некоторые параметры требующие дополнительного описания:

- COMMONINFO.name_regional_object - 'имя регионального объекта', модуль будет принимать и обрабатывать только запросы, в которых содержимое поля 'source' будет совпадать с этим значение, остальные запросы игнорируются.

- COMMONINFO.main_ftp_path_result_directory - путь к директории на ftp сервере MainFTP, в которую будут сохраняться переданные файлы.

- COMMONINFO.max_writing_file_limit - ограничение максимального размера файла, в мегабайтах, который будет передаваться на ftp MainFTP. При этом такой 'усечённый' файл будет иметь дополнительный суфикс '.limit', информирующий что файл был обрезан. Если данный параметр имеет значение 0, то уменьшение размера файла не выполняется.

- DATABASEWRITELOG.name_db - наименование базы данных, в случае реляционных БД и некоторых видов не реляционных БД.

- DATABASEWRITELOG.storage_name_db - наименование таблицы БД (MySQL, PostgresSQL), коллекции документов (MongoDB), индекса (Elasticsearch).

Основная переменная окружения для данного приложения - GO_PHFTP_MAIN. На основании значения этой переменной принимается решение какой из конфигурационных файлов config_dev.yaml или config_prod.yaml использовать. При GO_PHFTP_MAIN=development будет использоваться config_dev.yaml, во всех остальных случаях, в том числе и при отсутствии переменной окружения GO_PHFTP_MAIN будет использоваться конфигурационный файл config_prod.yaml. Перечень переменных окружения которые можно использовать для настройки приложения:

#### Переменная окружения отвечающая за тип запуска приложения "test", "development" или "production"

- GO_PHFTP_MAIN

#### Переменная окружения отвечающая за наименование регионального объекта

- GO_PHFTP_NAMEREGOBJ

#### Переменная окружения отвечающая за путь к папке на ftp сервере MainFTP, где хранятся загружаемые файлы

- GO_PHFTP_MAINFTPPATHRESDIR

#### Переменные окружения отвечающие за подключение к NATS

- GO_PHFTP_NPREFIX
- GO_PHFTP_NHOST
- GO_PHFTP_NPORT
- GO_PHFTP_NCACHETTL - данный параметр должен содержать время жизни записи
  кэша, по истечение которого запись автоматически удаляется, значение задается
  в секундах в диапазоне от 10 до 86400 секунд
- GO_PHFTP_NSUBLISTENERCOMMAND - канал для приема команд

#### Переменные окружения отвечающие за подключение к локальному FTP серверу

- GO_PHFTP_LOCALFTP_HOST
- GO_PHFTP_LOCALFTP_PORT
- GO_PHFTP_LOCALFTP_USERNAME
- GO_PHFTP_LOCALFTP_PASSWD

#### Переменные окружения отвечающие за подключение к FTP серверу агрегатору файлов

- GO_PHFTP_MAINFTP_HOST
- GO_PHFTP_MAINFTP_PORT
- GO_PHFTP_MAINFTP_USERNAME
- GO_PHFTP_MAINFTP_PASSWD

#### Переменные окружения отвечающие за настройки доступа к БД которая хранит логи приложения

- GO_PHFTP_DBWLOGHOST // доменное имя или ip БД
- GO_PHFTP_DBWLOGPORT // порт БД
- GO_PHFTP_DBWLOGNAME // наименование БД (при необходимости)
- GO_PHFTP_DBWLOGSTORAGENAME // наименование объекта хранения логов (таблица, документ, индекс и т.д. зависит от типа БД)
- GO_PHFTP_DBWLOGUSER // пользователь БД
- GO_PHFTP_DBWLOGPASSWD // пароль для доступа к БД

Настройки логирования данных в БД не являются обязательными и необходимы только если пользователь приложения желает хранить логи в базе данных

Приоритет значений заданных через переменные окружения выше чем значений полученных из конфигурационных файлов.

## Структура JSON запроса

```
{
  "task_id": "", //идентификатор задачи
  "source": "", //наименование регионального объекта к которому был адресован запрос
  "service": "test_service", //имя сервиса-инициатора команды
  "command": "convert_and_copy_file", //наименование команды
  "parameters": {
  "links": [
    "ftp://ftp.rcm.cloud.gcm/traff/test_pcap_file.pcap",
	  "ftp://ftp.rcm.cloud.gcm/traff/test_pcap_file_http.pcap",
	  "..."
	  ] //список ссылок на файлы, которые необходимо обработать
  }
}
```

Выполняются следующие типы команд:

- 'copy_file' просто выполняет копирование файла с одного ftp сервера на другой
- 'convert_and_copy_file' выполняется конвертирование pcap файла в текст и копирование с одного ftp сервера на другой

## Структура JSON ответа

Для получения ответа на запрос необходимо слушать 'временную тему' сообщения NATS.

```
{
  "request_id":"", //идентификатор задачи
  "source": "", //наименование регионального объекта к которому был адресован запрос
  "error": "", //содержит глобальные ошибки, такие как например, ошибка подключения к ftp серверу
  "processed_information": [
	{
	  "error": "" //ошибка возникшая при обработки файла
	  "link_old": "ftp://ftp.rcm.cloud.gcm/traff/test_pcap_file.pcap",
	  "link_new": "ftp://ftp.cloud.gcm/traff/test_pcap_file.pcap.txt"
    "size_befor_processing": int //размер файла до обработки
	  "size_after_processing": int //размер файла после обработки
	}
  ]
}
```
