package v20240812preview

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"github.com/Azure/ARO-RP/pkg/api"
)

// APIVersion contains a version string as it will be used by clients
const APIVersion = "2024-08-12-preview"

const (
	resourceProviderNamespace = "Microsoft.RedHatOpenShift"
	resourceType              = "openShiftClusters"
)

func init() {
	api.APIs[APIVersion] = &api.Version{
		OpenShiftClusterConverter:                openShiftClusterConverter{},
		OpenShiftClusterStaticValidator:          openShiftClusterStaticValidator{},
		OpenShiftClusterCredentialsConverter:     openShiftClusterCredentialsConverter{},
		OpenShiftClusterAdminKubeconfigConverter: openShiftClusterAdminKubeconfigConverter{},
		OpenShiftVersionConverter:                openShiftVersionConverter{},
		PlatformWorkloadIdentityRoleSetConverter: platformWorkloadIdentityRoleSetConverter{},
		OperationList: api.OperationList{
			Operations: []api.Operation{
				api.OperationResultsRead,
				api.OperationStatusRead,
				api.OperationRead,
				api.OperationOpenShiftClusterRead,
				api.OperationOpenShiftClusterWrite,
				api.OperationOpenShiftClusterDelete,
				api.OperationOpenShiftClusterListCredentials,
				api.OperationOpenShiftClusterListAdminCredentials,
				api.OperationListInstallVersions,
				api.OperationSyncSetsRead,
				api.OperationSyncSetsWrite,
				api.OperationSyncSetsDelete,
				api.OperationMachinePoolsRead,
				api.OperationMachinePoolsWrite,
				api.OperationMachinePoolsDelete,
				api.OperationSyncIdentityProvidersRead,
				api.OperationSyncIdentityProvidersWrite,
				api.OperationSyncIdentityProvidersDelete,
				api.OperationOpenShiftClusterGetDetectors,
			},
		},
		SyncSetConverter:              syncSetConverter{},
		MachinePoolConverter:          machinePoolConverter{},
		SyncIdentityProviderConverter: syncIdentityProviderConverter{},
		SecretConverter:               secretConverter{},
		ClusterManagerStaticValidator: clusterManagerStaticValidator{},
	}
}
