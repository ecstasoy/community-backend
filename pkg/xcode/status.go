package xcode

import (
	"community-backend/pkg/xcode/types"
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var _ XCode = (*Status)(nil)

type Status struct {
	sts *types.Status
}

func (s *Status) Error() string {
	return s.Message()
}

func (s *Status) Message() string {
	if s.sts.Message == "" {
		return strconv.Itoa(int(s.sts.Code))
	}

	return s.sts.Message
}

func (s *Status) Details() []interface{} {
	if s == nil || s.sts == nil {
		return nil
	}

	details := make([]interface{}, 0, len(s.sts.Details))
	for _, d := range s.sts.Details {
		detail := &anypb.Any{}

		if err := d.UnmarshalTo(detail); err != nil {
			details = append(details, err)
			continue
		}
		details = append(details, detail)
	}

	return details
}

func (s *Status) WithDetails(msgs ...proto.Message) (*Status, error) {
	for _, msg := range msgs {
		anyMsg, err := anypb.New(msg)
		if err != nil {
			return s, errors.Wrap(err, "failed to marshal message to any")
		}
		s.sts.Details = append(s.sts.Details, anyMsg)
	}

	return s, nil
}

func (s *Status) Code() int {
	return int(s.sts.Code)
}

func (s *Status) Proto() *types.Status {
	return s.sts
}

func Error(code Code) *Status {
	return &Status{sts: &types.Status{Code: int32(code.Code()), Message: code.Message()}}
}

func Errorf(code Code, format string, args ...interface{}) *Status {
	code.msg = fmt.Sprintf(format, args...)
	return Error(code)
}

// CodeFromError converts an error to an XCode
func CodeFromError(err error) XCode {
	err = errors.Cause(err)
	if xc, ok := err.(XCode); ok {
		return xc
	}

	switch err {
	case context.Canceled:
		return Canceled
	case context.DeadlineExceeded:
		return Deadline
	default:
		grpcStatus, _ := status.FromError(err)
		return GRPCStatusToXCode(grpcStatus)
	}
}

// StatusFromError converts an error to a gRPC status.Status
func StatusFromError(err error) *status.Status {
	err = errors.Cause(err)
	if xc, ok := err.(XCode); ok {
		grpcStatus, e := gRPCStatusFromXCode(xc)
		if e == nil {
			return grpcStatus
		}
	}

	var grpcStatus *status.Status
	switch err {
	case context.Canceled:
		grpcStatus, _ = gRPCStatusFromXCode(Canceled)
	case context.DeadlineExceeded:
		grpcStatus, _ = gRPCStatusFromXCode(Deadline)
	default:
		grpcStatus, _ = status.FromError(err)
	}

	return grpcStatus
}

// StatusFromCode converts an XCode to a Status
func StatusFromCode(code Code) *Status {
	return &Status{sts: &types.Status{Code: int32(code.Code()), Message: code.Message()}}
}

// CodeFromProto converts a protobuf message to an XCode
func CodeFromProto(pbMsg proto.Message) XCode {
	msg, ok := pbMsg.(*types.Status)
	if ok {
		if len(msg.Message) == 0 || msg.Message == strconv.FormatInt(int64(msg.Code), 10) {
			return Code{code: int(msg.Code)}
		}
		return &Status{sts: msg}
	}

	return Errorf(ServerErr, "cannot convert proto message to xcode status")
}

// gRPCStatusFromXCode converts an XCode to a gRPC status.Status
func gRPCStatusFromXCode(code XCode) (*status.Status, error) {
	var sts *Status
	switch v := code.(type) {
	case *Status:
		sts = v
	case Code:
		sts = StatusFromCode(v)
	default:
		sts = Error(Code{code.Code(), code.Message()})
		for _, detail := range code.Details() {
			if msg, ok := detail.(proto.Message); ok {
				_, _ = sts.WithDetails(msg)
			}
		}
	}

	stas := status.New(codes.Unknown, strconv.Itoa(sts.Code()))
	return stas.WithDetails(sts.Proto())
}

// GRPCStatusToXCode converts a gRPC status.Status to an XCode
func GRPCStatusToXCode(grpcStatus *status.Status) XCode {
	details := grpcStatus.Details()
	if len(details) > 0 {
		for i := len(details) - 1; i >= 0; i-- {
			if stProto, ok := details[i].(*types.Status); ok {
				return CodeFromProto(stProto)
			}
		}
	}

	return toXCode(grpcStatus)
}

func toXCode(grpcStatus *status.Status) Code {
	grpcCode := grpcStatus.Code()
	switch grpcCode {
	case codes.OK:
		return OK
	case codes.InvalidArgument:
		return RequestErr
	case codes.NotFound:
		return NotFound
	case codes.PermissionDenied:
		return AccessDenied
	case codes.Unauthenticated:
		return Unauthorized
	case codes.ResourceExhausted:
		return LimitExceed
	case codes.Unimplemented:
		return MethodNotAllowed
	case codes.DeadlineExceeded:
		return Deadline
	case codes.Unavailable:
		return ServiceUnavailable
	case codes.Unknown:
		return String(grpcStatus.Message())
	}

	return ServerErr
}
