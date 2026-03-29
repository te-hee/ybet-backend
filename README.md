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
```bash

make build-run

```
lub
```bash

make build
make run

```
Żeby sprawdzić logi:
```bash

make logs

```

Stop:
```bash

make down

```

Testy:
```bash

make test

```
---

Można użyć `make help`, żeby sprawdzić wszystkie komendy :3

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

## Gateway API (wszystko pod spodem ma routa /api)

### Wersja 1 (Wszystko pod spodem ma routa /v1)

#### Errory

Erorr wystąpi gdy status jest 4XX lub 5XX każdy response errora zwraca:
```json
{
  "error": "string"
}
```
w dalej części zwracane errory są pominięte

#### Endpointy:

* POST   /login
* GET    /messages
* POST   /messages
* PATCH  /messages
* DELETE /messages

#### Autoryzacja

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

---

#### Zmienne środowiskowe

| Zmienna            | Typ     | Opis                              |
| ------------------ | ------- | --------------------------------- |
| MESSAGE_SERVICE_IP | string  | adres message service (`ip:port`) |
| GATEWAY_PORT       | number  | port gatewaya                     |
| AUTH           | boolean | autoryzacja            |

---

#### /messages

##### GET

```
/messages?limit=<number>
```

Response success:

```json
{
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

##### POST

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

##### PATCH

Request:

```json
{
  "message_id": "string",
  "content": "string"
}
```

Status codes:
- **200** - udało się
- **400** - zły format json
- **401** - brak header `Authorization Bearer <token>`
- **403** - nie możesz edytować tej wiadomości
- **500** - internal server error lub unknown error

##### DELETE

Request:

```json
{
  "message_id": "string"
}
```

Status codes:
- **200** - udało się 
- **400** - zły format json 
- **401** - brak header `Authorization Bearer <token>`
- **403** - nie możesz usunąć tej wiadomości 
- **500** - internal server error lub unknown error
