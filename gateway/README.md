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
    {"Success": true, "Messages": [{"Uuid": string, "Content": string, "Timestamp": uint}, ...]}
    ```
    W przypadku niepowodzenia:

    ```json
    {"Success": false, "Error": wiadomość errora w stringu}
- POST:  
    W prypadku powodzenia:

    ```json
    {"Success": true}
    ```
    W przypadku niepowodzenia:

    ```json
    {"Success": false, "Error": wiadomość errora w stringu}
    ```


