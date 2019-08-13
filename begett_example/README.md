## Короткое описание генеработа сервиса

Сервис описывается файлом в формате **openapi 3.0** (_https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md_).
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
* серверную часть (HTTP)
    * хелпер для создания Mux объекта (*github.com/gorilla/mux*) - который на вход принимает тип с поддержкой интерфейса **IServer**. Результат можно дополнять и передавать в _http.ListenAndServe_ 
    * хелперы для создания ендпоинтов (_github.com/go-kit/kit/endpoint_)
* клиентскую часть (HTTP)
    * об]ект клиента с интерфейсом
    * хелперы _gokit_ ендпоинтов для клиента
    
TODO: AMQP Transport, GRPC Transport

## Описание типов и интерфейса _service.IService_

Типы генерируются на основе OpenAPI описания по пути _#/components/schemas/*_

При использовании OneOff генерируется несколько типов и в местах где используется этот тип подставляется тип _interface_, который может кастится к одному из этих типов так же этот интерфейс будет иметь хелперы для упращения логики  Is(Type), ...

Дополнительно генерируются типы для ошибок, описанных в _responses_

Интерфейс _IService_ включает в себя методы взятые из _#/paths/\<path\>/\<method\>/operationId_.
Входящие параметры метода берутся из _#/paths/\<path\>/\<method\>/parameters_.
Исходящих параметров два.

Первый - соответствует схеме описанной в _#/paths/\<path\>/\<method\>/responses/200_.

Второй - ошибка.
На основе других ключей из _#/paths/\<path\>/\<method\>/responses/200_ генерируются типы для ошибок. Например:

    type GetEmployeeResp400 struct {
    	ErrCode int    `json:"err_code"`
    	Msg     string `json:"msg"`
    }
    
    func (e *GetEmployeeResp400) Error() string {
    	return fmt.Sprintf("Bad request: error %d: message: %s", e.ErrCode, e.Msg); // ? как выводить кастомные ошибки
    }
    
Проверить на конкретную ошибку можно кастуя к типу:

    badRequest, ok := err.(*service.GetEmployeeResp400)
    if ok {
        panic(badRequest.Msg)
    }

Либо используя type switch:
    
    switch eType := err.(type) {
    case *service.GetEmployeeResp400:
        panic(badRequest.Msg)
    default:
        panic("Unknown error")
    }
    


## Серверная часть

* Возможность создать роутер который можно передать в _http.ListenAndServe_
* Возможность создать враппер над всеми ендпоинтами
* Возможность создать враппер над отдельным ендпоинтом
* Возможность определить свой encoder и decoder для каждого ендпоинта
* Возможность возвращать разные типы в отдном ответе (поддержка OneOf для OpenAPI)
* Возможность определить тип ошибки от серверной части
* Возможность инструментации запросов

### Возможность создать роутер который можно передать в _http.ListenAndServe_

    // инициализируем тип реализующий интерфейс service.IService
    businessLogicSvc := &business_logic.Service{}
    
    // Создаем стандартный роутер, на основе контракта
    router := service.NewRouter(businessLogicSvc)
 
    // Запускаем Сервис
 	err := http.ListenAndServe(":80", router)
 	if err != nil {
        panic(err)
 	}
 	os.Exit(0)
 	
_&business_logic.Service{}_ - должен реализовать _service.IService_ - интерфейс автоматически генерируется и нужно будет ручками создать тип и реализовать этот интерфейс

    router := service.NewRouter(businessLogicSvc)

На основе имеющегося интерфейса создаем _gorilla/mux_ объект, который реализует http.Handler интерфейс
далее передаем его в http.ListenAndServe либо в другой фреймворк поддерживающий стандартный http.Handler
 	
### Возможность кастомизации

	businessLogicSvc := &business_logic.Service{}

	getEmployeeCustomEndpoint := service.MakeGetEmployeeEndpoint(businessLogicSvc)

	getEmployeeCustomHandler := service.MakeGetEmployeeCustomHandler(
        getEmployeeCustomEndpoint,
        service.DecodeGetEmployeeRequest,
        service.EncodeResponse,
	)

	r := mux.NewRouter()
	service.AddGetEmployeeCustomRouter(r, getEmployeeCustomHandler)
	r.Use(loggingMiddleware)

	err := http.ListenAndServe(":80", r)
	if err != nil {
        panic(err)
	}
	os.Exit(0)
	
Тот же тип реализующий _service.IService_ интерфейс.

    getEmployeeCustomEndpoint := service.MakeGetEmployeeEndpoint(businessLogicSvc)

На его основании создается ендпоинт (_endpoint.Endpoint_)

	getEmployeeCustomHandler := service.MakeGetEmployeeCustomHandler(
        getEmployeeCustomEndpoint,
        service.DecodeGetEmployeeRequest,
        service.EncodeResponse,
	)

Можно создать кастомный хандлер. _service.MakeGetEmployeeCustomHandler_ - является простым маппером в _httptransport.NewServer_, который можно кастомизировать своими декодером и энкодером, плюс передать туда options (_options ...httptransport.ServerOption_) его так же можно враппить.

Либо создать простой хендлер на основе ендпоинта используя следующий хелпер:

    func MakeGetEmployeeHandler(ep endpoint.Endpoint, options ...httptransport.ServerOption) *httptransport.Server {
    
Далее создаем свой _mux_ объект и добавляем наш хендлер:

    r := mux.NewRouter()
	service.AddGetEmployeeCustomRouter(r, getEmployeeCustomHandler)
	
Тут так же можно добавить в сой роутер свои хандлеры либо добавить мидлварь:

    r.Use(loggingMiddleware)
    
Код функции будет иметь примерно такой вид:

    func loggingMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            log.Println(r.RequestURI)
            next.ServeHTTP(w, r)
        })
    }


## Клиентская часть

* Возможность создать простой клиент и вызвать метод, получив в ответ готовый тип и ошибку
* Возможность добавить wrappers (например логирование, рейтлимит, circuitbraker)
* Возможность асинхронного вызова, с получением результата через channel (?????? нужен или оставить пользователю самому реализовать)


    client, err := service.NewHTTPClient("http://127.0.0.1")
    if err != nil {
        panic(spew.Sprintf("Client error: %#v", err))
    }

    client.GetEmployeeEndpoint = logWrapper(client.GetEmployeeEndpoint)

    employee, err := client.GetEmployee(
        context.Background(),
        service.GetEmployeeRequest{
            Phone: "+79165177922",
        },
    )

    if err != nil {
        internalServerError, ok := err.(*service.GetEmployeeResp500)
        if ok {
            panic(internalServerError)
        }
        badRequest, ok := err.(*service.GetEmployeeResp400)
        if ok {
            panic(badRequest.Msg)
        }
        employeeNotFound, ok := err.(*service.GetEmployeeResp204)
        if ok {
            panic(employeeNotFound.Error())
        }
        panic("Unknown error")
    }
    
Создаем клиент:

    client, err := service.NewHTTPClient("http://127.0.0.1")
    if err != nil {
        panic(spew.Sprintf("Client error: %#v", err))
    }
    
Можем кастомно обернуть ендпоинт: 

    client.GetEmployeeEndpoint = logWrapper(client.GetEmployeeEndpoint)
    
Код враппера:

    func logWrapper(endpoint endpoint.Endpoint) endpoint.Endpoint {
        return func(ctx context.Context, request interface{}) (interface{}, error) {
            fmt.Println("log msg")
            resp, err := endpoint(ctx, request)
            // Also can log smth here
            return resp, err
        }
    }
    
Вызываем метод:

    employee, err := client.GetEmployee(
        context.Background(),
        service.GetEmployeeRequest{
            Phone: "+79165177922",
        },
    )
    
И далее обрабатываем ошибки.

Можно кастомизировать ендпоинт:

    func NewGetEmployeeCustomEndpoint(url *url.URL, enc httptransport.EncodeRequestFunc, dec httptransport.DecodeResponseFunc, options ...httptransport.ClientOption) endpoint.Endpoint {
        
Например:

    client.GetEmployeeEndpoint = NewGetEmployeeCustomEndpoint("http://example.com", customEncoder, customDecoder, options)