package service

import (
	"github.com/jacob-alt-del/explore-service/internal/dataaccess"
	pb "github.com/jacob-alt-del/explore-service/internal/proto"
)

type ExploreServiceServer struct {
	pb.UnimplementedExploreServiceServer
	Repo *dataaccess.Repository
}

func NewExploreServiceServer(db *dataaccess.Repository) *ExploreServiceServer {
	return &ExploreServiceServer{Repo: db}
}
