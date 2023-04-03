package env

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	mgmtauthorization "github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-09-01-preview/authorization"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/jongio/azidext/go/azidext"
	"github.com/sirupsen/logrus"

	"github.com/Azure/ARO-RP/pkg/util/azureclient/graphrbac"
	"github.com/Azure/ARO-RP/pkg/util/azureclient/mgmt/authorization"
	"github.com/Azure/ARO-RP/pkg/util/rbac"
	"github.com/Azure/ARO-RP/pkg/util/uuid"
)

// In INT or PROD, when the ARO RP is running behind ARM, ARM follows the RP's
// manifest configuration to grant the RP's first party service principal
// limited standing access into customers' subscriptions (when the RP is
// registered in the customer's AAD tenant).
//
// Our ARM manifest is set up such that the RP first party service principal can
// create new resource groups in any customer subscription (when the RP is
// registered in the customer's AAD tenant).  ARM invisibly then grants the RP
// Owner access on the created resource group.  For more information, read the
// ARM wiki ("ResourceGroup Scoped Service to Service Authorization") and the RP
// manifest.
//
// When running outside INT or PROD (i.e. in development or in CI), ARMHelper is
// used to fake up the above functionality with a separate development service
// principal (AZURE_ARM_CLIENT_ID).  We use a separate SP so that the RP's
// permissions are identical in dev and prod.  The advantage is that this helps
// prevent a developer rely on a permission in dev only to find that permission
// doesn't exist in prod.  The disadvantage is that an additional SP and this
// helper code is required in dev.
//
// There remains one other minor difference between running in dev and INT/PROD:
// in the former, the Owner role assignment created by the helper is visible.
// In INT/PROD I believe it is invisible.

type ARMHelper interface {
	EnsureARMResourceGroupRoleAssignment(context.Context, autorest.Authorizer, string) error
}

// noopARMHelper is used in INT and PROD.  It does nothing.
type noopARMHelper struct{}

func (*noopARMHelper) EnsureARMResourceGroupRoleAssignment(context.Context, autorest.Authorizer, string) error {
	return nil
}

// armHelper is used in dev.  It adds an Owner role assignment for the RP,
// faking up what ARM would do invisibly for us in INT/PROD.
type armHelper struct {
	log *logrus.Entry
	env Interface

	roleassignments authorization.RoleAssignmentsClient
	applications    graphrbac.ApplicationsClient
}

func newARMHelper(ctx context.Context, log *logrus.Entry, env Interface) (ARMHelper, error) {
	if os.Getenv("AZURE_ARM_CLIENT_ID") == "" {
		return &noopARMHelper{}, nil
	}

	var tokenCredential azcore.TokenCredential
	var err error

	if os.Getenv("AZURE_ARM_CLIENT_SECRET") != "" {
		// TODO: migrate away from AZURE_ARM_CLIENT_SECRET and remove this code
		// path

		tokenCredential, err = azidentity.NewClientSecretCredential(
			env.TenantID(),
			os.Getenv("AZURE_ARM_CLIENT_ID"),
			os.Getenv("AZURE_ARM_CLIENT_SECRET"),
			env.Environment().ClientSecretCredentialOptions())
		if err != nil {
			return nil, err
		}
	} else {
		key, certs, err := env.ServiceKeyvault().GetCertificateSecret(ctx, RPDevARMSecretName)
		if err != nil {
			return nil, err
		}

		options := env.Environment().ClientCertificateCredentialOptions()
		tokenCredential, err = azidentity.NewClientCertificateCredential(
			env.TenantID(), os.Getenv("AZURE_ARM_CLIENT_ID"), certs, key, options)
		if err != nil {
			return nil, err
		}
	}

	scopes := []string{env.Environment().ResourceManagerScope}
	armAuthorizer := azidext.NewTokenCredentialAdapter(tokenCredential, scopes)

	fpGraphAuthorizer, err := env.FPAuthorizer(env.TenantID(), env.Environment().ActiveDirectoryGraphScope)
	if err != nil {
		return nil, err
	}

	return &armHelper{
		log: log,
		env: env,

		roleassignments: authorization.NewRoleAssignmentsClient(env.Environment(), env.SubscriptionID(), armAuthorizer),
		applications:    graphrbac.NewApplicationsClient(env.Environment(), env.TenantID(), fpGraphAuthorizer),
	}, nil
}

func (ah *armHelper) EnsureARMResourceGroupRoleAssignment(ctx context.Context, fpAuthorizer autorest.Authorizer, resourceGroup string) error {
	ah.log.Print("ensuring resource group role assignment")

	res, err := ah.applications.GetServicePrincipalsIDByAppID(ctx, ah.env.FPClientID())
	if err != nil {
		return err
	}

	_, err = ah.roleassignments.Create(ctx, "/subscriptions/"+ah.env.SubscriptionID()+"/resourceGroups/"+resourceGroup, uuid.DefaultGenerator.Generate(), mgmtauthorization.RoleAssignmentCreateParameters{
		RoleAssignmentProperties: &mgmtauthorization.RoleAssignmentProperties{
			RoleDefinitionID: to.StringPtr("/subscriptions/" + ah.env.SubscriptionID() + "/providers/Microsoft.Authorization/roleDefinitions/" + rbac.RoleOwner),
			PrincipalID:      res.Value,
		},
	})
	if detailedErr, ok := err.(autorest.DetailedError); ok {
		if requestErr, ok := detailedErr.Original.(*azure.RequestError); ok &&
			requestErr.ServiceError != nil &&
			requestErr.ServiceError.Code == "RoleAssignmentExists" {
			err = nil
		}
	}

	return err
}
