/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package function

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/apple/pkl-go/pkl"
	"github.com/crossplane-contrib/function-pkl/input/v1beta1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/resource"
	//objectv1alpha2 "github.com/crossplane-contrib/provider-kubernetes/apis/object/v1alpha2"
)

var (
	//pklPackage = "package://pkg.pkl-lang.org/github.com/crossplane-contrib/function-pkl/crossplane-example@1.0.0#"
	workdir, _ = os.Getwd()
	pklPackage = fmt.Sprintf("%s/../../pkl/crossplane.contrib.example", workdir)
)

func TestRunFunction(t *testing.T) {

	type args struct {
		ctx context.Context
		req *fnv1beta1.RunFunctionRequest
	}
	type want struct {
		rsp *fnv1beta1.RunFunctionResponse
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"Full": {
			reason: "The Function should create a full functionResult",
			args: args{
				ctx: context.TODO(),
				req: &fnv1beta1.RunFunctionRequest{
					Meta: &fnv1beta1.RequestMeta{Tag: "extra"},
					Input: resource.MustStructObject(&v1beta1.Pkl{
						Spec: v1beta1.PklSpec{
							Type: "local",
							Local: &v1beta1.Local{
								ProjectDir: pklPackage,
								File:       pklPackage + "/full.pkl",
							},
						},
					}),
					Context: resource.MustStructJSON(`{
						"apiextensions.crossplane.io/environment": {
							"foo": "bar"
						}
					}`),
					Observed: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{
							Resource: resource.MustStructJSON(`{
								"apiVersion": "example.crossplane.io/v1",
								"kind": "XR",
								"metadata": {
									"name": "example-xr"
								},
								"spec": {}
							}`),
						},
					},
					ExtraResources: map[string]*fnv1beta1.Resources{
						"ineed": {
							Items: []*fnv1beta1.Resource{
								{
									Resource: resource.MustStructJSON(`{
										"apiVersion": "kubernetes.crossplane.io/v1alpha2",
										"kind": "Object",
										"metadata": {
											"name": "required"
										}
									}`),
								},
							},
						},
					},
				},
			},
			want: want{
				rsp: &fnv1beta1.RunFunctionResponse{
					Desired: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{
							Resource: resource.MustStructJSON(`{
								"apiVersion": "example.crossplane.io/v1",
								"kind": "XR",
								"status": {
									"someStatus": "pretty status"
								}
							}`),
						},
						Resources: map[string]*fnv1beta1.Resource{
							"cm-one": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "kubernetes.crossplane.io/v1alpha2",
									"kind": "Object",
									"metadata": {
										"name": "cm-one"
									},
									"spec": {
										"forProvider": {
											"manifest": {
												"apiVersion": "v1",
												"kind": "ConfigMap",
												"metadata": {
													"name": "cm-one",
													"namespace": "crossplane-system"
												},
												"data": {
													"foo": "example-xr",
													"required": "required"
												}
											}
										}
									}
								}`),
								Ready: fnv1beta1.Ready_READY_TRUE,
							},
						},
					},
					Requirements: &fnv1beta1.Requirements{
						ExtraResources: map[string]*fnv1beta1.ResourceSelector{
							"ineed": {
								ApiVersion: "kubernetes.crossplane.io/v1alpha2",
								Kind:       "Object",
								Match: &fnv1beta1.ResourceSelector_MatchName{
									MatchName: "required",
								},
							},
						},
					},
					Meta: &fnv1beta1.ResponseMeta{
						Tag: "extra",
						Ttl: &durationpb.Duration{
							Seconds: 60,
						},
					},
					Context: resource.MustStructJSON(`{
						"apiextensions.crossplane.io/environment": {
							"foo": "bar"
						},
						"greetings": "with <3 from function-pkl"
					}`),
					Results: []*fnv1beta1.Result{
						{
							Severity: fnv1beta1.Severity_SEVERITY_NORMAL,
							Message:  "welcome",
						},
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			evaluatorManager := pkl.NewEvaluatorManager()
			defer func(evaluatorManager pkl.EvaluatorManager) {
				err := evaluatorManager.Close()
				if err != nil {
					t.Error(err)
				}
			}(evaluatorManager)
			f := &Function{Log: logging.NewNopLogger(), EvaluatorManager: evaluatorManager}
			rsp, err := f.RunFunction(tc.args.ctx, tc.args.req)

			if diff := cmp.Diff(tc.want.rsp, rsp, protocmp.Transform()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want rsp, +got rsp:\n%s", tc.reason, diff)
			}

			if diff := cmp.Diff(tc.want.err, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want err, +got err:\n%s", tc.reason, diff)
			}
		})
	}
}
