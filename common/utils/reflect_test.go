package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestGames struct {
	CodeInt   int
	CodeInt8  int8
	CodeInt16 int16
	CodeInt32 int32
	CodeInt64 int64

	CodeUint   uint
	CodeUint8  uint8
	CodeUint16 uint16
	CodeUint32 uint32
	CodeUint64 uint64

	CodeFloat32 float32
	CodeFloat64 float64

	CodeBool bool

	CodeType      string
	CodeGameBrand []TestGameBrand

	TestMapInterface []map[string]interface{}
	TestMapStr       []map[string]string
}

type TestGameBrand struct {
	CodeStr      string
	CodeStrSlice []string
	CodeStrArray [3]string
}

type Games struct {
	CodeInt   int
	CodeInt8  int8
	CodeInt16 int16
	CodeInt32 int32
	CodeInt64 int64

	CodeUint   uint
	CodeUint8  uint8
	CodeUint16 uint16
	CodeUint32 uint32
	CodeUint64 uint64

	CodeFloat32 float32
	CodeFloat64 float64

	CodeBool bool

	CodeType      string
	CodeGameBrand []struct {
		CodeStr      string
		CodeStrSlice []string
		CodeStrArray [3]string
	}

	TestMapInterface []map[string]interface{}
	TestMapStr       []map[string]string
}

func TestConvSourceToData(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "測試 PromotionsContainGames conv Games",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourData := []TestGames{
				{
					CodeInt:   -123,
					CodeInt8:  -1,
					CodeInt16: -1,
					CodeInt32: -1,
					CodeInt64: -1,

					CodeUint:   456,
					CodeUint8:  4,
					CodeUint16: 4,
					CodeUint32: 4,
					CodeUint64: 4,

					CodeFloat32: 1.123,
					CodeFloat64: 1.123,

					CodeBool: true,

					CodeType: "live",
					CodeGameBrand: []TestGameBrand{
						{
							CodeStr: "ebet",
							CodeStrSlice: []string{
								"ebet_1",
								"ebet_2",
							},
							CodeStrArray: [3]string{
								"1", "2", "3",
							},
						},
					},

					TestMapInterface: []map[string]interface{}{
						{
							"1": 1,
							"2": true,
						},
						{
							"1": true,
							"3": 00.0,
						},
						{
							"1": 00.0,
							"4": "a",
						},
						{
							"1": "a",
						},
					},
					TestMapStr: []map[string]string{
						{
							"1": "1",
							"2": "2",
						},
						{
							"1": "1",
							"2": "2",
						},
					},
				},
				{
					CodeType:      "sport",
					CodeGameBrand: []TestGameBrand{},
				},
			}

			var resultData []Games

			if err := ConvSourceToData(&sourData, &resultData); (err != nil) != tt.wantErr {
				assert.NoError(t, err, "ConvSourceToData error = %v", err)
				return
			}

			assert.Equal(t, fmt.Sprintf("%v", sourData), fmt.Sprintf("%v", resultData))
		})
	}
}
