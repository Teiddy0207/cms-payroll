package mapper

import (
	"cal-salary/modules/payroll/dto"
	"testing"
)

func TestToDepartmentEntity(t *testing.T) {
	desc := "Phòng Hành chính Nhân sự"
	req := &dto.CreateDepartmentRequest{
		Code:        "HR",
		Name:        "Human Resources",
		Description: &desc,
	}

	entity := ToDepartmentEntity(req)

	if entity.Code != req.Code {
		t.Errorf("Expected Code %s, got %s", req.Code, entity.Code)
	}
	if entity.Name != req.Name {
		t.Errorf("Expected Name %s, got %s", req.Name, entity.Name)
	}
	if entity.Description == nil || *entity.Description != *req.Description {
		t.Errorf("Expected Description %s, got %v", *req.Description, entity.Description)
	}
}
