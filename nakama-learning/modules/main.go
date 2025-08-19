package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	// Register RPC functions
	initializer.RegisterRpc("hello_world", HelloWorldRPC)
	initializer.RegisterRpc("get_time", GetTimeRPC)
	
	logger.Info("Simple Nakama Module initialized successfully!")
	return nil
}

func HelloWorldRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	
	response := map[string]interface{}{
		"message": "Hello from Nakama!",
		"user_id": userID,
		"timestamp": time.Now().Unix(),
	}
	
	responseBytes, _ := json.Marshal(response)
	return string(responseBytes), nil
}

func GetTimeRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	response := map[string]interface{}{
		"server_time": time.Now().Format(time.RFC3339),
		"unix_time": time.Now().Unix(),
	}
	
	responseBytes, _ := json.Marshal(response)
	return string(responseBytes), nil
}
