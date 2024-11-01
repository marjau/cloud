package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const dynamodbSvcName = "DynamoDB"

type dynamoDBTester struct {
	awsTester
	svc *dynamodb.Client
}

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func newDynamoDBTester(cfg aws.Config) *dynamoDBTester {
	return &dynamoDBTester{
		awsTester: awsTester{
			cfg:    cfg,
			logger: newPrefixedLogger(dynamodbSvcName),
		},
		svc: dynamodb.NewFromConfig(cfg),
	}
}

func (dt *dynamoDBTester) RunTests() error {
	dt.logger.log("Running Testing...")
	tableName := "TestTable"
	item := Item{ID: "123", Name: "Sample Item"}

	if err := dt.createTable(tableName); err != nil {
		return err
	}
	defer dt.clean(tableName)

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

	dt.logger.log("Testing Completed")
	return nil
}

func (dt *dynamoDBTester) clean(tableName string) {
	dt.logger.log("Start cleaning tests...")
	if err := dt.deleteTable(tableName); err != nil {
		var notFoundErr *types.ResourceNotFoundException

		if errors.As(err, &notFoundErr) {
			dt.logger.logf("INFO: Table %v not found", tableName)
		} else {
			dt.logger.logf("ERROR: %v", err)
		}
	}
	dt.logger.log("Tests cleaning completed.")
}

func (dt *dynamoDBTester) createTable(tableName string) error {
	dt.logger.logf("Creating %q Table", tableName)

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

	dt.logger.logf("Table %q created successfully", tableName)
	return nil
}

func (dt *dynamoDBTester) listTables() error {
	dt.logger.log("Listing Tables")

	result, err := dt.svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		return fmt.Errorf("failed to list tables: %w", err)
	}

	if len(result.TableNames) == 0 {
		dt.logger.log("INFO: No tables found.")
		return nil
	}

	dt.logger.log("Tables:")
	for i, tableName := range result.TableNames {
		dt.logger.logf("  %d. %s", i+1, tableName)
	}
	return nil
}

func (dt *dynamoDBTester) putItem(tableName string, item Item) error {
	dt.logger.logf("Putting Item %+v to %v Table", item, tableName)

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

	dt.logger.log("Item put successfully")
	return nil
}

func (dt *dynamoDBTester) getItem(tableName, itemID string) error {
	dt.logger.logf("Getting Item %v from Table: %v", itemID, tableName)

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
		dt.logger.log("Item not found")
		return nil
	}

	item := Item{}
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return fmt.Errorf("failed to unmarshal item: %w", err)
	}

	dt.logger.log(fmt.Sprintf("Retrieved item: %+v", item))
	return nil
}

func (dt *dynamoDBTester) updateItem(tableName, itemID, newName string) error {
	dt.logger.logf("Updating Item %v With New Name %q in Table: %v", itemID, newName, tableName)

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

	dt.logger.log("Item updated successfully")
	return nil
}

func (dt *dynamoDBTester) deleteItem(tableName, itemID string) error {
	dt.logger.logf("Deleting Item %v from Table: %v", itemID, tableName)

	_, err := dt.svc.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: itemID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete item from table %v: %w", tableName, err)
	}

	dt.logger.logf("Item %v deleted successfully", itemID)
	return nil
}

func (dt *dynamoDBTester) deleteTable(tableName string) error {
	dt.logger.logf("Deleting Table: %v", tableName)

	_, err := dt.svc.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete table %v: %w", tableName, err)
	}

	dt.logger.logf("Table %v deleted successfully", tableName)
	return nil
}
