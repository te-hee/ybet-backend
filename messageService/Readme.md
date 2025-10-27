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

