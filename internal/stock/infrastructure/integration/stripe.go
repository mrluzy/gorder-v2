package integration

import (
	"context"
	_ "github.com/mrluzy/gorder-v2/common/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/product"
	"log"
)

type StripeAPI struct {
	apiKey string
}

func NewStripeAPI() *StripeAPI {
	key := viper.GetString("stripe-key")
	if key == "" {
		log.Fatal("stripe-key is empty")
	}
	return &StripeAPI{apiKey: key}
}

func (s *StripeAPI) GetPriceByProductID(ctx context.Context, pid string) (string, error) {
	stripe.Key = s.apiKey
	logrus.Infof("apiKey: %s, pid: %s", s.apiKey, pid)
	result, err := product.Get(pid, &stripe.ProductParams{})
	if err != nil {
		return "", err
	}
	return result.DefaultPrice.ID, nil
}

func (s *StripeAPI) GetProductByID(ctx context.Context, pid string) (*stripe.Product, error) {
	stripe.Key = s.apiKey
	return product.Get(pid, &stripe.ProductParams{})
}
