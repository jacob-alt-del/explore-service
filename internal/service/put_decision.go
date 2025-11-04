package service

import (
	"context"
	"strings"

	"github.com/jacob-alt-del/explore-service/internal/dataaccess"
	pb "github.com/jacob-alt-del/explore-service/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// allow stubbing data access repo for unit testing
var (
	fnUpsertDecision = func(repo *dataaccess.Repository, ctx context.Context, actorID, recipientID string, liked bool) error {
		return repo.UpsertDecision(ctx, actorID, recipientID, liked)
	}
	fnCheckMutualLike = func(repo *dataaccess.Repository, ctx context.Context, actorID, recipientID string) (bool, error) {
		return repo.CheckMutualLike(ctx, actorID, recipientID)
	}
)

func (s *ExploreServiceServer) PutDecision(ctx context.Context, req *pb.PutDecisionRequest) (*pb.PutDecisionResponse, error) {
	err := ValidatePutDecisionRequest(req)
	if err != nil {
		return nil, err
	}

	err = fnUpsertDecision(s.Repo, ctx, req.ActorUserId, req.RecipientUserId, req.LikedRecipient)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "UpsertDecision() error: %v", err)
	}

	mutualLike := false
	if req.LikedRecipient {
		mutualLike, err = fnCheckMutualLike(s.Repo, ctx, req.ActorUserId, req.RecipientUserId)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "CheckMutualLike() error: %v", err)
		}
	}

	return &pb.PutDecisionResponse{
		MutualLikes: mutualLike,
	}, nil
}

func ValidatePutDecisionRequest(req *pb.PutDecisionRequest) error {
	errList := []string{}

	actorUserID := req.GetActorUserId()
	if actorUserID == "" {
		errList = append(errList, "actor_user_id is required")
	}
	if !uuidRegex.MatchString(actorUserID) {
		errList = append(errList, "actor_user_id must be a valid UUID")
	}

	recipientUserID := req.GetRecipientUserId()
	if recipientUserID == "" {
		errList = append(errList, "recipient_user_id is required")
	}
	if !uuidRegex.MatchString(recipientUserID) {
		errList = append(errList, "recipient_user_id must be a valid UUID")
	}

	if actorUserID == recipientUserID {
		errList = append(errList, "recipient_user_id must not equal actor_user_id")
	}

	if len(errList) > 0 {
		return status.Errorf(codes.InvalidArgument, "request validation errors: [%v]", strings.Join(errList, ", "))
	}

	return nil
}
