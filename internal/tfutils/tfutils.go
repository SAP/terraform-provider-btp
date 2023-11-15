package tfutils

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

const btpcliTag = "btpcli"

const DefaultTimeout = 10 * time.Minute

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

		cliParameter, encoder, found := parseCliTag(fieldProps.Tag)

		if !found {
			continue
		}

		field := v.FieldByName(fieldProps.Name)

		if !field.IsValid() {
			return nil, fmt.Errorf("invalid field")
		}

		fieldValue, err := encoder.Encode(fieldProps.Type, field)

		if err != nil {
			return nil, fmt.Errorf("unable to encode '%s': %w", cliParameter, err)
		}

		if len(fieldValue) == 0 {
			continue
		}

		out[cliParameter] = fieldValue
	}

	return out, nil
}

func parseCliTag(tag reflect.StructTag) (cliParam string, encoder paramsEncoder, found bool) {
	tagValue := strings.Split(tag.Get(btpcliTag), ",")

	found = len(tagValue) > 0 && len(tagValue[0]) > 0

	if found {
		cliParam = tagValue[0]
	}

	if len(tagValue) > 1 && tagValue[1] == "json" {
		encoder = &jsonEncoder{}
	} else {
		encoder = &autoEncoder{}
	}

	return
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

func CalculateDelayAndMinTimeOut(timeout time.Duration) (delay time.Duration, minTimeout time.Duration) {
	// We define the polling interval as 1/100 of the timeout value in seconds
	// For 10 minutes this results in 6 seconds polling interval
	// For 1 hour this results in 36 seconds polling interval
	delay = time.Duration(math.Round(timeout.Seconds()/100)) * time.Second

	// We set the minTimeout equal to the polling interval
	minTimeout = delay
	return
}

type paramsEncoder interface {
	Encode(fieldType reflect.Type, field reflect.Value) (string, error)
}

type autoEncoder struct {
}

func (ae *autoEncoder) Encode(fieldType reflect.Type, field reflect.Value) (string, error) {
	switch fieldType.String() {
	case "int":
		return ae.encodeInt(field)
	case "string":
		return ae.encodeString(field)
	case "basetypes.StringValue":
		return ae.encodeStringValue(field)
	case "*string":
		return ae.encodeStringPointer(field)
	case "bool":
		return ae.encodeBool(field)
	case "basetypes.BoolValue":
		return ae.encodeBoolValue(field)
	case "*bool":
		return ae.encodeBoolPointer(field)
	case "map[string][]string":
		if field.IsNil() {
			return "", nil
		}

		valueArr, err := json.Marshal(field.Interface())
		if err != nil {
			return "", err
		}
		return string(valueArr), nil
	case "[]string":
		return ae.encodeStringSlice(field)
	default:
		return "", fmt.Errorf("unsupported type '%s'", fieldType)
	}
}

func (ae *autoEncoder) encodeInt(field reflect.Value) (string, error) {
	return fmt.Sprintf("%v", field.Interface().(int)), nil
}

func (ae *autoEncoder) encodeStringValue(field reflect.Value) (string, error) {
	fieldVal := field.Interface().(types.String)
	if fieldVal.IsUnknown() || fieldVal.IsNull() {
		return "", nil
	}

	return fieldVal.ValueString(), nil
}

func (ae *autoEncoder) encodeString(field reflect.Value) (string, error) {
	return field.Interface().(string), nil
}

func (ae *autoEncoder) encodeStringPointer(field reflect.Value) (string, error) {
	if field.IsNil() {
		return "", nil
	}

	return field.Elem().Interface().(string), nil
}

func (ae *autoEncoder) encodeBoolValue(field reflect.Value) (string, error) {
	fieldVal := field.Interface().(types.Bool)
	if fieldVal.IsUnknown() || fieldVal.IsNull() {
		return "", nil
	}

	return fmt.Sprintf("%v", fieldVal.ValueBool()), nil
}

func (ae *autoEncoder) encodeBool(field reflect.Value) (string, error) {
	return fmt.Sprintf("%v", field.Interface().(bool)), nil
}

func (ae *autoEncoder) encodeBoolPointer(field reflect.Value) (string, error) {
	if field.IsNil() {
		return "", nil
	}

	return fmt.Sprintf("%v", field.Elem().Interface().(bool)), nil
}

func (ae *autoEncoder) encodeStringSlice(field reflect.Value) (string, error) {
	if field.IsNil() {
		return "", nil
	}

	return fmt.Sprintf("%v", strings.Join(field.Interface().([]string), ",")), nil
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

type jsonEncoder struct {
}

func (je *jsonEncoder) Encode(_ reflect.Type, field reflect.Value) (string, error) {
	arr, err := json.Marshal(field.Interface())

	return string(arr), err
}
