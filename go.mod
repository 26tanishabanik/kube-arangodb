module github.com/arangodb/kube-arangodb

go 1.20

replace (
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring => github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.46.0
	github.com/prometheus-operator/prometheus-operator/pkg/client => github.com/prometheus-operator/prometheus-operator/pkg/client v0.46.0
	github.com/stretchr/testify => github.com/stretchr/testify v1.5.1
	github.com/ugorji/go => github.com/ugorji/go v0.0.0-20181209151446-772ced7fd4c2

	k8s.io/api => k8s.io/api v0.25.13
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.25.13
	k8s.io/apimachinery => k8s.io/apimachinery v0.25.13
	k8s.io/apiserver => k8s.io/apiserver v0.25.13
	k8s.io/client-go => k8s.io/client-go v0.25.13
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.25.13
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.25.13
	k8s.io/code-generator => ./deps/k8s.io/code-generator
	k8s.io/component-base => k8s.io/component-base v0.25.13
	k8s.io/kubernetes => k8s.io/kubernetes v0.25.13
	k8s.io/metrics => k8s.io/metrics v0.25.13
)

require (
	github.com/arangodb-helper/go-certificates v0.0.0-20180821055445-9fca24fc2680
	github.com/arangodb-helper/go-helper v0.4.2
	github.com/arangodb/arangosync-client v0.9.0
	github.com/arangodb/go-driver v1.6.1
	github.com/arangodb/go-driver/v2 v2.0.3
	github.com/arangodb/go-upgrade-rules v0.0.0-20180809110947-031b4774ff21
	//github.com/arangodb/rebalancer v0.1.1
	//github.com/arangodb/go-agency-helper v0.3.0
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/coreos/go-semver v0.3.1
	github.com/dchest/uniuri v0.0.0-20160212164326-8902c56451e9
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/jessevdk/go-assets v0.0.0-20160921144138-4f4301a06e15
	github.com/josephburnett/jd v1.6.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/magiconair/properties v1.8.5
	github.com/pkg/errors v0.9.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.44.1
	github.com/prometheus-operator/prometheus-operator/pkg/client v0.0.0-00010101000000-000000000000
	github.com/prometheus/client_golang v1.15.1
	github.com/prometheus/client_model v0.4.0
	github.com/prometheus/prom2json v1.3.3
	github.com/robfig/cron v1.2.0
	github.com/rs/zerolog v1.19.0
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.4
	golang.org/x/sync v0.1.0
	golang.org/x/sys v0.13.0
	golang.org/x/text v0.13.0
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8
	google.golang.org/grpc v1.56.3
	google.golang.org/protobuf v1.30.0
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/api v0.25.13
	k8s.io/apiextensions-apiserver v0.25.13
	k8s.io/apimachinery v0.25.13
	k8s.io/client-go v0.25.13
	k8s.io/kube-openapi v0.0.0-20220803162953-67bda5d908f1
	sigs.k8s.io/yaml v1.2.0
)

require (
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/arangodb/go-velocypack v0.0.0-20200318135517-5af53c29c67e // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dchest/siphash v1.2.2 // indirect
	github.com/emicklei/go-restful/v3 v3.8.0 // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kkdai/maglev v0.2.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pavel-v-chernykh/keystore-go v2.1.0+incompatible // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/oauth2 v0.8.0 // indirect
	golang.org/x/term v0.13.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/klog/v2 v2.70.1 // indirect
	k8s.io/utils v0.0.0-20220728103510-ee6ede2d64ed // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)
