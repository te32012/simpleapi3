[![builded](https://github.com/te32012/simpleapi3/actions/workflows/docker-image.yml/badge.svg)](https://github.com/te32012/simpleapi3/actions/workflows/docker-image.yml)



Тестовое задание pgpro 2024 

Формат запросов к api описан в файле api.yml при помощи openapi/swagger (для удобного просмотра можно - скопировать содержимое файла сюда https://editor.swagger.io/)

Общая схема взаимодействия следующая:
1) Отправить валидный bash скрипт при помощи /create endpoint 
2) Проверить, что скрипт корректно загружен в базу при помощи /get/{id} endpoint (опционально изучить все доступные для выполнения команды при помощи /getAll)
3) Запустить нужную команду на сервере при помощи /start endpoint (с параметрами и входными данными для скрипта) и получить id логов данного скрипта и id pid
4) Проверить статус запущенного скрипта при помощи /status enpoint (если скрипт завершился в ответе будет дата окончания запуска скрипта и код возврата) и получить логи данного скрипта.
5) Если скрипт завис - отправить команду /stop для его принудительного завершения. 


Для деплоя:
1) git clone репозиторий
2) cd папка с репозиторием
3) docker-compose build
4) docker-compose up

Работоспособность контейнера подтверждается сборкой при помощи github actions. 


Тестирование было осуществлено в ручном режиме без вмешательства в работоспособность системы (черный ящик методика) через postman, после написания программы (можно подразумевать - что оно приемочное)
Скришоты postman окон с запросами/ответами в папке examples

План тестирования:

Сначала создали два одинаковых скрипта: create1, create2

Запросили список всех доступных скриптов: getall

Проверили содержание одного из созданных скриптов (get1) и попытались получить информацию о не существующем скрипте get3 (получили 404 согластно api спецификации)

Запустили скрипты (start1, start2) и попытались запустить не существующий скрипт (start6). Получили pid и logid для запущенных скриптов (для не существующего получили 404 согластно api спецификации)

Получили статус запущенных скриптов до их завершения (т к скрипты просто флудилка) status1 и status2 - убедились, что скрипты не завершены т к нет даты окончания и кода возврата 
Стоит отметить что в одном скрипте есть выведеные данные а в другом еще пока нет т к в базу для того, чтобы ее не нагружать данные пишутся не в реальном времени а раз в минуту все что накопилось в буфере

Попытались остановить скрипты с pid 417,507, 506 - два остановились успешно т е 200 ок а одного не было найдено в списке запущенных т к мы его не запускали в связи с этим получили 404

Запросили стаус остановленных скриптов 417 и 507 и успешно получили логи и убедились в том, что они завершены т к есть код возврата (-1 потому что процесс был остановлен kill) и есть дата окончания запуска. 

Результаты тестирования:

Основной функционал приложения работоспособен для использования пользователем 
В ходе тестирования выявлены следующие замечания (по запросу может быть реализовано):
1) не хватает возможности посмотреть общую информацию о всех запущенных командах на данный момент
2) не хватает возможности посмотреть общую информацию обо всех командах которые ранее запускались и информация о них сохранена в базе
