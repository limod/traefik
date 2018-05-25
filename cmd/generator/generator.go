package generator

import (
	"fmt"
	"reflect"
)

const theAnswer = 42

type Hydration struct {
	ExcludedFieldNames  []string
	SlideItemNumber     int
	MapItemNumber       int
	DefaultString       string
	DefaultNumber       int
	DefaultBool         bool
	DefaultMapKeyPrefix string
}

func NewHydration() Hydration {
	return Hydration{
		SlideItemNumber:     2,
		MapItemNumber:       1,
		DefaultString:       "foobar",
		DefaultNumber:       theAnswer,
		DefaultBool:         true,
		DefaultMapKeyPrefix: "name",
	}
}

func (h Hydration) Fill(element interface{}) {
	field := reflect.ValueOf(element)
	name := reflect.TypeOf(element).Name()
	h.fill(name, field)
}

func (h Hydration) fill(previous string, field reflect.Value) error {
	//fmt.Println(previous)
	switch field.Kind() {
	case reflect.Struct:
		err := h.setStruct(previous, field)
		if err != nil {
			return err
		}
	case reflect.Ptr:
		err := h.setPointer(previous, field)
		if err != nil {
			return err
		}
	case reflect.Slice:
		err := h.setSlice(field)
		if err != nil {
			return err
		}
	case reflect.Map:
		err := h.setMap(field)
		if err != nil {
			return err
		}
	case reflect.Interface:
		if err := h.fill(previous, field.Elem()); err != nil {
			return err
		}
	case reflect.String:
		h.setTyped(field, h.DefaultString)
	case reflect.Int:
		h.setTyped(field, h.DefaultNumber)
	case reflect.Int8:
		h.setTyped(field, int8(h.DefaultNumber))
	case reflect.Int16:
		h.setTyped(field, int16(h.DefaultNumber))
	case reflect.Int32:
		h.setTyped(field, int32(h.DefaultNumber))
	case reflect.Int64:
		h.setTyped(field, int64(h.DefaultNumber))
	case reflect.Uint:
		h.setTyped(field, uint(h.DefaultNumber))
	case reflect.Uint8:
		h.setTyped(field, uint8(h.DefaultNumber))
	case reflect.Uint16:
		h.setTyped(field, uint16(h.DefaultNumber))
	case reflect.Uint32:
		h.setTyped(field, uint32(h.DefaultNumber))
	case reflect.Uint64:
		h.setTyped(field, uint64(h.DefaultNumber))
	case reflect.Bool:
		h.setTyped(field, h.DefaultBool)
	case reflect.Float32:
		h.setTyped(field, float32(h.DefaultNumber))
	case reflect.Float64:
		h.setTyped(field, float64(h.DefaultNumber))
	}

	//g := complex(float32(h.DefaultNumber),float32(h.DefaultNumber))

	// TODO ?
	//Uintptr
	//Complex64
	//Complex128
	//Array
	//Chan
	//Func
	//UnsafePointer
	return nil
}

func (h Hydration) setTyped(field reflect.Value, i interface{}) {
	baseValue := reflect.ValueOf(i)
	if field.Kind().String() == field.Type().String() {
		field.Set(baseValue)
	} else {
		field.Set(baseValue.Convert(field.Type()))
	}
}

func (h Hydration) setMap(field reflect.Value) error {
	field.Set(reflect.MakeMap(field.Type()))

	for i := 0; i < h.MapItemNumber; i++ {
		// TODO support only string... must be fixed
		//fmt.Println(field.Type().Key())
		baseKeyName := h.makeKeyName(field.Type().Elem())
		key := reflect.ValueOf(fmt.Sprintf("%s%d", baseKeyName, i))

		// generate value
		ptrType := reflect.PtrTo(field.Type().Elem())
		ptrValue := reflect.New(ptrType)
		if err := h.fill("", ptrValue); err != nil {
			return err
		}
		value := ptrValue.Elem().Elem()

		field.SetMapIndex(key, value)
	}
	return nil
}

func (h Hydration) makeKeyName(typ reflect.Type) string {
	switch typ.Kind() {
	case reflect.Ptr:
		return typ.Elem().Name()
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Bool, reflect.Float32, reflect.Float64:
		return h.DefaultMapKeyPrefix
	default:
		return typ.Name()
	}
}

func (h Hydration) setStruct(previous string, field reflect.Value) error {
	for i := 0; i < field.NumField(); i++ {
		fld := field.Field(i)
		stFld := field.Type().Field(i)

		if !isExported(stFld) || h.isExcluded(previous, stFld) || fld.Kind() == reflect.Func {
			continue
		}

		name := stFld.Name
		if len(previous) > 0 {
			name = previous + "." + stFld.Name
		}
		if h.isExcluded(name, stFld) {
			continue
		}

		if err := h.fill(name, fld); err != nil {
			return err
		}
	}
	return nil
}

func (h Hydration) setSlice(field reflect.Value) error {
	field.Set(reflect.MakeSlice(field.Type(), h.SlideItemNumber, h.SlideItemNumber))
	for j := 0; j < field.Len(); j++ {
		if err := h.fill("", field.Index(j)); err != nil {
			return err
		}
	}
	return nil
}

func (h Hydration) setPointer(previous string, field reflect.Value) error {
	if field.IsNil() {
		field.Set(reflect.New(field.Type().Elem()))
		if err := h.fill(previous, field.Elem()); err != nil {
			return err
		}
	} else {
		if err := h.fill(previous, field.Elem()); err != nil {
			return err
		}
	}
	return nil
}

// isExported return true is a struct field is exported, else false
func isExported(f reflect.StructField) bool {
	if f.PkgPath != "" && !f.Anonymous {
		return false
	}
	return true
}

func (h Hydration) isExcluded(fullName string, field reflect.StructField) bool {
	for _, name := range h.ExcludedFieldNames {
		if field.Name == name || fullName == name {
			return true
		}
	}
	return false
}
