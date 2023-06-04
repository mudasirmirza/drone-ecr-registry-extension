// Description: Drone plugin to provide ECR credentials to the Docker daemon.
// Author: Mudasir Mirza <github.com/mudasirmirza>

package main

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/registry"
	"github.com/sirupsen/logrus"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

// spec provides the plugin settings.
type spec struct {
	Bind      string `envconfig:"DRONE_BIND"`
	Debug     bool   `envconfig:"DRONE_DEBUG"`
	Secret    string `envconfig:"DRONE_SECRET"`
	ECRRegion string `envconfig:"DRONE_ECRREGION"`
}

type ECRPlugin struct {
	spec *spec
}

func (p *ECRPlugin) List(ctx context.Context, req *registry.Request) ([]*drone.Registry, error) {
	// Create a new session using the AWS SDK
	logrus.Infoln("Creating a new session using the AWS SDK")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(p.spec.ECRRegion),
	})
	if err != nil {
		logrus.Errorln("Failed to create a new session:", err)
		return nil, err
	}

	// Create a new ECR service client
	logrus.Infoln("Creating a new ECR service client")
	svc := ecr.New(sess)

	// Call the GetAuthorizationToken operation
	logrus.Infoln("Calling the GetAuthorizationToken operation")
	output, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		logrus.Errorln("Failed to call the GetAuthorizationToken operation:", err)
		return nil, err
	}

	// Parse the authorization data
	logrus.Infoln("Parsing the authorization data")
	data := output.AuthorizationData[0]
	logrus.Infoln("Authorization data received")
	token, err := base64.StdEncoding.DecodeString(*data.AuthorizationToken)
	if err != nil {
		logrus.Errorln("Failed to decode authorization token:", err)
		return nil, err
	}
	logrus.Infoln("Decoded authorization token successfully")
	parts := strings.SplitN(string(token), ":", 2)

	logrus.Infoln("Returning the registry credentials")
	// Return the registry credentials
	return []*drone.Registry{{
		Address:  *data.ProxyEndpoint,
		Username: parts[0],
		Password: parts[1],
	}}, nil
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("Missing secret key, please set DRONE_SECRET")
	}
	if spec.Bind == "" {
		logrus.Infoln("Missing bind address, defaulting to :3000")
		spec.Bind = ":3000"
	}
	if spec.ECRRegion == "" {
		logrus.Infoln("ECR region not set, defaulting to us-east-1")
		spec.ECRRegion = "us-east-1"
	}

	handler := registry.Handler(
		spec.Secret,
		&ECRPlugin{
			spec: spec,
		},
		logrus.StandardLogger(),
	)

	logrus.Infoln("Server listening on address %s", spec.Bind)

	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}
