package main

import (
	"context"
	"log"
	"net"

	pb "github.com/zerokkcoder/shippy/consignment-service/proto/consignment"
	"google.golang.org/grpc"
)

const (
	PORT = ":50051"
)

// 仓库接口
type IRepository interface {
	// Create 存放新货物
	Create(consignment *pb.Consignment) (*pb.Consignment, error)
	// GetAll 获取仓库中所有的货物
	GetAll() []*pb.Consignment
}

// 存放多批获取的仓库，实现了 IRepository 接口
type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// 定义微服务
type service struct {
	repo Repository
	pb.UnimplementedShippingServiceServer
}

// CreateConsignment 托运新的货物
// service 实现 consignment.pb.go 中 ShippingServiceServer 接口
// 使 service 作为 gRPC 的服务端
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	// 接收承运的获取
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}
	resp := &pb.Response{Created: true, Consignment: consignment}
	return resp, nil
}

// 获取目前所有托运的货物
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	consignments := s.repo.GetAll()
	resp := &pb.Response{Consignments: consignments}
	return resp, nil
}

func main() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening on: %s\n", PORT)

	server := grpc.NewServer()
	repo := Repository{}

	// 向 grpc 服务器注册微服务
	// 此时会把我们自己实现的微服务 service 与协议中的 ShippingServiceServer 进行绑定
	pb.RegisterShippingServiceServer(server, &service{repo: repo})

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
