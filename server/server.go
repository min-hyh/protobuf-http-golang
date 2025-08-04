package main

import (
	"context"
	"fmt"
	"log"
	pb "protobuf-http-golang/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// server implements the DiscoverServiceServer interface
type server struct {
	pb.UnimplementedDiscoverServiceServer
}

// GetParamInBody implements the GetParamInBody RPC method
func (s *server) GetParamInBody(ctx context.Context, req *pb.GetParamInBodyRequest) (*pb.Response, error) {
	log.Printf("GetParamInBody called with id: %s, content: %s", req.Id, req.Content)

	// Example error handling: validate required fields
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	if req.Content == "" {
		return nil, status.Errorf(codes.InvalidArgument, "content is required")
	}

	// Example error handling: simulate resource not found
	if req.Id == "not-found" {
		return nil, status.Errorf(codes.NotFound, "resource with id '%s' not found", req.Id)
	}

	// Example error handling: simulate permission denied
	if req.Id == "unauthorized" {
		return nil, status.Errorf(codes.PermissionDenied, "access denied for id '%s'", req.Id)
	}

	// Example error handling: simulate internal server error
	if req.Id == "error" {
		return nil, status.Errorf(codes.Internal, "internal server error occurred")
	}

	return &pb.Response{
		NewContent: fmt.Sprintf("Processed ID: %s, Content: %s", req.Id, req.Content),
	}, nil
}

// GetParamInHeader implements the GetParamInHeader RPC method
func (s *server) GetParamInHeader(ctx context.Context, req *pb.GetParamInHeaderRequest) (*pb.Response, error) {
	// Access HTTP headers from gRPC metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// Get specific headers - try both lowercase and original case
		if IDFromHeader := md.Get("x-custom-header-id"); len(IDFromHeader) > 0 {
			log.Printf("Found x-custom-header-id: %s", IDFromHeader[0])
			req.Id = IDFromHeader[0]
		} else {
			log.Printf("No x-custom-header-id found in metadata")
		}
	} else {
		log.Printf("No metadata found in context")
	}

	// Example error handling: validate required header
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "x-custom-header-id header is required")
	}

	// Example error handling: simulate authentication error
	if req.Id == "invalid-token" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authentication token")
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

	// Example error handling: validate required fields
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	if req.Data == nil {
		return nil, status.Errorf(codes.InvalidArgument, "data is required")
	}

	// Example error handling: simulate resource already exists
	if req.Id == "duplicate" {
		return nil, status.Errorf(codes.AlreadyExists, "resource with id '%s' already exists", req.Id)
	}

	// Example error handling: simulate rate limiting
	if req.Id == "rate-limit" {
		return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
	}

	// Echo back the received data
	return &pb.PostUnstructuredDataResponse{
		Id:   req.Id,
		Data: req.Data,
	}, nil
}
