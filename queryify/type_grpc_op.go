package queryify

import (
	"git.adapticode.com/libraries/golang/errors"
	"git.adapticode.com/libraries/golang/sfqb"
	pb_common_filter_v1 "git.adapticode.com/platform/contracts/gen/go/common/filter/v1"
)

const (
	badOperator = "bad operator: "
)

var ErrUnknownOperator = errors.New("unknown operator")

type Operator interface {
	String() string
}

// MapOperator takes an operator of any type and maps it to a corresponding method.
// It supports various types of operators including boolean, string, integer, float,
// and array. If the operator type is not recognized, it returns an error.
func MapOperator(op Operator) (sfqb.Method, error) {
	switch o := op.(type) {
	case pb_common_filter_v1.BoolFieldFilter_Operator:
		return mapBoolOperator(o)
	case pb_common_filter_v1.StringFieldFilter_Operator:
		return mapStringOperator(o)
	case pb_common_filter_v1.StringNumberFieldFilter_Operator:
		return mapStringNumberOperator(o)
	case pb_common_filter_v1.IntFieldFilter_Operator:
		return mapIntOperator(o)
	case pb_common_filter_v1.FloatFieldFilter_Operator:
		return mapFloatOperator(o)
	case pb_common_filter_v1.ArrayStringFieldFilter_Operator:
		return mapArrayStringOperator(o)
	case pb_common_filter_v1.ArrayIntFieldFilter_Operator:
		return mapArrayIntOperator(o)
	case pb_common_filter_v1.ArrayFloatFieldFilter_Operator:
		return mapArrayFloatOperator(o)
	case pb_common_filter_v1.RangeIntFieldFilter_Operator:
		return mapRangeIntOperator(o)
	case pb_common_filter_v1.RangeStringFieldFilter_Operator:
		return mapRangeStringOperator(o)
	default:
		return "", ErrUnknownOperator
	}
}

// mapBoolOperator takes a boolean operator and maps it to a corresponding method.
// It supports unspecified and equality operators. If the operator is not recognized,
// it returns an error.
func mapBoolOperator(op pb_common_filter_v1.BoolFieldFilter_Operator) (sfqb.Method, error) {
	switch op {
	case pb_common_filter_v1.BoolFieldFilter_OPERATOR_EQ:
		return sfqb.EQ, nil
	case pb_common_filter_v1.BoolFieldFilter_OPERATOR_NEQ:
		return sfqb.NE, nil
	case pb_common_filter_v1.BoolFieldFilter_OPERATOR_UNSPECIFIED:
		fallthrough
	default:
		return "", errors.New(badOperator + op.String())
	}
}

// mapStringOperator takes a string operator and maps it to a corresponding method.
// It supports unspecified, equality, non-equality, and like operators. If the operator
// is not recognized, it returns an error.
func mapStringOperator(op pb_common_filter_v1.StringFieldFilter_Operator) (sfqb.Method, error) {
	switch op {
	case pb_common_filter_v1.StringFieldFilter_OPERATOR_EQ:
		return sfqb.EQ, nil
	case pb_common_filter_v1.StringFieldFilter_OPERATOR_NEQ:
		return sfqb.NE, nil
	case pb_common_filter_v1.StringFieldFilter_OPERATOR_LIKE:
		return sfqb.LIKE, nil
	case pb_common_filter_v1.StringFieldFilter_OPERATOR_ILIKE:
		return sfqb.ILIKE, nil
	case pb_common_filter_v1.StringFieldFilter_OPERATOR_UNSPECIFIED:
		fallthrough
	default:
		return "", errors.New(badOperator + op.String())
	}
}

// mapStringNumberOperator takes a string operator and maps it to a corresponding method.
// It supports unspecified, equality, non-equality,
// If the operator is not recognized, it returns an error.
func mapStringNumberOperator(op pb_common_filter_v1.StringNumberFieldFilter_Operator) (sfqb.Method, error) {
	switch op {
	case pb_common_filter_v1.StringNumberFieldFilter_OPERATOR_EQ:
		return sfqb.EQ, nil
	case pb_common_filter_v1.StringNumberFieldFilter_OPERATOR_NEQ:
		return sfqb.NE, nil
	case pb_common_filter_v1.StringNumberFieldFilter_OPERATOR_LT:
		return sfqb.LT, nil
	case pb_common_filter_v1.StringNumberFieldFilter_OPERATOR_LTE:
		return sfqb.LTE, nil
	case pb_common_filter_v1.StringNumberFieldFilter_OPERATOR_GT:
		return sfqb.GT, nil
	case pb_common_filter_v1.StringNumberFieldFilter_OPERATOR_GTE:
		return sfqb.GTE, nil
	case pb_common_filter_v1.StringNumberFieldFilter_OPERATOR_UNSPECIFIED:
		fallthrough
	default:
		return "", errors.New(badOperator + op.String())
	}
}

// mapIntOperator takes an integer operator and maps it to a corresponding method.
// It supports unspecified, equality, non-equality, less than, less than or equal,
// greater than, and greater than or equal operators. If the operator is not recognized,
// it returns an error.
func mapIntOperator(op pb_common_filter_v1.IntFieldFilter_Operator) (sfqb.Method, error) {
	switch op {
	case pb_common_filter_v1.IntFieldFilter_OPERATOR_EQ:
		return sfqb.EQ, nil
	case pb_common_filter_v1.IntFieldFilter_OPERATOR_NEQ:
		return sfqb.NE, nil
	case pb_common_filter_v1.IntFieldFilter_OPERATOR_LT:
		return sfqb.LT, nil
	case pb_common_filter_v1.IntFieldFilter_OPERATOR_LTE:
		return sfqb.LTE, nil
	case pb_common_filter_v1.IntFieldFilter_OPERATOR_GT:
		return sfqb.GT, nil
	case pb_common_filter_v1.IntFieldFilter_OPERATOR_GTE:
		return sfqb.GTE, nil
	case pb_common_filter_v1.IntFieldFilter_OPERATOR_UNSPECIFIED:
		fallthrough
	default:
		return "", errors.New(badOperator + op.String())
	}
}

// mapFloatOperator takes a float operator and maps it to a corresponding method.
// It supports unspecified, equality, non-equality, less than, less than or equal,
// greater than, and greater than or equal operators. If the operator is not recognized,
// it returns an error.
func mapFloatOperator(op pb_common_filter_v1.FloatFieldFilter_Operator) (sfqb.Method, error) {
	switch op {
	case pb_common_filter_v1.FloatFieldFilter_OPERATOR_EQ:
		return sfqb.EQ, nil
	case pb_common_filter_v1.FloatFieldFilter_OPERATOR_NEQ:
		return sfqb.NE, nil
	case pb_common_filter_v1.FloatFieldFilter_OPERATOR_LT:
		return sfqb.LT, nil
	case pb_common_filter_v1.FloatFieldFilter_OPERATOR_LTE:
		return sfqb.LTE, nil
	case pb_common_filter_v1.FloatFieldFilter_OPERATOR_GT:
		return sfqb.GT, nil
	case pb_common_filter_v1.FloatFieldFilter_OPERATOR_GTE:
		return sfqb.GTE, nil
	case pb_common_filter_v1.FloatFieldFilter_OPERATOR_UNSPECIFIED:
		fallthrough
	default:
		return "", errors.New(badOperator + op.String())
	}
}

// mapArrayStringOperator takes an array string operator and maps it to a corresponding method.
// It supports unspecified, in, and not in operators. If the operator is not recognized,
// it returns an error.
func mapArrayStringOperator(op pb_common_filter_v1.ArrayStringFieldFilter_Operator) (sfqb.Method, error) {
	switch op {
	case pb_common_filter_v1.ArrayStringFieldFilter_OPERATOR_IN:
		return sfqb.IN, nil
	case pb_common_filter_v1.ArrayStringFieldFilter_OPERATOR_NIN:
		return sfqb.NIN, nil
	case pb_common_filter_v1.ArrayStringFieldFilter_OPERATOR_UNSPECIFIED:
		fallthrough
	default:
		return "", errors.New(badOperator + op.String())
	}
}

// mapArrayIntOperator takes an array integer operator and maps it to a corresponding method.
// It supports unspecified, in, and not in operators. If the operator is not recognized,
// it returns an error.
func mapArrayIntOperator(op pb_common_filter_v1.ArrayIntFieldFilter_Operator) (sfqb.Method, error) {
	switch op {
	case pb_common_filter_v1.ArrayIntFieldFilter_OPERATOR_IN:
		return sfqb.IN, nil
	case pb_common_filter_v1.ArrayIntFieldFilter_OPERATOR_NIN:
		return sfqb.NIN, nil
	case pb_common_filter_v1.ArrayIntFieldFilter_OPERATOR_UNSPECIFIED:
		fallthrough
	default:
		return "", errors.New(badOperator + op.String())
	}
}

// mapArrayFloatOperator takes an array float operator and maps it to a corresponding method.
// It supports unspecified, in, and not in operators. If the operator is not recognized,
// it returns an error.
func mapArrayFloatOperator(op pb_common_filter_v1.ArrayFloatFieldFilter_Operator) (sfqb.Method, error) {
	switch op {
	case pb_common_filter_v1.ArrayFloatFieldFilter_OPERATOR_IN:
		return sfqb.IN, nil
	case pb_common_filter_v1.ArrayFloatFieldFilter_OPERATOR_NIN:
		return sfqb.NIN, nil
	case pb_common_filter_v1.ArrayFloatFieldFilter_OPERATOR_UNSPECIFIED:
		fallthrough
	default:
		return "", errors.New(badOperator + op.String())
	}
}

// mapRangeIntOperator takes a range integer operator and maps it to a corresponding method.
func mapRangeIntOperator(op pb_common_filter_v1.RangeIntFieldFilter_Operator) (sfqb.Method, error) {
	switch op {
	case pb_common_filter_v1.RangeIntFieldFilter_OPERATOR_RANGE:
		return sfqb.Range, nil
	case pb_common_filter_v1.RangeIntFieldFilter_OPERATOR_UNSPECIFIED:
		fallthrough
	default:
		return "", errors.New(badOperator + op.String())
	}
}

func mapRangeStringOperator(op pb_common_filter_v1.RangeStringFieldFilter_Operator) (sfqb.Method, error) {
	switch op {
	case pb_common_filter_v1.RangeStringFieldFilter_OPERATOR_RANGE:
		return sfqb.Range, nil
	case pb_common_filter_v1.RangeStringFieldFilter_OPERATOR_UNSPECIFIED:
		fallthrough
	default:
		return "", errors.New(badOperator + op.String())
	}
}
