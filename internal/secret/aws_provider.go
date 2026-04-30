package secret

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// awsProvider implements Provider using AWS Secrets Manager.
type awsProvider struct {
	client *secretsmanager.Client
	prefix string
}

// NewAWSProvider creates a new AWS Secrets Manager provider.
// prefix is prepended to all secret names (e.g. "envchain/").
func NewAWSProvider(ctx context.Context, prefix string) (Provider, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("aws provider: load config: %w", err)
	}
	return &awsProvider{
		client: secretsmanager.NewFromConfig(cfg),
		prefix: prefix,
	}, nil
}

func (p *awsProvider) secretID(profile, key string) string {
	return fmt.Sprintf("%s%s/%s", p.prefix, profile, strings.ToLower(key))
}

// Set stores a key/value pair under the given profile as a JSON secret.
func (p *awsProvider) Set(ctx context.Context, profile, key, value string) error {
	id := p.secretID(profile, key)
	payload, _ := json.Marshal(map[string]string{"value": value})
	_, err := p.client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String(id),
		SecretString: aws.String(string(payload)),
	})
	if err != nil {
		// Attempt update if secret already exists.
		_, err = p.client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
			SecretId:     aws.String(id),
			SecretString: aws.String(string(payload)),
		})
		if err != nil {
			return fmt.Errorf("aws provider: set %q: %w", id, err)
		}
	}
	return nil
}

// Get retrieves the value for key under the given profile.
func (p *awsProvider) Get(ctx context.Context, profile, key string) (string, error) {
	id := p.secretID(profile, key)
	out, err := p.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(id),
	})
	if err != nil {
		return "", fmt.Errorf("aws provider: get %q: %w", id, err)
	}
	var payload map[string]string
	if err := json.Unmarshal([]byte(aws.ToString(out.SecretString)), &payload); err != nil {
		return "", fmt.Errorf("aws provider: decode %q: %w", id, err)
	}
	return payload["value"], nil
}

// Delete removes the secret for key under the given profile.
func (p *awsProvider) Delete(ctx context.Context, profile, key string) error {
	id := p.secretID(profile, key)
	_, err := p.client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
		SecretId:                   aws.String(id),
		ForceDeleteWithoutRecovery: aws.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("aws provider: delete %q: %w", id, err)
	}
	return nil
}
