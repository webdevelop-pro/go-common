package round

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type ValueFloat struct {
	v float64
}

func (m ValueFloat) GetFloatValue() float64 {
	return m.v
}

func (m *ValueFloat) SetFloatValue(value float64) {
	m.v = value
}

func TestRoundToLargestRemainder(t *testing.T) {
	type testCases struct {
		description    string
		input          []float64
		expectedResult []float64
		expectedError  error
	}

	for _, scenario := range []testCases{
		{
			description: "real case",
			input: []float64{
				48.648648648648646,
				40.54054054054054,
				10.81081081081081,
			},
			expectedResult: []float64{49, 40, 11},
			expectedError:  nil,
		},
		{
			description: "noramal test case 1",
			input: []float64{
				(15.0 / 30) * 100.0, // 50 %
				(10.0 / 30) * 100.0, // 33.3333 %
				(5.0 / 30) * 100.0,  // 16.666 %
			},
			expectedResult: []float64{50, 33, 17},
			expectedError:  nil,
		},
		{
			description: "noramal test case 2",
			input: []float64{
				(10.0 / 30) * 100.0, // 33.3333 %
				(10.0 / 30) * 100.0, // 33.3333 %
				(5.0 / 30) * 100.0,  // 16.666 %
				(5.0 / 30) * 100.0,  // 16.666 %
			},
			expectedResult: []float64{33, 33, 17, 17},
			expectedError:  nil,
		},
		{
			description: "noramal test case 3",
			input: []float64{
				(10.0 / 30) * 100.0, // 33.3333 %
				(6.0 / 30) * 100.0,  // 20 %
				(5.0 / 30) * 100.0,  // 16.666 %
				(5.0 / 30) * 100.0,  // 16.666 %
				(2.0 / 30) * 100.0,  // 6.666 %
				(1.0 / 30) * 100.0,  // 3.3333 %
				(1.0 / 30) * 100.0,  // 3.3333 %
			},
			expectedResult: []float64{33, 20, 17, 17, 7, 3, 3},
			expectedError:  nil,
		},
		{
			description: "noramal test case 4",
			input: []float64{
				(10.0 / 30) * 100.0, // 33.3333 %
				(5.0 / 30) * 100.0,  // 16.666 %
				(5.0 / 30) * 100.0,  // 16.666 %
				(5.0 / 30) * 100.0,  // 16.666 %
				(3.0 / 30) * 100.0,  // 10 %
				(1.0 / 30) * 100.0,  // 3.3333 %
				(1.0 / 30) * 100.0,  // 3.3333 %
			},
			expectedResult: []float64{33, 17, 17, 17, 10, 3, 3},
			expectedError:  nil,
		},
		{
			description: "noramal test case 5",
			input: []float64{
				(10.0 / 30) * 100.0, // 33.3333 %
				(6.0 / 30) * 100.0,  // 20 %
				(5.0 / 30) * 100.0,  // 16.666 %
				(5.0 / 30) * 100.0,  // 16.666 %
				(1.0 / 30) * 100.0,  // 3.3333 %
				(1.0 / 30) * 100.0,  // 3.3333 %
				(1.0 / 30) * 100.0,  // 3.3333 %
				(1.0 / 30) * 100.0,  // 3.3333 %
			},
			expectedResult: []float64{34, 20, 17, 17, 3, 3, 3, 3},
			expectedError:  nil,
		},

		{
			description: "normal test case 6",
			input: []float64{
				(3.0 / 7) * 100.0, // 42.857142857142855 %
				(2.0 / 7) * 100.0, // 28.5714285714285 %
				(1.0 / 7) * 100.0, // 14.285714285714285 %
				(1.0 / 7) * 100.0, // 14.285714285714285 %
			},
			expectedResult: []float64{43, 29, 14, 14},
			expectedError:  nil,
		},

		{
			description: "normal test case 7",
			input: []float64{
				(15.0 / 19) * 100.0, // 78.94736842105263 %
				(2.0 / 19) * 100.0,  // 10.526315789473684 %
				(1.0 / 19) * 100.0,  // 5.263157894736842 %
				(1.0 / 19) * 100.0,  // 5.263157894736842 %
			},
			expectedResult: []float64{79, 11, 5, 5},
			expectedError:  nil,
		},

		// save equal tests

		{
			description: "save equal for 16 and 16",
			input: []float64{
				(20.0 / 30) * 100.0, // 66.6666 %
				(5.0 / 30) * 100.0,  // 16.666 %
				(5.0 / 30) * 100.0,  // 16.666 %
			},
			expectedResult: []float64{66, 17, 17},
			expectedError:  nil,
		},

		{
			description: "save equal for 5, 5, 5. round 84 to 85",
			input: []float64{
				(16.0 / 19) * 100.0, // 84.21052631578947 %
				(1.0 / 19) * 100.0,  // 5.263157894736842 %
				(1.0 / 19) * 100.0,  // 5.263157894736842 %
				(1.0 / 19) * 100.0,  // 5.263157894736842 %
			},
			expectedResult: []float64{85, 5, 5, 5},
			expectedError:  nil,
		},

		// IMPORTANT //
		// this test case shows a case that
		// does not allow rounding to integer values
		// without losing the equality of numbers
		{
			description: "test case 2 problem",
			input: []float64{
				33.333333,
				33.333333,
				33.333333,
			},
			expectedResult: []float64{34, 33, 33},
			// ALTERNATIVE expectedResult: []float64{33.3, 33.3, 33.3},
			expectedError: nil,
		},

		// INVALID INPUT
		{
			description: "save equal for 5, 5, 5. round 84 to 85",
			input: []float64{
				1,
				8,
				9,
			},
			expectedResult: []float64{},
			expectedError:  errors.New("any error"),
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {

			values := make([]Value, len(scenario.input))
			expectedValues := make([]Value, len(scenario.expectedResult))

			for i := 0; i < len(scenario.input); i++ {
				values[i] = &ValueFloat{v: scenario.input[i]}
				if len(expectedValues) > i {
					expectedValues[i] = &ValueFloat{v: scenario.expectedResult[i]}
				}
			}

			err := SmartRound(values, 100)

			if scenario.expectedError != nil {
				require.NotNil(t, err)
			}
			if len(expectedValues) > 0 {
				require.Equal(t, values, expectedValues)
			}
		})
	}

}
