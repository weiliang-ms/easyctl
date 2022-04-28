package value

import (
	"fmt"
	"reflect"
)

// SetStructDefaultValue Set Struct Default value
/*
	Sets default values for structure objects.

	1.Get object types and values by reflection.
	2.
*/
func SetStructDefaultValue(obj interface{}, key string, value interface{}) error {

	dataType := reflect.TypeOf(obj)
	dataValue := reflect.ValueOf(obj)
	defaultValueType := reflect.ValueOf(value)

	if dataType.Kind() != reflect.Ptr {
		return fmt.Errorf("必需为指针类型")
	}

	if reflect.ValueOf(obj).Elem().Type().Kind() != reflect.Struct {
		return fmt.Errorf("必须为结构体类型")
	}

	dataValue = dataValue.Elem()
	dataType = dataType.Elem()

	for i := 0; i < dataType.NumField(); i++ {

		field := dataType.Field(i)
		fieldName := field.Name
		fieldValue := dataValue.FieldByName(fieldName)

		// todo: confirm this case
		//if !fieldValue.IsValid() {
		//	continue
		//}
		// todo: 默认值，key值必须为空，非空不赋值
		if fieldValue.CanInterface() && fieldValue.CanSet() && fieldValue.String() == "" {
			//fmt.Printf("exported fieldName:%v value:%v\n", fieldName, fieldValue.Interface())
			if defaultValueType.Kind() == fieldValue.Kind() && key == fieldName {
				switch defaultValueType.Kind() {
				case reflect.String:
					fieldValue.SetString(fmt.Sprintf("%s", value))
				case reflect.Int:
					v, _ := value.(int)
					fieldValue.SetInt(int64(v))
				case reflect.Int32:
					v, _ := value.(int32)
					fieldValue.SetInt(int64(v))
				case reflect.Int64:
					v, _ := value.(int64)
					fieldValue.SetInt(v)
				case reflect.Bool:
					v, _ := value.(bool)
					fieldValue.SetBool(v)
				}
			}

		}
		// todo: confirm this case
		//else {
		//	// 强行取址
		//	forceValue := reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()
		//	fmt.Printf("unexported fieldName:%v value:%v\n", fieldName, forceValue.Interface())
		//}

	}

	return nil
}
