package myaws

// ECSStatusOptions customize the behavior of the Ls command.
type ECSStatusOptions struct {
	Cluster string
}

// ECSStatus prints ECS status.
func (client *Client) ECSStatus(options ECSStatusOptions) error {
	return client.printECSStatus(options.Cluster)
}
