package keep

//go:generate go-enum -f "$GOFILE" --noprefix --names --nocase --mustparse

/*
	 ENUM(
		LAST
		SECOND
		MINUTE
		HOUR
		DAY
		WEEK
		MONTH
		QUARTER
		YEAR
		DECADE
		CENTURY
		MILLENIUM

)
*/
type TimeRange int8
