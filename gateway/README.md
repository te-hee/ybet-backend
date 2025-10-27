# Gateway

## Endpointy

- /messages

## /messages

### Metody

- GET - Pobieranie historii wiadomości
- POST - Wysyłanie wiadomości

### Ciała requestów

- GET:
 ```json
{"limit": dodatnia liczba}
```
- POST:
```json
{"content": treść wiadomości w stringu}
```

### Ciała responsów
- GET:  
    W prypadku powodzenia:

    ```json
    {"success": true, "messages": [{"uuid": string, "content": string, "timestamp": uint}, ...]}
    ```
    W przypadku niepowodzenia:

    ```json
    {"success": false, "error": wiadomość errora w stringu}
- POST:  
    W prypadku powodzenia:

    ```json
    {"success": true}
    ```
    W przypadku niepowodzenia:

    ```json
    {"success": false, "error": wiadomość errora w stringu}
    ```


