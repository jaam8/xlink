package pool

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"slices"
	"sync"
	"time"
	"xlink/common/logger"
)

type GrpcPool struct {
	idle   []*grpc.ClientConn
	active []*grpc.ClientConn
	mu     sync.Mutex
	config Config
}

func NewConnection(config Config) (*grpc.ClientConn, error) {
	_, cancel := context.WithTimeout(context.Background(), config.DialTimeout)
	defer cancel()
	return grpc.NewClient(config.Address, config.DialOptions...)
}

func NewGrpcPool(ctx context.Context, config Config) (*GrpcPool, error) {
	pool := &GrpcPool{
		idle:   make([]*grpc.ClientConn, config.MaxConnections),
		active: make([]*grpc.ClientConn, 0, config.MaxConnections),
		mu:     sync.Mutex{},
		config: config,
	}

	if pool.config.DialTimeout == 0 {
		pool.config.DialTimeout = time.Second * 5
	}

	successfulConnections := 0
	for i := 0; i < config.MaxConnections; i++ {
		grpcConn, err := NewConnection(pool.config)
		if err != nil {
			logger.GetOrCreateLoggerFromCtx(ctx).
				Info(ctx, "couldn't connect to grpc",
					zap.String("address", config.Address), zap.Error(err))
		}
		successfulConnections++
		pool.idle[i] = grpcConn
	}

	if successfulConnections < config.MinConnections {
		_ = pool.Close() //nolint:all
		return nil, fmt.Errorf("couldn't create at least %d connections: got %d",
			config.MinConnections, successfulConnections)
	}

	return pool, nil
}

func (p *GrpcPool) GetConn() (*grpc.ClientConn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var err error
	var grpcConn *grpc.ClientConn

	if len(p.idle) > 0 {
		grpcConn, p.idle = p.idle[0], p.idle[1:]
	} else {
		grpcConn, err = NewConnection(p.config)
		if err != nil {
			return nil, fmt.Errorf("couldn't create new grpc connection at address %s: %v", p.config.Address, err)
		}
	}

	p.active = append(p.active, grpcConn)

	// warning, but not returning on pool exhaust
	if len(p.active) > p.config.MaxConnections {
		loggerCtx, loggerErr := logger.New(context.Background())
		if loggerErr == nil {
			logger.GetLoggerFromCtx(loggerCtx).
				Warn(loggerCtx, "grpc pool overload",
					zap.Int("MaxConnections", p.config.MaxConnections),
					zap.Int("ActualConnections", len(p.active)),
				)
		}
	}

	return grpcConn, nil
}

func (p *GrpcPool) Restore(grpcConn *grpc.ClientConn) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for ind, connPointer := range p.active {
		if connPointer == grpcConn {
			p.active = slices.Delete(p.active, ind, ind+1)

			// append to idle only if connection is alive to use it later
			if grpcConn.GetState() != connectivity.Shutdown {
				p.idle = append(p.idle, grpcConn)
			}

			_ = grpcConn.Close() //nolint:all
			return nil
		}
	}

	return fmt.Errorf("conn not found in []active connections")
}

func (p *GrpcPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var errs []error
	for _, grpcConn := range p.idle {
		if err := grpcConn.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	for _, grpcConn := range p.active {
		if err := grpcConn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	p.active = nil
	p.idle = nil

	if len(errs) > 0 {
		return fmt.Errorf("couldn't close all connections: %v", errs)
	}
	return nil
}
