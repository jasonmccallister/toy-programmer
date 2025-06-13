// Code generated by dagger. DO NOT EDIT.

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"dagger/toy-workspace/internal/dagger"
	"dagger/toy-workspace/internal/querybuilder"
	"dagger/toy-workspace/internal/telemetry"
)

var dag = dagger.Connect()

func Tracer() trace.Tracer {
	return otel.Tracer("dagger.io/sdk.go")
}

// used for local MarshalJSON implementations
var marshalCtx = context.Background()

// called by main()
func setMarshalContext(ctx context.Context) {
	marshalCtx = ctx
	dagger.SetMarshalContext(ctx)
}

type DaggerObject = querybuilder.GraphQLMarshaller

type ExecError = dagger.ExecError

// ptr returns a pointer to the given value.
func ptr[T any](v T) *T {
	return &v
}

// convertSlice converts a slice of one type to a slice of another type using a
// converter function
func convertSlice[I any, O any](in []I, f func(I) O) []O {
	out := make([]O, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}

func (r ToyWorkspace) MarshalJSON() ([]byte, error) {
	var concrete struct {
		Container *dagger.Container
	}
	concrete.Container = r.Container
	return json.Marshal(&concrete)
}

func (r *ToyWorkspace) UnmarshalJSON(bs []byte) error {
	var concrete struct {
		Container *dagger.Container
	}
	err := json.Unmarshal(bs, &concrete)
	if err != nil {
		return err
	}
	r.Container = concrete.Container
	return nil
}

func main() {
	ctx := context.Background()

	// Direct slog to the new stderr. This is only for dev time debugging, and
	// runtime errors/warnings.
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))

	if err := dispatch(ctx); err != nil {
		os.Exit(2)
	}
}

func unwrapError(rerr error) string {
	var gqlErr *gqlerror.Error
	if errors.As(rerr, &gqlErr) {
		return gqlErr.Message
	}
	return rerr.Error()
}

func dispatch(ctx context.Context) (rerr error) {
	ctx = telemetry.InitEmbedded(ctx, resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("dagger-go-sdk"),
		// TODO version?
	))
	defer telemetry.Close()

	// A lot of the "work" actually happens when we're marshalling the return
	// value, which entails getting object IDs, which happens in MarshalJSON,
	// which has no ctx argument, so we use this lovely global variable.
	setMarshalContext(ctx)

	fnCall := dag.CurrentFunctionCall()
	defer func() {
		if rerr != nil {
			if err := fnCall.ReturnError(ctx, dag.Error(unwrapError(rerr))); err != nil {
				fmt.Println("failed to return error:", err)
			}
		}
	}()

	parentName, err := fnCall.ParentName(ctx)
	if err != nil {
		return fmt.Errorf("get parent name: %w", err)
	}
	fnName, err := fnCall.Name(ctx)
	if err != nil {
		return fmt.Errorf("get fn name: %w", err)
	}
	parentJson, err := fnCall.Parent(ctx)
	if err != nil {
		return fmt.Errorf("get fn parent: %w", err)
	}
	fnArgs, err := fnCall.InputArgs(ctx)
	if err != nil {
		return fmt.Errorf("get fn args: %w", err)
	}

	inputArgs := map[string][]byte{}
	for _, fnArg := range fnArgs {
		argName, err := fnArg.Name(ctx)
		if err != nil {
			return fmt.Errorf("get fn arg name: %w", err)
		}
		argValue, err := fnArg.Value(ctx)
		if err != nil {
			return fmt.Errorf("get fn arg value: %w", err)
		}
		inputArgs[argName] = []byte(argValue)
	}

	result, err := invoke(ctx, []byte(parentJson), parentName, fnName, inputArgs)
	if err != nil {
		var exec *dagger.ExecError
		if errors.As(err, &exec) {
			return exec.Unwrap()
		}
		return err
	}
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := fnCall.ReturnValue(ctx, dagger.JSON(resultBytes)); err != nil {
		return fmt.Errorf("store return value: %w", err)
	}
	return nil
}
func invoke(ctx context.Context, parentJSON []byte, parentName string, fnName string, inputArgs map[string][]byte) (_ any, err error) {
	_ = inputArgs
	switch parentName {
	case "ToyWorkspace":
		switch fnName {
		case "Read":
			var parent ToyWorkspace
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			var path string
			if inputArgs["path"] != nil {
				err = json.Unmarshal([]byte(inputArgs["path"]), &path)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg path", err))
				}
			}
			return (*ToyWorkspace).Read(&parent, ctx, path)
		case "Write":
			var parent ToyWorkspace
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			var path string
			if inputArgs["path"] != nil {
				err = json.Unmarshal([]byte(inputArgs["path"]), &path)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg path", err))
				}
			}
			var content string
			if inputArgs["content"] != nil {
				err = json.Unmarshal([]byte(inputArgs["content"]), &content)
				if err != nil {
					panic(fmt.Errorf("%s: %w", "failed to unmarshal input arg content", err))
				}
			}
			return (*ToyWorkspace).Write(&parent, path, content), nil
		case "Build":
			var parent ToyWorkspace
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			return (*ToyWorkspace).Build(&parent, ctx)
		case "":
			var parent ToyWorkspace
			err = json.Unmarshal(parentJSON, &parent)
			if err != nil {
				panic(fmt.Errorf("%s: %w", "failed to unmarshal parent object", err))
			}
			return New(), nil
		default:
			return nil, fmt.Errorf("unknown function %s", fnName)
		}
	case "":
		return dag.Module().
			WithDescription("A toy workspace for editing and building go programs\n").
			WithObject(
				dag.TypeDef().WithObject("ToyWorkspace", dagger.TypeDefWithObjectOpts{SourceMap: dag.SourceMap("main.go", 18, 6)}).
					WithFunction(
						dag.Function("Read",
							dag.TypeDef().WithKind(dagger.TypeDefKindStringKind)).
							WithDescription("Read a file").
							WithSourceMap(dag.SourceMap("main.go", 25, 1)).
							WithArg("path", dag.TypeDef().WithKind(dagger.TypeDefKindStringKind), dagger.FunctionWithArgOpts{Description: "The path of the file", SourceMap: dag.SourceMap("main.go", 28, 2)})).
					WithFunction(
						dag.Function("Write",
							dag.TypeDef().WithObject("ToyWorkspace")).
							WithDescription("Write a file").
							WithSourceMap(dag.SourceMap("main.go", 33, 1)).
							WithArg("path", dag.TypeDef().WithKind(dagger.TypeDefKindStringKind), dagger.FunctionWithArgOpts{Description: "The path of the file", SourceMap: dag.SourceMap("main.go", 35, 2)}).
							WithArg("content", dag.TypeDef().WithKind(dagger.TypeDefKindStringKind), dagger.FunctionWithArgOpts{Description: "The content to write", SourceMap: dag.SourceMap("main.go", 37, 2)})).
					WithFunction(
						dag.Function("Build",
							dag.TypeDef().WithKind(dagger.TypeDefKindStringKind)).
							WithDescription("Build the code at the current directory in the workspace").
							WithSourceMap(dag.SourceMap("main.go", 44, 1))).
					WithField("Container", dag.TypeDef().WithObject("Container"), dagger.TypeDefWithFieldOpts{Description: "The workspace container.", SourceMap: dag.SourceMap("main.go", 21, 2)}).
					WithConstructor(
						dag.Function("New",
							dag.TypeDef().WithObject("ToyWorkspace")).
							WithSourceMap(dag.SourceMap("main.go", 9, 1)))), nil
	default:
		return nil, fmt.Errorf("unknown object %s", parentName)
	}
}
