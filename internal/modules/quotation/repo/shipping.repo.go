package quotation_repos

// import (
// 	"context"
// 	"larsa-e-commerce-microservice-v2/internal/common/types"
// 	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
// 	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"

// 	"github.com/samber/do"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"

// 	"larsa-e-commerce-microservice-v2/internal/common/enums"
// 	"larsa-e-commerce-microservice-v2/internal/common/helpers"
// 	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
// 	dtos "larsa-e-commerce-microservice-v2/internal/modules/quotation/dto"
// 	models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"

// 	"larsa-e-commerce-microservice-v2/internal/common/util"
// 	utils "larsa-e-commerce-microservice-v2/internal/modules/quotation/utils"
// )

// type ShippingRepo struct {
// 	*base_repos.BaseRepo[models.Shipping, models.ShippingResponse]
// }

// func NewShippingRepo(i *do.Injector) (services.ShippingRepo, error) {
// 	db := do.MustInvoke[*mongo.Client](i)
// 	return &ShippingRepo{
// 		BaseRepo: base_repos.NewBaseRepo[models.Shipping, models.ShippingResponse](db, enums.SHIPPING),
// 	}, nil
// }

// type shippingOperation struct {
// 	shippingDto       *dtos.SubmitToAgentDto
// 	ctx               context.Context
// 	additionalOptions types.AdditionalOptions
// }

// func (repo *ShippingRepo) executeShippingOperation(op shippingOperation, updateId *primitive.ObjectID) (*models.ShippingResponse, *helpers.ApiError) {

// 	cfg, _ := util.GetReqAppCfg(op.ctx)

// 	shipping := utils.ConvertToSubmitToAgentModel(op.shippingDto, op.additionalOptions)

// 	// Start a session
// 	session, err := repo.GetDB().StartSession()
// 	if err != nil {
// 		return nil, helpers.DefaultError(err)
// 	}
// 	defer session.EndSession(op.ctx)

// 	// Start transaction
// 	err = session.StartTransaction()
// 	if err != nil {
// 		return nil, helpers.DefaultError(err)
// 	}

// 	// Create channels for goroutine results
// 	type result struct {
// 		response *models.ShippingResponse
// 		err      *helpers.ApiError
// 	}
// 	shippingChan := make(chan result, 1)
// 	quotationChan := make(chan *helpers.ApiError, 1)

// 	// Create/Update shipping cost in a goroutine
// 	go func() {
// 		var response *models.ShippingResponse
// 		var apiErr *helpers.ApiError
// 		if updateId != nil {
// 			response, apiErr = repo.BaseRepo.UpdateById(op.ctx, *updateId, *shipping, bson.M{})
// 		} else {
// 			response, apiErr = repo.BaseRepo.Create(op.ctx, *shipping)
// 		}
// 		shippingChan <- result{response: response, err: apiErr}
// 	}()

// 	// Update quotation in a goroutine
// 	go func() {
// 		quotationRepo := base_repos.NewBaseRepo[quotation_models.Quotation, quotation_models.QuotationResponse](repo.GetDB(), enums.QUOTATION)
// 		_, apiErr := quotationRepo.UpdateFieldsById(op.ctx, shipping.QuotationId, bson.M{"shippingId": shipping.Id}, cfg.User.Id)
// 		quotationChan <- apiErr
// 	}()

// 	// Wait for both operations to complete
// 	shippingResult := <-shippingChan
// 	quotationErr := <-quotationChan

// 	// Check for errors
// 	if shippingResult.err != nil {
// 		session.AbortTransaction(op.ctx)
// 		return nil, shippingResult.err
// 	}
// 	if quotationErr != nil {
// 		session.AbortTransaction(op.ctx)
// 		return nil, quotationErr
// 	}

// 	// Commit transaction
// 	err = session.CommitTransaction(op.ctx)
// 	if err != nil {
// 		session.AbortTransaction(op.ctx)
// 		return nil, helpers.DefaultError(err)
// 	}

// 	return shippingResult.response, nil
// }

// func (repo *ShippingRepo) SubmitToAgent(ctx context.Context, SubmitToAgentDto *dtos.SubmitToAgentDto) (*models.ShippingResponse, *helpers.ApiError) {

// 	cfg, _ := util.GetReqAppCfg(ctx)

// 	return repo.executeShippingOperation(shippingOperation{
// 		shippingDto: SubmitToAgentDto,
// 		ctx:         ctx,
// 		additionalOptions: types.AdditionalOptions{
// 			UserId: cfg.User.Id,
// 			Mode:   types.CREATE,
// 		},
// 	}, nil)
// }

// func (repo *ShippingRepo) UpdateSubmitToAgentById(ctx context.Context, id primitive.ObjectID, updateShippingDto *dtos.SubmitToAgentDto) (*models.ShippingResponse, *helpers.ApiError) {
// 	cfg, _ := util.GetReqAppCfg(ctx)
// 	return repo.executeShippingOperation(shippingOperation{
// 		shippingDto: updateShippingDto,
// 		ctx:         ctx,
// 		additionalOptions: types.AdditionalOptions{
// 			UserId: cfg.User.Id,
// 			Mode:   types.UPDATE,
// 		},
// 	}, &id)
// }

// func (repo *ShippingRepo) UpsertShippingFromAgent(ctx context.Context, shipping models.ShippingResponse, shippingAgent *dtos.UpsertShippingAgentDto) *helpers.ApiError {

// 	cfg, _ := util.GetReqAppCfg(ctx)

// 	From := utils.ConvertToUpsertShippingAgentModel(shippingAgent, types.AdditionalOptions{
// 		UserId: cfg.User.Id,
// 		Mode:   types.UPDATE,
// 	})

// 	updatedFields := bson.M{
// 		"from.method":           From.Method,
// 		"from.countryIso3":      From.CountryIso3,
// 		"from.customsClearance": From.CustomsClearance,
// 		"from.logistic":         From.Logistic,
// 		"from.insurance":        From.Insurance,
// 	}

// 	if shipping.From.Method != nil {
// 		updatedFields["status"] = quotation_enums.SHIPPING_APPROVED
// 	}

// 	_, err := repo.BaseRepo.UpdateFieldsById(ctx, shipping.Id, updatedFields, cfg.User.Id)
// 	if err != nil {
// 		return helpers.DefaultError(err)
// 	}

// 	return nil
// }

// func (repo *ShippingRepo) UpsertShippingToAgent(ctx context.Context, shipping models.ShippingResponse, shippingAgent *dtos.UpsertShippingAgentDto) *helpers.ApiError {

// 	cfg, _ := util.GetReqAppCfg(ctx)

// 	To := utils.ConvertToUpsertShippingAgentModel(shippingAgent, types.AdditionalOptions{
// 		UserId: cfg.User.Id,
// 		Mode:   types.UPDATE,
// 	})

// 	updatedFields := bson.M{
// 		"to.method":           To.Method,
// 		"to.countryIso3":      To.CountryIso3,
// 		"to.customsClearance": To.CustomsClearance,
// 		"to.logistic":         To.Logistic,
// 		"to.insurance":        To.Insurance,
// 	}

// 	if shipping.From.Method != nil {
// 		updatedFields["status"] = quotation_enums.SHIPPING_APPROVED
// 	}

// 	_, err := repo.BaseRepo.UpdateFieldsById(ctx, shipping.Id, updatedFields, cfg.User.Id)
// 	if err != nil {
// 		return helpers.DefaultError(err)
// 	}

// 	return nil
// }

// func (repo *ShippingRepo) SubmitToVendor(ctx context.Context, shipping models.ShippingResponse, shippingVendor *dtos.SubmitToVendorDto) *helpers.ApiError {

// 	cfg, _ := util.GetReqAppCfg(ctx)

// 	quotationRepo := base_repos.NewBaseRepo[quotation_models.Quotation, quotation_models.QuotationResponse](repo.GetDB(), enums.QUOTATION)

// 	platformCost := utils.ConvertToSubmitToVendorModel(shippingVendor, types.AdditionalOptions{
// 		UserId: cfg.User.Id,
// 		Mode:   types.UPDATE,
// 	})

// 	updatedFields := bson.M{
// 		"platformCost": platformCost,
// 		"status":       quotation_enums.SHIPPING_COMPLETED,
// 	}

// 	_, err := repo.BaseRepo.UpdateFieldsById(ctx, shipping.Id, updatedFields, cfg.User.Id)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = quotationRepo.UpdateFieldsById(ctx, shipping.QuotationId, bson.M{"status": quotation_enums.READY_TO_SEND}, cfg.User.Id)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
