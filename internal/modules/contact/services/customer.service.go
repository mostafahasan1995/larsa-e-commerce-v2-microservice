package contact_services

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"larsa-e-commerce-microservice-v2/internal/gateway"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/contact/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/contact/models"
	"net/http"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerRepo interface {
	Create(ctx context.Context, customer *dtos.CreateCustomerDto) (*models.CustomerResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CustomerResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CustomersWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, customer *dtos.CreateCustomerDto) (*models.CustomerResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type CustomerService interface {
	Create(ctx context.Context, customer *dtos.CreateCustomerDto) (*models.CustomerResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CustomerResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CustomersWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, customer *dtos.CreateCustomerDto) (*models.CustomerResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type customer struct {
	repo    CustomerRepo
	usersgw *gateway.UsersGw
	rolesgw *gateway.RolesGw
}

func NewCustomerService(i *do.Injector) (CustomerService, error) {
	return &customer{
		repo:    do.MustInvoke[CustomerRepo](i),
		usersgw: do.MustInvoke[*gateway.UsersGw](i),
		rolesgw: do.MustInvoke[*gateway.RolesGw](i),
	}, nil
}

func (p *customer) Create(ctx context.Context, customer *dtos.CreateCustomerDto) (*models.CustomerResponse, *helpers.ApiError) {

	role, err := p.rolesgw.GetRoleByName(ctx, "vendorEcommercev2")
	if err != nil {
		if err == types.ErrNotFound {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    "Role not found please create it first",
			}
		}

		return nil, &helpers.ApiError{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	user, err := p.usersgw.GetUserByEmail(ctx, customer.BusinessContactInfo.Email)
	if err != nil {
		if customer.BusinessContactInfo.Password == nil {
			Password := util.GeneratePassword(8, 2, 2, 2)
			customer.BusinessContactInfo.Password = &Password
		}
		fmt.Println(*customer.BusinessContactInfo.Password)

		if err == types.ErrUserNotFound {
			//no user with this email - create new one

			newUser := &types.PostUserData{
				FirstName: customer.BasicInformation.FullName["en"],
				LastName:  customer.BasicInformation.OwnerName["en"],
				Email:     customer.BusinessContactInfo.Email,
				Password:  *customer.BusinessContactInfo.Password,
				Roles:     []primitive.ObjectID{role.Id},
			}
			//TODO send password in email

			userId, err := p.usersgw.AddUser(ctx, newUser)

			fmt.Println(user)
			if err != nil {
				return nil, &helpers.ApiError{
					Message:    err.Error(),
					StatusCode: http.StatusInternalServerError,
				}
			}
			if userId.IsZero() {
				return nil, &helpers.ApiError{
					Message:    "User not created",
					StatusCode: http.StatusInternalServerError,
				}
			}

			customer.UserId = &userId

		} else {
			return nil, &helpers.ApiError{
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
			}
		}
	} else {
		userId, _ := primitive.ObjectIDFromHex(user.Id)
		customer.UserId = &userId
		//TODO add role to user

	}

	return p.repo.Create(ctx, customer)
}

func (p *customer) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CustomerResponse, *helpers.ApiError) {
	return p.repo.GetById(ctx, id, query)
}

func (p *customer) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CustomersWithPagination, *helpers.ApiError) {
	return p.repo.GetWithPagination(ctx, skip, limit, query)
}

func (p *customer) UpdateById(ctx context.Context, id primitive.ObjectID, customer *dtos.CreateCustomerDto) (*models.CustomerResponse, *helpers.ApiError) {
	//TODO prevent to update email
	return p.repo.UpdateById(ctx, id, customer)
}

func (p *customer) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return p.repo.DeleteById(ctx, id)
}

func (p *customer) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return p.repo.BulkDelete(ctx, ids)
}
