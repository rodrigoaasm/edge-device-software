package repositories

import (
	"database/sql"

	"opc.ua.agent/internal/domain/entities"
	domain_interfaces "opc.ua.agent/internal/domain/interfaces"
)

type DeviceRepository struct {
	con *sql.DB
}

func NewDeviceRepository(connection *sql.DB) *DeviceRepository {
	return &DeviceRepository{
		con: connection,
	}
}

func (r *DeviceRepository) Create(device entities.Device) (domain_interfaces.ITransactionManager, error) {
	tx, err := r.con.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("INSERT INTO device (device_id, device_ip) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(device.DeviceId, device.Ip)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return tx, nil
}

func (r *DeviceRepository) List() ([]entities.Device, error) {
	var devices []entities.Device
	rows, err := r.con.Query("SELECT device_id, device_ip, created_at, updated_at FROM device")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var device entities.Device
		if err := rows.Scan(&device.DeviceId, &device.Ip, &device.CreatedAt, &device.UpdatedAt); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func (r *DeviceRepository) Remove(deviceId string) (domain_interfaces.ITransactionManager, error) {
	tx, err := r.con.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("DELETE FROM device WHERE device_id = ?")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(deviceId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return tx, nil
}
