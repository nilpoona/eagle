package eagle

import (
	"net/url"
	"reflect"
	"testing"
)

func TestBindFormData(t *testing.T) {
	na := 1
	na8 := int8(na)
	na16 := int16(na)
	na32 := int32(na)
	na64 := int64(na)
	ui := uint(na)
	ui8 := uint8(na)
	ui16 := uint16(na)
	ui32 := uint32(na)
	ui64 := uint64(na)
	f32 := float32(na)
	// f64 := float64(na)
	name := "foo"

	type args struct {
		formData url.Values
		v        interface{}
	}

	type formData struct {
		Name            string   `form:"name"`
		NullableName    *string  `form:"nullablename"`
		Age             int      `form:"age"`
		Activate        bool     `form:"activate"`
		Rate            float64  `form:"rate"`
		Number          uint     `form:"number"`
		NullableAge     *int     `form:"nullableage"`
		Age8            int8     `form:"age8"`
		NullableAge8    *int8    `form:"nullableage8"`
		Age16           int16    `form:"age16"`
		NullableAge16   *int16   `form:"nullableage16"`
		Age32           int32    `form:"age32"`
		NullableAge32   *int32   `form:"nullableage32"`
		Age64           int64    `form:"age64"`
		NullableAge64   *int64   `form:"nullableage64"`
		Uint            uint     `form:"uint"`
		NullableUint    *uint    `form:"nullableuint"`
		Uint8           uint8    `form:"uint8"`
		NullableUint8   *uint8   `form:"nullableuint8"`
		Uint16          uint16   `form:"uint16"`
		NullableUint16  *uint16  `form:"nullableuint16"`
		Uint32          uint32   `form:"uint32"`
		NullableUint32  *uint32  `form:"nullableuint32"`
		Uint64          uint64   `form:"uint64"`
		NullableUint64  *uint64  `form:"nullableuint64"`
		Float32         float32  `form:"float32"`
		NullableFloat32 *float32 `form:"nullablefloat32"`
		Float64         float64  `form:"float64"`
		NullableFloat64 *float64 `form:"nullablefloat64"`
	}

	tests := []struct {
		name string
		args args
		want *formData
		err  error
	}{
		{
			name: "",
			args: args{
				formData: url.Values{
					"name": []string{"foo"},
					"age":  []string{"22"},
				},
				v: &formData{},
			},
			want: &formData{
				Name: "foo",
				Age:  22,
			},
			err: nil,
		},
		{
			name: "",
			args: args{
				formData: url.Values{
					"name":     []string{"foo"},
					"age":      []string{"22"},
					"activate": []string{"true"},
				},
				v: &formData{},
			},
			want: &formData{
				Name:     "foo",
				Age:      22,
				Activate: true,
			},
			err: nil,
		},
		{
			name: "",
			args: args{
				formData: url.Values{
					"name":     []string{"foo"},
					"age":      []string{"22"},
					"activate": []string{"true"},
					"rate":     []string{"1.2"},
				},
				v: &formData{},
			},
			want: &formData{
				Name:     "foo",
				Age:      22,
				Activate: true,
				Rate:     1.2,
			},
			err: nil,
		},
		{
			name: "",
			args: args{
				formData: url.Values{
					"name":     []string{"foo"},
					"age":      []string{"22"},
					"activate": []string{"true"},
					"rate":     []string{"1.2"},
					"number":   []string{"1"},
				},
				v: &formData{},
			},
			want: &formData{
				Name:     "foo",
				Age:      22,
				Activate: true,
				Rate:     1.2,
				Number:   1,
			},
			err: nil,
		},
		{
			name: "",
			args: args{
				formData: url.Values{
					"name":            []string{name},
					"nullablename":    []string{name},
					"age":             []string{"22"},
					"activate":        []string{"true"},
					"rate":            []string{"1.2"},
					"number":          []string{"1"},
					"nullableage":     []string{"1"},
					"age8":            []string{"1"},
					"nullableage8":    []string{"1"},
					"age16":           []string{"1"},
					"nullableage16":   []string{"1"},
					"age32":           []string{"1"},
					"nullableage32":   []string{"1"},
					"age64":           []string{"1"},
					"nullableage64":   []string{"1"},
					"uint":            []string{"1"},
					"nullableuint":    []string{"1"},
					"uint8":           []string{"1"},
					"nullableuint8":   []string{"1"},
					"uint16":          []string{"1"},
					"nullableuint16":  []string{"1"},
					"uint32":          []string{"1"},
					"nullableuint32":  []string{"1"},
					"uint64":          []string{"1"},
					"nullableuint64":  []string{"1"},
					"float32":         []string{"1"},
					"nullablefloat32": []string{"1"},
					/*
						"float64":         []string{"1"},
						"nullablefloat64": []string{"1"},
					*/
				},

				v: &formData{},
			},
			want: &formData{
				Name:            name,
				NullableName:    &name,
				Age:             22,
				Activate:        true,
				Rate:            1.2,
				Number:          1,
				NullableAge:     &na,
				Age8:            1,
				NullableAge8:    &na8,
				Age16:           1,
				NullableAge16:   &na16,
				Age32:           1,
				NullableAge32:   &na32,
				Age64:           1,
				NullableAge64:   &na64,
				Uint:            ui,
				NullableUint:    &ui,
				Uint8:           ui8,
				NullableUint8:   &ui8,
				Uint16:          ui16,
				NullableUint16:  &ui16,
				Uint32:          ui32,
				NullableUint32:  &ui32,
				Uint64:          ui64,
				NullableUint64:  &ui64,
				Float32:         f32,
				NullableFloat32: &f32,
				/*
					Float64:         f64,
					NullableFloat64: &f64,
				*/
			},
			err: nil,
		},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			err := bindFormData(td.args.formData, td.args.v)
			if err != td.err {
				t.Errorf("bindFormData failed err: %s expected error: %s", err, td.err)
			}

			if !reflect.DeepEqual(td.args.v, td.want) {
				t.Errorf("bindFormData failed err result: %+v expected: %+v", td.args.v, td.want)
			}
		})
	}
}
