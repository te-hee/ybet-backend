# Ybet - **Y**apping **b**ut **E**nd **T**o (End encrypted)

Our skilled team:
- Adas :*
- **hipo :333**
- **Dawid =^w^=**

## jak z tym pracować?
*włączanie*
`docker compose up --build -d`
*wyłączanie*
`docker compose down`

## architektura
`Client -> Gateway -> (Auth / Message) -> NATS -> Broadcast -> WebSocket`

## Web socket
- port: 8081
- endpoint: /ws

### /ws?token=jwt_token
#### Dane
**ogólny format**
```json
{
    "type": "messageType",
    "payload": "message"
}
```
type = payload type
**messageType**: `'systemMessage'|'userMessage'|'userListUpdate'|'editMessage'|'deleteMessage'`

**userMessage**
```json
{
    "message_id": "string",
    "user_id": "string",
    "username": "string",
    "content": "string",
    "timestamp": "uint"
}
```

**editMessage**
```json
{
    "message_id": "string",
    "content": "string"
}
```

**deleteMessage**
```json
{
    "message_id": "string",
}
```

**systemMessage**
```json
{
    "content": "string",
}
```

**userListUpdate**
```json
{
    "action": "'connect'|'disconnect'"
    "user_id": "string"
}
```

## gateway

- /messages
- /login

### autoryzacja
1. wyślij POST na `/login` `{"username":"cute"}`
2. dostajesz spowrotem `{"token":"token jwt"}`
3. przy każdym requeście ustaw header `Authorization` i wartość `Bearer <token>`

Żeby wyłączyć autoryzację użyj flagi `--noauth` lub zmienna środowiskowa `NO_AUTH` na `true`
wtedy każda wiadomość będzie od jedengo użytkownika testowego

### env

Przykładaowa env'ka znajudje się w katalogu gateway w przypadku jej braku są używane defultowe wartości

- MESSAGE_SERVICE_IP - ip do message servicu z portem (ip:port)
- GATEWAY_PORT - Port na jaki będzie nasłuchiwał gateway
- NO_AUTH - wyłącza autoryzację

### /login

#### Metody
- POST - zalogowanie się do systemu (tylko username)

#### Ciała requestów
- POST:
```json
{"username": "cutie"}
```
#### Ciała responsów
- POST:
```json
{"token": "token jwt"}
```

### /messages

#### Metody

- GET - Pobieranie historii wiadomości
- POST - Wysyłanie wiadomości
- PATCH - Edytowanie wiadomości
- DELETE - usuwanie wiadomości

#### Ciała requestów

- GET:
 ```
(protocol)://(domain):(port)/messages?limit=(num)
 ```
 so for running localy
 ```
 http://localhost:8080/messages?limit=(num)
 ```
- POST:
```json
{"content": "treść wiadomości w stringu"}
```
- PATCH
```json
{"message_id": "uuid", "content": "string"}
```
- DELETE
```json
{"message_id": "uuid"}
```

#### Ciała responsów
- GET:  
    W prypadku powodzenia:

    ```json
    {"success": true, "messages": [{"uuid": "string", "content": "string", "timestamp": "uint"}, ...]}
    ```
    W przypadku niepowodzenia:

    ```json
    {"success": false, "error": "wiadomość errora w stringu"}
    ```
- POST:  
    W prypadku powodzenia:

    ```json
    {"success": true}
    ```
    W przypadku niepowodzenia:

    ```json
    {"success": false, "error": "wiadomość errora w stringu"}
    ```
- PATCH:
    ```json
    {"success": true}
    ```
- DELETE:
    ```json
    {"success": true}
    ```
