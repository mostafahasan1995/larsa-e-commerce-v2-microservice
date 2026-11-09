package config_repos

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	config_dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	config_utils "larsa-e-commerce-microservice-v2/internal/modules/config/utils"
	seed_utils "larsa-e-commerce-microservice-v2/internal/modules/seed/utils"
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
)

type SettingRepo struct {
	*base_repos.BaseRepo[models.Setting, models.SettingResponse]
}

func NewSettingRepo(i *do.Injector) (services.SettingRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &SettingRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Setting, models.SettingResponse](db, enums.SETTING),
	}, nil
}

func (repo *SettingRepo) FindOneAndIncreaseSerialNumber(ctx context.Context, collectionName string) (*models.SettingResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	setting, err := repo.FindOne(ctx, []bson.M{{"$match": bson.M{"collection": collectionName}}})
	if err != nil {
		return nil, err
	}

	// If no setting exists, create one
	if setting == nil {
		userId := primitive.NilObjectID
		if cfg.User != nil {
			userId = cfg.User.Id
		}

		newSetting := &models.Setting{
			BaseSetting: models.BaseSetting{
				Collection:          collectionName,
				CurrentSerialNumber: 1, // Start with 1
				MetaDataBaseModel: base_models.MetaDataBaseModel{
					Id:        primitive.NewObjectID(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					CreatedBy: &userId,
					UpdatedBy: &userId,
				},
			},
		}

		createdSetting, createErr := repo.Create(ctx, *newSetting)
		if createErr != nil {
			return nil, createErr
		}

		return createdSetting, nil
	}

	// Increment the serial number
	setting.CurrentSerialNumber++
	userId := primitive.NilObjectID
	if cfg.User != nil {
		userId = cfg.User.Id
	}
	repo.UpdateFieldsById(ctx, setting.Id, map[string]interface{}{"currentSerialNumber": setting.CurrentSerialNumber}, userId)
	return setting, nil
}

func (repo *SettingRepo) Upsert() *helpers.ApiError {
	// fmt.Println("Starting settings upsert process...")

	dbs, err := repo.GetRelevantDbs()
	if err != nil {
		return err
	}
	for _, db := range dbs {

		if err := repo.upsertSettings(db); err != nil {
			fmt.Printf("Error upserting settings: %v\n", err)
			return err
		}

		// fmt.Println("Settings upsert completed successfully")

		if err := repo.UpsertIndexes(db); err != nil {
			fmt.Printf("Error upserting indexes: %v\n", err)
			return err
		}

		ctx := repo.AddContext(db)
		db := repo.GetDB().Database(db)
		if err := seed_utils.SeedAttributes(ctx, db, primitive.NilObjectID); err != nil {
			fmt.Printf("Error seeding attributes: %v\n", err)
			return helpers.DefaultError(err)
		}

	}
	// fmt.Println("Indexes upsert completed successfully")
	return nil
}

func (repo *SettingRepo) upsertSettings(databaseName string) *helpers.ApiError {
	collectionNameList := enums.GetCollectionNameList()

	ctx := repo.AddContext(databaseName)
	settingsModel := []models.Setting{}
	settingCollections := []string{}
	for _, collectionName := range collectionNameList {
		result, err := repo.FindOne(ctx, []bson.M{{"$match": bson.M{"collection": collectionName}}})
		if err != nil {
			return err
		}
		if result == nil {
			settingCollections = append(settingCollections, collectionName)
			settingsModel = append(settingsModel, models.Setting{
				BaseSetting: models.BaseSetting{
					Collection:          collectionName,
					CurrentSerialNumber: 0,
					MetaDataBaseModel: base_models.MetaDataBaseModel{
						Id:        primitive.NewObjectID(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						CreatedBy: nil,
						UpdatedBy: nil,
					},
				},
			})
		}

		if len(settingCollections) > 0 {
			// fmt.Printf("Upserting Setting for Db: %v, collections: %v\n", databaseName, settingCollections)
			repo.CreateMany(ctx, settingsModel)
		}
	}

	return nil
}

func (repo *SettingRepo) UpsertIndexes(databaseName string) *helpers.ApiError {
	collectionNameList := enums.GetCollectionNameList()

	ctx := repo.AddContext(databaseName)
	for _, collectionName := range collectionNameList {
		var indexes = config_utils.GetCollectionIndexes(collectionName)
		// fmt.Printf("Upserting Indexes For %v\n", collectionName)

		dbIndexes := []string{}
		myIndexes := []string{"_id_"}
		for _, index := range indexes {
			myIndexes = append(myIndexes, index.Inserted)
		}

		storedIndexes, _ := repo.GetDB().Database(databaseName).Collection(collectionName).Indexes().ListSpecifications(ctx)
		for _, index := range storedIndexes {
			dbIndexes = append(dbIndexes, index.Name)
		}

		if len(storedIndexes) > 0 {
			for _, index := range storedIndexes {
				if !util.Contains(myIndexes, index.Name) {
					fmt.Printf("Dropping Index %v\n", index.Name)
					_, err := repo.GetDB().Database(databaseName).Collection(collectionName).Indexes().DropOne(ctx, index.Name)
					if err != nil {
						fmt.Printf("Error dropping index %v: %v\n", index.Name, err)
					}
					dbIndexes = append(dbIndexes, index.Name)
				}
			}
		}

		if len(indexes) > 0 {
			for _, index := range indexes {
				if !util.Contains(dbIndexes, index.Inserted) {
					fmt.Printf("Creating Index %v\n", index.Inserted)
					_, err := repo.GetDB().Database(databaseName).Collection(collectionName).Indexes().CreateOne(ctx, index.ToBeInserted)
					if err != nil {
						fmt.Printf("Error creating index %v: %v\n", index.Inserted, err)
					}
				}
			}
		}

	}

	return nil
}

func (repo *SettingRepo) GetRelevantDbs() ([]string, *helpers.ApiError) {
	var dbs []string
	ch := make(chan string)
	go common.FetchSubscriptions(common.ClientSubscriptionInstance, ch)
	<-ch
	for _, v := range common.ClientSubscriptionInstance.SubscriptionData {
		microservices := v.Microservices
		for _, m := range microservices {
			if m == "e-commerce-v2" {
				dbs = append(dbs, v.Id)
				break
			}
		}
	}

	dbs = append(dbs, "604674966fe156529e3e03e9")
	return dbs, nil
}

func (repo *SettingRepo) AddContext(db string) context.Context {

	cfg := &types.AppCfg{
		Db:   db,
		Lang: string(enums.EN),
	}

	return context.WithValue(context.Background(), util.ReqAppCfg, cfg)
}

func (repo *SettingRepo) TruncateCollections(databaseName string, dto config_dtos.TruncateCollectionDtoList) *helpers.ApiError {
	ctx := repo.AddContext(databaseName)
	for _, item := range dto.Items {
		filter := bson.M{}
		if item.Filter != nil {
			filter = *item.Filter
		}
		repo.GetDB().Database(databaseName).Collection(item.CollectionName).DeleteMany(ctx, filter)
	}
	return nil
}
