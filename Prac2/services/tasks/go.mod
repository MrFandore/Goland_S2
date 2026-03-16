module github.com/MrFandore/Go_S2/Prac2/services/tasks

go 1.24.0

require (
	github.com/go-chi/chi/v5 v5.0.10
	google.golang.org/grpc v1.79.2
)

require github.com/MrFandore/Go_S2 v0.0.0-00010101000000-000000000000

require (
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/MrFandore/Go_S2 => ../../../
