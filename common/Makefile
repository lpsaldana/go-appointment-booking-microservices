gen:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       pb/auth.proto pb/professional.proto pb/client.proto \
		   pb/agenda.proto pb/notification.proto