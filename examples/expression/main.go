package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/expr-lang/expr"
)

type User struct {
	ID       int
	Username string
}

type Order struct {
	User User
}

func main() {
	successOrders := []Order{
		{User{ID: 1, Username: "alice"}},
		{User{ID: 2, Username: "bob"}},
	}

	failedOrders := []Order{
		{User{ID: 3, Username: "eve"}},
	}

	ctx1 := map[string]interface{}{
		"id":           "AAA",
		"balance":      100.00,
		"total_orders": 2,
		"active":       true,
		"devices": []map[string]interface{}{
			{
				"id":       "D1",
				"name":     "iphone",
				"location": "bangkok",
			},
			{
				"id":       "D2",
				"name":     "ipad",
				"location": "bangkok",
			},
		},
	}

	env := map[string]interface{}{
		"ctx1":                  ctx1,
		"success_orders":        successOrders,
		"failed_order":          failedOrders,
		"len":                   func(arr any) int { return getLength(arr) },
		"get_referred_username": func(id int) string { return fmt.Sprintf("ref_%d", id) },
		"uuidv7":                func() string { return fmt.Sprint(time.Now().Unix()) },
		"test": map[string]map[string]map[string]interface{}{
			"a1": {
				"output": {
					"value": "kk",
				},
			},
		},
	}

	mapping := map[string]string{
		"id":                "!uuidv7()",
		"name":              "!success_orders[0].User.Username + '_test1_' + \"_test2_\"",
		"organization":      "tech",
		"total_orders":      "!len(success_orders) + len(failed_order)",
		"order_failed_rate": "!len(failed_order)/(len(success_orders) + len(failed_order))",
		"referred_username": "!get_referred_username(success_orders[0].User.ID)",
		"user_devices":      "!ctx1.devices",
		"output":            "!test.a1.output.value",
		"cond":              "!test.a1.output.value == ctx1.id",
		"cond2":              "!test.a1.output.value == 'll'",
		"cond3":              "!test.a1.output.value == 'kk'",
		"cond4":              "!test.a1.output.value == 'kk' || test.a1.output.value == 'll'",
		"cond5":              "!test.a1.output.value == 'kk' && test.a1.output.value == 'll'",
		"cond6":              "!(test.a1.output.value == 'kk' && test.a1.output.value == 'll') || 1 == 1 && 1 < 2",
	}

	output, err := ex(env, mapping)

	if err != nil {
		panic(err)
	}

	for k, v := range output {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func getLength(v any) int {
	switch t := v.(type) {
	case []Order:
		return len(t)
	default:
		return 0
	}
}

func ex(env map[string]interface{}, mapping map[string]string) (map[string]string, error) {
	output := map[string]string{}

	for k, v := range mapping {

		if len(v) == 0 {
			output[k] = ""
			continue
		}

		if v[0] != '!' {
			output[k] = v
			continue
		}

		expression := v[1:]

		slog.Info("executing expression", slog.String("expression", expression))

		program, err := expr.Compile(expression, expr.Env(env))

		if err != nil {
			// output[k] = fmt.Sprintf("<compile error: %v>", err)
			// continue
			return nil, fmt.Errorf("error on expression %v: %s", expression, err.Error())
		}

		slog.Info("executing program", slog.String("disassemble", program.Disassemble()))

		result, err := expr.Run(program, env)

		if err != nil {
			// output[k] = fmt.Sprintf("<runtime error: %v>", err)
			// continue
			return nil, fmt.Errorf("error on expression %v: %s", expression, err.Error())
		}

		output[k] = fmt.Sprintf("%v", result)
	}

	return output, nil
}

func main2() {

}

type Ex struct {
}

func (e *Ex) RegisterContext(key string, value interface{}) {
	//
}

// func (e *Ex) RegisterContextJSON(key string, value interface{}) {
// 	//
// }

func (e *Ex) Bind(mapping map[string]string, output interface{}) error {
	//

	return nil
}
