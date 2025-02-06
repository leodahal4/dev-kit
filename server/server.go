package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/leodahal4/dev-kit/config"
	pb "github.com/leodahal4/dev-kit/protos"
	"github.com/leodahal4/dev-kit/server/models"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedConfigServiceServer
	mu     sync.RWMutex
	config *config.GlobalConfig
	repo   models.RepoImpl
}

func loadConfig() (*config.GlobalConfig, error) {
	file, err := os.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config config.GlobalConfig
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}

func saveConfig(config *config.GlobalConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile("config.json", data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func convertToProtoProject(p config.ProjectConfig) *pb.ProjectConfig {
	environments := make([]*pb.EnvironmentConfig, len(p.Environments))
	for i, env := range p.Environments {
		environments[i] = &pb.EnvironmentConfig{
			Name:        env.Name,
			Description: env.Description,
			Language:    env.Language,
			Path:        env.Path,
		}
	}

	return &pb.ProjectConfig{
		Id:             p.ID,
		Name:           p.Name,
		Description:    p.Description,
		IsMicroservice: p.IsMicroservice,
		Environments:   environments,
	}
}

func (s *Server) GetGlobalConfig(ctx context.Context, _ *pb.Empty) (*pb.GlobalConfigResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projects := make([]*pb.ProjectConfig, len(s.config.Projects))
	for i, p := range s.config.Projects {
		projects[i] = convertToProtoProject(p)
	}

	return &pb.GlobalConfigResponse{
		Debug:           s.config.DEBUG,
		PprofEnabled:    s.config.PPROF_ENABLED,
		PprofAddAndPort: s.config.PPROF_ADD_AND_PORT,
		LogFormat:       s.config.LOG_FORMAT,
		Kubeconfig:      s.config.KUBECONFIG,
		CheckedTools:    s.config.CHECKED_TOOLS,
		Projects:        projects,
		CurrentCmd:      s.config.CURRENT_CMD,
	}, nil
}

func (s *Server) UpdateGlobalConfig(ctx context.Context, req *pb.GlobalConfigRequest) (*pb.GlobalConfigResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.config.DEBUG = req.Config.Debug
	s.config.PPROF_ENABLED = req.Config.PprofEnabled
	s.config.PPROF_ADD_AND_PORT = req.Config.PprofAddAndPort
	s.config.LOG_FORMAT = req.Config.LogFormat
	s.config.KUBECONFIG = req.Config.Kubeconfig
	s.config.CHECKED_TOOLS = req.Config.CheckedTools
	s.config.CURRENT_CMD = req.Config.CurrentCmd

	if err := saveConfig(s.config); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save config: %v", err)
	}

	return req.Config, nil
}

// CreateEnvironment handles the creation of a new environment
func (s *Server) CreateEnvironment(ctx context.Context, req *pb.CreateEnvironmentRequest) (*pb.Empty, error) {
	logrus.Infof("got the request :%+v", req)
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find the project by ID
	for _, project := range s.config.Projects {
		if project.ID == req.ProjectId {
			// Check for duplicate environment names
			for _, env := range project.Environments {
				if env.Name == req.Environment.Name {
					return nil, status.Errorf(codes.AlreadyExists, "environment with this name already exists")
				}
			}

			// Append the new environment
			project.Environments = append(project.Environments, config.EnvironmentConfig{
				Name:        req.Environment.Name,
				Description: req.Environment.Description,
				Path:        req.Environment.Path,
			})

			// Save the updated configuration
			if err := saveConfig(s.config); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to save config: %v", err)
			}

			return &pb.Empty{}, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "project not found")
}

func main() {
	configPath := flag.String("c", "", "Path to the config file")
	flag.Parse()

	// Load the config from the specified path
	config, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize the database connection
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	repo := models.NewConfigRepository(config)
	pb.RegisterConfigServiceServer(s, &Server{
		config: config,
		repo:   repo,
	})

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
