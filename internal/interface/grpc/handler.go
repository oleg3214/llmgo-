package grpc

import (
	"context"

	"github.com/Ablyamitov/userbot-core/internal/usecase/userbot"
	pb "github.com/Ablyamitov/userbot-protobuf"
	"google.golang.org/grpc"
)

type UserbotGrpcHandler struct {
	pb.UnimplementedUserbotServiceServer
	uc *userbot.Usecase
}

func NewUserbotGrpcHandler(uc *userbot.Usecase) *UserbotGrpcHandler {
	return &UserbotGrpcHandler{uc: uc}
}

func RegisterServices(server *grpc.Server, uc *userbot.Usecase) {
	userbotHandler := NewUserbotGrpcHandler(uc)
	pb.RegisterUserbotServiceServer(server, userbotHandler)
}
func (h *UserbotGrpcHandler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	err := h.uc.SendMessage(ctx, req.GetUsername(), req.GetText())
	if err != nil {
		return nil, err
	}
	return &pb.SendMessageResponse{Success: true}, nil
}

func (h *UserbotGrpcHandler) GetMessages(ctx context.Context, req *pb.MessagesRequest) (*pb.MessagesResponse, error) {
	messages, err := h.uc.GetMessages(ctx, req.GetUsername(), int(req.GetLimit()))
	if err != nil {
		return nil, err
	}

	var pbMessages []*pb.Message
	for _, m := range messages {
		pbMessages = append(pbMessages, &pb.Message{
			MessageId: int64(m.ID),
			Text:      m.Text,
			Date:      int64(m.Date),
		})
	}

	return &pb.MessagesResponse{Messages: pbMessages}, nil
}
