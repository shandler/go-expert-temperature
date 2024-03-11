package domain

import "github.com/shandler/go-expert-temperature/internal/dto"

type ZipCode interface {
	Search(request dto.SearchRequest) dto.SearchResponse
}
