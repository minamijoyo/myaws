package myaws

import (
	"fmt"
	"strings"

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

	for _, db := range response.DBInstances {
		fmt.Println(formatDBInstance(client, options, db))
	}

	return nil
}

func formatDBInstance(client *Client, options RDSLsOptions, db *rds.DBInstance) string {
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

	var outputFields []string
	if options.Quiet {
		outputFields = []string{"DBInstanceIdentifier"}
	} else {
		outputFields = options.Fields
	}

	output := []string{}

	for _, field := range outputFields {
		value := formatFuncs[field](client, options, db)
		output = append(output, value)
	}

	return strings.Join(output[:], "\t")
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
	return fmt.Sprintf("%-15s", fmt.Sprintf("%s:%s", *db.Engine, *db.EngineVersion))
}

func formatRDSAllocatedStorage(client *Client, options RDSLsOptions, db *rds.DBInstance) string {
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

	return fmt.Sprintf("%-8s", fmt.Sprintf("%s:%s", *db.StorageType, iops))
}

func formatRDSReadReplicaSource(client *Client, options RDSLsOptions, db *rds.DBInstance) string {
	if db.ReadReplicaSourceDBInstanceIdentifier == nil {
		return "source:---"
	}
	return fmt.Sprintf("source:%s", *db.ReadReplicaSourceDBInstanceIdentifier)
}

func formatRDSInstanceCreateTime(client *Client, options RDSLsOptions, db *rds.DBInstance) string {
	return client.FormatTime(db.InstanceCreateTime)
}
