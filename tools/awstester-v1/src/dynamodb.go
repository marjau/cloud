package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	testTableName = "Test"
)

var (
	testItem Item = Item{ID: "123", Name: "TestItem"}
)

type DynamoDBTester struct {
	tester
	svc *dynamodb.DynamoDB
}

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewDynamoDBTester(cfg AWSConfig, sess *session.Session) *DynamoDBTester {
	return &DynamoDBTester{
		tester: tester{
			config: cfg,
			logger: newPrefixedLogger("DynamoDB"),
		},
		svc: dynamodb.New(sess),
	}
}

func (dt DynamoDBTester) RunTests() error {
	dt.logger.log("Start testing...")
	if err := dt.listTables(); err != nil {
		return err
	}
	if err := dt.createTable(); err != nil {
		return err
	}
	// if err := dt.insertItem(); err != nil {
	// 	return err
	// }
	// if err := dt.readItem(); err != nil {
	// 	return err
	// }
	// if err := dt.updateItem(); err != nil {
	// 	return err
	// }
	// if err := dt.deleteItem(); err != nil {
	// 	return err
	// }
	if err := dt.DeleteTable(); err != nil {
		return err
	}

	dt.logger.log("Testing completed")
	return nil
}

func (dt DynamoDBTester) Clean() {
	dt.logger.log("Cleaning tests...")
	if err := dt.DeleteTable(); err != nil {
		dt.logger.logf("Warning: error deleting table: %v", err)
	}
	dt.logger.log("Tests cleaned")
}

func (dt DynamoDBTester) listTables() error {
	dt.logger.log("List Tables...")

	result, err := dt.svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return err
	}

	if len(result.TableNames) == 0 {
		dt.logger.log("No tables found.")
		return nil
	}

	dt.logger.logf("Tables(%d):", len(result.TableNames))
	for i, tableName := range result.TableNames {
		dt.logger.logf("  %v. %s", (i + 1), aws.StringValue(tableName))
	}
	return nil
}

func (dt DynamoDBTester) createTable() error {
	dt.logger.logf("Create %v Table...", testTableName)

	_, err := dt.svc.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(testTableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"), // Partition key
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"), // String type
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
	})
	if err != nil {
		return err
	}

	dt.logger.logf("Table %v created successfully", testTableName)
	return nil
}

func (dt DynamoDBTester) insertItem() error {
	dt.logger.logf("Insert %v Item...", testItem.Name)

	av, err := dynamodbattribute.MarshalMap(testItem)
	if err != nil {
		return err
	}

	_, err = dt.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(testTableName),
		Item:      av,
	})
	if err != nil {
		return err
	}
	dt.logger.logf("Item %v inserted successfully", testItem.Name)
	return nil
}

func (dt DynamoDBTester) readItem() error {
	dt.logger.logf("Read %v Item...", testItem.Name)

	result, err := dt.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(testTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(testItem.ID),
			},
		},
	})
	if err != nil {
		return err
	}

	if result.Item == nil {
		dt.logger.logf("Item %v not found", testItem.Name)
		return nil
	}

	item := Item{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return err
	}

	dt.logger.logf("Item retrieved: %+v\n", item)
	return nil
}

func (dt DynamoDBTester) updateItem() error {
	dt.logger.log("Update Item...")

	_, err := dt.svc.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(testTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String("123"),
			},
		},
		UpdateExpression: aws.String("set #n = :name"),
		ExpressionAttributeNames: map[string]*string{
			"#n": aws.String("Name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {
				S: aws.String("Updated Item Name"),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	})
	if err != nil {
		return err
	}

	dt.logger.log("Item updated successfully")
	return nil
}

func (dt DynamoDBTester) deleteItem() error {
	dt.logger.logf("Delete %v Item...", testItem.Name)

	_, err := dt.svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(testTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(testItem.ID),
			},
		},
	})
	if err != nil {
		return err
	}
	dt.logger.logf("Item %v deleted successfully", testItem.Name)
	return nil
}

func (dt DynamoDBTester) DeleteTable() error {
	dt.logger.logf("Delete %v Table...", testTableName)

	_, err := dt.svc.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String(testTableName),
	})
	if err != nil {
		return err
	}

	dt.logger.logf("Table %s deleted successfully\n", testTableName)
	return nil
}
