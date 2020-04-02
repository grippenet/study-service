package studyengine

import (
	"testing"

	"github.com/influenzanet/study-service/models"
)

// Reference/Lookup methods
func TestEvalGetEventType(t *testing.T) {
	exp := models.Expression{Name: "getEventType"}

	t.Run("for enter", func(t *testing.T) {
		evalContext := evalContext{
			event: models.StudyEvent{Type: "enter"},
		}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(string) != "enter" {
			t.Errorf("unexpected type or str value: %s", ret)
		}
	})
}

func TestEvalGetParticipantState(t *testing.T) {
	t.Run("with normal state", func(t *testing.T) { t.Error("test unimplemented") })
}

// Comparisons
func TestEvalEq(t *testing.T) {
	t.Run("for eq numbers", func(t *testing.T) {
		exp := models.Expression{Name: "eq", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 23},
			models.ExpressionArg{DType: "num", Num: 23},
		}}
		evalContext := evalContext{
			event: models.StudyEvent{Type: "TIMER"},
		}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("for not equal numbers", func(t *testing.T) {
		exp := models.Expression{Name: "eq", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 13},
			models.ExpressionArg{DType: "num", Num: 23},
		}}
		evalContext := evalContext{
			event: models.StudyEvent{Type: "enter"},
		}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("for equal strings", func(t *testing.T) {
		exp := models.Expression{Name: "eq", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "exp", Exp: models.Expression{Name: "getEventType"}},
			models.ExpressionArg{DType: "str", Str: "enter"},
		}}
		evalContext := evalContext{
			event: models.StudyEvent{Type: "enter"},
		}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("for not equal strings", func(t *testing.T) {
		exp := models.Expression{Name: "eq", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "exp", Exp: models.Expression{Name: "getEventType"}},
			models.ExpressionArg{DType: "str", Str: "time..."},
		}}
		evalContext := evalContext{
			event: models.StudyEvent{Type: "enter"},
		}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

func TestEvalLT(t *testing.T) {
	t.Run("2 < 2", func(t *testing.T) {
		exp := models.Expression{Name: "lt", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 2},
			models.ExpressionArg{DType: "num", Num: 2},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("2 < 1", func(t *testing.T) {
		exp := models.Expression{Name: "lt", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 2},
			models.ExpressionArg{DType: "num", Num: 1},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 < 2", func(t *testing.T) {
		exp := models.Expression{Name: "lt", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 1},
			models.ExpressionArg{DType: "num", Num: 2},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("a < b", func(t *testing.T) {
		exp := models.Expression{Name: "lt", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "str", Str: "a"},
			models.ExpressionArg{DType: "str", Str: "b"},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b < b", func(t *testing.T) {
		exp := models.Expression{Name: "lt", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "str", Str: "b"},
			models.ExpressionArg{DType: "str", Str: "b"},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b < a", func(t *testing.T) {
		exp := models.Expression{Name: "lt", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "str", Str: "b"},
			models.ExpressionArg{DType: "str", Str: "a"},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

func TestEvalLTE(t *testing.T) {
	t.Run("2 <= 2", func(t *testing.T) {
		exp := models.Expression{Name: "lte", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 2},
			models.ExpressionArg{DType: "num", Num: 2},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("2 <= 1", func(t *testing.T) {
		exp := models.Expression{Name: "lte", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 2},
			models.ExpressionArg{DType: "num", Num: 1},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 <= 2", func(t *testing.T) {
		exp := models.Expression{Name: "lte", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 1},
			models.ExpressionArg{DType: "num", Num: 2},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("a <= b", func(t *testing.T) {
		exp := models.Expression{Name: "lte", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "str", Str: "a"},
			models.ExpressionArg{DType: "str", Str: "b"},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b <= b", func(t *testing.T) {
		exp := models.Expression{Name: "lte", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "str", Str: "b"},
			models.ExpressionArg{DType: "str", Str: "b"},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b <= a", func(t *testing.T) {
		exp := models.Expression{Name: "lte", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "str", Str: "b"},
			models.ExpressionArg{DType: "str", Str: "a"},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

func TestEvalGT(t *testing.T) {
	t.Run("2 > 2", func(t *testing.T) {
		exp := models.Expression{Name: "gt", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 2},
			models.ExpressionArg{DType: "num", Num: 2},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("2 > 1", func(t *testing.T) {
		exp := models.Expression{Name: "gt", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 2},
			models.ExpressionArg{DType: "num", Num: 1},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 > 2", func(t *testing.T) {
		exp := models.Expression{Name: "gt", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 1},
			models.ExpressionArg{DType: "num", Num: 2},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("a > b", func(t *testing.T) {
		exp := models.Expression{Name: "gt", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "str", Str: "a"},
			models.ExpressionArg{DType: "str", Str: "b"},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b > b", func(t *testing.T) {
		exp := models.Expression{Name: "gt", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "str", Str: "b"},
			models.ExpressionArg{DType: "str", Str: "b"},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b > a", func(t *testing.T) {
		exp := models.Expression{Name: "gt", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "str", Str: "b"},
			models.ExpressionArg{DType: "str", Str: "a"},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

func TestEvalGTE(t *testing.T) {
	t.Run("2 >= 2", func(t *testing.T) {
		exp := models.Expression{Name: "gte", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 2},
			models.ExpressionArg{DType: "num", Num: 2},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("2 >= 1", func(t *testing.T) {
		exp := models.Expression{Name: "gte", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 2},
			models.ExpressionArg{DType: "num", Num: 1},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 >= 2", func(t *testing.T) {
		exp := models.Expression{Name: "gte", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 1},
			models.ExpressionArg{DType: "num", Num: 2},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("a >= b", func(t *testing.T) {
		exp := models.Expression{Name: "gte", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "str", Str: "a"},
			models.ExpressionArg{DType: "str", Str: "b"},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b >= b", func(t *testing.T) {
		exp := models.Expression{Name: "gte", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "str", Str: "b"},
			models.ExpressionArg{DType: "str", Str: "b"},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("b >= a", func(t *testing.T) {
		exp := models.Expression{Name: "gte", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "str", Str: "b"},
			models.ExpressionArg{DType: "str", Str: "a"},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

// Logic operators
func TestEvalAND(t *testing.T) {
	t.Run("0 && 0 ", func(t *testing.T) {
		exp := models.Expression{Name: "and", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 0},
			models.ExpressionArg{DType: "num", Num: 0},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 && 0 ", func(t *testing.T) {
		exp := models.Expression{Name: "and", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 1},
			models.ExpressionArg{DType: "num", Num: 0},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("0 && 1 ", func(t *testing.T) {
		exp := models.Expression{Name: "and", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 0},
			models.ExpressionArg{DType: "num", Num: 1},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 && 1 ", func(t *testing.T) {
		exp := models.Expression{Name: "and", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 1},
			models.ExpressionArg{DType: "num", Num: 1},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

func TestEvalOR(t *testing.T) {
	t.Run("0 || 0 ", func(t *testing.T) {
		exp := models.Expression{Name: "or", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 0},
			models.ExpressionArg{DType: "num", Num: 0},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 || 0 ", func(t *testing.T) {
		exp := models.Expression{Name: "or", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 1},
			models.ExpressionArg{DType: "num", Num: 0},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("0 || 1 ", func(t *testing.T) {
		exp := models.Expression{Name: "or", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 0},
			models.ExpressionArg{DType: "num", Num: 1},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})

	t.Run("1 || 1 ", func(t *testing.T) {
		exp := models.Expression{Name: "or", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 1},
			models.ExpressionArg{DType: "num", Num: 1},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}

func TestEvalNOT(t *testing.T) {
	t.Run("0", func(t *testing.T) {
		exp := models.Expression{Name: "not", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 0},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if !ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
	t.Run("1", func(t *testing.T) {
		exp := models.Expression{Name: "not", Data: []models.ExpressionArg{
			models.ExpressionArg{DType: "num", Num: 1},
		}}
		evalContext := evalContext{}
		ret, err := ExpressionEval(exp, evalContext)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
			return
		}
		if ret.(bool) {
			t.Errorf("unexpected value: %b", ret)
		}
	})
}