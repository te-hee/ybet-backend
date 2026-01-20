package utils

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCToHTTPResponse(err error) (int, map[string]string) {
	if err == nil {
		return http.StatusOK, nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, map[string]string{
			"error": "Internal Gateway Error: " + err.Error(),
		}
	}

	switch st.Code() {
	case codes.OK:
		return http.StatusOK, nil

	case codes.InvalidArgument:
		return http.StatusBadRequest, map[string]string{"error": st.Message()}

	case codes.NotFound:
		return http.StatusNotFound, map[string]string{"error": st.Message()}

	case codes.AlreadyExists:
		return http.StatusConflict, map[string]string{"error": st.Message()}

	case codes.PermissionDenied:
		return http.StatusForbidden, map[string]string{"error": st.Message()}

	case codes.Unauthenticated:
		return http.StatusUnauthorized, map[string]string{"error": st.Message()}

	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout, map[string]string{"error": "Service timeout"}

	case codes.Unavailable:
		return http.StatusServiceUnavailable, map[string]string{"error": "Service unavailable"}

	default:
		return http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"}
	}
}
