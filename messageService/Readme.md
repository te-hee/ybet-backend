# Message Service

## Typy

```
Message{
    uuid string
    content string
    timestamp int64
}

SendMessageRequest{
    content string
}

SendMessageResponse {
success bool
    error string (pusty kiedy success jest true)
}

GetHistoryRequest{
    limit uint32
}

GetHistoryResponse{
    messages []Message
}
```

## RPC's

`SendMessage(SendMessageRequest) returns (SendMessageResponse)`
wysyłanie wiadomości :3

`GetHistory(GetHistoryRequest) returns (GetHistoryResponse)`
historia wiadomości z limitem :O

`StreamMessages() returns (stream Message)`
stream wiadomości :3

## config
**env**
string
in what enviroment app is running. (used to turn on gRPC reflection for easier debuging)
recommended to set to "dev"
default - ""
env var - `ENV`
flag --env dev

**custom buffer**
int
used to set custom message buffer
default - 100
env var - `BUFFER_SIZE`
flag --buffer 100

**No auth**
bool
used to disable service to service authorization. if set to true no need to pass API token through gRPC metadata.
**If authoriaztion is turned on you have to provide auth key (next config option)**
default - false
env var - `NO_AUTH`
flag --noauth true

**service api key**
string
key that message service expects when auth is turned on
env var - `SERVICE_API_KEY`

**NATS address**
string
address on which NATS is running
NATS default port is `4222`
default - `localhost:4222`
env var - `NATS_ADDRESS`
