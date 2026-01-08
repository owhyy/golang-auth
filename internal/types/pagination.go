package types

type PaginationData struct {
	CurrentPage int
	PerPage     int
	TotalPages  int
	Prev        int
	Next        int
}
