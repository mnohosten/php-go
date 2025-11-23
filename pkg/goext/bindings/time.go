package bindings

import (
	"fmt"
	"time"

	"github.com/krizos/php-go/pkg/goext"
	"github.com/krizos/php-go/pkg/types"
	"github.com/krizos/php-go/pkg/vm"
)

// TimeExtension provides Go time/date functions for PHP.
type TimeExtension struct {
	*goext.BaseExtension
}

// NewTimeExtension creates a new time extension.
func NewTimeExtension() *TimeExtension {
	ext := &TimeExtension{
		BaseExtension: goext.NewBaseExtension("go_time", "1.0.0"),
	}

	// Register functions
	ext.AddFunction("go_time_now", ext.timeNow)
	ext.AddFunction("go_time_unix", ext.timeUnix)
	ext.AddFunction("go_time_format", ext.timeFormat)
	ext.AddFunction("go_time_parse", ext.timeParse)
	ext.AddFunction("go_time_sleep", ext.timeSleep)
	ext.AddFunction("go_time_add", ext.timeAdd)
	ext.AddFunction("go_time_diff", ext.timeDiff)

	// Register constants (common time formats)
	ext.AddConstant("TIME_RFC3339", types.NewString(time.RFC3339))
	ext.AddConstant("TIME_RFC822", types.NewString(time.RFC822))
	ext.AddConstant("TIME_ANSIC", types.NewString(time.ANSIC))
	ext.AddConstant("TIME_UNIX_DATE", types.NewString(time.UnixDate))
	ext.AddConstant("TIME_RUBY_DATE", types.NewString(time.RubyDate))
	ext.AddConstant("TIME_KITCHEN", types.NewString(time.Kitchen))

	return ext
}

// timeNow implements go_time_now()
//
// Returns: current Unix timestamp
func (e *TimeExtension) timeNow(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("go_time_now() expects 0 arguments, got %d", len(args))
	}

	return types.NewInt(time.Now().Unix()), nil
}

// timeUnix implements go_time_unix($timestamp)
//
// Returns: array with time components
//   - year, month, day, hour, minute, second, weekday
func (e *TimeExtension) timeUnix(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_time_unix() expects 1 argument (timestamp), got %d", len(args))
	}

	timestamp := args[0].ToInt()
	t := time.Unix(timestamp, 0)

	result := types.NewEmptyArray()
	result.Set(types.NewString("year"), types.NewInt(int64(t.Year())))
	result.Set(types.NewString("month"), types.NewInt(int64(t.Month())))
	result.Set(types.NewString("day"), types.NewInt(int64(t.Day())))
	result.Set(types.NewString("hour"), types.NewInt(int64(t.Hour())))
	result.Set(types.NewString("minute"), types.NewInt(int64(t.Minute())))
	result.Set(types.NewString("second"), types.NewInt(int64(t.Second())))
	result.Set(types.NewString("weekday"), types.NewInt(int64(t.Weekday())))

	return types.NewArray(result), nil
}

// timeFormat implements go_time_format($timestamp, $format)
//
// Returns: formatted time string
func (e *TimeExtension) timeFormat(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("go_time_format() expects 2 arguments (timestamp, format), got %d", len(args))
	}

	timestamp := args[0].ToInt()
	format := args[1].ToString()

	t := time.Unix(timestamp, 0)
	formatted := t.Format(format)

	return types.NewString(formatted), nil
}

// timeParse implements go_time_parse($format, $value)
//
// Returns: Unix timestamp or false on error
func (e *TimeExtension) timeParse(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("go_time_parse() expects 2 arguments (format, value), got %d", len(args))
	}

	format := args[0].ToString()
	value := args[1].ToString()

	t, err := time.Parse(format, value)
	if err != nil {
		return types.NewBool(false), nil
	}

	return types.NewInt(t.Unix()), nil
}

// timeSleep implements go_time_sleep($seconds)
//
// Sleeps for the specified number of seconds
func (e *TimeExtension) timeSleep(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("go_time_sleep() expects 1 argument (seconds), got %d", len(args))
	}

	seconds := args[0].ToInt()
	time.Sleep(time.Duration(seconds) * time.Second)

	return types.NewNull(), nil
}

// timeAdd implements go_time_add($timestamp, $duration)
//
// Adds duration (in seconds) to timestamp
func (e *TimeExtension) timeAdd(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("go_time_add() expects 2 arguments (timestamp, duration), got %d", len(args))
	}

	timestamp := args[0].ToInt()
	duration := args[1].ToInt()

	t := time.Unix(timestamp, 0)
	newTime := t.Add(time.Duration(duration) * time.Second)

	return types.NewInt(newTime.Unix()), nil
}

// timeDiff implements go_time_diff($timestamp1, $timestamp2)
//
// Returns: difference in seconds between two timestamps
func (e *TimeExtension) timeDiff(v *vm.VM, args []*types.Value) (*types.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("go_time_diff() expects 2 arguments (timestamp1, timestamp2), got %d", len(args))
	}

	timestamp1 := args[0].ToInt()
	timestamp2 := args[1].ToInt()

	t1 := time.Unix(timestamp1, 0)
	t2 := time.Unix(timestamp2, 0)

	diff := t2.Sub(t1).Seconds()

	return types.NewInt(int64(diff)), nil
}
