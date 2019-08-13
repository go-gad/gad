#  Короткое описание генеработа сервиса

Сервис описывается файлом в формате **openapi 3.0**.
Файл размещается в папке **./contract/openapi.yaml**.

Генератор считывает контракт и создает библиотеку в папке **service**.
Библиотека содержит:
* описание типов взятых из _#/components/schemas/*_
* описание интерфейса **IServer**
    * названия методов **\<method_name\>** берутся из **#/paths/\<path\>/\<method\>/operationId** _Example: GetEmployee_
    * входящие параметры:
        * context.Context
        * **\<method_name\>Request** _Example: GetEmployeeRequest_
    * исходящие параметры
        * **\<method_name\>Response** _Example: GetEmployeeResponse_
        * error
* серверную часть
    * хелпер для создания Mux объекта (*github.com/gorilla/mux*) - который на вход принимает тип с поддержкой интерфейса **IServer**. Результат можно дополнять и передавать в _http.ListenAndServe_ 
    * хелперы для создания ендпоинтов (_github.com/go-kit/kit/endpoint_)
* клиентскую часть
    