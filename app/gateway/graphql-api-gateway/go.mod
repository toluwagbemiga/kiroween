module github.com/haunted-saas/graphql-api-gateway

go 1.21

require (
	github.com/99designs/gqlgen v0.17.43
	github.com/graph-gophers/dataloader/v7 v7.1.0
	github.com/rs/cors v1.10.1
	github.com/vektah/gqlparser/v2 v2.5.11
	go.uber.org/zap v1.26.0
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/sosodev/duration v1.2.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231212172506-995d672761c0 // indirect
)

// Local proto dependencies
replace (
	github.com/haunted-saas/analytics-service => ../../services/analytics-service
	github.com/haunted-saas/billing-service => ../../services/billing-service
	github.com/haunted-saas/feature-flags-service => ../../services/feature-flags-service
	github.com/haunted-saas/llm-gateway-service => ../../services/llm-gateway-service
	github.com/haunted-saas/notifications-service => ../../services/notifications-service
	github.com/haunted-saas/user-auth-service => ../../services/user-auth-service
)
