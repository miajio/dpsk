package dto

import "time"

// Timestamp 时间
type Timestamp uint64

func (t Timestamp) Time() time.Time {
	return time.Unix(int64(t), 0)
}

func (t Timestamp) String() string {
	if t == 0 {
		return ""
	}
	return t.Time().Format("2006-01-02 15:04:05")
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}
