package handler

import (
	"context"
	"demo/internal/middleware"
	"demo/internal/model"
	"demo/internal/service"
	newerror "demo/pkg/error"
	newlog "demo/pkg/log"
	"encoding/json/v2"
	"errors"
	"uuid"

	"github.com/gorilla/websocket"
)

type Websocket struct {
	logger  newlog.Logger
	conn    *websocket.Conn
	cancel  context.CancelFunc
	SendCh  chan model.Message
	service *service.Service
	Uuid    uuid.UUID
}

func NewWebsocket(URL string, logger newlog.Logger, Service *service.Service) (newHandler *Websocket, err error) {
	conn, resp, err := websocket.DefaultDialer.Dial(URL, nil)
	if err != nil {
		if resp != nil {
			err = newerror.MakeError(newerror.CodeWebSocketUpgradeFailed, err, newerror.LevelError).SetHttpError(resp.StatusCode, "websocket dial error")
			resp.Body.Close()
			return nil, err
		}
		return nil, newerror.MakeError(newerror.CodeWebSocketUpgradeFailed, err, newerror.LevelError).SetHttpError(-1, "tcp connect error")
	}
	newlog.SetLabel("Lack Config To limit Channel Capacity")
	ch := make(chan model.Message, 100)
	service.SetOperatorHub(ch)
	return &Websocket{
		logger:  logger,
		conn:    conn,
		SendCh:  ch,
		service: Service,
		Uuid:    uuid.NewV4(),
	}, nil
}

//	func (h *Websocket) SendCloseMessage() error {
//		err := h.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
//		if err != nil {
//
//		}
//	}
func (h *Websocket) Close() error {
	if h.cancel != nil {
		h.cancel()
	}
	return h.conn.Close()
}
func (h *Websocket) Run() {
	if h.conn == nil {
		newlog.SetLabel(newlog.LackLog)
	}
	if h.service == nil {
		newlog.SetLabel(newlog.LackLog)
	}
	if h.SendCh == nil {
		newlog.SetLabel(newlog.LackLog)
	}
	ctx, cancel := context.WithCancel(context.Background())
	h.cancel = cancel
	defer cancel()
	errCh := make(chan error, 100)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case logger := <-h.service.GetLogCh():
				middleware.WebEndLog(logger)
			}
		}
	}()
	go func() {
		var err error
		for {
			select {
			case <-ctx.Done():
				return
			case err = <-errCh:
				newlog.SetLabel(newlog.LackLog)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-h.SendCh:
				err := h.conn.WriteMessage(msg.MessageType, msg.Message)
				if err != nil {
					errCh <- newerror.MakeError(newerror.CodeWebSocketWriteError, h.Close(), newerror.LevelError)
					return
				}
			}
		}
	}()
	for {
		messageType, message, err := h.conn.ReadMessage()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				errCh <- newerror.MakeError(newerror.CodeWebSocketReadError, h.Close(), newerror.LevelError)
				return
			}
		}
		var oneBotResp *model.OneBotResp
		if err = json.Unmarshal(message, &oneBotResp); err != nil {
			errCh <- newerror.MakeError(newerror.CodeInvalidJSON, err, newerror.LevelError)
			continue
		}
		traceID := h.Uuid.String()
		logger := middleware.WebBeginLog(h.logger, traceID)
		if messageType == websocket.TextMessage {
			h.service.DistributeMessage(logger, traceID, oneBotResp)
		}
		if messageType == websocket.BinaryMessage {
			errCh <- newerror.MakeError(newerror.CodeUnsupportedMedia, errors.New("recive binary message"), newerror.LevelInfo)
		}

	}
}
