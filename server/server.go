package main

import (
	"context"
	"fmt"
	"log"
	pb "protobuf-http-golang/pb"

	"google.golang.org/grpc/metadata"
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

	// Access HTTP headers from gRPC metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Printf("Metadata detected: %+v", md)

		// Get specific headers - try both lowercase and original case
		if IDFromHeader := md.Get("x-custom-header-id"); len(IDFromHeader) > 0 {
			log.Printf("Found x-custom-header-id: %s", IDFromHeader[0])
			req.Id = IDFromHeader[0]
		} else if IDFromHeader := md.Get("X-Custom-Header-Id"); len(IDFromHeader) > 0 {
			log.Printf("Found X-Custom-Header-Id: %s", IDFromHeader[0])
			req.Id = IDFromHeader[0]
		} else {
			log.Printf("No X-Custom-Header-Id found in metadata")
		}

		// Log all available metadata keys for debugging
		for key, values := range md {
			log.Printf("Metadata key: %s, values: %v", key, values)
		}
	} else {
		log.Printf("No metadata found in context")
	}

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
