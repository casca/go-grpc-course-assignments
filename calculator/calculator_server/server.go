package main

import (
	"context"
	"fmt"
	"go-grpc-course/calculator/calculatorpb"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)
	firstNumber := req.GetFirstNumber()
	secondNumber := req.GetSecondNumber()
	res := &calculatorpb.SumResponse{
		SumResult: firstNumber + secondNumber,
	}
	return res, nil
}

func (*server) DecomposePrimeNumber(req *calculatorpb.DecomposePrimeNumberRequest, stream calculatorpb.CalculatorService_DecomposePrimeNumberServer) error {
	divisor := int32(2)
	number := req.GetNumber()

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.DecomposePrimeNumberResponse{
				PrimeFactor: divisor,
			})
			time.Sleep(time.Second)
			number /= divisor
		} else {
			divisor++
		}
	}

	return nil
}

// func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
// 	var sum int32
// 	var count int32
// 	for {
// 		req, err := stream.Recv()
// 		if err == io.EOF {
// 			stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
// 				Average: float64(sum) / float64(count),
// 			})
// 		}
// 		if err != nil {
// 			log.Fatalf("error while reading client stream: %v\n", err)
// 		}
// 		sum += req.GetNumber()
// 		count++
// 	}
// }

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("Received ComputeAverage RPC\n")

	sum := int32(0)
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		sum += req.GetNumber()
		count++
	}

}

func main() {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
