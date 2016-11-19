package rds

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/minamijoyo/myaws/myaws"
)

// Ls describes RDSs.
func Ls(*cobra.Command, []string) error {
	client := newRDSClient()
	params := &rds.DescribeDBInstancesInput{}

	response, err := client.DescribeDBInstances(params)
	if err != nil {
		return errors.Wrap(err, "DescribeDBInstances failed:")
	}

	for _, db := range response.DBInstances {
		fmt.Println(formatDBInstance(db))
	}

	return nil
}

func formatDBInstance(db *rds.DBInstance) string {
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

	var fields []string
	if viper.GetBool("rds.ls.quiet") {
		fields = []string{"DBInstanceIdentifier"}
	} else {
		fields = viper.GetStringSlice("rds.ls.fields")
	}

	output := []string{}

	for _, field := range fields {
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
	} else {
		return fmt.Sprintf("%s:single", *db.DBInstanceClass)
	}
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
