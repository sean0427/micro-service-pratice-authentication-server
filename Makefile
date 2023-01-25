protoco_gen: 
	protoc proto/*.proto --go_out=${PWD} --go-grpc_out=${PWD} --experimental_allow_proto3_optional

mock_gen:
	mockgen -package=mock --source=redis_repo/redis.go > mock/redis_client_mock.go
	mockgen -package=mock --source=auth.go > mock/auth_mock.go