
fmt:
	go fmt ./...

test:
	time go test ./...

cover_middlewares:
	go test github.com/deliverous/cocktails/middlewares -coverprofile=/tmp/coverage.out
	go tool cover -html=/tmp/coverage.out

cover_context:
	go test github.com/deliverous/cocktails/httpcontext -coverprofile=/tmp/coverage.out
	go tool cover -html=/tmp/coverage.out

cover_render:
	go test github.com/deliverous/cocktails/render -coverprofile=/tmp/coverage.out
	go tool cover -html=/tmp/coverage.out

docserve:
	godoc -http=:6060 &
