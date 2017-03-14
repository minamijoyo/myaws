package myaws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

// FindEC2ReservedInstances return an array of reserved instances matching the conditions.
func (client *Client) FindEC2ReservedInstances(all bool) ([]*ec2.ReservedInstances, error) {
	params := &ec2.DescribeReservedInstancesInput{
		Filters: []*ec2.Filter{
			buildEC2RIStateFilter(all),
		},
	}

	response, err := client.EC2.DescribeReservedInstances(params)
	if err != nil {
		return nil, errors.Wrap(err, "DescribeReservedInstances failed")
	}

	var reservedInstances []*ec2.ReservedInstances
	for _, ri := range response.ReservedInstances {
		reservedInstances = append(reservedInstances, ri)
	}

	return reservedInstances, nil
}

func buildEC2RIStateFilter(all bool) *ec2.Filter {
	var stateFilter *ec2.Filter
	if !all {
		stateFilter = &ec2.Filter{
			Name: aws.String("state"),
			Values: []*string{
				aws.String("active"),
			},
		}
	}
	return stateFilter
}
