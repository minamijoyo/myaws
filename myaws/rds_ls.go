package myaws

import (
	"fmt"
	"strings"

	"encoding/json"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/pkg/errors"
)

// RDSLsOptions customize the behavior of the Ls command.
type RDSLsOptions struct {
	Quiet  bool
	Fields []string
}

// RDSLs describes RDSs.
func (client *Client) RDSLs(options RDSLsOptions) error {
	params := &rds.DescribeDBInstancesInput{}

	response, err := client.RDS.DescribeDBInstances(params)
	if err != nil {
		return errors.Wrap(err, "DescribeDBInstances failed:")
	}

	switch client.format {
	case "json":
		fmt.Fprintln(client.stdout, formatJSONDBInstances(client, options, response.DBInstances))
	default:
		for _, db := range response.DBInstances {
			fmt.Fprintln(client.stdout, formatTsvDBInstance(client, options, db))
		}
	}

	return nil
}

func dbInstanceValues(client *Client, options RDSLsOptions, db *rds.DBInstance) map[string]string {
	formatFuncs := map[string]func(client *Client, options RDSLsOptions, db *rds.DBInstance) string{
		"DBInstanceClass":      formatRDSDBInstanceClass,
		"Engine":               formatRDSEngine,
		"AllocatedStorage":     formatRDSAllocatedStorage,
		"StorageType":          formatRDSStorageType,
		"StorageTypeIops":      formatRDSStorageTypeIops,
		"DBInstanceIdentifier": formatRDSDBInstanceIdentifier,
		"ReadReplicaSource":    formatRDSReadReplicaSource,
		"InstanceCreateTime":   formatRDSInstanceCreateTime,
	}

	values := map[string]string{}

	for _, field := range options.Fields {
		value := formatFuncs[field](client, options, db)
		values[field] = value
	}

	return values
}

func formatTsvDBInstance(client *Client, options RDSLsOptions, db *rds.DBInstance) string {
	values := dbInstanceValues(client, options, db)

	var outputFields []string
	if options.Quiet {
		outputFields = []string{"DBInstanceIdentifier"}
	} else {
		outputFields = options.Fields
	}

	output := []string{}
	for _, field := range outputFields {
		output = append(output, values[field])
	}

	return strings.Join(output[:], "\t")
}

func formatJSONDBInstances(client *Client, options RDSLsOptions, dbInstances []*rds.DBInstance) string {
	outputs := []map[string]string{}
	for _, db := range dbInstances {
		outputs = append(outputs, dbInstanceValues(client, options, db))
	}

	bytes, err := json.Marshal(outputs)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", bytes)
}

func formatRDSDBInstanceIdentifier(client *Client, options RDSLsOptions, db *rds.DBInstance) string {
	return *db.DBInstanceIdentifier
}

func formatRDSDBInstanceClass(client *Client, options RDSLsOptions, db *rds.DBInstance) string {
	if *db.MultiAZ {
		return fmt.Sprintf("%s:multi", *db.DBInstanceClass)
	}
	return fmt.Sprintf("%s:single", *db.DBInstanceClass)
}

func formatRDSEngine(client *Client, options RDSLsOptions, db *rds.DBInstance) string {
	engine := fmt.Sprintf("%s:%s", *db.Engine, *db.EngineVersion)
	if client.format == "json" {
		return engine
	}
	return fmt.Sprintf("%-15s", engine)
}

func formatRDSAllocatedStorage(client *Client, options RDSLsOptions, db *rds.DBInstance) string {
	if client.format == "json" {
		return fmt.Sprintf("%dGB", *db.AllocatedStorage)
	}
	return fmt.Sprintf("%4dGB", *db.AllocatedStorage)
}

func formatRDSStorageType(client *Client, options RDSLsOptions, db *rds.DBInstance) string {
	return *db.StorageType
}

func formatRDSStorageTypeIops(client *Client, options RDSLsOptions, db *rds.DBInstance) string {
	iops := "-"
	if db.Iops != nil {
		iops = fmt.Sprint(*db.Iops)
	}
	storage := fmt.Sprintf("%s:%s", *db.StorageType, iops)
	if client.format == "json" {
		return storage
	}

	return fmt.Sprintf("%-8s", storage)
}

func formatRDSReadReplicaSource(client *Client, options RDSLsOptions, db *rds.DBInstance) string {
	if db.ReadReplicaSourceDBInstanceIdentifier == nil {
		if client.format == "json" {
			return ""
		}
		return "source:---"
	}
	return fmt.Sprintf("source:%s", *db.ReadReplicaSourceDBInstanceIdentifier)
}

func formatRDSInstanceCreateTime(client *Client, options RDSLsOptions, db *rds.DBInstance) string {
	return client.FormatTime(db.InstanceCreateTime)
}
