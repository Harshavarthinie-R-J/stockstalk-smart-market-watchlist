package graphql

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	gql "github.com/graphql-go/graphql"
	gqlhandler "github.com/graphql-go/handler"

	"stockstalk/internal/auth"
)

type Server struct {
	Schema   gql.Schema
	Auth     *auth.Service
	Resolver *Resolver
}

func NewServer(resolver *Resolver, authService *auth.Service) (*Server, error) {
	stringType := gql.String
	intType := gql.Int
	floatType := gql.Float
	idType := gql.ID

	// ------------------------------------------------------------
	// Instrument
	// ------------------------------------------------------------
	instrumentType := gql.NewObject(gql.ObjectConfig{
		Name: "Instrument",
		Fields: gql.Fields{
			"id":       &gql.Field{Type: idType},
			"symbol":   &gql.Field{Type: stringType},
			"name":     &gql.Field{Type: stringType},
			"exchange": &gql.Field{Type: stringType},
			"segment":  &gql.Field{Type: stringType},
			"isin":     &gql.Field{Type: stringType},
			"currency": &gql.Field{Type: stringType},
			"sector":   &gql.Field{Type: stringType},
			"industry": &gql.Field{Type: stringType},
		},
	})

	// ------------------------------------------------------------
	// Quote
	// ------------------------------------------------------------
	quoteType := gql.NewObject(gql.ObjectConfig{
		Name: "Quote",
		Fields: gql.Fields{
			"instrumentId":      &gql.Field{Type: idType},
			"symbol":            &gql.Field{Type: stringType},
			"price":             &gql.Field{Type: floatType},
			"previousClose":     &gql.Field{Type: floatType},
			"open":              &gql.Field{Type: floatType},
			"high":              &gql.Field{Type: floatType},
			"low":               &gql.Field{Type: floatType},
			"volume":            &gql.Field{Type: intType},
			"change":            &gql.Field{Type: floatType},
			"changePercent":     &gql.Field{Type: floatType},
			"week52High":        &gql.Field{Type: floatType},
			"week52Low":         &gql.Field{Type: floatType},
			"marketStatus":      &gql.Field{Type: stringType},
			"marketTimestamp":   &gql.Field{Type: stringType},
			"receivedTimestamp": &gql.Field{Type: stringType},
			"source":            &gql.Field{Type: stringType},

			"reliabilityState": &gql.Field{
				Type: stringType,
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(map[string]interface{})
					if !ok {
						return "UNAVAILABLE", nil
					}

					marketTimestamp, ok := source["marketTimestamp"].(string)
					if !ok || marketTimestamp == "" {
						return "UNAVAILABLE", nil
					}

					t, err := time.Parse(time.RFC3339, marketTimestamp)
					if err != nil {
						return "UNAVAILABLE", nil
					}

					age := time.Since(t)

					if age < 0 {
						return "CONFLICTING", nil
					}

					switch {
					case age <= 5*time.Minute:
						return "LIVE", nil
					case age <= 30*time.Minute:
						return "DELAYED", nil
					default:
						return "STALE", nil
					}
				},
			},
		},
	})

	// ------------------------------------------------------------
	// Market Data
	// ------------------------------------------------------------
	marketType := gql.NewObject(gql.ObjectConfig{
		Name: "MarketData",
		Fields: gql.Fields{
			"market":      &gql.Field{Type: stringType},
			"lastUpdated": &gql.Field{Type: stringType},
			"instruments": &gql.Field{Type: gql.NewList(instrumentType)},
			"quotes":      &gql.Field{Type: gql.NewList(quoteType)},
		},
	})

	// ------------------------------------------------------------
	// Watch Item
	// ------------------------------------------------------------
	watchItemType := gql.NewObject(gql.ObjectConfig{
		Name: "WatchItem",
		Fields: gql.Fields{
			"instrumentId": &gql.Field{Type: idType},
			"position":     &gql.Field{Type: intType},
			"addedAt":      &gql.Field{Type: stringType},
		},
	})

	// ------------------------------------------------------------
	// Watchlist
	// ------------------------------------------------------------
	watchlistType := gql.NewObject(gql.ObjectConfig{
		Name: "Watchlist",
		Fields: gql.Fields{
			"id":        &gql.Field{Type: idType},
			"name":      &gql.Field{Type: stringType},
			"createdAt": &gql.Field{Type: stringType},
			"updatedAt": &gql.Field{Type: stringType},
			"stocks":    &gql.Field{Type: gql.NewList(watchItemType)},
		},
	})

	// ------------------------------------------------------------
	// User
	// ------------------------------------------------------------
	userType := gql.NewObject(gql.ObjectConfig{
		Name: "User",
		Fields: gql.Fields{
			"id":        &gql.Field{Type: idType},
			"name":      &gql.Field{Type: stringType},
			"email":     &gql.Field{Type: stringType},
			"createdAt": &gql.Field{Type: stringType},
		},
	})

	// ------------------------------------------------------------
	// Checkpoint
	// ------------------------------------------------------------
	checkpointType := gql.NewObject(gql.ObjectConfig{
		Name: "Checkpoint",
		Fields: gql.Fields{
			"id":          &gql.Field{Type: idType},
			"watchlistId": &gql.Field{Type: idType},
			"createdAt":   &gql.Field{Type: stringType},
		},
	})

	// ------------------------------------------------------------
	// Market Change / Event
	// ------------------------------------------------------------
	changeType := gql.NewObject(gql.ObjectConfig{
		Name: "MarketChange",
		Fields: gql.Fields{
			"id":            &gql.Field{Type: idType},
			"instrumentId":  &gql.Field{Type: idType},
			"symbol":        &gql.Field{Type: stringType},
			"type":          &gql.Field{Type: stringType},
			"severity":      &gql.Field{Type: stringType},
			"previousValue": &gql.Field{Type: floatType},
			"currentValue":  &gql.Field{Type: floatType},
			"changePercent": &gql.Field{Type: floatType},

			"currentVolume":  &gql.Field{Type: intType},
			"averageVolume":  &gql.Field{Type: floatType},
			"volumeRatio":    &gql.Field{Type: floatType},
			"volumeAnomaly":  &gql.Field{Type: gql.Boolean},
			"priceAnomaly":   &gql.Field{Type: gql.Boolean},
			"near52WeekHigh": &gql.Field{Type: gql.Boolean},
			"near52WeekLow":  &gql.Field{Type: gql.Boolean},

			"signals":               &gql.Field{Type: gql.NewList(stringType)},
			"context":               &gql.Field{Type: stringType},
			"confidence":            &gql.Field{Type: stringType},
			"contextType":           &gql.Field{Type: stringType},
			"attribution":           &gql.Field{Type: stringType},
			"attributionConfidence": &gql.Field{Type: stringType},
			"attributionSource":     &gql.Field{Type: stringType},

			"attentionScore": &gql.Field{Type: intType},
			"detectedAt":     &gql.Field{Type: stringType},

			"reliabilityState": &gql.Field{
				Type: stringType,
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(map[string]interface{})
					if !ok {
						return "UNAVAILABLE", nil
					}

					value, ok := source["reliabilityState"].(string)
					if !ok || value == "" {
						return "UNAVAILABLE", nil
					}

					return value, nil
				},
			},
		},
	})

	// ------------------------------------------------------------
	// Change Summary
	// ------------------------------------------------------------
	changeSummaryType := gql.NewObject(gql.ObjectConfig{
		Name: "ChangeSummary",
		Fields: gql.Fields{
			"watchlistId":        &gql.Field{Type: idType},
			"previousCheckpoint": &gql.Field{Type: stringType},
			"currentTime":        &gql.Field{Type: stringType},
			"totalStocks":        &gql.Field{Type: intType},
			"meaningfulChanges":  &gql.Field{Type: intType},
			"notableChanges":     &gql.Field{Type: intType},
			"unchanged":          &gql.Field{Type: intType},
			"changes":            &gql.Field{Type: gql.NewList(changeType)},
		},
	})

	// ------------------------------------------------------------
	// Authentication Payload
	// ------------------------------------------------------------
	authType := gql.NewObject(gql.ObjectConfig{
		Name: "AuthPayload",
		Fields: gql.Fields{
			"user":  &gql.Field{Type: userType},
			"token": &gql.Field{Type: stringType},
		},
	})

	// ============================================================
	// QUERY
	// ============================================================
	query := gql.NewObject(gql.ObjectConfig{
		Name: "Query",
		Fields: gql.Fields{

			// -------------------------
			// Current User
			// -------------------------
			"me": &gql.Field{
				Type: gql.NewNonNull(userType),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return resolver.Me(withContext(p.Context))
				},
			},

			// -------------------------
			// Market
			// -------------------------
			"market": &gql.Field{
				Type: marketType,
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return resolver.Market(p.Context)
				},
			},

			// -------------------------
			// Quote
			// -------------------------
			"quote": &gql.Field{
				Type: quoteType,
				Args: gql.FieldConfigArgument{
					"symbol": &gql.ArgumentConfig{
						Type: gql.NewNonNull(stringType),
					},
				},
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return resolver.Quote(p.Context, p.Args)
				},
			},

			// -------------------------
			// Instrument Search
			// -------------------------
			"searchInstruments": &gql.Field{
				Type: gql.NewList(instrumentType),
				Args: gql.FieldConfigArgument{
					"query": &gql.ArgumentConfig{
						Type: gql.NewNonNull(stringType),
					},
				},
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return resolver.SearchInstruments(p.Context, p.Args)
				},
			},

			// -------------------------
			// Watchlists
			// -------------------------
			"watchlists": &gql.Field{
				Type: gql.NewList(watchlistType),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return resolver.Watchlists(p.Context)
				},
			},

			// -------------------------
			// Single Watchlist
			// -------------------------
			"watchlist": &gql.Field{
				Type: watchlistType,
				Args: gql.FieldConfigArgument{
					"id": &gql.ArgumentConfig{
						Type: gql.NewNonNull(idType),
					},
				},
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return resolver.Watchlist(p.Context, p.Args)
				},
			},

			// -------------------------
			// Changes Since Last Visit
			// -------------------------
			"changesSinceLastVisit": &gql.Field{
				Type: changeSummaryType,
				Args: gql.FieldConfigArgument{
					"watchlistId": &gql.ArgumentConfig{
						Type: gql.NewNonNull(idType),
					},
				},
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return resolver.ChangesSinceLastVisit(p.Context, p.Args)
				},
			},

			// -------------------------
			// Event History
			// -------------------------
			"eventHistory": &gql.Field{
				Type: gql.NewList(changeType),
				Args: gql.FieldConfigArgument{
					"instrumentId": &gql.ArgumentConfig{
						Type: gql.NewNonNull(idType),
					},
					"limit": &gql.ArgumentConfig{
						Type: gql.Int,
					},
				},
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return resolver.EventHistory(p.Context, p.Args)
				},
			},
		},
	})

	// ============================================================
	// MUTATION
	// ============================================================
	mutation := gql.NewObject(gql.ObjectConfig{
		Name: "Mutation",
		Fields: gql.Fields{

			// -------------------------
			// Register
			// -------------------------
			"register": &gql.Field{
				Type: authType,
				Args: args("name", "email", "password"),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return resolver.Register(p.Context, p.Args)
				},
			},

			// -------------------------
			// Login
			// -------------------------
			"login": &gql.Field{
				Type: authType,
				Args: args("email", "password"),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return resolver.Login(p.Context, p.Args)
				},
			},

			// -------------------------
			// Create Watchlist
			// -------------------------
			"createWatchlist": &gql.Field{
				Type: watchlistType,
				Args: args("name"),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return resolver.CreateWatchlist(p.Context, p.Args)
				},
			},

			// -------------------------
			// Add Stock
			// -------------------------
			"addStock": &gql.Field{
				Type: gql.Boolean,
				Args: args("watchlistId", "instrumentId"),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return resolver.AddStock(p.Context, p.Args)
				},
			},

			// -------------------------
			// Remove Stock
			// -------------------------
			"removeStock": &gql.Field{
				Type: gql.Boolean,
				Args: args("watchlistId", "instrumentId"),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return resolver.RemoveStock(p.Context, p.Args)
				},
			},

			// -------------------------
			// Create Checkpoint
			// -------------------------
			"createCheckpoint": &gql.Field{
				Type: gql.NewNonNull(checkpointType),
				Args: args("watchlistId"),
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					return resolver.CreateCheckpoint(p.Context, p.Args)
				},
			},
		},
	})

	// ============================================================
	// SUBSCRIPTION
	// ============================================================
	subscription := gql.NewObject(gql.ObjectConfig{
		Name: "Subscription",
		Fields: gql.Fields{

			"meaningfulChange": &gql.Field{
				Type: changeType,

				Args: gql.FieldConfigArgument{
					"watchlistId": &gql.ArgumentConfig{
						Type: gql.NewNonNull(idType),
					},
				},

				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					userID, ok := auth.UserID(p.Context)
					if !ok {
						return nil, errors.New("authentication required")
					}

					watchlistID, ok := p.Args["watchlistId"].(string)
					if !ok || watchlistID == "" {
						return nil, errors.New("watchlistId is required")
					}

					if resolver.Broker == nil {
						return nil, errors.New("event broker is not initialized")
					}

					return resolver.Broker.Subscribe(
						p.Context,
						userID,
						watchlistID,
					), nil
				},
			},
		},
	})

	// ============================================================
	// BUILD GRAPHQL SCHEMA
	// ============================================================
	schema, err := gql.NewSchema(gql.SchemaConfig{
		Query:        query,
		Mutation:     mutation,
		Subscription: subscription,
	})

	if err != nil {
		return nil, err
	}

	return &Server{
		Schema:   schema,
		Auth:     authService,
		Resolver: resolver,
	}, nil
}

// ================================================================
// GRAPHQL ARGUMENT HELPER
// ================================================================

func args(names ...string) gql.FieldConfigArgument {
	result := gql.FieldConfigArgument{}

	for _, name := range names {
		t := gql.String

		if name == "watchlistId" || name == "instrumentId" {
			t = gql.ID
		}

		result[name] = &gql.ArgumentConfig{
			Type: gql.NewNonNull(t),
		}
	}

	return result
}

// ================================================================
// CONTEXT HELPER
// ================================================================

func withContext(ctx context.Context) context.Context {
	return ctx
}

// ================================================================
// HTTP HANDLER
// ================================================================

func (s *Server) Handler() http.Handler {
	return corsMiddleware(
		authMiddleware(
			s.Auth,
			gqlhandler.New(
				&gqlhandler.Config{
					Schema:   &s.Schema,
					Pretty:   true,
					GraphiQL: true,
				},
			),
		),
	)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")

		if origin == "http://localhost:5173" ||
			origin == "https://stockstalk-frontend.onrender.com" {

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization",
			)
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ================================================================
// AUTHENTICATION MIDDLEWARE
// ================================================================

func authMiddleware(service *auth.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")

		if len(header) > 7 && header[:7] == "Bearer " {
			if id, err := service.Parse(header[7:]); err == nil {
				r = r.WithContext(
					auth.WithUserID(r.Context(), id),
				)
			}
		}

		next.ServeHTTP(w, r)
	})
}

// ================================================================
// TOKEN PAYLOAD HELPER
// ================================================================

func parseTokenPayload(tokenString string) (map[string]interface{}, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}

	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}

	return data, nil
}
