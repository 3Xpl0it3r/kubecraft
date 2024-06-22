package install

import (
	"github.com/3Xpl0it3r/kubecraft/pkg/api/legacyscheme"
	"github.com/3Xpl0it3r/kubecraft/pkg/apis/example"
	v1 "github.com/3Xpl0it3r/kubecraft/pkg/apis/example/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func init() {
	Install(legacyscheme.Scheme)
}

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(example.AddToScheme(scheme))
	utilruntime.Must(v1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(v1.SchemeGroupVersion))
}
