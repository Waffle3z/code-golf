package pretty

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"
)

const day = 24 * time.Hour

func Bytes(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d.0 B", b)
	}
	div, exp := int(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KM"[exp])
}

// Comma returns a string of a int with thousand separators. Only 0 - 999,999.
func Comma(i int) string {
	if i > 999 {
		return fmt.Sprintf("%d,%03d", i/1000, i%1000)
	}

	return strconv.Itoa(i)
}

// Ordinal returns the ordinal of a int.
func Ordinal(i int) string {
	switch i % 10 {
	case 1:
		if i%100 != 11 {
			return "st"
		}
	case 2:
		if i%100 != 12 {
			return "nd"
		}
	case 3:
		if i%100 != 13 {
			return "rd"
		}
	}

	return "th"
}

// Time returns a fuzzy HTML <time> tag of a time.Time.
//
//    a min ago
//   2 mins ago
//          ...
//  59 mins ago
//  an hour ago
//  2 hours ago
//          ...
// 23 hours ago
//    a day ago
//   2 days ago
//          ...
//  28 days ago
//   1 Jan 2020
//          ...
//  31 Dec 2020
func Time(t time.Time) template.HTML {
	var sb strings.Builder

	rfc := t.Format(time.RFC3339)

	sb.WriteString("<time datetime=")
	sb.WriteString(rfc)
	sb.WriteString(" title=")
	sb.WriteString(rfc)
	sb.WriteRune('>')

	switch diff := time.Since(t); {
	case diff < 2*time.Minute:
		sb.WriteString("a min ago")
	case diff < time.Hour:
		fmt.Fprintf(&sb, "%d mins ago", diff/time.Minute)
	case diff < 2*time.Hour:
		sb.WriteString("an hour ago")
	case diff < day:
		fmt.Fprintf(&sb, "%d hours ago", diff/time.Hour)
	case diff < 2*day:
		sb.WriteString("a day ago")
	case diff < 28*day:
		fmt.Fprintf(&sb, "%d days ago", diff/day)
	default:
		sb.WriteString(t.Format("2 Jan 2006"))
	}

	sb.WriteString("</time>")

	return template.HTML(sb.String())
}
