package storage

import (
	"context"

	"sample-project/structs"
)

func (s *Storage) CreateDrone(d *structs.Drone) error {
	_, err := s.Pool.Exec(context.TODO(), `INSERT INTO drones
			(name,
			description,
			user_uuid,
			frame,
			motors,
			esc,
			propellers,
			fpv_camera,
			vtx,
			vtx_antenna,
			rx,
			flight_controller,
			extra_equipment) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		d.Name,
		d.Description,
		d.UserUUID,
		d.Frame,
		d.Motors,
		d.ESC,
		d.Propellers,
		d.FPVCamera,
		d.VTX,
		d.VTXAntenna,
		d.RX,
		d.FlightController,
		d.ExtraEquipment)
	return err
}

func (s *Storage) GetAllDrones() ([]structs.Drone, error) {
	rows, err := s.Pool.Query(context.TODO(), `SELECT d.id,
       d.name,
       d.description,
       u.name,
       d.frame,
       d.motors,
       d.esc,
       d.propellers,
       d.fpv_camera,
       d.vtx,
       d.vtx_antenna,
       d.rx,
       d.flight_controller,
       d.extra_equipment,
       d.created_at::varchar,
       d.updated_at::varchar
		FROM drones d
		LEFT JOIN users u ON u.uuid = d.user_uuid
		ORDER BY d.id`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	drones := []structs.Drone{}
	for rows.Next() {
		d := structs.Drone{}
		err = rows.Scan(
			&d.ID,
			&d.Name,
			&d.Description,
			&d.Author,
			&d.Frame,
			&d.Motors,
			&d.ESC,
			&d.Propellers,
			&d.FPVCamera,
			&d.VTX,
			&d.VTXAntenna,
			&d.RX,
			&d.FlightController,
			&d.ExtraEquipment,
			&d.CreatedAt,
			&d.UpdatedAt)
		if err != nil {
			return nil, err
		}

		drones = append(drones, d)
	}

	return drones, err
}
