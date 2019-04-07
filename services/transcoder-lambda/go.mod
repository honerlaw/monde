module services/transcoder-lambda

replace services/server v0.0.0 => ../server

require (
	github.com/aws/aws-lambda-go v1.10.0
	github.com/aws/aws-sdk-go v1.19.11
	github.com/joho/godotenv v1.3.0
	github.com/stretchr/testify v1.3.0 // indirect
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3 // indirect
	services/server v0.0.0
)
