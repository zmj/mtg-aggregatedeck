package main

import (
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/zmj/mtg-aggregatedeck/internal/logic"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	mr, err := multipartReader(req)
	if err != nil {
		return errResp(http.StatusBadRequest, err)
	}
	decks, err := parseDecks(mr)
	if err != nil {
		return errResp(http.StatusBadRequest, err)
	}
	aggregated, err := aggregate(decks)
	if err != nil {
		return errResp(http.StatusInternalServerError, err)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       aggregated,
	}, nil
}

func errResp(status int, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       err.Error(),
	}, err
}

func dbg(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	s := ""
	for k, v := range req.Headers {
		s += fmt.Sprintf("%v: %v\n", k, v)
	}
	s += req.Body
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       s,
	}, nil
}

func multipartReader(req events.APIGatewayProxyRequest) (*multipart.Reader, error) {
	if req.HTTPMethod != http.MethodPost {
		return nil, fmt.Errorf("Expected POST")
	}
	ctype := req.Headers["content-type"]
	mediaType, mimeParams, err := mime.ParseMediaType(ctype)
	if err != nil {
		return nil, fmt.Errorf("Parse content-type '%v': %v", ctype, err)
	}
	if !strings.HasPrefix(mediaType, "multipart/") {
		return nil, fmt.Errorf("Expected Content-Type multipart: %v", mediaType)
	}
	mr := multipart.NewReader(strings.NewReader(req.Body), mimeParams["boundary"])
	return mr, nil
}

func parseDecks(mr *multipart.Reader) ([]*logic.Deck, error) {
	var decks []*logic.Deck
	for {
		file, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if file.FileName() == "" {
			continue
		}
		deck, err := logic.NewDeck(file)
		if err != nil {
			return nil, fmt.Errorf("Parse deck: %v", err)
		}
		decks = append(decks, deck)
	}
	return decks, nil
}

func aggregate(decks []*logic.Deck) (string, error) {
	d, err := logic.Aggregate(decks)
	return d.String(), err
}
