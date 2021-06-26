package partunmarhal

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type timestamp int64

func (ts *timestamp) MarshalJSON() ([]byte, error) {
	t := time.Unix(int64(*ts), 0)
	return t.MarshalJSON()
}

func (ts *timestamp) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	*ts = timestamp(t.Unix())
	return nil
}

type fruit struct {
	Name      string    `json:"name"`                // not updatable
	Color     string    `json:"color" u:"true"`      // updatable
	Price     int       `json:"price" u:"true"`      // updatable
	CreatedAt timestamp `json:"created_at" u:"true"` // special marshalling
}

func watermelon() *fruit {
	return &fruit{
		Name:      "watermelon",
		Color:     "green",
		Price:     100,
		CreatedAt: 1624701600, //  Saturday, June 26, 2021 10:00:00 AM GMT
	}
}

func TestJSON(t *testing.T) {
	testCases := []struct {
		msg         string
		obj         interface{}
		input       []byte
		expectedObj interface{}
		expectedErr error
	}{
		{
			msg:         "return error if not pointer",
			obj:         fruit{},
			input:       []byte(`{}`),
			expectedErr: ErrPointerExpected,
		},
		{
			msg: "return error if field not settable",
			obj: &struct {
				x int `u:"true"`
			}{1},
			input:       []byte(`{"x":10}`),
			expectedErr: ErrFieldCannotBeSet,
		},
		{
			msg: "return error if field does not have json",
			obj: &struct {
				X int `u:"true"`
			}{1},
			input:       []byte(`{"x":10}`),
			expectedErr: ErrNoTagJSON,
		},
		{
			msg: "only update updatable fields",
			obj: watermelon(),
			input: []byte(`{
				"name":  "orange",
				"color": "green,red",
				"price":   200,
				"created_at":   "2021-06-26T11:00:00Z"
				}`),
			expectedObj: &fruit{
				Name:      "watermelon",
				Color:     "green,red",
				Price:     200,
				CreatedAt: 1624705200,
			},
		},
		{
			msg:   "only update fields present in input",
			obj:   watermelon(),
			input: []byte(`{"color": "green,red"}`),
			expectedObj: &fruit{
				Name:      "watermelon",
				Color:     "green,red",
				Price:     100,
				CreatedAt: 1624701600,
			},
		},
	}

	for _, tc := range testCases {
		err := JSON(tc.obj, tc.input)
		assert.Equal(t, tc.expectedErr, err, tc.msg)
		if tc.expectedErr == nil {
			assert.Equal(t, tc.expectedObj, tc.obj, tc.msg)
		}
	}
}
