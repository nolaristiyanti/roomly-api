package postgres

import (
	"context"

	"gorm.io/gorm"

	"roomly-api/domain/room/entity"
	"roomly-api/domain/room/model"
	"roomly-api/domain/room/repository"
)

type RoomRepositoryStruct struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) repository.RoomRepository {
	return &RoomRepositoryStruct{
		db: db,
	}
}

// FindByID mengambil satu room berdasarkan ID
func (r *RoomRepositoryStruct) FindByID(ctx context.Context, id int64) (*entity.Room, error) {
	var roomModel model.Room

	err := r.db.WithContext(ctx).First(&roomModel, id).Error
	if err != nil {
		return nil, err
	}

	roomEntity := &entity.Room{
		ID:          roomModel.ID,
		Name:        roomModel.Name,
		Description: roomModel.Description,
		Capacity:    roomModel.Capacity,
		Location:    roomModel.Location,
		IsActive:    roomModel.IsActive,
	}

	return roomEntity, nil
}

// FindList mengambil seluruh data room
func (r *RoomRepositoryStruct) FindList(ctx context.Context) ([]entity.Room, error) {
	var roomModels []model.Room

	err := r.db.WithContext(ctx).Order("id ASC").Find(&roomModels).Error
	if err != nil {
		return nil, err
	}

	var roomEntities []entity.Room

	for _, roomModel := range roomModels {
		roomEntity := entity.Room{
			ID:          roomModel.ID,
			Name:        roomModel.Name,
			Description: roomModel.Description,
			Capacity:    roomModel.Capacity,
			Location:    roomModel.Location,
			IsActive:    roomModel.IsActive,
		}

		roomEntities = append(roomEntities, roomEntity)
	}

	return roomEntities, nil
}
