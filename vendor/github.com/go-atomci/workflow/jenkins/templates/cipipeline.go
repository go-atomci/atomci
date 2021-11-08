/*
Copyright 2021 The AtomCI Group Authors.

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

package templates

// CIPipeline defained default jenkins pipeline
const CIPipeline = `
pipeline {
    agent {
        kubernetes {
            defaultContainer 'jnlp'
            yaml """
apiVersion: v1
kind: Pod
metadata:
  namespace: devops
spec:
  containers:
  {{- range $i, $item := .ContainerTemplates }}
  - name: {{ $item.Name }}
    image: {{ $item.Image }}
    workingDir: {{ $item.WorkingDir }}
    command: 
    {{- range $cmd := $item.CommandArr }}
    - {{ $cmd }}
    {{- end }}
    args:
    {{- range $arg := $item.ArgsArr }}
    - {{ $arg }}
    {{- end }}
    tty: true
  {{- end }}
"""          
        }
    }
    environment {
        {{- range $i, $item := .EnvVars }}
        def {{ $item.Key }} = '{{ $item.Value }}'
        {{- end }}
    }
    stages {
        {{ .Stages }}

        stage('Callback') {
            steps {
                retry(count: 5) {
                    httpRequest acceptType: 'APPLICATION_JSON', contentType: 'APPLICATION_JSON', customHeaders: [[maskValue: true, name: 'Authorization', value: 'Bearer {{ .CallBack.Token }}']], httpMode: 'POST', requestBody: '''{{ .CallBack.Body }}''', responseHandle: 'NONE', timeout: 10, url: '{{ .CallBack.URL }}'
                }
            }
        }
    }
}
`

// Checkout ..
const Checkout = `
stage('Checkout') {
    {{if .CheckoutItems }}
    parallel {
        {{- range $i, $item := .CheckoutItems }}
        stage('{{ $item.Name }}') {
            steps {
                {{ $item.Command }}
            }
        }
        {{- end }}
    }
    {{ else }}
        steps {
            sh "echo 'there was no checkout items'"
        }
    {{ end }}
}
`

// Compile stage
const Compile = `
stage('Builds') {
    {{if .BuildItems }}
    parallel {
        {{- range $i, $item := .BuildItems }}
        stage('{{ $item.Name }}') {
            steps {
                container('{{ $item.ContainerName }}') {
                    {{ $item.Command }}
                }
            }
        }
        {{- end }}
    }
    {{ else }}
        steps {
            sh "echo 'there was no build items'"
        }
    {{ end }}
}
`

// BuildImage stage
const BuildImage = `
stage('Images') {
    {{if .ImageItems }}
    parallel {
        {{- range $i, $item := .ImageItems }}
        stage('{{ $item.Name }}') {
            steps {
                container("kaniko") {
                    sh "[ -d $DOCKER_CONFIG ] || mkdir -pv $DOCKER_CONFIG"

                    sh """
                    echo '{"auths": {"'$REGISTRY_ADDR'": {"auth": "'$DOCKER_AUTH'"}}}' > $DOCKER_CONFIG/config.json
                    """
                    {{ $item.Command }}
                }
            }
        }
        {{- end }}
    }
    {{ else }}
        steps {
            sh "echo 'there was no images items'"
        }
    {{ end }}
}
`

// CustomScript stage
const CustomScript = `
stage({{ .CustomScriptItem.Name }}) {
    steps {
        {{ .CustomScriptItem.Command }}
    }
}
`
