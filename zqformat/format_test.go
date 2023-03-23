package zqformat

import (
	"testing"

	"github.com/Kotodian/format"
	"github.com/stretchr/testify/assert"
)

func TestParseFormat(t *testing.T) {
	testFormat := `{"imei": "%i", "iccid": "%c", "time": %t, %d}`
	frt, err := ParseFormat(testFormat)

	assert.NoError(t, err)
	f := frt.(*formatter)
	assert.Equal(t, "imei", f.imeiKey)
	assert.Equal(t, "iccid", f.iccidKey)
	assert.Equal(t, "time", f.timeKey)

	err = f.Match([]byte(`{"imei": "123456789012345", "iccid": "98765432109876543210", "time": "2023-03-23 10:30:00", "customData1": "value1", "customData2": "value2"}`))
	assert.NoError(t, err)

	err = f.Match([]byte(`{"imei": "123456789012345", "iccid": "98765432109876543210", "time": "invalid-time-format", "customData1": "value1", "customData2": "value2"}`))
	assert.Error(t, err)

	formatted, err := f.Format("123456789012345", "98765432109876543210", "2023-03-23 10:30:00", map[string]interface{}{"customData1": "value1", "customData2": "value2"})
	assert.NoError(t, err)
	assert.JSONEq(t, `{"imei": "123456789012345", "iccid": "98765432109876543210", "time": "2023-03-23 10:30:00", "customData1": "value1", "customData2": "value2"}`, string(formatted))
}

func TestParseInvalidFormat(t *testing.T) {
	testFormat := `{"imei": "%i", "iccid": "%c", "time": %t}`
	_, err := ParseFormat(testFormat)
	assert.Error(t, err)
}

func TestFormatterImplementation(t *testing.T) {
	var _ format.Formatter = (*formatter)(nil)
}
