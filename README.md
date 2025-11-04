# Ybet - **Y**apping **b**ut **E**nd **T**o (End encrypted)

Our skilled team:
- Radiant-ちゃん
- Julek <3
- Adas :*
- **Krzyś >.<**
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
    "target": "string",
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
{"content": treść wiadomości w stringu}
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

