package translator_test

import (
	"embed"
	"khanhanhtr/sample/translator"
	"testing"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

//go:embed test_data/*.toml
var test_folder embed.FS

type TestCase struct {
	lang          string
	param         translator.ParamConfig
	plurals_count any
	expected      string
}

func TestTranslateMessage_Unit(t *testing.T) {
	test_specs_name := "test_specs"
	ut, err := translator.NewUtTrans(test_folder, "test_data")
	require.NoError(t, err)

	app := fiber.New()

	tcs := map[string]TestCase{
		"Test Vietnamese Single Case": {
			lang:          "vi",
			param:         translator.ParamConfig{"Count": 1},
			plurals_count: 1,
			expected:      "Đây là 1 apple",
		},

		"Test English Single Case": {
			lang:          "en",
			param:         translator.ParamConfig{"Count": 1},
			plurals_count: 1,
			expected:      "This is 1 apple",
		},

		"Test Vietnamese Plurals Case": {
			lang:          "vi",
			param:         translator.ParamConfig{"Count": 2},
			plurals_count: 2,
			expected:      "Đây là 2 apples",
		},

		"Test English Plurals Case": {
			lang:          "en",
			param:         translator.ParamConfig{"Count": 2},
			plurals_count: 2,
			expected:      "This is 2 apples",
		},
	}

	t.Parallel()
	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {

			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(ctx)
			ctx.Locals("locale", tc.lang)
			result := ut.TranslateMessage(ctx, test_specs_name, tc.param, tc.plurals_count)
			// TODO:  singular doesn't seem to work correctly
			require.Contains(t, result, tc.expected)

		})

	}
}

type (
	TestUserStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	TestRequest struct {
		Locale     string
		Expect     []string
		TestStruct *TestUserStruct
	}
)

func TestTranslateValidateError_Unit(t *testing.T) {
	ut, err := translator.NewUtTrans(test_folder, "test_data")
	require.NoError(t, err)
	app := fiber.New()
	v := validator.New()

	tcs := map[string]TestRequest{
		"Invalid email format error in English": {
			Locale:     "en",
			Expect:     []string{"Invalid Email Address"},
			TestStruct: &TestUserStruct{Name: "ok", Email: "com@com"},
		},

		"Invalid email format error in Vietnamese": {
			Locale:     "vi",
			Expect:     []string{"vui lòng nhập dung email"},
			TestStruct: &TestUserStruct{Name: "ok", Email: "com@com"},
		},

		"Missing name error in English": {
			Locale:     "en",
			Expect:     []string{"Please input name"},
			TestStruct: &TestUserStruct{Name: "", Email: "goodemail@gmail.com"},
		},

		"Invalid name error in Vietnamese": {
			Locale:     "vi",
			Expect:     []string{"vui lòng nhập tên"},
			TestStruct: &TestUserStruct{Name: "", Email: "goodemail@gmail.com"},
		},
	}

	t.Parallel()

	for name, tc := range tcs {
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(ctx)
			ctx.Locals("locale", tc.Locale)

			errs := ut.ValidateRequest(ctx, v, tc.TestStruct)
			require.NotEmpty(t, errs)
			require.Equal(t, tc.Expect, errs)

		})

	}

}
