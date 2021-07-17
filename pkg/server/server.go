package server

import (
	"github.com/dobriychelpozitivniy/telegram-go-pocket-bot/pkg/repository"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
	"net/http"
	"strconv"
)

type AuthorizationServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string
}

func (s *AuthorizationServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatIDParam := request.URL.Query().Get("chat_id")
	if chatIDParam == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDParam, 10, 64)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	requestToken, err := s.tokenRepository.Get(chatID, repository.RequestTokens)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	authResponse, err := s.pocketClient.Authorize(request.Context(), requestToken)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.tokenRepository.Save(chatID, authResponse.AccessToken, repository.AccessTokens)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("chat_id: %d \nrequest_toket: %s \naccess_token: %s", chatID, requestToken, authResponse.AccessToken)

	writer.Header().Add("Location", s.redirectURL)
	writer.WriteHeader(http.StatusMovedPermanently)
}

func NewAuthorizationServer(pocketClient *pocket.Client, tokenRepository repository.TokenRepository, redirectURL string) *AuthorizationServer {
	return &AuthorizationServer{pocketClient: pocketClient, tokenRepository: tokenRepository, redirectURL: redirectURL}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":80",
		Handler: s,
	}

	return s.server.ListenAndServe()
}
