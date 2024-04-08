package powervs

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/power-go-client/power/models"
)

const (
	networkTypeName = "network"
)

// listDHCPNetworks lists previously found DHCP networks in found instances in the vpc.
func (o *ClusterUninstaller) listNetworks() (cloudResources, error) {
	var networks *models.Networks
	var network *models.NetworkReference
	var err error

	o.Logger.Debugf("Listing networks")

	if o.networkClient == nil {
		o.Logger.Infof("Skipping deleting networks because no service instance was found")
		result := []cloudResource{}
		return cloudResources{}.insert(result...), nil
	}

	networks, err = o.networkClient.GetAll()
	if err != nil {
		o.Logger.Fatalf("Failed to list networks: %v", err)
	}

	var foundOne = false

	result := []cloudResource{}
	for _, network = range networks.Networks {
		if network.Name == nil {
			o.Logger.Debugf("listNetworks: Network has empty Network.Name: %s", *network.NetworkID)
			continue
		}

		if strings.Contains(*network.Name, o.InfraID) {
			o.Logger.Debugf("listNetworks: FOUND: %s (%s)", *network.Name, *network.NetworkID)
			foundOne = true
			result = append(result, cloudResource{
				key:      *network.NetworkID,
				name:     *network.Name,
				status:   "",
				typeName: networkTypeName,
				id:       *network.NetworkID,
			})
		}
	}
	if !foundOne {
		o.Logger.Debugf("listNetworks: NO matching network found in:")
		for _, network = range networks.Networks {
			if network.NetworkID == nil {
				continue
			}
			if network.Name == nil {
				continue
			}
			o.Logger.Debugf("listNetworks: only found Network: %s", *network.Name)
		}
	}

	return cloudResources{}.insert(result...), nil
}

// func (o *ClusterUninstaller) destroyNetwork(item cloudResource) error {
// 	var err error

// 	_, err = o.dhcpClient.Get(item.id)
// 	if err != nil {
// 		o.deletePendingItems(item.typeName, []cloudResource{item})
// 		o.Logger.Infof("Deleted DHCP Network %q", item.name)
// 		return nil
// 	}

// 	o.Logger.Debugf("Deleting DHCP network %q", item.name)

// 	err = o.dhcpClient.Delete(item.id)
// 	if err != nil {
// 		o.Logger.Infof("Error: o.dhcpClient.Delete: %q", err)
// 		return err
// 	}

// 	o.deletePendingItems(item.typeName, []cloudResource{item})
// 	o.Logger.Infof("Deleted DHCP Network %q", item.name)

// 	return nil
// }

// destroyDHCPNetworks searches for DHCP networks that are in a previous list
// the cluster's infra ID.
func (o *ClusterUninstaller) destroyNetworks() error {
	firstPassList, err := o.listNetworks()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}
	fmt.Println(firstPassList.list())
	// items := o.insertPendingItems(networkTypeName, firstPassList.list())

	// ctx, cancel := o.contextWithTimeout()
	// defer cancel()

	// for _, item := range items {
	// 	select {
	// 	case <-ctx.Done():
	// 		o.Logger.Debugf("destroyNetworks: case <-ctx.Done()")
	// 		return o.Context.Err() // we're cancelled, abort
	// 	default:
	// 	}

	// 	backoff := wait.Backoff{
	// 		Duration: 15 * time.Second,
	// 		Factor:   1.1,
	// 		Cap:      leftInContext(ctx),
	// 		Steps:    math.MaxInt32}
	// 	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
	// 		err2 := o.destroyNetwork(item)
	// 		if err2 == nil {
	// 			return true, err2
	// 		}
	// 		o.errorTracker.suppressWarning(item.key, err2, o.Logger)
	// 		return false, err2
	// 	})
	// 	if err != nil {
	// 		o.Logger.Fatal("destroyNetworks: ExponentialBackoffWithContext (destroy) returns ", err)
	// 	}
	// }

	// if items = o.getPendingItems(networkTypeName); len(items) > 0 {
	// 	for _, item := range items {
	// 		o.Logger.Debugf("destroyNetworks: found %s in pending items", item.name)
	// 	}
	// 	return fmt.Errorf("destroyNetworks: %d undeleted items pending", len(items))
	// }

	// backoff := wait.Backoff{
	// 	Duration: 15 * time.Second,
	// 	Factor:   1.1,
	// 	Cap:      leftInContext(ctx),
	// 	Steps:    math.MaxInt32}
	// err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
	// 	secondPassList, err2 := o.listDHCPNetworks()
	// 	if err2 != nil {
	// 		return false, err2
	// 	}
	// 	if len(secondPassList) == 0 {
	// 		// We finally don't see any remaining instances!
	// 		return true, nil
	// 	}
	// 	for _, item := range secondPassList {
	// 		o.Logger.Debugf("destroyDHCPNetworks: found %s in second pass", item.name)
	// 	}
	// 	return false, nil
	// })
	// if err != nil {
	// 	o.Logger.Fatal("destroyDHCPNetworks: ExponentialBackoffWithContext (list) returns ", err)
	// }

	return nil
}
