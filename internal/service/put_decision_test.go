package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jacob-alt-del/explore-service/internal/dataaccess"
	pb "github.com/jacob-alt-del/explore-service/internal/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPutDecision(t *testing.T) {
	// Keep originals so we can restore them after
	origUpsert := fnUpsertDecision
	origCheck := fnCheckMutualLike
	defer func() {
		fnUpsertDecision = origUpsert
		fnCheckMutualLike = origCheck
	}()

	validUUID1 := "550e8400-e29b-41d4-a716-446655440000"
	validUUID2 := "123e4567-e89b-12d3-a456-426614174000"

	tests := []struct {
		name           string
		req            *pb.PutDecisionRequest
		mockUpsertErr  error
		mockCheckLike  bool
		mockCheckErr   error
		wantErrCode    codes.Code
		wantMutualLike bool
	}{
		{
			name: "success with mutual like",
			req: &pb.PutDecisionRequest{
				ActorUserId:     validUUID1,
				RecipientUserId: validUUID2,
				LikedRecipient:  true,
			},
			mockCheckLike:  true,
			wantMutualLike: true,
		},
		{
			name: "success no mutual like",
			req: &pb.PutDecisionRequest{
				ActorUserId:     validUUID1,
				RecipientUserId: validUUID2,
				LikedRecipient:  true,
			},
			mockCheckLike:  false,
			wantMutualLike: false,
		},
		{
			name: "UpsertDecision returns error",
			req: &pb.PutDecisionRequest{
				ActorUserId:     validUUID1,
				RecipientUserId: validUUID2,
				LikedRecipient:  true,
			},
			mockUpsertErr: errors.New("db error"),
			wantErrCode:   codes.Unknown,
		},
		{
			name: "CheckMutualLike returns error",
			req: &pb.PutDecisionRequest{
				ActorUserId:     validUUID1,
				RecipientUserId: validUUID2,
				LikedRecipient:  true,
			},
			mockCheckErr: errors.New("check error"),
			wantErrCode:  codes.Unknown,
		},
		{
			name: "invalid request (missing actor ID)",
			req: &pb.PutDecisionRequest{
				ActorUserId:     "",
				RecipientUserId: validUUID2,
				LikedRecipient:  true,
			},
			wantErrCode: codes.InvalidArgument, // assuming your ValidatePutDecisionRequest returns this
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override the function variables for test
			fnUpsertDecision = func(repo *dataaccess.Repository, ctx context.Context, actorID, recipientID string, liked bool) error {
				return tt.mockUpsertErr
			}
			fnCheckMutualLike = func(repo *dataaccess.Repository, ctx context.Context, actorID, recipientID string) (bool, error) {
				return tt.mockCheckLike, tt.mockCheckErr
			}

			s := &ExploreServiceServer{Repo: &dataaccess.Repository{}}

			resp, err := s.PutDecision(context.Background(), tt.req)

			if tt.wantErrCode != codes.OK {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantErrCode, st.Code(), "unexpected gRPC error code")
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, tt.wantMutualLike, resp.MutualLikes)
		})
	}
}

func Test_validatePutDecisionRequest(t *testing.T) {
	validUUID1 := "550e8400-e29b-41d4-a716-446655440000"
	validUUID2 := "123e4567-e89b-12d3-a456-426614174000"

	tests := []struct {
		name    string
		req     *pb.PutDecisionRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: &pb.PutDecisionRequest{
				ActorUserId:     validUUID1,
				RecipientUserId: validUUID2,
			},
			wantErr: false,
		},
		{
			name: "missing actor_user_id",
			req: &pb.PutDecisionRequest{
				ActorUserId:     "",
				RecipientUserId: validUUID2,
			},
			wantErr: true,
			errMsg:  "actor_user_id is required",
		},
		{
			name: "missing recipient_user_id",
			req: &pb.PutDecisionRequest{
				ActorUserId:     validUUID1,
				RecipientUserId: "",
			},
			wantErr: true,
			errMsg:  "recipient_user_id is required",
		},
		{
			name: "invalid actor_user_id format",
			req: &pb.PutDecisionRequest{
				ActorUserId:     "not-a-uuid",
				RecipientUserId: validUUID2,
			},
			wantErr: true,
			errMsg:  "actor_user_id must be a valid UUID",
		},
		{
			name: "invalid recipient_user_id format",
			req: &pb.PutDecisionRequest{
				ActorUserId:     validUUID1,
				RecipientUserId: "invalid-uuid",
			},
			wantErr: true,
			errMsg:  "recipient_user_id must be a valid UUID",
		},
		{
			name: "actor_user_id equals recipient_user_id",
			req: &pb.PutDecisionRequest{
				ActorUserId:     validUUID1,
				RecipientUserId: validUUID1,
			},
			wantErr: true,
			errMsg:  "recipient_user_id must not equal actor_user_id",
		},
		{
			name: "multiple validation errors",
			req: &pb.PutDecisionRequest{
				ActorUserId:     "",
				RecipientUserId: "not-a-uuid",
			},
			wantErr: true,
			errMsg:  "actor_user_id is required, actor_user_id must be a valid UUID, recipient_user_id must be a valid UUID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePutDecisionRequest(tt.req)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				st, ok := status.FromError(err)
				if !ok || st.Code() != codes.InvalidArgument {
					t.Fatalf("expected InvalidArgument error, got %v", err)
				}
				if !strings.Contains(st.Message(), tt.errMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errMsg, st.Message())
				}
			} else if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
