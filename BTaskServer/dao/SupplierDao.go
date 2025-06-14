package dao

import (
	"BTaskServer/global"
	"BTaskServer/model"
)

// AddSupplier creates a new supplier record in the database
func AddSupplier(supplier *model.Supplier) error {
	result := global.GVA_DB.Create(supplier)
	return result.Error
}

// GetSuppliers retrieves all supplier records from the database
func GetSuppliers() ([]model.Supplier, error) {
	var suppliers []model.Supplier
	result := global.GVA_DB.Find(&suppliers)
	return suppliers, result.Error
}

// UpdateSupplier updates an existing supplier record in the database
func UpdateSupplier(supplier *model.Supplier) error {
	result := global.GVA_DB.Save(supplier)
	return result.Error
}

// DeleteSupplier deletes a supplier record from the database by ID
func DeleteSupplier(id uint) error {
	result := global.GVA_DB.Delete(&model.Supplier{}, id)
	return result.Error
}
