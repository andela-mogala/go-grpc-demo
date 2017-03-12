package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"flag"

	"github.com/andela-mogala/go-grpc-demo/src/grpcdemo/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const port = ":9000"

func main() {
	option := flag.Int("o", 1, "Command to run")
	flag.Parse()
	creds, err := credentials.NewClientTLSFromFile("server.crt", "")
	if err != nil {
		log.Fatal(err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial("localhost"+port, opts...)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	client := pb.NewEmployeeServiceClient(conn)

	switch *option {
	case 1:
		sendMetaData(client)
	case 2:
		getByBadgeNumber(client)
	case 3:
		getAll(client)
	case 4:
		addPhoto(client)
	case 5
		saveAll(client)
	}
}

func saveAll(client pb.EmployeeServiceClient) {
	employees := []pb.Employee{
		pb.Employee{
			BadgeNumber: 123,
			FirstName: "John",
			LastName: "Smith",
			VacationAccrualRate: 1.2,
			VacationAccrued: 0,
		},
		pb.Employee{
			BadgeNumber: 234,
			FirstName: "Lisa",
			LastName: "Wu",
			VacationAccrualRate: 1.7,
			VacationAccrued: 10,
		},
	}
	stream, err := client.SaveAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	doneCh := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				doneCh <- struct{}{}
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(res.Employee)
		}
	}()

	for _, e := range employees {
		err := stream.Send(&pb.EmployeeRequest{Employee: &e})
		if err != nil {
			log.Fatal(err)
		}
	}
	stream.CloseSend()
	<- doneCh
}

func addPhoto(client pb.EmployeeServiceClient) {
	f, err := os.Open("file.mp3")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	md := metadata.New(map[string]string{"badgenumber": "2080"})
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, md)
	stream, err := client.AddPhoto(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for {
		chunk := make([]byte, 64*1024)
		n, err := f.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if n < len(chunk) {
			chunk = chunk[:n]
		}
		stream.Send(&pb.AddPhotoRequest{Data: chunk})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.IsOk)
}

func getAll(client pb.EmployeeServiceClient) {
	stream, err := client.GetAll(context.Background(), &pb.GetAllRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res.Employee)
	}
}

func getByBadgeNumber(client pb.EmployeeServiceClient) {
	res, err := client.GetByBadgeNumber(context.Background(), &pb.GetByBadgeNumberRequest{BadgeNumber: 2080})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Employee)
}

func sendMetaData(client pb.EmployeeServiceClient) {
	md := metadata.MD{}
	md["user"] = []string{"michael"}
	md["password"] = []string{"pass"}
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, md)
	client.GetByBadgeNumber(ctx, &pb.GetByBadgeNumberRequest{})
}
