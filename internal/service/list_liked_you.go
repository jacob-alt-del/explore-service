package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/jacob-alt-del/explore-service/internal/pagination"
	pb "github.com/jacob-alt-del/explore-service/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ExploreServiceServer) ListLikedYou(ctx context.Context, req *pb.ListLikedYouRequest) (*pb.ListLikedYouResponse, error) {
	err := validateListLikedYouRequest(req)
	if err != nil {
		return nil, err
	}

	pageSize := int(req.GetPageSize())
	if pageSize == 0 {
		pageSize = likedYouDefaultPageSize
	}

	paginationUnix, err := pagination.Decode(req.GetPaginationToken())
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Decode() error: %v", err)
	}

	decisions, err := s.Repo.ListLikedYou(ctx, req.RecipientUserId, paginationUnix, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "ListLikedYou() error: %v", err)
	}

	var likers []*pb.ListLikedYouResponse_Liker
	var nextToken string

	// check if next page is needed
	if len(decisions) > pageSize {
		nextToken = pagination.Encode(decisions[pageSize-1].UpdatedAtUnix)
		decisions = decisions[:pageSize]
	}

	for _, d := range decisions {
		likers = append(likers, &pb.ListLikedYouResponse_Liker{
			ActorId:       d.ActorID,
			UnixTimestamp: uint64(d.UpdatedAtUnix),
		})
	}

	return &pb.ListLikedYouResponse{
		Likers:              likers,
		NextPaginationToken: &nextToken,
	}, nil
}

func validateListLikedYouRequest(req *pb.ListLikedYouRequest) error {
	errList := []string{}

	recipientUserID := req.GetRecipientUserId()
	if recipientUserID == "" {
		errList = append(errList, "recipient_user_id is required")
	}
	if !uuidRegex.MatchString(recipientUserID) {
		errList = append(errList, "recipient_user_id must be a valid UUID")
	}

	pageSize := req.GetPageSize()
	if pageSize > likedYouMaxPageSize {
		errList = append(errList, fmt.Sprintf("page_size cannot exceed %v", likedYouMaxPageSize))
	}

	paginationToken := req.GetPaginationToken()
	if paginationToken != "" {
		if _, err := base64.StdEncoding.DecodeString(paginationToken); err != nil {
			errList = append(errList, "pagination_token must be valid base64")
		}
	}

	if len(errList) > 0 {
		return status.Errorf(codes.InvalidArgument, "request validation errors: [%v]", strings.Join(errList, ", "))
	}

	return nil
}
