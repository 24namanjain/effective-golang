package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

type EchoRequest struct {
	Message string `json:"message"`
}

type EchoResponse struct {
	Message   string `json:"message"`
	Echoed    bool   `json:"echoed"`
	Timestamp string `json:"timestamp"`
}

/*
InitModule is the entry point for the module.
It registers the RPC functions and initializes the module.
*/
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("Cricket Nakama module initialized successfully!")

	if error := initializer.RegisterRpc("echo", rpcEcho); error != nil {
		logger.Error("Failed to register RPC: %v", error)
		return error
	}

	logger.Info("Nakama Practice module intialized successfully")
	return nil
}

/*
rpcEcho is the function that is called when the client calls the echo RPC.
It echoes the message back to the client.
*/
func rpcEcho(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Info("Echo RPC called with the payload: %s", payload)

	var request EchoRequest
	if error := json.Unmarshal([]byte(payload), &request); error != nil {
		logger.Error("Failed to unmarshal request: %v", error)
		return "", error
	}

	// Create Response
	response := EchoResponse{
		Message:   request.Message,
		Echoed:    true,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Convert response to JSON
	responseBytes, err := json.Marshal(response)
	if err != nil {
		logger.Error("Failed to marshal response: %v", err)
		return "", err
	}

	return string(responseBytes), nil
}
