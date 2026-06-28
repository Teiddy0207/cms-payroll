package mapper

import (
	"cal-salary/core/logger"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func PQArrayPtrToUUIDSlicePtr(arr *pq.StringArray) *[]uuid.UUID {
	if arr == nil || len(*arr) == 0 {
		return nil
	}

	result := make([]uuid.UUID, 0, len(*arr))
	for _, s := range *arr {
		if id, err := uuid.Parse(s); err == nil {
			result = append(result, id)
		} else {
			logger.Warn("Invalid UUID in array:", s)
		}
	}

	if len(result) == 0 {
		return nil
	}

	return &result
}

func UUIDSlicePtrToPQArrayPtr(slice *[]uuid.UUID) *pq.StringArray {
	if slice == nil || len(*slice) == 0 {
		return nil
	}

	result := make([]string, len(*slice))
	for i, id := range *slice {
		result[i] = id.String()
	}

	arr := pq.StringArray(result)
	return &arr
}

func ApplyPatchPtr[T any](patch **T, existing *T) *T {
	if patch == nil {
		return existing
	}
	return *patch
}

type Patch[T any] struct {
	Set   bool // field có xuất hiện trong JSON không
	Value *T   // nil => JSON null, != nil => có giá trị
}

func (p *Patch[T]) UnmarshalJSON(b []byte) error {
	p.Set = true

	// null
	if string(b) == "null" {
		p.Value = nil
		return nil
	}

	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	p.Value = &v
	return nil
}

func ApplyPatch[T any](patch Patch[T], existing *T) *T {
	if !patch.Set {
		return existing
	}
	return patch.Value // nil => clear, != nil => set
}
