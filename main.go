package main

import (
	"context"
	"time"

	"pex-backend/http_server"
)

const serviceAddress = ":10101"

func main() {
	ctx := context.Background()
	cancellable, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	ctx, cancelFunc := context.WithTimeout(cancellable, 5*time.Second)
	defer cancelFunc()

	http_server.InstantiateServer(serviceAddress)
}
