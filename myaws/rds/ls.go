package rds

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/pkg/errors"

	"github.com/minamijoyo/myaws/myaws"
)

// LsOptions customize the behavior of the Ls command.
type LsOptions struct {
	Quiet  bool
	Fields []string
}

// Ls describes RDSs.
func Ls(client *myaws.Client, options LsOptions) error {
	params := &rds.DescribeDBInstancesInput{}

	response, err := client.RDS.DescribeDBInstances(params)
	if err != nil {
		return errors.Wrap(err, "DescribeDBInstances failed:")
	}

	for _, db := range response.DBInstances {
		fmt.Println(formatDBInstance(db, options.Fields, options.Quiet))
	}

	return nil
}

func formatDBInstance(db *rds.DBInstance, fields []string, quiet bool) string {
	formatFuncs := map[string]func(db *rds.DBInstance) string{
		"DBInstanceClass":      formatDBInstanceClass,
		"Engine":               formatEngine,
		"AllocatedStorage":     formatAllocatedStorage,
		"StorageType":          formatStorageType,
		"StorageTypeIops":      formatStorageTypeIops,
		"DBInstanceIdentifier": formatDBInstanceIdentifier,
		"ReadReplicaSource":    formatReadReplicaSource,
		"InstanceCreateTime":   formatInstanceCreateTime,
	}

	var outputFields []string
	if quiet {
		outputFields = []string{"DBInstanceIdentifier"}
	} else {
		outputFields = fields
	}

	output := []string{}

	for _, field := range outputFields {
		value := formatFuncs[field](db)
		output = append(output, value)
	}

	return strings.Join(output[:], "\t")
}

func formatDBInstanceIdentifier(db *rds.DBInstance) string {
	return *db.DBInstanceIdentifier
}

func formatDBInstanceClass(db *rds.DBInstance) string {
	if *db.MultiAZ {
		return fmt.Sprintf("%s:multi", *db.DBInstanceClass)
	}
	return fmt.Sprintf("%s:single", *db.DBInstanceClass)
}

func formatEngine(db *rds.DBInstance) string {
	return fmt.Sprintf("%-15s", fmt.Sprintf("%s:%s", *db.Engine, *db.EngineVersion))
}

func formatAllocatedStorage(db *rds.DBInstance) string {
	return fmt.Sprintf("%4dGB", *db.AllocatedStorage)
}

func formatStorageType(db *rds.DBInstance) string {
	return *db.StorageType
}

func formatStorageTypeIops(db *rds.DBInstance) string {
	iops := "-"
	if db.Iops != nil {
		iops = fmt.Sprint(*db.Iops)
	}

	return fmt.Sprintf("%-8s", fmt.Sprintf("%s:%s", *db.StorageType, iops))
}

func formatReadReplicaSource(db *rds.DBInstance) string {
	if db.ReadReplicaSourceDBInstanceIdentifier == nil {
		return "source:---"
	}
	return fmt.Sprintf("source:%s", *db.ReadReplicaSourceDBInstanceIdentifier)
}

func formatInstanceCreateTime(db *rds.DBInstance) string {
	return myaws.FormatTime(db.InstanceCreateTime)
}
