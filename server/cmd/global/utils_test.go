package utils_test

import (
	"math"
	"math/big"
	"os"
	"teniditter-server/cmd/global/utils"
	utils_enc "teniditter-server/cmd/global/utils/encryption"
	"testing"
)

func init() {
	os.Setenv("TEST", "1")
	// utils_env.LoadEnv()
}

type TestCase[T any, E any] struct {
	input    T
	expected E
}

func TestIsEmptyString(t *testing.T) {
	ref := "notemtpy"
	testCases := []TestCase[any, bool]{
		{input: "", expected: true},
		{input: " ", expected: true},
		{input: "test", expected: false},
		{input: " test ", expected: false},
		{input: true, expected: true},
		{input: []string{"test", ""}, expected: true},
		{input: string([]byte{' '}), expected: true},
		{input: string([]byte{'i', 'l'}), expected: false},
		{input: []byte{'i', 'l'}, expected: true},
		{input: ref, expected: false},
		{input: &ref, expected: true},
		{input: 12, expected: true},
	}
	for _, test := range testCases {
		result := utils.IsEmptyString(test.input)
		if result != test.expected {
			t.Errorf("IsEmptyString(%s), expected: %t, got: %t", test.input, test.expected, result)
		}
	}
}
func TestIsSafeString(t *testing.T) {
	testCases := []TestCase[string, bool]{
		{input: "", expected: true},
		{input: " ", expected: false},
		{input: "test", expected: true},
		{input: "%", expected: false},
		{input: "test test", expected: false},
		{input: "15Ilingu28", expected: true},
		{input: "Test\nTest", expected: false},
		{
			input: `Test
							Test`,
			expected: false,
		},
		{input: "*I&l^i%n#g$u", expected: false},
		{input: `<script>alert("hacked")</script>`, expected: false},
	}
	for _, test := range testCases {
		result := utils.IsSafeString(test.input)
		if result != test.expected {
			t.Errorf("IsSafeString(%s), expected: %t, got: %t", test.input, test.expected, result)
		}
	}
}

func TestSafeString(t *testing.T) {
	testCases := []TestCase[string, string]{
		{input: "CAPITAL", expected: "capital"},
		{input: "teSt123", expected: "test123"},
		{input: " tesst test ", expected: "tesst+test"},
		{
			input:    "[1, test, ILOVE_MATH, #i-Hate-testing, exp(ln(23)+ln(3))]",
			expected: "%5B1%2C+test%2C+ilove_math%2C+%23i-hate-testing%2C+exp%28ln%2823%29%2Bln%283%29%29%5D",
		},
	}
	for _, test := range testCases {
		result := utils.SafeString(test.input)
		if result != test.expected {
			t.Errorf("SafeString(%s), expected: %s, got: %s", test.input, test.expected, result)
		}
	}
}

func TestFormatUsername(t *testing.T) {
	testCases := []TestCase[string, string]{
		{input: "CAPITAL", expected: "capital"},
		{input: "12IWGu773gej09d2eu", expected: "iwgugejdeu"},
		{input: "%20", expected: ""},
		{input: "ili gu", expected: "iligu"},
	}
	for _, test := range testCases {
		result := utils.FormatUsername(test.input)
		if result != test.expected {
			t.Errorf("FormatUsername(%s), expected: %s, got: %s", test.input, test.expected, result)
		}
	}
}

func TestContainsScript(t *testing.T) {
	testCases := []TestCase[string, bool]{
		{input: "noscripttag", expected: false},
		{
			input: `
			<html>
					<head>
					</head>
					<body>
						<a href="#">
							<p><script>alert("hacked");</script></p>
						</a>
					</body>
			</html>`,
			expected: true,
		},
		{
			input: `
		<!DOCTYPE html>
<!--[if lt IE 7]> <html class="ie ie6 lt-ie9 lt-ie8 lt-ie7" lang="en"> <![endif]-->
<!--[if IE 7]>    <html class="ie ie7 lt-ie9 lt-ie8"        lang="en"> <![endif]-->
<!--[if IE 8]>    <html class="ie ie8 lt-ie9"               lang="en"> <![endif]-->
<!--[if IE 9]>    <html class="ie ie9"                      lang="en"> <![endif]-->
<!--[if !IE]><!-->
<html lang="en" class="no-ie">
  <!--<![endif]-->

  <head>
    <title>Gophercises - Coding exercises for budding gophers</title>
  </head>

  <body>
    <section class="header-section">
      <div class="jumbo-content">
        <div class="pull-right login-section">
          Already have an account?
          <a href="#" class="btn btn-login"
            >Login <i class="fa fa-sign-in" aria-hidden="true"></i
          ></a>
        </div>
        <center>
          <img
            src="https://gophercises.com/img/gophercises_logo.png"
            style="max-width: 85%; z-index: 3"
          />
          <h1>coding exercises for budding gophers</h1>
          <br />
          <form action="/do-stuff" method="post">
            <div class="input-group">
              <input
                type="email"
                id="drip-email"
                name="fields[email]"
                class="btn-input"
                placeholder="Email Address"
                required
              />
              <button class="btn btn-success btn-lg" type="submit">
                Sign me up!
              </button>
              <a href="/lost"
                >Lost? Need help?
                <script></script>
              </a>
            </div>
          </form>
          <p class="disclaimer disclaimer-box">
            Gophercises is 100% FREE, but is currently in beta. There will be
            bugs, and things will be changing significantly over the coming
            weeks.
          </p>
        </center>
      </div>
    </section>
    <section class="footer-section">
      <div class="row">
        <div class="col-md-6 col-md-offset-1 vcenter">
          <div class="quote">
            "Success is no accident. It is hard work, perseverance, learning,
            studying, sacrifice and most of all, love of what you are doing or
            learning to do." - Pele
          </div>
        </div>
        <div class="col-md-4 col-md-offset-0 vcenter">
          <center>
            <img
              src="https://gophercises.com/img/gophercises_lifting.gif"
              style="width: 80%"
            />
            <br />
            <br />
          </center>
        </div>
      </div>
      <div class="row">
        <div class="col-md-10 col-md-offset-1">
          <center>
            <p class="disclaimer">
              Artwork created by Marcus Olsson (<a
                href="https://twitter.com/marcusolsson"
                >@marcusolsson</a
              >), animated by Jon Calhoun (that's me!), and inspired by the
              original Go Gopher created by Renee French.
            </p>
          </center>
        </div>
      </div>
    </section>
  </body>
</html>
`,
			expected: true,
		},
		{
			input: `<html>
  <head>
    <link
      rel="stylesheet"
      href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css"
    />
  </head>
  <body>
    <h1>Social stuffs</h1>
    <div>
      <a href="https://www.twitter.com/joncalhoun">
        Check me out on twitter
        <i class="fa fa-twitter" aria-hidden="true"></i>
      </a>
      <a href="https://github.com/gophercises">
        Gophercises is on <strong>Github</strong>!
      </a>
    </div>
  </body>
</html>
`,
			expected: false,
		},
	}
	for _, test := range testCases {
		result := utils.ContainsScript(test.input)
		if result != test.expected {
			t.Errorf("ContainsScript(%s), expected: %t, got: %t", test.input, test.expected, result)
		}
	}
}

func TestBigIntToInt(t *testing.T) {
	testCases := []TestCase[*big.Int, int64]{
		{input: big.NewInt(25121999), expected: 25121999},
		{input: big.NewInt(int64(1e18)), expected: 1e18},
		{input: big.NewInt(-int64(1e18)), expected: -1e18},
	}
	for _, test := range testCases {
		result, err := utils.BigIntToInt(test.input, 64)
		if err != nil {
			t.Errorf("BigIntToInt(%s), got error: %s", test.input.String(), err.Error())
		} else if result != test.expected {
			t.Errorf("BigIntToInt(%s), expected: %d, got: %d", test.input.String(), test.expected, result)
		}
	}
}

func TestShuffleSlice(t *testing.T) {
	testCases := []TestCase[[]any, error]{
		{input: []any{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{input: []any{"apple", "microsoft", "google", "facebook", "amazon"}},
	}
	for _, test := range testCases {
		hashSumBefore := utils_enc.GenerateHashFromArgs(test.input...)
		utils.ShuffleSlice(test.input)
		hashSumAfter := utils_enc.GenerateHashFromArgs(test.input...)

		if hashSumBefore == hashSumAfter {
			t.Errorf("ShuffleSlice(%s), array not shuffled", test.input)
		}
	}
}

func TestIsValidURL(t *testing.T) {
	testCases := []TestCase[string, bool]{
		{input: "", expected: false},
		{input: "not an url", expected: false},
		{input: "exemple.com", expected: false},
		{input: "http exemple.com", expected: false},
		{input: "://ack.vercel.app", expected: false},

		{input: "/some/route", expected: true},
		{input: "http://ack.vercel.app", expected: true},
		{input: "https://ack.vercel.app/some/route", expected: true},
		{input: "https://ack.vercel.app/some/route?test=true&vitest634%20great", expected: true},
	}
	for _, test := range testCases {
		result := utils.IsValidURL(test.input)
		if result != test.expected {
			t.Errorf("IsValidURL(%s), expected: %t, got: %t", test.input, test.expected, result)
		}
	}
}
func TestGenerateRandomChars(t *testing.T) {
	testCases := []TestCase[uint, error]{}
	for i := 0; i <= 8; i++ {
		testCases = append(testCases, TestCase[uint, error]{input: uint(math.Pow(2, float64(i)))})
	}

	for _, test := range testCases {
		result, err := utils.GenerateRandomChars(test.input)
		if err != nil {
			t.Errorf("GenerateRandomChars(%d), got an error", test.input)
		} else if len(result) != int(test.input) {
			t.Errorf("GenerateRandomChars(%d), expected: %s, got: %s", test.input, test.expected, result)
		}
	}
}
func TestIsStrongPassword(t *testing.T) {
	testCases := []TestCase[string, bool]{}

	for _, test := range testCases {
		result := utils.IsStrongPassword(test.input)
		if result != test.expected {
			t.Errorf("GenerateRandomChars(%d), expected: %t, got: %t", test.input, test.expected, result)
		}
	}
}
