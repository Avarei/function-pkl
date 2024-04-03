package function

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"github.com/avarei/function-pkl/input/v1beta1"
	"github.com/avarei/function-pkl/internal/pkl/reader"
	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/response"
	"sigs.k8s.io/yaml"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer

	Log logging.Logger
}

// RunFunction runs the Function.
func (f *Function) RunFunction(ctx context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	f.Log.Info("Running function", "tag", req.GetMeta().GetTag())

	rsp := response.To(req, response.DefaultTTL)

	in := &v1beta1.Pkl{}
	if err := request.GetInput(req, in); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get Function input from %T", req))
		return rsp, nil
	}
	/*
		xr, err := request.GetObservedCompositeResource(req)
		myxr := req.Observed.GetComposite()
	*/
	evaluatorManager := pkl.NewEvaluatorManager()
	defer evaluatorManager.Close()
	evaluator, err := evaluatorManager.NewEvaluator(ctx,
		pkl.PreconfiguredOptions,
		reader.WithCrossplane(req, "crossplane"),
	) // TODO disallow FS access
	if err != nil {
		evaluator.Close()
		response.Fatal(rsp, errors.Wrap(err, "could not create Pkl Evaluater"))
	}
	defer evaluator.Close()

	var outResources map[string]*fnv1beta1.Resource = make(map[string]*fnv1beta1.Resource)

	var sources map[string]*pkl.ModuleSource = make(map[string]*pkl.ModuleSource)

	for _, pklFileRef := range in.Spec.PklManifests {
		if pklFileRef.Type == "uri" {
			sources[pklFileRef.Name] = pkl.UriSource(pklFileRef.Uri)
		} else if pklFileRef.Type == "inline" {
			sources[pklFileRef.Name] = pkl.TextSource(pklFileRef.Inline)
		} else if pklFileRef.Type == "configMap" {
			// TODO configMap
		} else {
			response.Fatal(rsp, errors.New("unknown PklFileRef type"))
			return rsp, nil
		}
	}

	for name, source := range sources {
		resource, err := parseFile(ctx, evaluator, source)
		if err == nil {
			outResources[name] = resource
		} else {
			f.Log.Debug("Could not parse file \"" + name + "\": " + err.Error())
		}
	}

	rsp.Desired.Resources = outResources
	return rsp, nil
}

func parseFile(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*fnv1beta1.Resource, error) {
	// TODO request a new Function to EvaluateOutputValue which does not require a Struct Tag
	renderedManifest, err := evaluator.EvaluateOutputText(ctx, source)
	if err != nil {
		return nil, errors.Wrap(err, "could not evaluate Pkl file")
	}

	resource := &fnv1beta1.Resource{}
	if err := yaml.Unmarshal([]byte(renderedManifest), resource); err != nil {
		return nil, errors.Wrap(err, "could not parse yaml to Resource")
	}

	return resource, nil
}
