package ios

import (
	"context"
	"errors"

	pb "github.com/u03013112/ss-pb/ios"
)

// Srv ：服务
type Srv struct{}

// Purchase : 支付确认
func (s *Srv) Purchase(ctx context.Context, in *pb.PurchaseRequest) (*pb.PurchaseReply, error) {
	return &pb.PurchaseReply{}, errors.New("auth failed")
}
