package auth

import (
	"context"
	"testing"

	"github.com/sean0427/micro-service-pratice-auth-domain/model"
)

type mockRepo struct{}

func (r *mockRepo) Get(ctx context.Context, params *model.GetProductsParams) ([]*model.Product, error) {
	return []*model.Product{{ID: "test"}}, nil
}

func (r *mockRepo) GetByID(ctx context.Context, id string) (*model.Product, error) {
	return &model.Product{
		ID:   "testfjeia",
		Name: id}, nil
}

func createMockRepo() *mockRepo {
	// TODO
	return &mockRepo{}
}

var testService *ProductService

func TestMain(m *testing.M) {
	testService = New(createMockRepo())
}

func TestProductService_Get(t *testing.T) {
	t.Run("Should success get auth", func(t *testing.T) {
		list, err := testService.Get(context.TODO(), nil)

		if len(list) == 0 {
			t.Errorf("Get auth list is empty")
		}

		if err != nil {
			t.Error(err)
		}
	})
}

func TestProductService_GetByID(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		const testName = "test"
		item, err := testService.GetByID(context.TODO(), testName)
		if err != nil {
			t.Error(err)
		}

		if item.Name != testName {
			t.Errorf("Get auth by name is not equal")
		}

		if item.ID == "" {
			t.Errorf("Returned auth id should not be empty")
		}
	})
}
