package models

import "strconv"

type WorkEntry struct {
	Company  string
	Role     string
	YearsExp int
}

type Person struct {
	Name       string
	Age        int64
	Contact    ContactInfo
	Address    Address
	Employment Employment
}

// Level 2
type ContactInfo struct {
	Email  string
	Phone  string
	Social SocialMedia
}

type Address struct {
	Street      string
	City        string
	Country     string
	ZipCode     string
	Region      Region
	Coordinates Coordinates
}

type Employment struct {
	Company     string
	Role        string
	Experience  int
	Skills      []string
	Manager     Manager
	Salary      Salary
	WorkHistory WorkHistory
}

// Level 3
type SocialMedia struct {
	Twitter  string
	LinkedIn string
	Stats    SocialStats
}

type Region struct {
	State    string
	TimeZone string
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

type Manager struct {
	Name    string
	Contact ContactInfo
}

type Salary struct {
	Total     float64
	Currency  string
	Breakdown SalaryBreakdown
}

// Level 4
type SocialStats struct {
	Followers int64
	Posts     int64
	Verified  bool
}

type SalaryBreakdown struct {
	Base      float64
	Bonus     float64
	TaxRegion TaxRegion
}

// Level 5
type TaxRegion struct {
	Code string
	Rate float64
}

// WorkHistory is a slice of WorkEntry that implements ArrayMarshal,
// writing a compact JSON array directly into the encoder buffer.
type WorkHistory []WorkEntry

func (w WorkHistory) MarshalArray(b []byte) ([]byte, error) {
	b = append(b, '[')
	for i, e := range w {
		if i > 0 {
			b = append(b, ',')
		}
		b = append(b, `{"company":"`...)
		b = append(b, e.Company...)
		b = append(b, `","role":"`...)
		b = append(b, e.Role...)
		b = append(b, `","years":`...)
		b = strconv.AppendInt(b, int64(e.YearsExp), 10)
		b = append(b, '}')
	}
	b = append(b, ']')
	return b, nil
}
