# Gateway

## Endpointy

- /messages

## /messages

`GET /messages?limit=10`
Pobiera historię wiadomości
`limit` określa maksymalną liczbę zwróconych wiadomości

**Example Request**
curl -X GET 'https://api.example.com/v1/messages?limit=10'

**Example response**
Success (200 OK)
```json
{
    "Success": true,
    "Messages": [
        {"Uuid": "a1b2c3", "Content": "Hello world!", "Timestamp": 1730049060},
        {"Uuid": "d4e5f6", "Content": "meowwww", "Timestamp": 1730049060},
        {"Uuid": "g7h8i9", "Content": "Hi cutie~", "Timestamp": 1730049060}
    ]
}
```

Error (400 Bad Request)
```json
{
    "Success": false,
    "Error": "Failed reading message history"
}
```

`POST /messages`
Wysyła wiadomość

**Example request body (JSON)**
```json
{ "content": "Hello world!" }
```

**Example response**
Success (200 OK)
```json
{ "Success": true }
```

Error (500 Internal Server Error)
```json
{
    "Success": false,
    "Error": "Failed to send message"
}
```


