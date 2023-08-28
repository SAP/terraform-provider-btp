package tfutils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

const btpcliTag = "btpcli"

type any interface{}
type equalityPredicate[E any] func(E, E) bool

func ToBTPCLIParamsMap(a any) (map[string]string, error) {
	out := map[string]string{}

	v, done, err := unwrapToStruct(a)
	if done {
		return out, err
	}

	for i := 0; i < v.NumField(); i++ {
		fieldProps := v.Type().Field(i)
		tagValue := fieldProps.Tag.Get(btpcliTag)
		if len(tagValue) == 0 {
			continue
		}

		field := v.FieldByName(fieldProps.Name)
		if !field.IsValid() {
			return nil, fmt.Errorf("invalid field")
		}

		switch fieldProps.Type.String() {
		case "string":
			setString(field, tagValue, out)
		case "basetypes.StringValue":
			setStringValue(field, tagValue, out)
		case "*string":
			setStringPointer(field, tagValue, out)
		case "bool":
			setBool(field, tagValue, out)
		case "basetypes.BoolValue":
			setBoolValue(field, tagValue, out)
		case "map[string][]string": // TODO would be nice to have `encodethisasjson` tag, instead of an explicit type mapping
			if !field.IsNil() {
				valueArr, err := json.Marshal(field.Interface())
				if err != nil {
					return nil, err
				}
				out[tagValue] = string(valueArr)
			}
		case "[]string":
			setStringSlice(field, tagValue, out)
		default:
			return nil, fmt.Errorf("the type '%s' assigned to '%s' is not yet supported", fieldProps.Type.String(), tagValue)
		}
	}

	return out, nil
}

func unwrapToStruct(a any) (reflect.Value, bool, error) {
	v := reflect.ValueOf(a)

	if !v.IsValid() {
		return reflect.Value{}, true, fmt.Errorf("invalid value")
	}

loop:
	for {
		switch v.Kind() {
		case reflect.Pointer:
			v = v.Elem()
		case reflect.Interface:
			if v.IsNil() {
				return reflect.Value{}, true, nil
			}
			v = v.Elem()
		default:
			break loop
		}
	}

	if v.Kind() != reflect.Struct || !v.IsValid() {
		return reflect.Value{}, true, fmt.Errorf("unsupported type: %s", v.Kind())
	}

	return v, false, nil
}

func setStringValue(field reflect.Value, tagValue string, out map[string]string) {
	fieldVal := field.Interface().(types.String)
	if !(fieldVal.IsUnknown() || fieldVal.IsNull()) {
		out[tagValue] = fieldVal.ValueString()
	}
}

func setString(field reflect.Value, tagValue string, out map[string]string) {
	fieldVal := field.Interface().(string)
	if fieldVal != "" {
		out[tagValue] = fieldVal
	}
}

func setStringPointer(field reflect.Value, tagValue string, out map[string]string) {
	if !field.IsNil() {
		out[tagValue] = field.Elem().Interface().(string)
	}
}

func setBoolValue(field reflect.Value, tagValue string, out map[string]string) {
	fieldVal := field.Interface().(types.Bool)
	if !(fieldVal.IsUnknown() || fieldVal.IsNull()) {
		out[tagValue] = fmt.Sprintf("%v", fieldVal.ValueBool())
	}
}

func setBool(field reflect.Value, tagValue string, out map[string]string) {
	fieldVal := field.Interface().(bool)
	out[tagValue] = fmt.Sprintf("%v", fieldVal)
}

func setStringSlice(field reflect.Value, tagValue string, out map[string]string) {
	if !field.IsNil() {
		valueString := fmt.Sprintf("%v", strings.Join(field.Interface().([]string), ","))
		out[tagValue] = valueString
	}
}

// TODO This is a utility function to compute to be removed and to be added substructures in resource configurations.
// TODO This is required since terraform only computes required CRUD operations on resource level. Changes in inner
// TODO configurations need to be computed based on the state and plan data by the update operation of a provider.
// TODO Should the terraform plugin framework support this functionality in the future, e.g. as a part of the Set
// TODO datatype, we can remove this code.
func SetDifference[S ~[]E, E any](setA, setB S, isEqual equalityPredicate[E]) (result S) {
	for _, element := range setA {
		if !setContains(setB, element, isEqual) {
			result = append(result, element)
		}
	}
	return
}

func setContains[S ~[]E, E any](set S, element E, isEqual equalityPredicate[E]) bool {
	for _, setElement := range set {
		if isEqual(setElement, element) {
			return true
		}
	}
	return false
}
