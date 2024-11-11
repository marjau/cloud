package service

import (
	"context"
	"fmt"
	"myapps/awstester/internal/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const secretsManagerSvcName = "Secrets-Manager"

type SecretsManagerTester struct {
	awsTester
	svc *secretsmanager.Client
}

func NewSecretsManagerTester(cfg aws.Config) *SecretsManagerTester {
	return &SecretsManagerTester{
		awsTester: awsTester{
			cfg:    cfg,
			logger: logger.NewPrefixedLogger(secretsManagerSvcName),
		},
		svc: secretsmanager.NewFromConfig(cfg),
	}
}

func (SecretsManagerTester) GetName() string {
	return secretsManagerSvcName
}

func (sm SecretsManagerTester) Run() error {

	sm.logger.Log("Running Secrets Manager testing...")

	// Create a secret
	secretName := "TestSecret"
	secretValue := "MySecretValue"
	if err := sm.CreateSecret(secretName, secretValue); err != nil {
		return err
	}

	if err := sm.ListSecrets(); err != nil {
		return err
	}

	// Retrieve the secret
	_, err := sm.GetSecretValue(secretName)
	if err != nil {
		return err
	}

	// Delete the secret
	if err := sm.DeleteSecret(secretName); err != nil {
		return err
	}

	return nil
}

func (sm *SecretsManagerTester) ListSecrets() error {
	sm.logger.Log("Listing Secrets")

	// Use the ListSecrets method to retrieve all secrets
	output, err := sm.svc.ListSecrets(context.TODO(), &secretsmanager.ListSecretsInput{})
	if err != nil {
		return fmt.Errorf("failed to list secrets: %w", err)
	}

	// Display each secret's name
	for i, secret := range output.SecretList {
		sm.logger.Logf("  %d. Name: %s\n", i+1, aws.ToString(secret.Name))
	}
	return nil
}

func (sm *SecretsManagerTester) CreateSecret(secretName, secretValue string) error {
	sm.logger.Logf("Creating Secret: %v", secretName)

	_, err := sm.svc.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
		Name:         aws.String(secretName),
		SecretString: aws.String(secretValue),
	})
	if err != nil {
		return fmt.Errorf("failed to create secret %v: %w", secretName, err)
	}

	sm.logger.Logf("Secret %v created successfully", secretName)
	return nil
}

func (sm *SecretsManagerTester) GetSecretValue(secretName string) (string, error) {
	sm.logger.Logf("Retrieving Secret: %v", secretName)
	output, err := sm.svc.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get secret value for %v: %w", secretName, err)
	}
	sm.logger.Logf("Secret Value: %v", aws.ToString(output.SecretString))
	return aws.ToString(output.SecretString), nil
}

// func (sm *SecretsManagerTester) GetSecretWithCaching(secretName string) error {
// 	sm.logger.Logf("Retrieving Secret with Caching: %v", secretName)

// 	// Create a CacheConfig with the client and TTL settings
// 	// cacheConfig := &secretcache.CacheConfig{
// 	// 	Client:       svc,
// 	// 	CacheItemTTL: 5 * time.Minute, // Set TTL for cached items
// 	// }

// 	// Initialize the cache with a custom configuration
// 	// cache, err := secretcache.New(&secretcache.CacheConfig{
// 	// 	Client:       sm.svc,
// 	// 	CacheItemTTL: 5 * time.Minute, // Adjust TTL as needed
// 	// })
// 	// cache, err := secretcache.New(secretcache.Options{
// 	// 	CacheExpiration: 5 * time.Minute,                      // Set TTL as needed
// 	// 	SecretsManager:  secretsmanager.NewFromConfig(sm.cfg), // Pass the SecretsManager client
// 	// })

// 	cache, err := secretcache.New(func(c *secretcache.Cache) {
// 		c.Client = sm.svc
// 	})

// 	// Initialize the cache with a custom configuration
// 	// cache, err := secretcache.New(
// 	// 	secretcache.WithClient(svc),
// 	// 	secretcache.WithCacheItemTTL(5*time.Minute), // Set TTL for cached items
// 	// )
// 	if err != nil {
// 		sm.logger.Logf("failed to initialize secret cache: %v", err)
// 		return err
// 	}

// 	// Retrieve the secret with caching
// 	secret, err := cache.GetSecretString(secretName)
// 	if err != nil {
// 		return fmt.Errorf("failed to retrieve secret with caching: %w", err)
// 	}

// 	sm.logger.Logf("Retrieved secret value with caching: %v", secret)
// 	return nil
// }

func (sm *SecretsManagerTester) DeleteSecret(secretName string) error {
	sm.logger.Logf("Deleting Secret: %v", secretName)

	_, err := sm.svc.DeleteSecret(context.TODO(), &secretsmanager.DeleteSecretInput{
		SecretId:                   aws.String(secretName),
		ForceDeleteWithoutRecovery: aws.Bool(true), // Bypass recovery for immediate deletion
	})
	if err != nil {
		return fmt.Errorf("failed to delete secret %v: %w", secretName, err)
	}

	sm.logger.Log("Secret deleted successfully")
	return nil
}
