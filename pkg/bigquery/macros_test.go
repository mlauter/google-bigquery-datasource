package bigquery

import (
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/sqlds/v3"
	"github.com/pkg/errors"
)

func Test_macros(t *testing.T) {
	tests := []struct {
		description      string
		macro            string
		query            *sqlds.Query
		args             []string
		expected         string
		expectedErr      error
		expectedFillMode *data.FillMissing
	}{
		{
			"time groups 1w",
			"timeGroup",
			&sqlds.Query{},
			[]string{"created_at", "1w"},
			"TIMESTAMP_SECONDS(DIV(UNIX_SECONDS(created_at), 604800) * 604800)",
			nil,
			nil,
		},
		{
			"time groups 1d",
			"timeGroup",
			&sqlds.Query{},
			[]string{"created_at", "1d"},
			"TIMESTAMP_SECONDS(DIV(UNIX_SECONDS(created_at), 86400) * 86400)",
			nil,
			nil,
		},
		{
			"time groups 1M",
			"timeGroup",
			&sqlds.Query{},
			[]string{"created_at", "1M"},
			"TIMESTAMP((PARSE_DATE(\"%Y-%m-%d\",CONCAT( CAST((EXTRACT(YEAR FROM created_at)) AS STRING),'-',CAST((EXTRACT(MONTH FROM created_at)) AS STRING),'-','01'))))",
			nil,
			nil,
		},
		{
			"time groups '1M'",
			"timeGroup",
			&sqlds.Query{},
			[]string{"created_at", "'1M'"},
			"TIMESTAMP((PARSE_DATE(\"%Y-%m-%d\",CONCAT( CAST((EXTRACT(YEAR FROM created_at)) AS STRING),'-',CAST((EXTRACT(MONTH FROM created_at)) AS STRING),'-','01'))))",
			nil,
			nil,
		},
		{
			"time groups \"1M\"",
			"timeGroup",
			&sqlds.Query{},
			[]string{"created_at", "\"1M\""},
			"TIMESTAMP((PARSE_DATE(\"%Y-%m-%d\",CONCAT( CAST((EXTRACT(YEAR FROM created_at)) AS STRING),'-',CAST((EXTRACT(MONTH FROM created_at)) AS STRING),'-','01'))))",
			nil,
			nil,
		},
		{
			"time groups fill 0",
			"timeGroup",
			&sqlds.Query{},
			[]string{"created_at", "1w", "0"},
			"TIMESTAMP_SECONDS(DIV(UNIX_SECONDS(created_at), 604800) * 604800)",
			nil,
			&data.FillMissing{
				Mode:  data.FillModeValue,
				Value: 0,
			},
		},
		{
			"time groups fill null",
			"timeGroup",
			&sqlds.Query{},
			[]string{"created_at", "1w", "NULL"},
			"TIMESTAMP_SECONDS(DIV(UNIX_SECONDS(created_at), 604800) * 604800)",
			nil,
			&data.FillMissing{
				Mode:  data.FillModeNull,
			},
		},
		{
			"time groups fill previous",
			"timeGroup",
			&sqlds.Query{},
			[]string{"created_at", "1w", "previous"},
			"TIMESTAMP_SECONDS(DIV(UNIX_SECONDS(created_at), 604800) * 604800)",
			nil,
			&data.FillMissing{
				Mode:  data.FillModePrevious,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			res, err := macros[tt.macro](tt.query, tt.args)
			if (err != nil || tt.expectedErr != nil) && !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error %v, expecting %v", err, tt.expectedErr)
			}
			if res != tt.expected {
				t.Errorf("unexpected result %v, expecting %v", res, tt.expected)
			}
			if (!reflect.DeepEqual(tt.query.FillMissing, tt.expectedFillMode)) {
				t.Errorf("unexpected fill mode %v, expecting %v", tt.query.FillMissing, tt.expectedFillMode)
			}
		})
	}
}
