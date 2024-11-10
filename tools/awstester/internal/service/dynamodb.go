package service

import (
	"context"
	"errors"
	"fmt"

	"myapps/awstester/internal/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const dynamodbSvcName = "DynamoDB"

type DynamoDBTester struct {
	awsTester
	svc *dynamodb.Client
}

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewDynamoDBTester(cfg aws.Config) *DynamoDBTester {
	return &DynamoDBTester{
		awsTester: awsTester{
			cfg:    cfg,
			logger: logger.NewPrefixedLogger(dynamodbSvcName),
		},
		svc: dynamodb.NewFromConfig(cfg),
	}
}

func (dt DynamoDBTester) GetName() string {
	return dynamodbSvcName
}

func (dt DynamoDBTester) Run() error {
	dt.logger.Log("Running Testing...")
	tableName := "TestTable"
	item := Item{ID: "123", Name: "Sample Item"}
	defer dt.clean(tableName)

	if err := dt.createTable(tableName); err != nil {
		return err
	}
	if err := dt.listTables(); err != nil {
		return err
	}
	if err := dt.putItem(tableName, item); err != nil {
		return err
	}
	if err := dt.getItem(tableName, item.ID); err != nil {
		return err
	}
	if err := dt.updateItem(tableName, item.ID, "Sample2"); err != nil {
		return err
	}
	if err := dt.deleteItem(tableName, item.ID); err != nil {
		return err
	}
	if err := dt.deleteTable(tableName); err != nil {
		return err
	}

	dt.logger.Log("Testing Completed")
	return nil
}

func (dt DynamoDBTester) clean(tableName string) {
	dt.logger.Log("Start cleaning tests...")
	if err := dt.deleteTable(tableName); err != nil {
		var notFoundErr *types.ResourceNotFoundException

		if errors.As(err, &notFoundErr) {
			dt.logger.Logf("INFO: Table %v not found", tableName)
		} else {
			dt.logger.Logf("ERROR: %v", err)
		}
	}
	dt.logger.Log("Tests cleaning completed.")
}

func (dt DynamoDBTester) createTable(tableName string) error {
	dt.logger.Logf("Creating %q Table", tableName)

	_, err := dt.svc.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       types.KeyTypeHash,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		return fmt.Errorf("failed to create table %v: %w", tableName, err)
	}

	dt.logger.Logf("Table %q created successfully", tableName)
	return nil
}

func (dt DynamoDBTester) listTables() error {
	dt.logger.Log("Listing Tables")

	result, err := dt.svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		return fmt.Errorf("failed to list tables: %w", err)
	}

	if len(result.TableNames) == 0 {
		dt.logger.Log("INFO: No tables found.")
		return nil
	}

	dt.logger.Log("Tables:")
	for i, tableName := range result.TableNames {
		dt.logger.Logf("  %d. %s", i+1, tableName)
	}
	return nil
}

func (dt DynamoDBTester) putItem(tableName string, item Item) error {
	dt.logger.Logf("Putting Item %+v to %v Table", item, tableName)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	_, err = dt.svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to put item in table %v: %w", tableName, err)
	}

	dt.logger.Log("Item put successfully")
	return nil
}

func (dt DynamoDBTester) getItem(tableName, itemID string) error {
	dt.logger.Logf("Getting Item %v from Table: %v", itemID, tableName)

	result, err := dt.svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: itemID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get item from table %v: %w", tableName, err)
	}

	if result.Item == nil {
		dt.logger.Log("Item not found")
		return nil
	}

	item := Item{}
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return fmt.Errorf("failed to unmarshal item: %w", err)
	}

	dt.logger.Log(fmt.Sprintf("Retrieved item: %+v", item))
	return nil
}

func (dt DynamoDBTester) updateItem(tableName, itemID, newName string) error {
	dt.logger.Logf("Updating Item %v With New Name %q in Table: %v", itemID, newName, tableName)

	_, err := dt.svc.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: itemID},
		},
		UpdateExpression: aws.String("set #n = :name"),
		ExpressionAttributeNames: map[string]string{
			"#n": "Name",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name": &types.AttributeValueMemberS{Value: newName},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	})
	if err != nil {
		return fmt.Errorf("failed to update item in table %v: %w", tableName, err)
	}

	dt.logger.Log("Item updated successfully")
	return nil
}

func (dt DynamoDBTester) deleteItem(tableName, itemID string) error {
	dt.logger.Logf("Deleting Item %v from Table: %v", itemID, tableName)

	_, err := dt.svc.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: itemID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete item from table %v: %w", tableName, err)
	}

	dt.logger.Logf("Item %v deleted successfully", itemID)
	return nil
}

func (dt DynamoDBTester) deleteTable(tableName string) error {
	dt.logger.Logf("Deleting Table: %v", tableName)

	_, err := dt.svc.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete table %v: %w", tableName, err)
	}

	dt.logger.Logf("Table %v deleted successfully", tableName)
	return nil
}
