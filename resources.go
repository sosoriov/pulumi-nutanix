// Copyright 2016-2018, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nutanix

import (
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/pulumi/pulumi-terraform-bridge/pkg/tfbridge"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/tokens"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

// all of the token components used below.
const (
	// packages:
	mainPkg = "nutanix"
	// modules:
	mainMod = "index" // the y module
)

// makeMember manufactures a type token for the package and the given module and type.
func makeMember(mod string, mem string) tokens.ModuleMember {
	return tokens.ModuleMember(mainPkg + ":" + mod + ":" + mem)
}

// makeType manufactures a type token for the package and the given module and type.
func makeType(mod string, typ string) tokens.Type {
	return tokens.Type(makeMember(mod, typ))
}

// makeDataSource manufactures a standard resource token given a module and resource name.  It
// automatically uses the main package and names the file by simply lower casing the data source's
// first character.
func makeDataSource(mod string, res string) tokens.ModuleMember {
	fn := string(unicode.ToLower(rune(res[0]))) + res[1:]
	return makeMember(mod+"/"+fn, res)
}

// makeResource manufactures a standard resource token given a module and resource name.  It
// automatically uses the main package and names the file by simply lower casing the resource's
// first character.
func makeResource(mod string, res string) tokens.Type {
	fn := string(unicode.ToLower(rune(res[0]))) + res[1:]
	return makeType(mod+"/"+fn, res)
}

// boolRef returns a reference to the bool argument.
func boolRef(b bool) *bool {
	return &b
}

// stringValue gets a string value from a property map if present, else ""
func stringValue(vars resource.PropertyMap, prop resource.PropertyKey) string {
	val, ok := vars[prop]
	if ok && val.IsString() {
		return val.StringValue()
	}
	return ""
}

// preConfigureCallback is called before the providerConfigure function of the underlying provider.
// It should validate that the provider can be configured, and provide actionable errors in the case
// it cannot be. Configuration variables can be read from `vars` using the `stringValue` function -
// for example `stringValue(vars, "accessKey")`.
func preConfigureCallback(vars resource.PropertyMap, c *terraform.ResourceConfig) error {
	return nil
}

// managedByPulumi is a default used for some managed resources, in the absence of something more meaningful.
var managedByPulumi = &tfbridge.DefaultInfo{Value: "Managed by Pulumi"}

// Provider returns additional overlaid schema and metadata associated with the provider..
func Provider() tfbridge.ProviderInfo {
	// Instantiate the Terraform provider
	p := nutanix.Provider().(*schema.Provider)

	// Create a Pulumi provider mapping
	prov := tfbridge.ProviderInfo{
		P:           p,
		Name:        "nutanix",
		Description: "A Pulumi package for creating and managing nutanix cloud resources.",
		Keywords:    []string{"pulumi", "nutanix"},
		License:     "Apache-2.0",
		Homepage:    "https://pulumi.io",
		Repository:  "https://github.com/pulumi/pulumi-nutanix",
		Config:      map[string]*tfbridge.SchemaInfo{
			// Add any required configuration here, or remove the example below if
			// no additional points are required.
			"username": {
				Type: makeType("username", "Username"),
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"NUTANIX_USERNAME"},
				},
			},
			"password": {
				Type: makeType("password", "Password"),
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"NUTANIX_PASSWORD"},
				},
			},
			"endpoint": {
				Type: makeType("endpoint", "Endpoint"),
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"NUTANIX_ENDPOINT"},
				},
			},
			"port": {
				Type: makeType("port", "Port"),
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"NUTANIX_PORT"},
				},
			},
			"proxy_url": {
				Type: makeType("proxy_url", "Proxy_url"),
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"NUTANIX_PROXY_URL"},
				},
			},
			"insecure": {
				Type: makeType("insecure", "Insecure"),
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"NUTANIX_INSECURE"},
				},
			},
			"wait_timeout": {
				Type: makeType("wait_timeout", "wait_timeout"),
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"NUTANIX_WAIT_TIMEOUT"},
				},
			},
		},
		PreConfigureCallback: preConfigureCallback,
		Resources: map[string]*tfbridge.ResourceInfo{
			// Map each resource in the Terraform provider to a Pulumi type. Two examples
			// are below - the single line form is the common case. The multi-line form is
			// needed only if you wish to override types or other default options.
			//
			"nutanix_virtual_machine": {Tok: makeResource(mainMod, "VirtualMachine")},
			"nutanix_image": {Tok: makeResource(mainMod, "Image")},
			"nutanix_subnet": {Tok: makeResource(mainMod, "Subnet")},
			"nutanix_category_key": {Tok: makeResource(mainMod, "CategoryKey")},
			"nutanix_category_value": {Tok: makeResource(mainMod, "CategoryValue")},
			"nutanix_network_security_rule": {Tok: makeResource(mainMod, "NetworkSecurityGroup")},
			
		},
		DataSources: map[string]*tfbridge.DataSourceInfo{
			// Map each resource in the Terraform provider to a Pulumi function. An example
			// is below.
			// "aws_ami": {Tok: makeDataSource(mainMod, "getAmi")},
			"nutanix_virtual_machine": {Tok: makeDataSource(mainMod, "getVirtualMachine")},
			"nutanix_cluster": {Tok: makeDataSource(mainMod, "getCluster")},
			"nutanix_clusters": {Tok: makeDataSource(mainMod, "getClusters")},
			"nutanix_image": {Tok: makeDataSource(mainMod, "getImage")},
			"nutanix_subnet": {Tok: makeDataSource(mainMod, "getSubnet")},
			"nutanix_category_key": {Tok: makeDataSource(mainMod, "getCategoryKey")},
			"nutanix_network_security_rule": {Tok: makeDataSource(mainMod, "getNetworkSecurityGroup")},
		},
		JavaScript: &tfbridge.JavaScriptInfo{
			// List any npm dependencies and their versions
			Dependencies: map[string]string{
				"@pulumi/pulumi": "latest",
			},
			DevDependencies: map[string]string{
				"@types/node": "^8.0.25", // so we can access strongly typed node definitions.
				"@types/mime": "^2.0.0",
			},
			// See the documentation for tfbridge.OverlayInfo for how to lay out this
			// section, or refer to the AWS provider. Delete this section if there are
			// no overlay files.
			//Overlay: &tfbridge.OverlayInfo{},
		},
		Python: &tfbridge.PythonInfo{
			// List any Python dependencies and their version ranges
			Requires: map[string]string{
				"pulumi": ">=1.0.0,<2.0.0",
			},
		},
		CSharp: &tfbridge.CSharpInfo{
			PackageReferences: map[string]string{
				"Pulumi":                       "1.7.0-preview",
				"System.Collections.Immutable": "1.6.0",
			},
		},
	}

	// For all resources with name properties, we will add an auto-name property.  Make sure to skip those that
	// already have a name mapping entry, since those may have custom overrides set above (e.g., for length).
	const nameProperty = "name"
	for resname, res := range prov.Resources {
		if schema := p.ResourcesMap[resname]; schema != nil {
			// Only apply auto-name to input properties (Optional || Required) named `name`
			if tfs, has := schema.Schema[nameProperty]; has && (tfs.Optional || tfs.Required) {
				if _, hasfield := res.Fields[nameProperty]; !hasfield {
					if res.Fields == nil {
						res.Fields = make(map[string]*tfbridge.SchemaInfo)
					}
					res.Fields[nameProperty] = tfbridge.AutoName(nameProperty, 255)
				}
			}
		}
	}

	return prov
}
