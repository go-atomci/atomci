module github.com/go-atomci/atomci

go 1.15

replace (
	k8s.io/api => github.com/kubernetes/api v0.17.0
	k8s.io/apiextensions-apiserver => github.com/kubernetes/apiextensions-apiserver v0.17.0
	k8s.io/apimachiner => github.com/kubernetes/apimachinery v0.17.0
	k8s.io/apimachinery => github.com/kubernetes/apimachinery v0.17.0
	k8s.io/apiserver => github.com/kubernetes/apiserver v0.17.0
	k8s.io/cli-runtime => github.com/kubernetes/cli-runtime v0.17.0
	k8s.io/client-go => github.com/kubernetes/client-go v0.17.0
	k8s.io/cloud-provider => github.com/kubernetes/cloud-provider v0.17.0
	k8s.io/cluster-bootstrap => github.com/kubernetes/cluster-bootstrap v0.17.0
	k8s.io/code-generator => github.com/kubernetes/code-generator v0.17.0
	k8s.io/component-base => github.com/kubernetes/component-base v0.17.0
	k8s.io/cri-api => github.com/kubernetes/cri-api v0.17.0
	k8s.io/csi-translation-lib => github.com/kubernetes/csi-translation-lib v0.17.0
	// helm v2.11.0
	k8s.io/helm => k8s.io/helm v2.11.0+incompatible
	k8s.io/kube-aggregator => github.com/kubernetes/kube-aggregator v0.17.0
	k8s.io/kube-controller-manager => github.com/kubernetes/kube-controller-manager v0.17.0
	k8s.io/kube-proxy => github.com/kubernetes/kube-proxy v0.17.0
	k8s.io/kube-scheduler => github.com/kubernetes/kube-scheduler v0.17.0
	k8s.io/kubectl => github.com/kubernetes/kubectl v0.17.0
	k8s.io/kubelet => github.com/kubernetes/kubelet v0.17.0
	k8s.io/kubernetes => k8s.io/kubernetes v1.17.0
	k8s.io/legacy-cloud-providers => github.com/kubernetes/legacy-cloud-providers v0.17.0
	k8s.io/metrics => github.com/kubernetes/metrics v0.17.0
	k8s.io/sample-apiserver => github.com/kubernetes/sample-apiserver v0.17.0
)

require (
	github.com/astaxie/beego v1.12.0
	github.com/casbin/casbin/v2 v2.19.4
	github.com/colynn/go-ldap-client/v3 v3.0.0-20201016034829-4c1455a490de
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-atomci/go-scm v1.13.2-0.20210629010829-147be8a9bdd3
	github.com/go-atomci/workflow v0.0.0-20211126090842-208f180b47ab
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/isbrick/tools v0.0.0-20211027093338-a3a0ded37175
	github.com/kr/text v0.2.0 // indirect
	github.com/nbio/st v0.0.0-20140626010706-e9e8d9816f32 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pborman/uuid v1.2.0
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/net v0.0.0-20210224082022-3d97a244fca7 // indirect
	golang.org/x/oauth2 v0.0.0-20210126194326-f9ce19ea3013 // indirect
	golang.org/x/sys v0.0.0-20210225134936-a50acf3fe073 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.18.0
	k8s.io/kubernetes v0.0.0-00010101000000-000000000000
)
