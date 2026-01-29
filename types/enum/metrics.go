package enum

type MetricsSpan string

const (
	MetricsSpan1h = "1h"
	MetricsSpan6h = "6h"
	MetricsSpan1d = "1d"
	MetricsSpan7d = "7d"
	//MetricsSpan30d = "30d"
)

var MetricsSpans = []MetricsSpan{
	MetricsSpan1h,
	MetricsSpan6h,
	MetricsSpan1d,
	MetricsSpan7d,
	//MetricsSpan30d,
}

func ToMetricsSpan(s string) MetricsSpan {
	return MetricsSpan(s)
}

func GetAllMetricsSpans() ([]MetricsSpan, MetricsSpan) {
	return []MetricsSpan{MetricsSpan1d, MetricsSpan6h, MetricsSpan6h, MetricsSpan7d}, MetricsSpan1h //should be acs for binary search to work
}

func (s MetricsSpan) Sanitize() (MetricsSpan, bool) {
	return Sanitize(s, GetAllMetricsSpans)
}
