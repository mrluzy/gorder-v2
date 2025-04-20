package adapters

import (
	"context"
	_ "github.com/mrluzy/gorder-v2/common/config"
	domain "github.com/mrluzy/gorder-v2/order/domain/order"
	"github.com/mrluzy/gorder-v2/order/entity"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	dbName   = viper.GetString("mongo.db-name")
	collName = viper.GetString("mongo.coll-name")
)

type OrderRepositoryMongo struct {
	db *mongo.Client
}

func NewOrderRepositoryMongo(db *mongo.Client) *OrderRepositoryMongo {
	return &OrderRepositoryMongo{db: db}
}

func (r *OrderRepositoryMongo) Collection() *mongo.Collection {
	return r.db.Database(dbName).Collection(collName)
}

type orderModel struct {
	MongoID     primitive.ObjectID `bson:"_id"`
	ID          string             `bson:"id"`
	CustomerID  string             `bson:"customer_id"`
	Status      string             `bson:"status"`
	PaymentLink string             `bson:"payment_link"`
	Items       []*entity.Item     `bson:"items"`
}

func (r *OrderRepositoryMongo) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
	defer r.logWithTag("create", err, created)
	write := r.marshalToModel(order)
	resp, err := r.Collection().InsertOne(ctx, write)
	if err != nil {
		return nil, err
	}
	created = order
	created.ID = resp.InsertedID.(primitive.ObjectID).Hex()
	return
}

func (r *OrderRepositoryMongo) logWithTag(tag string, err error, result interface{}) {
	l := logrus.WithFields(logrus.Fields{
		"tag":             "order_repository_mongo",
		"performed _time": time.Now().Unix(),
		"err":             err,
		"result":          result,
	})
	if err != nil {
		l.Infof("%s_fail", tag)
	} else {
		l.Infof("%s_success", tag)
	}
}

func (r *OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
	defer r.logWithTag("get", err, got)
	read := &orderModel{}
	mongoID, _ := primitive.ObjectIDFromHex(id)
	// condition
	cond := bson.M{"_id": mongoID}
	if err = r.Collection().FindOne(ctx, cond).Decode(read); err != nil {
		return
	}
	if read == nil {
		return nil, domain.NotFoundError{OrderID: id}
	}
	return r.unmarshal(read), nil
}

func (r *OrderRepositoryMongo) Update(
	ctx context.Context,
	order *domain.Order,
	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error),
) (err error) {

	defer r.logWithTag("update", err, nil)
	if order == nil {
		panic("nil order")
	}
	session, err := r.db.StartSession()
	if err != nil {
		return
	}
	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		return err
	}
	defer func() {
		if err == nil {
			_ = session.CommitTransaction(ctx)
		} else {
			_ = session.AbortTransaction(ctx)
		}
	}()

	// inside transaction
	oldOrder, err := r.Get(ctx, order.ID, order.CustomerID)
	if err != nil {
		return
	}
	updated, err := UpdateFn(ctx, order)
	if err != nil {
		return
	}
	MongoID, _ := primitive.ObjectIDFromHex(oldOrder.ID)
	res, err := r.Collection().UpdateOne(
		ctx,
		bson.M{"_id": MongoID, "customer_id": oldOrder.CustomerID},
		bson.M{"$set": bson.M{
			"status":       updated.Status,
			"payment_link": updated.PaymentLink,
		}},
	)
	if err != nil {
		return
	}
	r.logWithTag("finish_update", err, res)
	return
}

func (r *OrderRepositoryMongo) marshalToModel(order *domain.Order) *orderModel {
	return &orderModel{
		MongoID:     primitive.NewObjectID(),
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
}

func (r *OrderRepositoryMongo) unmarshal(m *orderModel) *domain.Order {
	return &domain.Order{
		ID:          m.MongoID.Hex(),
		CustomerID:  m.CustomerID,
		Status:      m.Status,
		PaymentLink: m.PaymentLink,
		Items:       m.Items,
	}
}
