// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// CallFunction calls the CallFunction method of the underlying provider
// serving the function.
func (s *muxServer) CallFunction(ctx context.Context, req *tfprotov5.CallFunctionRequest) (*tfprotov5.CallFunctionResponse, error) {
	rpc := "CallFunction"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getFunctionServer(ctx, req.Name)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov5.CallFunctionResponse{
			Diagnostics: diags,
		}, nil
	}

	ctx = logging.Tfprotov5ProviderServerContext(ctx, server)

	// Remove and call server.CallFunction below directly.
	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/210
	functionServer, ok := server.(tfprotov5.FunctionServer)

	if !ok {
		resp := &tfprotov5.CallFunctionResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Provider Functions Not Implemented",
					Detail: "A provider-defined function call was received by the provider, however the provider does not implement functions. " +
						"Either upgrade the provider to a version that implements provider-defined functions or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return resp, nil
	}

	logging.MuxTrace(ctx, "calling downstream server")

	// return server.CallFunction(ctx, req)
	return functionServer.CallFunction(ctx, req)
}
