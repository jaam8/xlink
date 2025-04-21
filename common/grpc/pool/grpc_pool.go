package pool

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"slices"
	"sync"
	"xlink/common/logger"
)

type GrpcPool struct {
	idle   []*grpc.ClientConn
	active []*grpc.ClientConn
	mu     sync.Mutex
	config Config
}

func NewGrpcPool(ctx context.Context, config Config) (*GrpcPool, error) {
	pool := make([]*grpc.ClientConn, config.MaxConnections)

	connectionsLeft := config.MaxConnections
	for i := 0; i < config.MaxConnections; i++ {
		grpcConn, err := grpc.NewClient(config.Address, config.DialOptions...)
		if err != nil {
			logger.GetOrCreateLoggerFromCtx(ctx).
				Info(ctx, "couldn't connect to grpc",
					zap.String("address", config.Address), zap.Error(err))
		} else {
			connectionsLeft--
		}

		pool[i] = grpcConn
	}

	if connectionsLeft > config.MaxConnections-config.MinConnections {
		wg := sync.WaitGroup{}
		wg.Add(connectionsLeft)
		for _, grpcConn := range pool {
			go func() {
				defer wg.Done()
				if grpcConn != nil {
					grpcConn.Close()
				}
			}()
		}

		wg.Wait()
		return nil, fmt.Errorf("couldn't create at least %d connections: got %d",
			config.MinConnections, config.MaxConnections-connectionsLeft)
	}

	return &GrpcPool{
		idle:   pool,
		config: config,
	}, nil
}

func (p *GrpcPool) GetConn() (*grpc.ClientConn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var err error
	var grpcConn *grpc.ClientConn

	if len(p.idle) == 0 {
		grpcConn, err = grpc.NewClient(p.config.Address, p.config.DialOptions...)
		if err != nil {
			return nil, fmt.Errorf("couldn't connect to grpc at address %s: %v", p.config.Address, err)
		}
	} else {
		grpcConn, p.idle = p.idle[0], p.idle[1:]
	}

	p.active = append(p.active, grpcConn)
	return grpcConn, nil
}

func (p *GrpcPool) Restore(grpcConn *grpc.ClientConn) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for ind, connPointer := range p.active {
		if connPointer == grpcConn {
			p.active = slices.Delete(p.active, ind, ind+1)
			p.idle = append(p.idle, grpcConn)
			return nil
		}
	}

	return fmt.Errorf("couldn't restore connection: conn not found in []active")
}
