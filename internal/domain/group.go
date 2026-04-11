package domain

type Group uint8

const (
	TenthGradeGroup = iota
	EleventhGradeGroup
	TwelfthGradeGroup
	CollegeGroup
)

func (g Group) String() string {
	switch g {
	case TenthGradeGroup:
		return "10th Grade"
	case EleventhGradeGroup:
		return "11th Grade"
	case TwelfthGradeGroup:
		return "12th Grade"
	case CollegeGroup:
		return "College"
	default:
		return "Unknown Group"
	}
}
