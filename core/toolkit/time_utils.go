package toolkit

import "time"

var TimeUtils = &timeUtil{
	SecondStamp:   1e3,
	MinuteStamp:   60e3,
	HourStamp:     3600e3,
	DayStamp:      86400e3,
	Year2100Stamp: 4102416000000,

	Nanosecond:  time.Nanosecond,
	Microsecond: time.Microsecond,
	Millisecond: time.Millisecond,
	Second:      time.Second,
	Minute:      time.Minute,
	Hour:        time.Hour,
	Day:         time.Hour * 24,
	Week:        time.Hour * 24 * 7,

	Layout:           "2006-01-02 15:04:05",
	LayoutDate:       "2006-01-02",
	LayoutTime:       "15:04:05",
	LayoutTimeMinute: "15:04",
	LayoutDateMinute: "2006-01-02 15:04",
}

// 时间工具
//
// 这个工具内的所有时间戳均以毫秒为单位
type timeUtil struct {
	SecondStamp   int64 // 一秒的毫秒数
	MinuteStamp   int64 // 一分的毫秒数
	HourStamp     int64 // 一小时的毫秒数
	DayStamp      int64 // 一天的毫秒数
	Year2100Stamp int64 // 2100年的时间戳

	Nanosecond  time.Duration // 纳秒
	Microsecond time.Duration // 微秒
	Millisecond time.Duration // 毫秒
	Second      time.Duration // 秒
	Minute      time.Duration // 分
	Hour        time.Duration // 时
	Day         time.Duration // 天
	Week        time.Duration // 周

	Layout           string // 默认时间样式 YYYY-MM-DD hh:mm:ss
	LayoutDate       string // 日期样式 YYYY-MM-DD
	LayoutTime       string // 时间样式 hh:mm:ss
	LayoutTimeMinute string // 分钟样式 hh:mm
	LayoutDateMinute string // 带日期的分钟样式 YYYY-MM-DD hh:mm
}

// region 基础
// 获取当前毫秒级时间戳
func (*timeUtil) GetStamp() int64 {
	return time.Now().UnixNano() / 1e6
}

// 获取当前时间
func (*timeUtil) GetTime() time.Time {
	return time.Now()
}

// 获取当前时间默认样式, YYYY-MM-DD hh:mm:ss
func (u *timeUtil) GetText() string {
	return time.Now().Format(u.Layout)
}

// 获取当前时间日期样式, YYYY-MM-DD
func (u *timeUtil) GetDateText() string {
	return time.Now().Format(u.LayoutDate)
}

// 获取当前时间时间样式, hh:mm:ss
func (u *timeUtil) GetTimeText() string {
	return time.Now().Format(u.LayoutTime)
}

// 获取当前时间分钟样式, hh:mm
func (u *timeUtil) GetTimeMinuteText() string {
	return time.Now().Format(u.LayoutTimeMinute)
}

// 获取当前时间带日期的分钟样式, YYYY-MM-DD hh:mm
func (u *timeUtil) GetDateMinuteTextHour() string {
	return time.Now().Format(u.LayoutDateMinute)
}

// endregion

// region 转换

// 将时间转为毫秒级时间戳
func (u *timeUtil) TimeToStamp(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

// 将时间转为毫秒级时间戳
func (u *timeUtil) TimeToSecond(t time.Time) uint32 {
	return uint32(t.UnixNano() / 1e6 /1000)
}

// 将时间转为默认样式的字符串
func (u *timeUtil) TimeToText(t time.Time) string {
	return t.Format(u.Layout)
}

// 将时间转为日期样式的字符串
func (u *timeUtil) TimeToDateText(t time.Time) string {
	return t.Format(u.LayoutDate)
}

// 将时间转为指定样式的字符串
func (*timeUtil) TimeToTextOfLayout(t time.Time, layout string) string {
	return t.Format(layout)
}

// 毫秒级时间戳转时间
func (u *timeUtil) StampToTime(stamp int64) time.Time {
	return time.Unix(0, stamp*1e6)
}

// 将毫秒级时间戳转为默认样式的字符串
func (u *timeUtil) StampToText(stamp int64) string {
	return time.Unix(0, stamp*1e6).Format(u.Layout)
}

// 将毫秒级时间戳转为日期样式的字符串
func (u *timeUtil) StampToDateText(stamp int64) string {
	return time.Unix(0, stamp*1e6).Format(u.LayoutDate)
}

// 将毫秒级时间戳转为指定样式的字符串
func (*timeUtil) StampToTextOfLayout(stamp int64, layout string) string {
	return time.Unix(0, stamp*1e6).Format(layout)
}

// 将标准样式时间文本转为时间
func (u *timeUtil) TextToTime(text string) (time.Time, error) {
	return time.ParseInLocation(u.Layout, text, time.Local)
}

// 将日期样式时间文本转为时间
func (u *timeUtil) DateTextToTime(text string) (time.Time, error) {
	return time.ParseInLocation(u.LayoutDate, text, time.Local)
}

// 将指定样式时间文本转为时间
func (*timeUtil) TextToTimeOfLayout(text, layout string) (time.Time, error) {
	return time.ParseInLocation(layout, text, time.Local)
}

// 将标准样式时间文本转为毫秒级时间戳
func (u *timeUtil) TextToStamp(text string) (int64, error) {
	return u.TextToStampOfLayout(text, u.Layout)
}

// 将日期样式时间文本转为毫秒级时间戳
func (u *timeUtil) DateTextToStamp(text string) (int64, error) {
	return u.TextToStampOfLayout(text, u.LayoutDate)
}

// 将时间文本转为毫秒级时间戳
func (*timeUtil) TextToStampOfLayout(text, layout string) (int64, error) {
	t, e := time.ParseInLocation(layout, text, time.Local)
	if e != nil {
		return 0, e
	}
	return t.UnixNano() / 1e6, nil
}

// endregion

// 获取当天开始时的毫秒级时间戳(0时0分0秒)
func (*timeUtil) GetDayStartTime() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// 获取当天开始时的毫秒级时间戳(0时0分0秒)
func (*timeUtil) GetDayStartStamp() int64 {
	t := time.Now()
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return t.Unix() * 1e3
}

// 获取传入时间戳当天的开始时间戳(0时0分0秒)
func (*timeUtil) GetDayStartTimeOfTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// 获取传入时间戳当天的开始时间戳(0时0分0秒)
func (*timeUtil) GetDayStartStampOfStamp(stamp int64) int64 {
	t := time.Unix(0, stamp*1e6)
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return t.Unix() * 1e3
}