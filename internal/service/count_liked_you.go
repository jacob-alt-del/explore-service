package service

import (
	"context"
	"strings"

	pb "github.com/jacob-alt-del/explore-service/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ExploreServiceServer) CountLikedYou(ctx context.Context, req *pb.CountLikedYouRequest) (*pb.CountLikedYouResponse, error) {
	err := validateCountLikedYouRequest(req)
	if err != nil {
		return nil, err
	}

	count, err := s.Repo.CountLikedYou(ctx, req.RecipientUserId)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "CountLikedYou error: %v", err)
	}

	return &pb.CountLikedYouResponse{
		Count: count,
	}, nil
}

func validateCountLikedYouRequest(req *pb.CountLikedYouRequest) error {
	errList := []string{}

	recipientUserID := req.GetRecipientUserId()
	if recipientUserID == "" {
		errList = append(errList, "recipient_user_id is required")
	}
	if !uuidRegex.MatchString(recipientUserID) {
		errList = append(errList, "recipient_user_id must be a valid UUID")
	}

	if len(errList) > 0 {
		return status.Errorf(codes.InvalidArgument, "request validation errors: [%v]", strings.Join(errList, ", "))
	}

	return nil
}
