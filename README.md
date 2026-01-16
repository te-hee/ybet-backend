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

## Web socket
- port: 8081
- endpoint: /ws

### /ws?token=jwt_token
#### Dane
**ogólny format**
```json
{
    "type": "messageType",
    "payload": "systemMessage | userMessage | userListUpdate"
}
```
type = payload type
**messageType**: `'systemMessage'|'userMessage'|'userListUpdate'`

**systemMessage**
```json
{
    "content": "string",
}
```

**userMessage**
```json
{
    "uuid": "string",
    "user_id": "string",
    "username": "string",
    "content": "string",
    "timestamp": "uint"
}

```

**userListUpdate**
```json
{
    "action": "'connect'|'disconnect'"
    "uuid": "string"
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
```
{"username": "cutie"}
```
#### Ciała responsów
- POST:
```
{"token": "token jwt"}
```

### /messages

#### Metody

- GET - Pobieranie historii wiadomości
- POST - Wysyłanie wiadomości

#### Ciała requestów

- GET:
 ```
(protocol)://(domain):(port)/messgaes?limit=(num)
 ```
 so for running localy
 ```
 http://localhost:8080/messages?limit=(num)
 ```
- POST:
```json
{"content": "treść wiadomości w stringu"}
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
- POST:  
    W prypadku powodzenia:

    ```json
    {"success": true}
    ```
    W przypadku niepowodzenia:

    ```json
    {"success": false, "error": "wiadomość errora w stringu"}
    ```

