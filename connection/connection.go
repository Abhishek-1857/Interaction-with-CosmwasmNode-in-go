package connection
import "google.golang.org/grpc"

func Connection()*grpc.ClientConn {
	grpcConn, err := grpc.Dial(
		"127.0.0.1:9090",
		grpc.WithInsecure(),
	)
	if err != nil {
		panic(err)

	}
	return grpcConn
}