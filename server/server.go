package main

import (
	"context"
	"fmt"
	"log"
	pb "protobuf-http-golang/pb"
)

// server implements the DiscoverServiceServer interface
type server struct {
	pb.UnimplementedDiscoverServiceServer
}

// GetParamInBody implements the GetParamInBody RPC method
func (s *server) GetParamInBody(ctx context.Context, req *pb.GetParamInBodyRequest) (*pb.Response, error) {
	log.Printf("GetParamInBody called with id: %s, content: %s", req.Id, req.Content)

	return &pb.Response{
		NewContent: fmt.Sprintf("Processed ID: %s, Content: %s", req.Id, req.Content),
	}, nil
}

// GetParamInHeader implements the GetParamInHeader RPC method
func (s *server) GetParamInHeader(ctx context.Context, req *pb.GetParamInHeaderRequest) (*pb.Response, error) {
	log.Printf("GetParamInHeader called with id: %s, content: %s", req.Id, req.Content)

	return &pb.Response{
		NewContent: fmt.Sprintf("Header processed - ID: %s, Content: %s", req.Id, req.Content),
	}, nil
}

// PostUnstructuredData implements the PostUnstructuredData RPC method
func (s *server) PostUnstructuredData(
	ctx context.Context,
	req *pb.PostUnstructuredDataRequest,
) (*pb.PostUnstructuredDataResponse, error) {
	log.Printf("PostUnstructuredData called with id: %s", req.Id)

	// Echo back the received data
	return &pb.PostUnstructuredDataResponse{
		Id:   req.Id,
		Data: req.Data,
	}, nil
}
