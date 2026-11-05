package utils

import (
	"errors"
	"roomService/internal/core/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AppErrorToGrpcError(incomingError error) error {

	var appError *domain.AppError
	if !errors.As(incomingError, &appError) {
		return status.Error(codes.Unknown, incomingError.Error())
	}

	switch appError.Code {
	case domain.CodeNotFound:
		return status.Error(codes.NotFound, appError.Message)
	case domain.CodeInvalidArgument:
		return status.Error(codes.InvalidArgument, appError.Message)
	case domain.CodePermissionDenied:
		return status.Error(codes.PermissionDenied, appError.Message)
	case domain.CodeAlreadyExists:
		return status.Error(codes.AlreadyExists, appError.Message)
	case domain.CodeUnauthenticated:
		return status.Error(codes.Unauthenticated, appError.Message)
	case domain.CodeInternal:
		return status.Error(codes.Internal, appError.Message)
	}

	return status.Error(codes.Unknown, appError.Message)
}
