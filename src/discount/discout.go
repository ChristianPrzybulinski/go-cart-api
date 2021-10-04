package discount

import (
	context "context"
	"time"

	log "github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
)

func DescountPercentage(port string, product int32) float32 {

	var percentageResponse float32 = 0
	var conn *grpc.ClientConn

	log.Infoln("Connecting in gRPC server in port: " + port)

	conn, err := grpc.Dial(port, grpc.WithInsecure())

	if err != nil {
		log.Errorln("did not connect: ", err)
	} else {

		client := NewDiscountClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		resp, err := client.GetDiscount(ctx, &GetDiscountRequest{ProductID: product})

		if err != nil {
			log.Errorln("Error when calling GetDiscount: ", err)
		} else {
			log.Infoln("Response from server: ", resp.Percentage)
			percentageResponse = resp.Percentage
		}
	}

	defer conn.Close()
	return percentageResponse
}
