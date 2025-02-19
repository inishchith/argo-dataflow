// +build test

package http_stress

import (
	"testing"

	. "github.com/argoproj-labs/argo-dataflow/test/stress"

	. "github.com/argoproj-labs/argo-dataflow/api/v1alpha1"
	. "github.com/argoproj-labs/argo-dataflow/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestHTTPSourceStress(t *testing.T) {
	defer Setup(t)()

	CreatePipeline(Pipeline{
		ObjectMeta: metav1.ObjectMeta{Name: "http"},
		Spec: PipelineSpec{
			Steps: []StepSpec{{
				Name:     "main",
				Cat:      &Cat{},
				Replicas: Params.Replicas,
				Sources:  []Source{{HTTP: &HTTPSource{}}},
				Sinks:    []Sink{DefaultLogSink},
			}},
		},
	})

	WaitForPipeline()
	WaitForPod()

	defer StartPortForward("http-main-0")()

	WaitForService()

	n := Params.N
	prefix := "my-msg"

	defer StartTPSReporter(t, "main", prefix, n)()

	go PumpHTTP("https://http-main/sources/default", prefix, n, Params.MessageSize)
	WaitForStep(TotalSunkMessages(n), Params.Timeout)
}

func TestHTTPSinkStress(t *testing.T) {
	defer Setup(t)()

	CreatePipeline(Pipeline{
		ObjectMeta: metav1.ObjectMeta{Name: "http"},
		Spec: PipelineSpec{
			Steps: []StepSpec{{
				Name:     "main",
				Cat:      &Cat{},
				Replicas: Params.Replicas,
				Sources:  []Source{{HTTP: &HTTPSource{}}},
				Sinks:    []Sink{{HTTP: &HTTPSink{URL: "http://testapi/count/incr"}}, DefaultLogSink},
			}},
		},
	})

	WaitForPipeline()

	defer StartPortForward("http-main-0")()

	WaitForService()

	n := Params.N
	prefix := "my-msg"

	defer StartTPSReporter(t, "main", prefix, n)()

	go PumpHTTP("https://http-main/sources/default", prefix, n, Params.MessageSize)
	WaitForStep(TotalSunkMessages(n*2), Params.Timeout)
}
