package utils

import "time"

// 中文：DateTimeLayout 声明当前包使用的常量。
// English: DateTimeLayout declares constants used by this package.
const DateTimeLayout = "2006-01-02 15:04:05"

// 中文：NowUnixMilli 执行当前包中的对应流程。
// English: NowUnixMilli executes the corresponding workflow in this package.
func NowUnixMilli() int64 {
	return time.Now().UnixMilli()
}

// 中文：FormatTime 执行当前包中的对应流程。
// English: FormatTime executes the corresponding workflow in this package.
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DateTimeLayout)
}

// 中文：ParseTime 执行当前包中的对应流程。
// English: ParseTime executes the corresponding workflow in this package.
func ParseTime(value string, layouts ...string) (time.Time, error) {
	if len(layouts) == 0 {
		layouts = []string{time.RFC3339, DateTimeLayout, time.DateOnly}
	}
	var lastErr error
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

// 中文：StartOfDay 执行当前包中的对应流程。
// English: StartOfDay executes the corresponding workflow in this package.
func StartOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// 中文：EndOfDay 执行当前包中的对应流程。
// English: EndOfDay executes the corresponding workflow in this package.
func EndOfDay(t time.Time) time.Time {
	return StartOfDay(t).Add(24*time.Hour - time.Nanosecond)
}
