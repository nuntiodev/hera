package initializer

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

const (
	EngineKubernetes = "kubernetes"
	EngineMemory     = "memory"
)

type Initializer interface {
	CreateSecrets(ctx context.Context) error
}

func New(zapLog *zap.Logger, engine string) (Initializer, error) {
	zapLog.Info("initializing system with encryption secrets and public/private keys")
	if engine == EngineKubernetes {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		tokenPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"
		if _, err := os.Stat(tokenPath); err != nil {
			return nil, err
		}
		config.BearerTokenFile = tokenPath
		clientSet, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}
		bytesNamespace, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		if err != nil {
			return nil, err
		}
		return &kubernetesInitializer{
			zapLog:    zapLog,
			k8s:       clientSet,
			redLog:    color.New(color.FgRed),
			blueLog:   color.New(color.FgBlue),
			namespace: string(bytesNamespace),
		}, nil
	} else if engine == EngineMemory {
		return &memoryInitializer{
			zapLog: zapLog,
			redLog: color.New(color.FgRed),
		}, nil
	}
	return nil, fmt.Errorf("invalid engine %s", engine)
}
