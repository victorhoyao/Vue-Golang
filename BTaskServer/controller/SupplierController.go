package controller

import (
	"BTaskServer/dao"
	"BTaskServer/model"
	"BTaskServer/util/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SupplierController handles supplier-related requests
type SupplierController struct {
}

// NewSupplierController creates and returns a new SupplierController
func NewSupplierController() *SupplierController {
	return &SupplierController{}
}

// AddSupplier handles the creation of a new supplier
func (s *SupplierController) AddSupplier(c *gin.Context) {
	var form model.SupplierCreateForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Fail(c, nil, "Invalid parameters")
		return
	}

	supplier := &model.Supplier{
		SupplierID:    form.SupplierID,
		SupplierName:  form.SupplierName,
		Contact:       form.Contact,
		Status:        form.Status,
		TasksAssigned: form.TasksAssigned,
	}

	err := dao.AddSupplier(supplier)
	if err != nil {
		response.ServerBad(c, nil, "Failed to add supplier")
		return
	}

	response.Success(c, nil, "Supplier added successfully")
}

// GetSuppliers retrieves a list of suppliers
func (s *SupplierController) GetSuppliers(c *gin.Context) {
	suppliers, err := dao.GetSuppliers()
	if err != nil {
		response.ServerBad(c, nil, "Failed to retrieve suppliers")
		return
	}

	response.Success(c, gin.H{"suppliers": suppliers}, "Suppliers retrieved successfully")
}

// UpdateSupplier handles the update of an existing supplier
func (s *SupplierController) UpdateSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, nil, "Invalid supplier ID")
		return
	}

	var form model.SupplierUpdateForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Fail(c, nil, "Invalid parameters")
		return
	}

	supplier := &model.Supplier{
		ID:            uint(id),
		SupplierID:    form.SupplierID,
		SupplierName:  form.SupplierName,
		Contact:       form.Contact,
		Status:        form.Status,
		TasksAssigned: form.TasksAssigned,
	}

	err = dao.UpdateSupplier(supplier)
	if err != nil {
		response.ServerBad(c, nil, "Failed to update supplier")
		return
	}

	response.Success(c, nil, "Supplier updated successfully")
}

// DeleteSupplier handles the deletion of a supplier
func (s *SupplierController) DeleteSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Fail(c, nil, "Invalid supplier ID")
		return
	}

	err = dao.DeleteSupplier(uint(id))
	if err != nil {
		response.ServerBad(c, nil, "Failed to delete supplier")
		return
	}

	response.Success(c, nil, "Supplier deleted successfully")
}
