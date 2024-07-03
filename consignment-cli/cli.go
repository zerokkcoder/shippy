package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	pb "github.com/zerokkcoder/shippy/consignment-service/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ADDRESS           = "localhost:50051"
	DEFAULT_INFO_FILE = "consignment.json"
)

// 读取 consignment.json 文件
func parseFile(fileName string) (*pb.Consignment, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var consignment *pb.Consignment
	err = json.Unmarshal(data, &consignment)
	if err != nil {
		log.Printf("parse info file error: %v", err)
		return nil, errors.New("consignment.json file content error")
	}
	return consignment, err
}

func main() {
	// 连接 grpc 服务器
	conn, err := grpc.NewClient(ADDRESS, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}
	defer conn.Close()

	// 初始化 gRPC 客户端
	client := pb.NewShippingServiceClient(conn)

	// 在命令行中指定新的货物信息 json 文件
	infoFile := DEFAULT_INFO_FILE
	if len(os.Args) > 1 {
		infoFile = os.Args[1]
	}

	// 解析货物信息
	consignment, err := parseFile(infoFile)
	if err != nil {
		log.Fatalf("parse info file error: %v", err)
	}

	// 调用 grpc
	// 将获取存储到仓库
	resp, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("create consignment error: %v", err)
	}

	// 新获取是否托运成功
	log.Printf("created: %t", resp.Created)

	// 列出目前所有托运的货物
	resp, err = client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("failed to list consignments: %v", err)
	}
	for _, c := range resp.Consignments {
		log.Printf("%+v", c)
	}
}
