# Ybet - **Y**apping **b**ut **E**nd **T**o (End encrypted)

Our skilled team:
- Radiant-ちゃん
- Julek <3
- Adas :*
- **Krzyś >.<**
- **hipo :333**
- **Dawid =^w^=**

## Web socket
- port: 8081
- endpoint: /ws

### /ws

#### responsy
w przypadku niepowodzenia w połączeniu do websocketa
```json
{"error": treść błędu}

```

## gateway

- /messages

### /messages

#### Metody

- GET - Pobieranie historii wiadomości
- POST - Wysyłanie wiadomości

#### Ciała requestów

- GET:
 ```json
{"limit": dodatnia liczba}
```
- POST:
```json
{"content": treść wiadomości w stringu}
```

#### Ciała responsów
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

