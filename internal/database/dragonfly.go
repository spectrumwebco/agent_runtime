package database

import (
	"context"
	"fmt"
	"strings"
	"time"
	
	"github.com/go-redis/redis/v8"
)

type DragonflyAdapter struct {
	connStr string
	client  *redis.Client
}

func NewDragonflyAdapter(connStr string) (*DragonflyAdapter, error) {
	var host, port, password string
	var db int
	
	host = "localhost"
	port = "6379"
	password = ""
	db = 0
	
	if strings.HasPrefix(connStr, "redis://") {
		connStr = strings.TrimPrefix(connStr, "redis://")
		
		if userPassEndIdx := strings.Index(connStr, "@"); userPassEndIdx != -1 {
			userPass := connStr[:userPassEndIdx]
			connStr = connStr[userPassEndIdx+1:]
			
			if passIdx := strings.Index(userPass, ":"); passIdx != -1 {
				password = userPass[passIdx+1:]
			}
		}
		
		if hostPortEndIdx := strings.Index(connStr, "/"); hostPortEndIdx != -1 {
			hostPort := connStr[:hostPortEndIdx]
			dbStr := connStr[hostPortEndIdx+1:]
			
			if dbStr != "" {
				fmt.Sscanf(dbStr, "%d", &db)
			}
			
			if portIdx := strings.Index(hostPort, ":"); portIdx != -1 {
				host = hostPort[:portIdx]
				port = hostPort[portIdx+1:]
			} else {
				host = hostPort
			}
		} else {
			if portIdx := strings.Index(connStr, ":"); portIdx != -1 {
				host = connStr[:portIdx]
				port = connStr[portIdx+1:]
			} else {
				host = connStr
			}
		}
	}
	
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to DragonflyDB: %w", err)
	}
	
	return &DragonflyAdapter{
		connStr: connStr,
		client:  client,
	}, nil
}

func (a *DragonflyAdapter) Query(ctx context.Context, command string) (interface{}, error) {
	args := strings.Fields(command)
	if len(args) == 0 {
		return nil, fmt.Errorf("empty command")
	}
	
	cmd := strings.ToUpper(args[0])
	
	switch cmd {
	case "GET":
		if len(args) < 2 {
			return nil, fmt.Errorf("GET command requires a key")
		}
		return a.client.Get(ctx, args[1]).Result()
		
	case "MGET":
		if len(args) < 2 {
			return nil, fmt.Errorf("MGET command requires at least one key")
		}
		return a.client.MGet(ctx, args[1:]...).Result()
		
	case "HGETALL":
		if len(args) < 2 {
			return nil, fmt.Errorf("HGETALL command requires a key")
		}
		return a.client.HGetAll(ctx, args[1]).Result()
		
	case "KEYS":
		if len(args) < 2 {
			return nil, fmt.Errorf("KEYS command requires a pattern")
		}
		return a.client.Keys(ctx, args[1]).Result()
		
	case "SCAN":
		if len(args) < 2 {
			return nil, fmt.Errorf("SCAN command requires a cursor")
		}
		
		var cursor uint64
		fmt.Sscanf(args[1], "%d", &cursor)
		
		var match string
		var count int64
		
		for i := 2; i < len(args); i++ {
			if strings.ToUpper(args[i]) == "MATCH" && i+1 < len(args) {
				match = args[i+1]
				i++
			} else if strings.ToUpper(args[i]) == "COUNT" && i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &count)
				i++
			}
		}
		
		return a.client.Scan(ctx, cursor, match, count).Result()
		
	case "LRANGE":
		if len(args) < 4 {
			return nil, fmt.Errorf("LRANGE command requires a key, start, and stop")
		}
		
		var start, stop int64
		fmt.Sscanf(args[2], "%d", &start)
		fmt.Sscanf(args[3], "%d", &stop)
		
		return a.client.LRange(ctx, args[1], start, stop).Result()
		
	case "ZRANGE":
		if len(args) < 4 {
			return nil, fmt.Errorf("ZRANGE command requires a key, start, and stop")
		}
		
		var start, stop int64
		fmt.Sscanf(args[2], "%d", &start)
		fmt.Sscanf(args[3], "%d", &stop)
		
		return a.client.ZRange(ctx, args[1], start, stop).Result()
		
	case "INFO":
		section := ""
		if len(args) > 1 {
			section = args[1]
		}
		return a.client.Info(ctx, section).Result()
		
	default:
		return nil, fmt.Errorf("unsupported command: %s", cmd)
	}
}

func (a *DragonflyAdapter) Execute(ctx context.Context, command string) (interface{}, error) {
	args := strings.Fields(command)
	if len(args) == 0 {
		return nil, fmt.Errorf("empty command")
	}
	
	cmd := strings.ToUpper(args[0])
	
	switch cmd {
	case "SET":
		if len(args) < 3 {
			return nil, fmt.Errorf("SET command requires a key and value")
		}
		
		var expiration time.Duration
		var setNX, setXX bool
		
		for i := 3; i < len(args); i++ {
			arg := strings.ToUpper(args[i])
			if arg == "EX" && i+1 < len(args) {
				var seconds int64
				fmt.Sscanf(args[i+1], "%d", &seconds)
				expiration = time.Duration(seconds) * time.Second
				i++
			} else if arg == "PX" && i+1 < len(args) {
				var milliseconds int64
				fmt.Sscanf(args[i+1], "%d", &milliseconds)
				expiration = time.Duration(milliseconds) * time.Millisecond
				i++
			} else if arg == "NX" {
				setNX = true
			} else if arg == "XX" {
				setXX = true
			}
		}
		
		if setNX {
			return a.client.SetNX(ctx, args[1], args[2], expiration).Result()
		} else if setXX {
			return a.client.SetXX(ctx, args[1], args[2], expiration).Result()
		} else {
			return a.client.Set(ctx, args[1], args[2], expiration).Result()
		}
		
	case "DEL":
		if len(args) < 2 {
			return nil, fmt.Errorf("DEL command requires at least one key")
		}
		return a.client.Del(ctx, args[1:]...).Result()
		
	case "HSET":
		if len(args) < 4 || len(args)%2 != 0 {
			return nil, fmt.Errorf("HSET command requires a key and at least one field-value pair")
		}
		
		key := args[1]
		fieldValues := make(map[string]interface{})
		
		for i := 2; i < len(args); i += 2 {
			fieldValues[args[i]] = args[i+1]
		}
		
		return a.client.HSet(ctx, key, fieldValues).Result()
		
	case "LPUSH":
		if len(args) < 3 {
			return nil, fmt.Errorf("LPUSH command requires a key and at least one value")
		}
		
		values := make([]interface{}, len(args)-2)
		for i := 2; i < len(args); i++ {
			values[i-2] = args[i]
		}
		
		return a.client.LPush(ctx, args[1], values...).Result()
		
	case "RPUSH":
		if len(args) < 3 {
			return nil, fmt.Errorf("RPUSH command requires a key and at least one value")
		}
		
		values := make([]interface{}, len(args)-2)
		for i := 2; i < len(args); i++ {
			values[i-2] = args[i]
		}
		
		return a.client.RPush(ctx, args[1], values...).Result()
		
	case "ZADD":
		if len(args) < 4 || (len(args)-2)%2 != 0 {
			return nil, fmt.Errorf("ZADD command requires a key and at least one score-member pair")
		}
		
		key := args[1]
		members := make([]*redis.Z, 0, (len(args)-2)/2)
		
		for i := 2; i < len(args); i += 2 {
			var score float64
			fmt.Sscanf(args[i], "%f", &score)
			members = append(members, &redis.Z{
				Score:  score,
				Member: args[i+1],
			})
		}
		
		return a.client.ZAdd(ctx, key, members...).Result()
		
	case "EXPIRE":
		if len(args) < 3 {
			return nil, fmt.Errorf("EXPIRE command requires a key and seconds")
		}
		
		var seconds int64
		fmt.Sscanf(args[2], "%d", &seconds)
		
		return a.client.Expire(ctx, args[1], time.Duration(seconds)*time.Second).Result()
		
	default:
		return nil, fmt.Errorf("unsupported command: %s", cmd)
	}
}

func (a *DragonflyAdapter) GetSchema(ctx context.Context) (interface{}, error) {
	info, err := a.client.Info(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get server info: %w", err)
	}
	
	dbSize, err := a.client.DBSize(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get database size: %w", err)
	}
	
	keys, err := a.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}
	
	if len(keys) > 1000 {
		keys = keys[:1000]
	}
	
	keyTypes := make(map[string]string)
	for _, key := range keys {
		keyType, err := a.client.Type(ctx, key).Result()
		if err != nil {
			continue
		}
		keyTypes[key] = keyType
	}
	
	return map[string]interface{}{
		"info":      info,
		"db_size":   dbSize,
		"key_count": len(keys),
		"key_types": keyTypes,
	}, nil
}

func (a *DragonflyAdapter) Close() error {
	return a.client.Close()
}
