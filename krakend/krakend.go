// Package krakend registers a bloomfilter given a config and registers the service with consul.
package krakend

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	redisbloom "github.com/RedisBloom/redisbloom-go"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
)

// Namespace for bloomfilter
const Namespace = "github_com/devopsfaith/bloomfilter"

var (
	errNoConfig    = errors.New("no config for the bloomfilter")
	errWrongConfig = errors.New("invalid config for the bloomfilter")
)

type Config struct {
	HashName string
	// Address   string
	TokenKeys []string
	Headers   []string
}

// Register registers a bloomfilter given a config and registers the service with consul
func Register(ctx context.Context, serviceName string, cfg config.ServiceConfig,
	logger logging.Logger, register func(n string, p int)) (Rejecter, error) {
	data, ok := cfg.ExtraConfig[Namespace]
	if !ok {
		logger.Debug(errNoConfig.Error())
		return nopRejecter, errNoConfig
	}

	raw, err := json.Marshal(data)
	if err != nil {
		logger.Debug(errWrongConfig.Error())
		return nopRejecter, errWrongConfig
	}

	var rpcConfig Config
	if err := json.Unmarshal(raw, &rpcConfig); err != nil {
		logger.Debug(err.Error(), string(raw))
		return nopRejecter, err
	}

	rejecter := Rejecter{
		TokenKeys: rpcConfig.TokenKeys,
		Headers:   rpcConfig.Headers,
		HashName:  rpcConfig.HashName,
	}

	redis_pass := os.Getenv("REDIS_PASSWORD")
	redis_ddress := os.Getenv("REDIS_ADDRESS")
	if redis_ddress == "" {
		return nopRejecter, &RedisAddressEmpyErr{}
	}
	if redis_pass == "" {
		rejecter.redis_client = redisbloom.NewClient(redis_ddress, serviceName, nil)
	} else {
		rejecter.redis_client = redisbloom.NewClient(redis_ddress, serviceName, &redis_pass)
	}

	return rejecter, nil
}

type Rejecter struct {
	HashName     string
	redis_client *redisbloom.Client
	TokenKeys    []string
	Headers      []string
}

func (r *Rejecter) RejectToken(claims map[string]interface{}) bool {
	for _, k := range r.TokenKeys {
		v, ok := claims[k]
		if !ok {
			continue
		}
		data, ok := v.(string)
		if !ok {
			continue
		}
		exists, _ := r.redis_client.CfExists(r.HashName, (k + "-" + data))
		if exists {
			return true
		}
	}
	return false
}

func (r *Rejecter) RejectHeader(header http.Header) bool {
	for _, k := range r.Headers {
		data := header.Get(k)
		if data == "" {
			continue
		}
		exists, _ := r.redis_client.CfExists(r.HashName, (k + "-" + data))
		if exists {
			return true
		}
	}
	return false
}

var nopRejecter = Rejecter{HashName: ""}

type RedisAddressEmpyErr struct {
}

func (n *RedisAddressEmpyErr) Error() string {
	return "The 'REDIS_ADDRESS' environment variable dose not exist!"
}
