package ocr

import (
	context "context"
	"log"

	"github.com/otiai10/gosseract/v2"
)

type Server struct {
	UnimplementedOCRServer
}

func (s *Server) Read(ctx context.Context, in *ReadRequest) (*ReadReply, error) {
	filePath := in.GetFilePath()
	log.Printf("ocr.Read filePath: %v", filePath)

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(filePath)
	client.Languages = []string{"eng"}
	client.SetWhitelist("*xX×+%-.0123456789")
	text, _ := client.Text()

	log.Printf("ocr.Read text: %v", text)
	return &ReadReply{Text: text}, nil
}
