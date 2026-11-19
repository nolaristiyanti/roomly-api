package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
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

func (r *RoomRepositoryStruct) Create(ctx context.Context, room entity.Room) (*entity.Room, error) {
	// Mapping Entity menjadi GORM Model Room
	roomModel := model.Room{
		Name:        room.Name,
		Description: room.Description,
		Capacity:    room.Capacity,
		Location:    room.Location,
		IsActive:    true,
	}

	// Simpan Room dan relasi Facility dalam satu transaction
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// // Simpan room ke PostgreSQL
		if err := tx.Create(&roomModel).Error; err != nil {
			return err
		}

		// INSERT INTO rooms (name, capacity, location, is_active)
		// VALUES ('Bali', 12, 'Lantai 10', true)
		// RETURNING id;

		// Jika terdapat Facility, relasinya disimpan ke tabel room_facilities
		if len(room.Facilities) > 0 {
			var facilityModels []facilitymodel.Facility

			for _, facility := range room.Facilities {
				facilityModels = append(facilityModels, facilitymodel.Facility{
					ID: facility.ID,
				})
			}

			if err := tx.Model(&roomModel).Association("Facilities").Append(&facilityModels); err != nil {
				return err
			}

			// Model(&roomModel) : Room yang menjadi sumber relasi
			// Association("Facilities") : Mengakses relasi Facilities pada Room
			// Append(&facilityModels) : Menambahkan Facility ke relasi Room

			// INSERT INTO room_facilities (room_id, facility_id)
			// VALUES
			// (1, 1),
			// (1, 2);
		}

		return nil
	}) // commit / rollback

	if err != nil {
		return nil, err
	}

	// Jika berhasil, get Room beserta Facility nya
	return r.FindByID(ctx, roomModel.ID)
}

func (r *RoomRepositoryStruct) Update(ctx context.Context, room entity.Room) (*entity.Room, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Cari Room berdasarkan ID
		var roomModel model.Room

		if err := tx.First(&roomModel, room.ID).Error; err != nil {
			return err
		}

		// SELECT * FROM rooms WHERE id = 1;

		// 2. Mapping Entity ke Model
		roomModel.Name = room.Name
		roomModel.Description = room.Description
		roomModel.Capacity = room.Capacity
		roomModel.Location = room.Location
		roomModel.IsActive = room.IsActive

		// 3. Update data Room
		if err := tx.Model(&roomModel).
			Select("Name", "Description", "Capacity", "Location", "IsActive").
			Updates(&roomModel).Error; err != nil {
			return err
		}

		// UPDATE rooms
		// SET
		//     name = 'Bali',
		//     description = 'Meeting room lantai 10',
		//     capacity = 20,
		//     location = 'Lantai 12',
		//     is_active = true,
		//     updated_at = NOW()
		// WHERE id = 1;

		// 4. Mapping Facility Entity ke Model
		var facilityModels []facilitymodel.Facility

		for _, facility := range room.Facilities {
			facilityModels = append(facilityModels, facilitymodel.Facility{
				ID: facility.ID,
			})
		}

		// 5. Replace relasi Facility
		if err := tx.Model(&roomModel).
			Association("Facilities").
			Replace(&facilityModels); err != nil {
			return err
		}

		// DELETE FROM room_facilities WHERE room_id = 1;
		// INSERT INTO room_facilities (room_id, facility_id)
		// VALUES
		// (1, 1),
		// (1, 3);

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 6. Ambil kembali Room beserta Facilities
	return r.FindByID(ctx, room.ID)
}

func (r *RoomRepositoryStruct) Delete(ctx context.Context, id int64) error {
	err := r.db.WithContext(ctx).Delete(&model.Room{}, id).Error

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return repository.ErrRoomHasBookings
		}
		// errors.As() memeriksa apakah error berasal dari PostgreSQL
		// pgErr.Code == "23503" adalah kode PostgreSQL untuk foreign key violation.

		return err
	}

	return nil
}
