package gapi

import (
	"context"

	"github.com/FrostJ143/simplebank/internal/query"
	"github.com/FrostJ143/simplebank/internal/val"
	"github.com/FrostJ143/simplebank/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	txResult, err := server.store.VerifyEmailTx(ctx, query.VerifyEmailTxParams{
		EmailID:    req.GetId(),
		SecretCode: req.GetSecretCode(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can not verify email: %s", err)
	}

	res := &pb.VerifyEmailResponse{
		IsVerified: txResult.User.IsEmailVerified,
	}

	return res, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	err := val.ValidateEmailID(req.GetId())
	if err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}

	err = val.ValidateSecretCode(req.GetSecretCode())
	if err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}

	return
}
