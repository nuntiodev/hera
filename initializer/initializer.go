package initializer

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"strings"
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
	redLog := color.New(color.FgRed)
	blueLog := color.New(color.FgBlue)
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
			redLog:    redLog,
			blueLog:   blueLog,
			namespace: string(bytesNamespace),
		}, nil
	} else if engine == EngineMemory {
		return &memoryInitializer{
			zapLog:  zapLog,
			redLog:  redLog,
			blueLog: blueLog,
		}, nil
	}
	if strings.TrimSpace(engine) == "" {
		redLog.Println("Hera is running without an engine, which means Hera will not automatically be able to generate any public/private keys or symmetric keys. Specify an engine in the environment.")
		return nil, errors.New("no engine specified")
	} else {
		redLog.Println(fmt.Sprintf("Hera is customized with an invalid engine: %s\n", engine))
		return nil, errors.New("invalid engine")
	}
	return nil, nil
}
