package frontend

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/Azure/ARO-RP/pkg/api"
	"github.com/Azure/ARO-RP/pkg/database/cosmosdb"
	"github.com/Azure/ARO-RP/pkg/frontend/middleware"
)

func (f *frontend) getAsyncOperationResult(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(middleware.ContextKeyLog).(*logrus.Entry)
	vars := mux.Vars(r)

	header := http.Header{}
	b, err := f._getAsyncOperationResult(r, header, api.APIs[vars["api-version"]]["OpenShiftCluster"].(api.OpenShiftClusterToExternal))

	reply(log, w, header, b, err)
}

func (f *frontend) _getAsyncOperationResult(r *http.Request, header http.Header, external api.OpenShiftClusterToExternal) ([]byte, error) {
	vars := mux.Vars(r)

	asyncdoc, err := f.db.AsyncOperations.Get(vars["operationId"])
	switch {
	case cosmosdb.IsErrorStatusCode(err, http.StatusNotFound):
		return nil, api.NewCloudError(http.StatusNotFound, api.CloudErrorCodeNotFound, "", "The entity was not found.")
	case err != nil:
		return nil, err
	}

	doc, err := f.db.OpenShiftClusters.Get(asyncdoc.OpenShiftClusterKey)
	if err != nil && !cosmosdb.IsErrorStatusCode(err, http.StatusNotFound) {
		return nil, err
	}

	// don't give away the final operation status until it's committed to the
	// database
	if doc != nil && doc.AsyncOperationID == vars["operationId"] {
		header["Location"] = r.Header["Referer"]
		return nil, statusCodeError(http.StatusAccepted)
	}

	if asyncdoc.OpenShiftCluster == nil {
		return nil, statusCodeError(http.StatusNoContent)
	}

	asyncdoc.OpenShiftCluster.Properties.ServicePrincipalProfile.ClientSecret = ""

	return json.MarshalIndent(external.OpenShiftClusterToExternal(asyncdoc.OpenShiftCluster), "", "    ")
}
