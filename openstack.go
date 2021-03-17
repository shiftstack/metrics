package metrics

import (
	"github.com/gophercloud/gophercloud"
	baremetalVersions "github.com/gophercloud/gophercloud/openstack/baremetal/apiversions"
	blockstorageVersions "github.com/gophercloud/gophercloud/openstack/blockstorage/apiversions"
	computeVersions "github.com/gophercloud/gophercloud/openstack/compute/apiversions"
)

func baremetalService(client *gophercloud.ServiceClient) service {
	return NewService(
		"openstack_service_baremetal",
		"OpenStack Baremetal service version availability",
		func() ([]versionInformation, error) {
			allVersions, err := baremetalVersions.List(client).Extract()
			if err != nil {
				return nil, err
			}

			versions := make([]versionInformation, len(allVersions.Versions))
			for i, version := range allVersions.Versions {
				versions[i] = versionInformation{
					ID:         version.ID,
					MinVersion: version.MinVersion,
					Version:    version.Version,
				}
			}

			return versions, nil
		},
	)
}

func blockstorageService(client *gophercloud.ServiceClient) service {
	return NewService(
		"openstack_service_blockstorage",
		"OpenStack Block storage service version availability",
		func() ([]versionInformation, error) {
			allPages, err := blockstorageVersions.List(client).AllPages()
			if err != nil {
				return nil, err
			}

			allVersions, err := blockstorageVersions.ExtractAPIVersions(allPages)
			if err != nil {
				return nil, err
			}

			versions := make([]versionInformation, len(allVersions))
			for i, version := range allVersions {
				versions[i] = versionInformation{
					ID:         version.ID,
					MinVersion: version.MinVersion,
					Version:    version.Version,
				}
			}

			return versions, nil
		},
	)
}

func computeService(client *gophercloud.ServiceClient) service {
	return NewService(
		"openstack_service_compute",
		"OpenStack Compute service version availability",

		func() ([]versionInformation, error) {
			allPages, err := computeVersions.List(client).AllPages()
			if err != nil {
				return nil, err
			}

			allVersions, err := computeVersions.ExtractAPIVersions(allPages)
			if err != nil {
				return nil, err
			}

			versions := make([]versionInformation, len(allVersions))
			for i, version := range allVersions {
				versions[i] = versionInformation{
					ID:         version.ID,
					MinVersion: version.MinVersion,
					Version:    version.Version,
				}
			}

			return versions, nil
		},
	)
}

// openstack_service_baremetal {id, min_version, version}
// openstack_service_blockstorage {id, min_version, version}
// openstack_service_cloudformation {id}
// openstack_service_compute {id, min_version, version}
// openstack_service_dns{id}
// openstack_service_identity {id}
// openstack_service_image {id}
// openstack_service_loadbalancer {id}
// openstack_service_network {id}
// openstack_service_objectstore {id}
// openstack_service_orchestration {id}
// openstack_service_placement {id, min_version, version}
// openstack_service_sharedfilesystem {id, min_version, version}
