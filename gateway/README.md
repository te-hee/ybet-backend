# Gateway

## Endpointy

- /messages

## /messages

### Metody

- GET - Pobieranie historii wiadomości
- POST - Wysyłanie wiadomości

### Ciała requestów

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

### Ciała responsów
- GET:  
    W prypadku powodzenia:

    ```json
    {"success": true, "messages": [{"uuid": "string", "Content": "string", "Timestamp": "string"}, ...]}
    ```
    W przypadku niepowodzenia:

    ```json
    {"Success": false, "Error": "wiadomość errora w stringu"}
- POST:  
    W prypadku powodzenia:

    ```json
    {"Success": true}
    ```
    W przypadku niepowodzenia:

    ```json
    {"Success": false, "Error": "wiadomość errora w stringu"}
    ```


