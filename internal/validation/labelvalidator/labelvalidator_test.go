package labelvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestLabelValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        types.Map
		expErrors int
	}

	var validlabels = map[string][]string{
		"costcenter":  {"12345"},
		"environment": {"production"},
		"features":    {"featureset1", "featureset2"},
	}

	validLabelsTypesMap, _ := types.MapValueFrom(context.TODO(), types.SetType{ElemType: types.StringType}, validlabels)

	var invalidLabels = map[string][]int{
		"costcenter":  {12345},
	}

	invalidLabelsTypesMap, _ := types.MapValueFrom(context.TODO(), types.SetType{ElemType: types.Int32Type}, invalidLabels)

	var invalidLabels_LabelType = map[string]string{
		"costcenter":  "12345",
	}

	invalidLabels_LabelTypeTypesMap, _ := types.MapValueFrom(context.TODO(), types.StringType, invalidLabels_LabelType)

	var invalidLabelsLength = map[string][]string{
		"costcenter":  {"12345"},
		"environment": {"production"},
		"features":    {"featureset1, featureset2, featureset3, featureset4, featureset5, featureset6"},
	}

	invalidLabelsLengthTypesMap, _ := types.MapValueFrom(context.TODO(), types.SetType{ElemType: types.StringType}, invalidLabelsLength)

	var invalidLabelsNumber = map[string][]string{
		"costcenter":   {"12345"},
		"environment":  {"production"},
		"features":     {"featureset1", "featureset2"},
		"costcenter2":  {"12345"},
		"environment2": {"production"},
		"features2":    {"featureset1", "featureset2"},
		"costcenter3":  {"12345"},
		"environment3": {"production"},
		"features3":    {"featureset1", "featureset2"},
		"costcenter4":  {"12345"},
		"environment4": {"production"},
		"features4":    {"featureset1", "featureset2"},
		"costcenter5":  {"12345"},
		"environment5": {"production"},
		"features5":    {"featureset1", "featureset2"},
	}

	invalidLabelsNumberTypesMap, _ := types.MapValueFrom(context.TODO(), types.SetType{ElemType: types.StringType}, invalidLabelsNumber)

	testCases := map[string]testCase{
		"valid-labels": {
			in:        validLabelsTypesMap,
			expErrors: 0,
		},
		"invalid-labels": {
			in: 	  invalidLabelsTypesMap,
			expErrors: 1,
		},
		"invalid-labels-lableType": {
			in: 	  invalidLabels_LabelTypeTypesMap,
			expErrors: 1,
		},
		"too-many-labels": {
			in:        invalidLabelsNumberTypesMap,
			expErrors: 1,
		},
		"too-long-label": {
			in:        invalidLabelsLengthTypesMap,
			expErrors: 1,
		},
	}

	for name, test := range testCases {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			req := validator.MapRequest{
				ConfigValue: test.in,
			}
			res := validator.MapResponse{}
			ValidLabels().ValidateMap(context.TODO(), req, &res)

			if test.expErrors > 0 && !res.Diagnostics.HasError() {
				t.Fatalf("expected %d error(s), got none", test.expErrors)
			}

			if test.expErrors > 0 && test.expErrors != res.Diagnostics.ErrorsCount() {
				t.Fatalf("expected %d error(s), got %d: %v", test.expErrors, res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}

			if test.expErrors == 0 && res.Diagnostics.HasError() {
				t.Fatalf("expected no error(s), got %d: %v", res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}
		})
	}
}
