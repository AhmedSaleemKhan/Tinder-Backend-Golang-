package gapi

// import (
// 	"context"
// 	"database/sql"
// 	"encoding/json"

// 	"github.com/tabbed/pqtype"
// 	db "github.com/techschool/simplebank/db/sqlc"
// 	"github.com/techschool/simplebank/pb"
// 	"github.com/techschool/simplebank/val"
// 	"google.golang.org/genproto/googleapis/rpc/errdetails"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
// 	authPayload, err := server.authorizeUser(ctx)
// 	if err != nil {
// 		return nil, unauthenticatedError(err)
// 	}

// 	violations := validateUpdateUserRequest(req)
// 	if violations != nil {
// 		return nil, invalidArgumentError(violations)
// 	}

// 	if authPayload.ID != req.GetId() {
// 		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
// 	}

// 	arg := db.UpdateUserParams{
// 		ID: req.GetId(),
// 		FirstName: sql.NullString{
// 			String: req.GetFirstName(),
// 			Valid:  req.FirstName != nil,
// 		},
// 		LastName: sql.NullString{
// 			String: req.GetLastName(),
// 			Valid:  req.LastName != nil,
// 		},
// 		Email: sql.NullString{
// 			String: req.GetEmail(),
// 			Valid:  req.Email != nil,
// 		},
// 		Phone: sql.NullString{
// 			String: req.GetPhone(),
// 			Valid:  req.Phone != nil,
// 		},
// 		Age: sql.NullInt64{
// 			Int64: req.GetAge(),
// 			Valid: req.Age != nil,
// 		},
// 		Gender: sql.NullString{
// 			String: req.GetGender(),
// 			Valid:  req.Gender != nil,
// 		},
// 		Ethnicity: req.Ethnicity,
// 		Nsfw: sql.NullBool{
// 			Bool:  req.GetNsfw(),
// 			Valid: req.Nsfw != nil,
// 		},
// 		Metadata: pqtype.NullRawMessage{
// 			RawMessage: json.RawMessage(req.GetMetadata()),
// 			Valid:      req.Metadata != nil,
// 		},
// 	}

// 	user, err := server.store.UpdateUser(ctx, arg)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, status.Errorf(codes.NotFound, "user not found")
// 		}
// 		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
// 	}

// 	rsp := &pb.UpdateUserResponse{
// 		User: convertUser(user),
// 	}
// 	return rsp, nil
// }

// func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
// 	if req.Email != nil {
// 		if err := val.ValidateEmail(req.GetEmail()); err != nil {
// 			violations = append(violations, fieldViolation("email", err))
// 		}
// 	}

// 	return violations
// }
