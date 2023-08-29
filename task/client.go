package main
import (    
	pb "task/helloworld"    
	"context"    
	"fmt"    
	"log"    
	"google.golang.org/grpc"
)
func main() {    
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())    
	if err != nil {        
		log.Fatalf("Failed to connect: %v", err)    
	}    
	defer conn.Close()    
	client := pb.NewTaskServiceClient(conn)    // Enter account details    
	var title string    
	var completed bool    
	fmt.Print("Enter account title: ")    
	fmt.Scanln(&title)    
	fmt.Print("Is account completed (true/false): ")    
	fmt.Scanln(&completed)    
	task := &pb.Task{        
		Title:     title,       
		 Completed: completed,    
	}    // Add account    
	addResp, err := client.AddTask(context.Background(), task)    
	if err != nil {        
		log.Fatalf("Failed to add account: %v", err)    
	}    
	fmt.Printf("Added account with ID: %s\n", addResp.Id)    // Get list of accounts    
	tasksResp, err := client.GetTasks(context.Background(), &pb.Empty{})    
	if err != nil {        
		log.Fatalf("Failed to retrieve accounts: %v", err)    
	}    
	fmt.Println("Accounts:")    
	for _, task := range tasksResp.Tasks {        
		fmt.Printf("ID: %s, Title: %s, Completed: %v\n", task.Id, task.Title, task.Completed)    
}
}