# Ybet - **Y**apping **b**ut **E**nd **T**o
> End-to-end encrypted chat system

## Nasz zespół
- Adas :*
- hipo :333
- Dawid =^w^=
- ishu uwu

---

## Uruchamianie

Start:
```

docker compose up --build -d

```

Stop:
```

docker compose down

```

---

## Architektura

```

Client
→ Gateway
→ Auth / Message
→ NATS
→ Broadcast
→ WebSocket

```

---

## WebSocket

Endpoint:
```

ws://<host>:8081/ws

````
### Potrzebny header
- `Authorization: Bearer <token>`

Format wiadomości:
```json
{
  "type": "<type>",
  "payload": { ... }
}
````

Typy `type`:

* `userMessage`
* `editMessage`
* `deleteMessage`
* `systemMessage`
* `userListUpdate`

### userMessage

```json
{
  "message_id": "string",
  "user_id": "string",
  "username": "string",
  "content": "string",
  "timestamp": number
}
```

### editMessage

```json
{
  "message_id": "string",
  "content": "string"
}
```

### deleteMessage

```json
{
  "message_id": "string"
}
```

### systemMessage

```json
{
  "content": "string"
}
```

### userListUpdate

```json
{
  "action": "connect | disconnect",
  "user_id": "string",
  "username": "string"
}
```

---

## Gateway API

Endpointy:

* POST   /login
* GET    /messages
* POST   /messages
* PATCH  /messages
* DELETE /messages

### Autoryzacja

POST `/login`
Request:

```json
{
  "username": "string"
}
```

Response:

```json
{
  "token": "<jwt>"
}
```

Header dla requestów (poza /login):

```
Authorization: Bearer <jwt>
```

#### Tryb testowy (bez auth)

Flaga `--noauth` lub env `NO_AUTH=true`
Wszystkie akcje wykonywane jako użytkownik testowy.

---

## Zmienne środowiskowe

| Zmienna            | Typ     | Opis                              |
| ------------------ | ------- | --------------------------------- |
| MESSAGE_SERVICE_IP | string  | adres message service (`ip:port`) |
| GATEWAY_PORT       | number  | port gatewaya                     |
| NO_AUTH            | boolean | wyłączenie autoryzacji            |

---

## /messages

### GET

```
/messages?limit=<number>
```

Response success:

```json
{
  "success": true,
  "messages": [
    {
        "message_id": "string",
        "user_id": "string",
        "username": "string,"
        "content": "string",
        "timestamp": number
    }
  ]
}
```

Response error:

```json
{
  "success": false,
  "error": "string"
}
```

### POST

Request:

```json
{
  "content": "string"
}
```

Response:

```json
{
    "message_id": "string",
    "timestamp": number
}
```

### PATCH

Request:

```json
{
  "message_id": "string",
  "content": "string"
}
```

Status codes:
- **200** - udało się :3
- **400** - zły format json :c
- **401** - brak header `Authorization Bearer <token>` >~<
- **403** - nie możesz edytować tej wiadomości >:c 
- **500** - internal server error lub unknown error TwT
### DELETE

Request:

```json
{
  "message_id": "string"
}
```

Status codes:
- **200** - udało się :3
- **400** - zły format json :c
- **401** - brak header `Authorization Bearer <token>` >~<
- **403** - nie możesz usunąć tej wiadomości >:c 
- **500** - internal server error lub unknown error TwT
