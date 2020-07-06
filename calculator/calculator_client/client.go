package main

import (
	"context"
	"fmt"
	"go-grpc-course/calculator/calculatorpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnary(c)

	// doServerStreaming(c)

	doClientStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Sum Unary RPC")
	req := &calculatorpb.SumRequest{
		FirstNumber:  2,
		SecondNumber: 40,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	req := &calculatorpb.DecomposePrimeNumberRequest{Number: 213122}
	resStream, err := c.DecomposePrimeNumber(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling DecomposePrimeNumber RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from DecomposePrimeNumber: %v\n", msg.GetPrimeFactor())
	}
}

// func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
// 	stream, err := c.ComputeAverage(context.Background())
// 	if err != nil {
// 		log.Fatalf("error while calling ComputeAverage: %v\n", err)
// 	}

// 	// numbers := []*calculatorpb.ComputeAverageRequest{
// 	// 	&calculatorpb.ComputeAverageRequest{
// 	// 		Number: 1,
// 	// 	},
// 	// 	&calculatorpb.ComputeAverageRequest{
// 	// 		Number: 2,
// 	// 	},
// 	// 	&calculatorpb.ComputeAverageRequest{
// 	// 		Number: 3,
// 	// 	},
// 	// 	&calculatorpb.ComputeAverageRequest{
// 	// 		Number: 4,
// 	// 	},
// 	// }

// 	// for _, req := range numbers {
// 	// 	stream.Send(req)
// 	// }

// 	numbers := []int32{3, 5, 9, 54, 23}

// 	for _, number := range numbers {
// 		stream.Send(&calculatorpb.ComputeAverageRequest{
// 			Number: number,
// 		})
// 	}

// 	res, err := stream.CloseAndRecv()
// 	if err != nil {
// 		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
// 	}
// 	fmt.Printf("The Average is: %v\n", res.GetAverage())

// 	// res, err := stream.CloseAndRecv()
// 	// if err != nil {
// 	// 	log.Fatalf("error while receiving response from ComputeAverage: %v", err)
// 	// }
// 	// fmt.Printf("Average: %v\n", res.Average)

// }

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a ComputeAverage Client Streaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	fmt.Printf("The Average is: %v\n", res.GetAverage())
}
