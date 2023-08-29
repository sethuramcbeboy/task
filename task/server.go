package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"task/constants"
	pb "task/helloworld"
	"log"
	"net"
	"sync"
)

type TaskServiceServer struct {
	mu     sync.Mutex
	tasks  map[string]*pb.Task
	client *mongo.Client // MongoDB client
	pb.UnimplementedTaskServiceServer
}

func (s *TaskServiceServer) AddTask(ctx context.Context, req *pb.Task) (*pb.TaskResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock() // Generate a new account ID
	accountID := generateAccountID(s.client)
	req.Id = accountID
	s.tasks[accountID] = req // Insert account details into MongoDB
	collection := s.client.Database("acc").Collection("account")
	_, err := collection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.TaskResponse{Id: accountID}, nil
} // ...
func generateAccountID(client *mongo.Client) string {
	ctx := context.Background()
	collection := client.Database("acc").Collection("account")
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.M{"_id": "accountID"}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	result := collection.FindOneAndUpdate(ctx, filter, update, opts)
	if result.Err() != nil {
		log.Fatal(result.Err())
	}
	var counter struct {
		Seq int `bson:"seq"`
	}
	err := result.Decode(&counter)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("ACC%05d", counter.Seq)
}
func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	} // Connect to MongoDB
	clientOptions := options.Client().ApplyURI(constants.Connectionstring)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		return
	}
	server := grpc.NewServer()
	pb.RegisterTaskServiceServer(server, &TaskServiceServer{
		tasks: make(map[string]*pb.Task), client: client})
	fmt.Println("Server listening on :50051")
	if err := server.Serve(lis); err != nil {
		fmt.Printf("Failed to serve: %v", err)
	}
}
