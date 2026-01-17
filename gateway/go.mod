module gateway

go 1.25.3

require (
	backend/proto v0.0.0
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/schema v1.4.1
	github.com/joho/godotenv v1.5.1
	github.com/rs/cors v1.11.1
	google.golang.org/grpc v1.78.0
)

require (
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260114163908-3f89685c29c3 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace backend/proto => ../proto
