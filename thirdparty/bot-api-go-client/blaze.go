package bot

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/fox-one/pkg/uuid"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

const (
	keepAlivePeriod = 3 * time.Second
	writeWait       = 10 * time.Second
	pongWait        = 10 * time.Second
	pingPeriod      = (pongWait * 9) / 10

	createMessageAction = "CREATE_MESSAGE"
)

const (
	MessageCategoryPlainText             = "PLAIN_TEXT"
	MessageCategoryPlainImage            = "PLAIN_IMAGE"
	MessageCategoryPlainData             = "PLAIN_DATA"
	MessageCategoryPlainSticker          = "PLAIN_STICKER"
	MessageCategoryPlainLive             = "PLAIN_LIVE"
	MessageCategoryPlainContact          = "PLAIN_CONTACT"
	MessageCategorySystemConversation    = "SYSTEM_CONVERSATION"
	MessageCategorySystemAccountSnapshot = "SYSTEM_ACCOUNT_SNAPSHOT"
	MessageCategoryMessageRecall         = "MESSAGE_RECALL"
	MessageCategoryAppButtonGroup        = "APP_BUTTON_GROUP"
)

type BlazeMessage struct {
	Id     string                 `json:"id"`
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params,omitempty"`
	Data   json.RawMessage        `json:"data,omitempty"`
	Error  *Error                 `json:"error,omitempty"`
}

type MessageView struct {
	ConversationId   string    `json:"conversation_id"`
	UserId           string    `json:"user_id"`
	MessageId        string    `json:"message_id"`
	Category         string    `json:"category"`
	Data             string    `json:"data"`
	RepresentativeId string    `json:"representative_id"`
	QuoteMessageId   string    `json:"quote_message_id"`
	Status           string    `json:"status"`
	Source           string    `json:"source"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type TransferView struct {
	Type          string    `json:"type"`
	SnapshotId    string    `json:"snapshot_id"`
	CounterUserId string    `json:"counter_user_id"`
	AssetId       string    `json:"asset_id"`
	Amount        string    `json:"amount"`
	TraceId       string    `json:"trace_id"`
	Memo          string    `json:"memo"`
	CreatedAt     time.Time `json:"created_at"`
}

type systemConversationPayload struct {
	Action        string `json:"action"`
	ParticipantId string `json:"participant_id"`
	UserId        string `json:"user_id,omitempty"`
	Role          string `json:"role,omitempty"`
}

type BlazeClient struct {
	c *Credential
}

type BlazeListener interface {
	OnMessage(ctx context.Context, msg *MessageView, userId string) error
}

func NewBlazeClient(c *Credential) *BlazeClient {
	client := BlazeClient{
		c: c,
	}
	return &client
}

func (b *BlazeClient) Loop(ctx context.Context, listener BlazeListener) error {
	conn, err := connectMixinBlaze(b.c)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go tick(ctx, conn)

	if err = writeMessage(conn, "LIST_PENDING_MESSAGES", nil); err != nil {
		return fmt.Errorf("write LIST_PENDING_MESSAGES failed: %w", err)
	}

	var (
		blazeMessage BlazeMessage
		message      MessageView
	)

	for {
		typ, r, err := conn.NextReader()
		if err != nil {
			return err
		}

		if typ != websocket.BinaryMessage {
			continue
		}

		if err := parseBlazeMessage(r, &blazeMessage); err != nil {
			return err
		}

		if blazeMessage.Error != nil {
			return err
		}

		if blazeMessage.Action != createMessageAction {
			continue
		}

		if err := jsoniter.Unmarshal(blazeMessage.Data, &message); err != nil {
			return err
		}

		if err := listener.OnMessage(ctx, &message, b.c.uid); err != nil {
			return err
		}
	}
}

func connectMixinBlaze(c *Credential) (*websocket.Conn, error) {
	token, err := SignAuthenticationTokenByCredential(c, "GET", "/", "")
	if err != nil {
		return nil, err
	}
	header := make(http.Header)
	header.Add("Authorization", "Bearer "+token)
	u := url.URL{Scheme: "wss", Host: "mixin-blaze.zeromesh.net", Path: "/"}
	dialer := &websocket.Dialer{
		Subprotocols: []string{"Mixin-Blaze-1"},
	}
	conn, _, err := dialer.Dial(u.String(), header)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func tick(ctx context.Context, conn *websocket.Conn) error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return fmt.Errorf("write PING failed: %w", err)
			}
		}
	}
}

func writeMessage(coon *websocket.Conn, action string, params map[string]interface{}) error {
	id := uuid.New()
	blazeMessage, err := jsoniter.Marshal(BlazeMessage{Id: id, Action: action, Params: params})
	if err != nil {
		return err
	}

	if err := writeGzipToConn(coon, blazeMessage); err != nil {
		return err
	}

	return nil
}

func writeGzipToConn(conn *websocket.Conn, msg []byte) error {
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	wsWriter, err := conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}
	gzWriter, err := gzip.NewWriterLevel(wsWriter, 3)
	if err != nil {
		return err
	}
	if _, err := gzWriter.Write(msg); err != nil {
		return err
	}

	if err := gzWriter.Close(); err != nil {
		return err
	}
	return wsWriter.Close()
}

func parseBlazeMessage(r io.Reader, msg *BlazeMessage) error {
	gzReader, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	if err = jsoniter.NewDecoder(gzReader).Decode(msg); err != nil {
		return err
	}

	return nil
}
