package gapi

// func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
// 	violations := validateLoginUserRequest(req)
// 	if violations != nil {
// 		return nil, invalidArgumentError(violations)
// 	}

// 	user, err := server.store.GetUser(ctx, req.GetEmail())
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, status.Errorf(codes.NotFound, "user not found")
// 		}
// 		return nil, status.Errorf(codes.Internal, "failed to find user")
// 	}

// 	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
// 		user.Email,
// 		server.config.AccessTokenDuration,
// 	)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to create access token")
// 	}

// 	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
// 		user.Email,
// 		server.config.RefreshTokenDuration,
// 	)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
// 	}

// 	mtdt := server.extractMetadata(ctx)
// 	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
// 		ID:           refreshPayload.ID,
// 		UserID:       user.ID,
// 		RefreshToken: refreshToken,
// 		UserAgent:    mtdt.UserAgent,
// 		ClientIp:     mtdt.ClientIP,
// 		IsBlocked:    false,
// 		ExpiresAt:    refreshPayload.ExpiredAt,
// 	})
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to create session")
// 	}

// 	rsp := &pb.LoginUserResponse{
// 		User:                  convertUser(user),
// 		SessionId:             session.ID.String(),
// 		AccessToken:           accessToken,
// 		RefreshToken:          refreshToken,
// 		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
// 		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
// 	}
// 	return rsp, nil
// }

// func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
// 	if err := val.ValidatePassword(req.GetPassword()); err != nil {
// 		violations = append(violations, fieldViolation("password", err))
// 	}

// 	return violations
// }
