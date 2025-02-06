package main

import (
	"context"

	"github.com/leodahal4/dev-kit/config"
	pb "github.com/leodahal4/dev-kit/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetProject(ctx context.Context, req *pb.ProjectRequest) (*pb.ProjectResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, p := range s.config.Projects {
		if p.ID == req.ProjectId {
			return &pb.ProjectResponse{
				Project: convertToProtoProject(p),
			}, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "project not found")
}

func (s *Server) UpdateProject(ctx context.Context, req *pb.ProjectRequest) (*pb.ProjectResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.config.Projects {
		if p.ID == req.ProjectId {
			// Update project
			s.config.Projects[i] = config.ProjectConfig{
				ID:             req.Project.Id,
				Name:           req.Project.Name,
				Description:    req.Project.Description,
				IsMicroservice: req.Project.IsMicroservice,
				IsValid:        true, // Set default value
				Environments:   make([]config.EnvironmentConfig, len(req.Project.Environments)),
			}

			for j, env := range req.Project.Environments {
				s.config.Projects[i].Environments[j] = config.EnvironmentConfig{
					Name:        env.Name,
					Description: env.Description,
					Language:    env.Language,
					Path:        env.Path,
				}
			}

			if err := saveConfig(s.config); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to save config: %v", err)
			}

			return &pb.ProjectResponse{
				Project: convertToProtoProject(s.config.Projects[i]),
			}, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "project not found")
}

func (s *Server) ListProjects(ctx context.Context, _ *pb.Empty) (*pb.ListProjectsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projects := make([]*pb.ProjectConfig, len(s.config.Projects))
	for i, p := range s.config.Projects {
		projects[i] = convertToProtoProject(p)
	}

	return &pb.ListProjectsResponse{
		Projects: projects,
	}, nil
}
