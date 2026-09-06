package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Development/hackathon configuration.
		// Restrict this to the frontend origin in production.
		return true
	},
}

type wsMessage struct {
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type wsSubscribePayload struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                 `json:"operationName,omitempty"`
}

type wsConnectionPayload struct {
	Authorization string `json:"authorization,omitempty"`
	Token         string `json:"token,omitempty"`
}

type wsResponse struct {
	ID      string      `json:"id,omitempty"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

type wsErrorPayload struct {
	Message string `json:"message"`
}

func (s *Server) WebSocketHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		// A WebSocket connection can have multiple goroutines writing to it.
		// Gorilla WebSocket does not allow concurrent writes, so every write
		// goes through this mutex.
		var writeMu sync.Mutex

		writeJSON := func(response wsResponse) error {
			writeMu.Lock()
			defer writeMu.Unlock()

			return conn.WriteJSON(response)
		}

		userID := ""

		// Support Authorization header when the WebSocket client can send it.
		if token := bearerToken(r.Header.Get("Authorization")); token != "" {
			userID, err = s.authenticateToken(ctx, token)
			if err != nil {
				_ = writeJSON(wsResponse{
					Type: "error",
					Payload: wsErrorPayload{
						Message: "authentication failed",
					},
				})
				return
			}
		}

		// subscription ID -> cancel function
		subscriptions := make(map[string]context.CancelFunc)

		defer func() {
			for _, cancelSubscription := range subscriptions {
				cancelSubscription()
			}
		}()

		for {
			var msg wsMessage

			if err := conn.ReadJSON(&msg); err != nil {
				return
			}

			switch msg.Type {

			// ---------------------------------------------------------
			// CONNECTION INITIALIZATION
			// ---------------------------------------------------------
			case "connection_init":
				var payload wsConnectionPayload

				if len(msg.Payload) > 0 {
					if err := json.Unmarshal(msg.Payload, &payload); err != nil {
						_ = writeJSON(wsResponse{
							Type: "error",
							Payload: wsErrorPayload{
								Message: "invalid connection_init payload",
							},
						})
						return
					}
				}

				// If authentication was not supplied through the HTTP
				// Authorization header, try the connection_init payload.
				if userID == "" {
					token := payload.Token

					if token == "" {
						token = bearerToken(payload.Authorization)
					}

					if token != "" {
						userID, err = s.authenticateToken(ctx, token)
						if err != nil {
							_ = writeJSON(wsResponse{
								Type: "error",
								Payload: wsErrorPayload{
									Message: "authentication failed",
								},
							})
							return
						}
					}
				}

				if userID == "" {
					_ = writeJSON(wsResponse{
						Type: "error",
						Payload: wsErrorPayload{
							Message: "authentication required",
						},
					})
					return
				}

				// Tell the client that the WebSocket connection is ready.
				if err := writeJSON(wsResponse{
					Type: "connection_ack",
				}); err != nil {
					return
				}

			// ---------------------------------------------------------
			// START SUBSCRIPTION
			// ---------------------------------------------------------
			case "subscribe":
				if userID == "" {
					_ = writeJSON(wsResponse{
						ID:   msg.ID,
						Type: "error",
						Payload: wsErrorPayload{
							Message: "authentication required",
						},
					})
					continue
				}

				var payload wsSubscribePayload

				if err := json.Unmarshal(msg.Payload, &payload); err != nil {
					_ = writeJSON(wsResponse{
						ID:   msg.ID,
						Type: "error",
						Payload: wsErrorPayload{
							Message: "invalid subscription payload",
						},
					})
					continue
				}

				// Extract watchlist ID from GraphQL variables.
				watchlistID, err := extractWatchlistID(payload)
				if err != nil {
					_ = writeJSON(wsResponse{
						ID:   msg.ID,
						Type: "error",
						Payload: wsErrorPayload{
							Message: err.Error(),
						},
					})
					continue
				}

				if s.Resolver == nil || s.Resolver.Broker == nil {
					_ = writeJSON(wsResponse{
						ID:   msg.ID,
						Type: "error",
						Payload: wsErrorPayload{
							Message: "event broker unavailable",
						},
					})
					continue
				}

				// If this subscription ID already exists, cancel it first.
				if oldCancel, ok := subscriptions[msg.ID]; ok {
					oldCancel()
					delete(subscriptions, msg.ID)
				}

				subCtx, subCancel := context.WithCancel(ctx)
				subscriptions[msg.ID] = subCancel

				eventCh := s.Resolver.Broker.Subscribe(
					subCtx,
					userID,
					watchlistID,
				)

				// Listen for events from the broker.
				go func(
					subID string,
					ch <-chan interface{},
					subContext context.Context,
				) {
					for {
						select {
						case <-subContext.Done():
							return

						case event, ok := <-ch:
							if !ok {
								return
							}

							err := writeJSON(wsResponse{
								ID:   subID,
								Type: "next",
								Payload: map[string]interface{}{
									"data": map[string]interface{}{
										"meaningfulChange": event,
									},
								},
							})

							if err != nil {
								// The connection is probably closed.
								return
							}
						}
					}
				}(msg.ID, eventCh, subCtx)

			// ---------------------------------------------------------
			// END SUBSCRIPTION
			// ---------------------------------------------------------
			case "complete":
				if cancelSubscription, ok := subscriptions[msg.ID]; ok {
					cancelSubscription()
					delete(subscriptions, msg.ID)
				}

			// ---------------------------------------------------------
			// PING
			// ---------------------------------------------------------
			case "ping":
				if err := writeJSON(wsResponse{
					Type: "pong",
				}); err != nil {
					return
				}

			// ---------------------------------------------------------
			// PONG
			// ---------------------------------------------------------
			case "pong":
				// Nothing to do.

			// ---------------------------------------------------------
			// UNKNOWN MESSAGE
			// ---------------------------------------------------------
			default:
				if err := writeJSON(wsResponse{
					ID:   msg.ID,
					Type: "error",
					Payload: wsErrorPayload{
						Message: "unsupported WebSocket message type",
					},
				}); err != nil {
					return
				}
			}
		}
	})
}

// authenticateToken validates the JWT and confirms that the user still exists.
func (s *Server) authenticateToken(
    ctx context.Context,
    token string,
) (string, error) {

    if s.Auth == nil {
        return "", errors.New("auth service unavailable")
    }

    userID, err := s.Auth.Parse(token)
    if err != nil {
        return "", errors.New("invalid token")
    }

    if s.Resolver == nil || s.Resolver.UserService == nil {
        return "", errors.New("user service unavailable")
    }

    if _, err := s.Resolver.UserService.GetByID(ctx, userID); err != nil {
        return "", errors.New("user not found")
    }

    return userID, nil
}
// bearerToken extracts a token from an Authorization value.
//
// Supported:
//   - "Bearer <token>"
//   - "<token>"
func bearerToken(value string) string {
	value = strings.TrimSpace(value)

	if value == "" {
		return ""
	}

	parts := strings.Fields(value)

	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}

	return value
}

// extractWatchlistID gets the watchlist ID from the GraphQL subscription
// variables.
//
// Expected client payload:
//
//	{
//	  "query": "subscription($watchlistId: ID!) { meaningfulChange(watchlistId: $watchlistId) { ... } }",
//	  "variables": {
//	    "watchlistId": "..."
//
//	  }
//	}
func extractWatchlistID(
	payload wsSubscribePayload,
) (string, error) {
	if payload.Variables != nil {
		if value, ok := payload.Variables["watchlistId"]; ok {
			if id, ok := value.(string); ok && id != "" {
				return id, nil
			}
		}
	}

	query := strings.TrimSpace(payload.Query)

	if query == "" {
		return "", errors.New("subscription query is required")
	}

	if !strings.Contains(query, "meaningfulChange") {
		return "", errors.New(
			"only meaningfulChange subscription is supported",
		)
	}

	return "", errors.New(
		"watchlistId variable is required",
	)
}
