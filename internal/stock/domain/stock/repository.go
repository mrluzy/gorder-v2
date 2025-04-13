package stock

import (
	"context"
	"fmt"
	"github.com/mrluzy/gorder-v2/common/genproto/orderpb"
	"strings"
)

type Repository interface {
	GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error)
}

type NotFoundError struct {
	Missing []string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("not found in stock:%s", strings.Join(e.Missing, ","))
}
