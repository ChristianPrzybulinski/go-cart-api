package discount

import (
	context "context"
	"time"

	log "github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
)

func DescountPercentage(port string, product int32, timeout int) float32 {

	var percentageResponse float32 = 0
	var conn *grpc.ClientConn

	log.Infoln("Connecting in gRPC server in port: " + port)

	conn, err := grpc.Dial(port, grpc.WithInsecure())

	if err == nil {
		client := NewDiscountClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second) //default 1 second timeout for tests

		defer cancel()

		resp, err := client.GetDiscount(ctx, &GetDiscountRequest{ProductID: product})

		if err != nil {
			log.Errorln("Error when calling GetDiscount: ", err)
		} else {
			percentageResponse = resp.Percentage
		}
	}

	defer conn.Close()
	return percentageResponse
}
