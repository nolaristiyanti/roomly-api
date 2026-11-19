package postgres

import (
	"context"

	"gorm.io/gorm"

	"roomly-api/domain/room/entity"
	"roomly-api/domain/room/model"
	"roomly-api/domain/room/repository"

	facilityentity "roomly-api/domain/facility/entity"
	facilitymodel "roomly-api/domain/facility/model"
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

	err := r.db.WithContext(ctx).Preload("Facilities").First(&roomModel, id).Error

	// SELECT f.*
	// FROM facilities AS f
	// JOIN room_facilities AS rf
	//     ON rf.facility_id = f.id
	// WHERE rf.room_id = 1;

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
		Facilities:  mapFacilities(roomModel.Facilities),
	}

	return roomEntity, nil
}

// FindList mengambil seluruh data room
func (r *RoomRepositoryStruct) FindList(ctx context.Context) ([]entity.Room, error) {
	var roomModels []model.Room

	err := r.db.WithContext(ctx).Preload("Facilities").Order("id ASC").Find(&roomModels).Error

	// SELECT *
	// FROM room_facilities AS rf
	// JOIN facilities AS f
	//     ON f.id = rf.facility_id
	// ORDER BY rf.room_id ASC;

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
			Facilities:  mapFacilities(roomModel.Facilities),
		}

		roomEntities = append(roomEntities, roomEntity)
	}

	return roomEntities, nil
}

func mapFacilities(facilities []facilitymodel.Facility) []facilityentity.Facility {
	result := make([]facilityentity.Facility, 0, len(facilities))

	for _, facility := range facilities {
		result = append(result, facilityentity.Facility{
			ID:          facility.ID,
			Name:        facility.Name,
			Description: facility.Description,
		})
	}

	return result
}
