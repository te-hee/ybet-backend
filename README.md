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

ws://<host>:8081/ws?token=<jwt>

````

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
  "user_id": "string"
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
      "uuid": "string",
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
  "success": true
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

Response:

```json
{
  "success": true
}
```

### DELETE

Request:

```json
{
  "message_id": "string"
}
```

Response:

```json
{
  "success": true
}
```
