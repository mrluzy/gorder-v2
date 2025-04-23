package adapters

import (
	"context"
	"fmt"
	_ "github.com/mrluzy/gorder-v2/common/config"
	"github.com/mrluzy/gorder-v2/stock/entity"
	"github.com/mrluzy/gorder-v2/stock/infrastructure/persistent"
	"github.com/mrluzy/gorder-v2/stock/infrastructure/persistent/builder"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"testing"
)

func setupTestDB(t *testing.T) *persistent.MySQL {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		"",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	testDB := viper.GetString("mysql.dbname") + "_shadow"

	assert.NoError(t, db.Exec("DROP DATABASE IF EXISTS "+testDB).Error)
	assert.NoError(t, db.Exec("CREATE DATABASE IF NOT EXISTS "+testDB).Error)

	dsn = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		testDB,
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&persistent.StockModel{}))
	return persistent.NewMySQLWithDB(db)
}

func TestMySQLStockRepository_UpdateStock_Race(t *testing.T) {
	db := setupTestDB(t)
	t.Parallel()

	var (
		ctx          = context.Background()
		testItem     = "test-race-item"
		initialStock = 100
	)

	err := db.Create(ctx, &persistent.StockModel{
		ProductID: testItem,
		Quantity:  int32(initialStock),
	})
	assert.NoError(t, err)

	repo := NewMySQLStockRepository(db)
	var wg sync.WaitGroup
	goroutineCount := 10
	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.UpdateStock(
				ctx,
				[]*entity.ItemWithQuantity{
					{ID: testItem, Quantity: 1},
				},
				func(ctx context.Context, existing []*entity.ItemWithQuantity, query []*entity.ItemWithQuantity) ([]*entity.ItemWithQuantity, error) {
					var newItems []*entity.ItemWithQuantity
					for _, e := range existing {
						for _, q := range query {
							if e.ID == q.ID {
								newItems = append(newItems, &entity.ItemWithQuantity{
									ID:       e.ID,
									Quantity: e.Quantity - q.Quantity,
								})
							}
						}
					}
					return newItems, nil
				},
			)
			assert.NoError(t, err)
		}()
	}
	wg.Wait()

	res, err := db.BatchGetStockByID(ctx, builder.NewStock().ProductIDs(testItem))
	assert.NoError(t, err)
	assert.NotEmpty(t, res, "res should not be empty")

	expected := initialStock - goroutineCount
	assert.EqualValues(t, expected, res[0].Quantity)
}

func TestMySQLStockRepository_UpdateStock_OverSell(t *testing.T) {
	db := setupTestDB(t)
	t.Parallel()

	var (
		ctx          = context.Background()
		testItem     = "test-race-item"
		initialStock = 5
	)

	err := db.Create(ctx, &persistent.StockModel{
		ProductID: testItem,
		Quantity:  int32(initialStock),
	})
	assert.NoError(t, err)

	repo := NewMySQLStockRepository(db)
	var wg sync.WaitGroup
	goroutineCount := 50
	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.UpdateStock(
				ctx,
				[]*entity.ItemWithQuantity{
					{ID: testItem, Quantity: 1},
				},
				func(ctx context.Context, existing []*entity.ItemWithQuantity, query []*entity.ItemWithQuantity) ([]*entity.ItemWithQuantity, error) {
					var newItems []*entity.ItemWithQuantity
					for _, e := range existing {
						for _, q := range query {
							if e.ID == q.ID {
								newItems = append(newItems, &entity.ItemWithQuantity{
									ID:       e.ID,
									Quantity: e.Quantity - q.Quantity,
								})
							}
						}
					}
					return newItems, nil
				},
			)
			assert.NoError(t, err)
		}()
	}
	wg.Wait()
	res, err := db.BatchGetStockByID(ctx, builder.NewStock().ProductIDs([]string{testItem}...))
	assert.NoError(t, err)
	assert.NotEmpty(t, res, "res should not be empty")

	//assert.EqualValues(t, int32(0), res[0].Quantity)
	assert.GreaterOrEqual(t, res[0].Quantity, int32(0))
}
