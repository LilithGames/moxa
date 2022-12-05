package service

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/LilithGames/protoc-gen-dragonboat/runtime"
)

const DragonboatGrpcErrorCodeBase int32 = 20000000
const DragonboatGrpcErrorCodeEnd  int32 = 30000000

func GetGrpcErrorCode(err error) codes.Code {
	if s, ok := status.FromError(err); !ok {
		return codes.Unknown
	} else {
		return s.Code()
	}
}

func DragonboatErrorToGrpcError(err error) error {
	if err == nil {
		return nil
	}
	code := runtime.GetDragonboatErrorCode(err)
	if code == runtime.ErrCodeOK {
		return nil
	}
	return grpc.Errorf(DragonboatErrCodeToGrpcErrCode(code), err.Error())
}

func GrpcErrorToDragonboatError(err error) error {
	if e, ok := status.FromError(err); ok {
		dcode := GrpcErrCodeToDragonboatErrCode(e.Code())
		if dcode == nil {
			if e.Code() == codes.OK {
				return nil
			} else if e.Code() == codes.InvalidArgument {
				return runtime.NewDragonboatError(runtime.ErrCodeBadReqeust, err.Error())
			}
			return runtime.NewDragonboatError(runtime.ErrCodeInternal, err.Error())
		}
		return runtime.NewDragonboatError(*dcode, err.Error())
	}
	return runtime.NewDragonboatError(runtime.ErrCodeInternal, err.Error())
}

func DragonboatErrCodeToGrpcErrCode(code int32) codes.Code {
	return codes.Code(DragonboatGrpcErrorCodeBase+code)
}

func GrpcErrCodeToDragonboatErrCode(code codes.Code) *int32 {
	if DragonboatGrpcErrorCodeBase <= int32(code) && int32(code) < DragonboatGrpcErrorCodeEnd {
		result := int32(code) - DragonboatGrpcErrorCodeBase
		return &result
	}
	return nil
}
