#!/bin/bash
# description: this script is used to build some neceressury files/scripts for init an crd controller
# Copyright 2021 l0calh0st
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#      https://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# PROJECT_NAME is directory of project
SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")
OUTDIR="${SCRIPT_ROOT}/../"

PROJECT_NAME=$(grep  "module" go.mod |awk '{print $NF}')
if [ $# -lt 3 ]; then
    echo "Usage: add_type.sh <group_name> <group_version> <resource>"
    exit 1
fi

GROUP_NAME=$1
GROUP_VERSION=$2
RESOURCE_KIND=$3



# 所有字符串大写
function fn_strings_to_upper(){
    echo $(echo $1|tr '[:lower:]' '[:upper:]')
}
# 所有字符串小写
function fn_strings_to_lower(){
    echo $(echo $1|tr '[:upper:]' '[:lower:]')
}
# 去除特殊符号
function fn_strings_strip_special_charts(){
  echo $(echo ${1}|sed 's/-//'|sed 's/_//')
}

# 首字母大写
function fn_strings_first_upper(){
    str=$1
    firstLetter=`echo ${str:0:1} | awk '{print toupper($0)}'`
    otherLetter=${str:1}
    result=$firstLetter$otherLetter
    echo $result
}

# 生成 go mod 名称
function fn_project_to_gomod(){
    echo "${GIT_DOMAIN}/${PROJECT_AUTHOR}/${PROJECT_NAME}"
}

function fn_group_name() {
    echo $(echo ${PROJECT_NAME}|sed 's/-//'|sed 's/_//').${PROJECT_AUTHOR}.cn
}




####################################################################################################
#  全局 相关的
####################################################################################################


####################################################################################################
#                            资源类型定义
####################################################################################################
# auto generate doc.go
function fn_gen_gofile_group_internal_doc(){
    PROJECT_NAME=$1
    GROUP_NAME=$2
    GROUP_VERSION=$3
    RESOURCE_KIND=$4
    if [ -e "pkg/apis/${GROUP_NAME}/doc.go" ]; then
        return 0
    fi
    gendir="${OUTDIR}/pkg/apis/${GROUP_NAME}"
    mkdir -pv ${gendir}
    # mkdir -pv  pkg/apis/${GROUP_NAME}/
    echo "Generate ${gendir}/doc.go"
    cat >> ${gendir}/doc.go << EOF
/*
Copyright `date "+%Y"` The ${PROJECT_NAME} Authors.
Licensed under the Apache License, PROJECT_VERSION 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

    // +k8s:deepcopy-gen=package
    // +groupName=${GROUP_NAME}

    // Package ${GROUP_NAME} is the INTERNAL version of the API.
    package ${GROUP_NAME} // import "$(fn_project_to_gomod ${PROJECT_NAME})/pkg/apis/${GROUP_NAME}/${GROUP_VERSION}"
EOF
    # gofmt -w pkg/apis/${GROUP_NAME}/doc.go
    gofmt -w ${gendir}/doc.go
}

#
function fn_gen_gofile_group_version_doc(){
    PROJECT_NAME=$1
    GROUP_NAME=$2
    GROUP_VERSION=$3
    RESOURCE_KIND=$4
    if [ -e "pkg/apis/${GROUP_NAME}/${GROUP_VERSION}/doc.go" ]; then
        return 0
    fi
    
    gendir="${OUTDIR}/pkg/apis/${GROUP_NAME}/${GROUP_VERSION}"
    echo "Generate ${gendir}/doc.go"
    cat >>${gendir}/doc.go << EOF
/*
Copyright `date "+%Y"` The ${PROJECT_NAME} Authors.
Licensed under the Apache License, PROJECT_VERSION 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

    // +k8s:openapi-gen=true
    // +k8s:deepcopy-gen=package
    // +k8s:conversion-gen=${PROJECT_NAME}/pkg/apis/${GROUP_NAME}
    // +k8s:defaulter-gen=TypeMeta
    // +groupName=${GROUP_NAME}

    // Package ${GROUP_VERSION} is the ${GROUP_VERSION} version of the API.
    package ${GROUP_VERSION} // import "$(fn_project_to_gomod ${PROJECT_NAME})/pkg/apis/${GROUP_NAME}/${GROUP_VERSION}"
EOF
    gofmt -w ${gendir}/doc.go
    # gofmt -w pkg/apis/${GROUP_NAME}/${GROUP_VERSION}/doc.go
}

# auto geneate types.go
function fn_gen_gofile_group_versioned_types(){
    PROJECT_NAME=$1
    GROUP_NAME=$2
    GROUP_VERSION=$3
    RESOURCE_KIND=$4
    RESOURCE_NAME=$(fn_strings_first_upper ${RESOURCE_KIND})    #RESOURCE_NAME 名称，首字母要大写
    gendir="${OUTDIR}/pkg/apis/${GROUP_NAME}/${GROUP_VERSION}"
    # mkdir -pv pkg/apis/${GROUP_NAME}/${GROUP_VERSION}
    mkdir -pv ${gendir}
    echo "Generate ${gendir}/${RESOURCE_KIND}.go"
    cat >> ${gendir}/${RESOURCE_KIND}.go << EOF
/*
Copyright `date "+%Y"` The ${PROJECT_NAME} Authors.
Licensed under the Apache License, PROJECT_VERSION 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

    package ${GROUP_VERSION}

    import (
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    )

    // +genclient
    // +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
    // +k8s:defaulter-gen=true

    // ${RESOURCE_NAME} defines ${RESOURCE_NAME} deployment
    type ${RESOURCE_NAME} struct {
        metav1.TypeMeta \`json:",inline"\`
        metav1.ObjectMeta \`json:"metadata,omitempty"\`

        Spec ${RESOURCE_NAME}Spec \`json:"spec"\`
        Status ${RESOURCE_NAME}Status \`json:"status"\`
    }



    // ${RESOURCE_NAME}Spec describes the specification of ${RESOURCE_NAME} applications using kubernetes as a cluster manager
    type ${RESOURCE_NAME}Spec struct {
        // TODO, write your code
    }

    // ${RESOURCE_NAME}Status describes the current status of ${RESOURCE_NAME} applications
    type ${RESOURCE_NAME}Status struct {
        // TODO, write your code
    }

    // +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

    // ${RESOURCE_NAME}List carries a list of ${RESOURCE_NAME} objects
    type ${RESOURCE_NAME}List struct {
        metav1.TypeMeta \`json:",inline"\`
        metav1.ListMeta \`json:"metadata,omitempty"\`

        Items []$RESOURCE_NAME \`json:"items"\`
    }
EOF
    gofmt -w ${gendir}/${RESOURCE_KIND}.go
}


# auto geneate types.go
function fn_gen_gofile_group_internal_types(){
    PROJECT_NAME=$1
    GROUP_NAME=$2
    GROUP_VERSION=$3
    RESOURCE_KIND=$4
    RESOURCE_NAME=$(fn_strings_first_upper ${RESOURCE_KIND})    #RESOURCE_NAME 名称，首字母要大写
    gendir="${OUTDIR}/pkg/apis/${GROUP_NAME}"
    echo "Generate ${gendir}/${RESOURCE_KIND}.go"
    cat >> ${gendir}/${RESOURCE_KIND}.go << EOF
/*
Copyright `date "+%Y"` The ${PROJECT_NAME} Authors.
Licensed under the Apache License, PROJECT_VERSION 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

    package ${GROUP_NAME}

    import (
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    )

    // +genclient
    // +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
    // +k8s:defaulter-gen=true

    // ${RESOURCE_NAME} defines ${RESOURCE_NAME} deployment
    type ${RESOURCE_NAME} struct {
        metav1.TypeMeta \`json:",inline"\`
        metav1.ObjectMeta \`json:"metadata,omitempty"\`

        Spec ${RESOURCE_NAME}Spec \`json:"spec"\`
        Status ${RESOURCE_NAME}Status \`json:"status"\`
    }



    // ${RESOURCE_NAME}Spec describes the specification of ${RESOURCE_NAME} applications using kubernetes as a cluster manager
    type ${RESOURCE_NAME}Spec struct {
        // TODO, write your code
    }

    // ${RESOURCE_NAME}Status describes the current status of ${RESOURCE_NAME} applications
    type ${RESOURCE_NAME}Status struct {
        // TODO, write your code
    }

    // +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

    // ${RESOURCE_NAME}List carries a list of ${RESOURCE_NAME} objects
    type ${RESOURCE_NAME}List struct {
        metav1.TypeMeta \`json:",inline"\`
        metav1.ListMeta \`json:"metadata,omitempty"\`

        Items []$RESOURCE_NAME \`json:"items"\`
    }
EOF
    gofmt -w ${gendir}/${RESOURCE_KIND}.go
}


# generate regiser.go
function fn_gen_gofile_group_versioned_register(){
    PROJECT_NAME=$1
    GROUP_NAME=$2
    GROUP_VERSION=$3
    RESOURCE_KIND=$4
    if [ -e "pkg/apis/${GROUP_NAME}/${GROUP_VERSION}/register.go" ]; then
        return 0
    fi
    gendir="${OURDIR}/pkg/apis/${GROUP_NAME}/${GROUP_VERSION}/"
    echo "Generate ${gendir}/register.go"
    cat >> ${gendir}/register.go << EOF
/*
Copyright `date "+%Y"` The ${PROJECT_NAME} Authors.
Licensed under the Apache License, PROJECT_VERSION 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

    package ${GROUP_VERSION}

    import (

      metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
      "k8s.io/apimachinery/pkg/runtime"
      "k8s.io/apimachinery/pkg/runtime/schema"
    )
    const (
        Version = "${GROUP_VERSION}"
        GroupName = "${GROUP_NAME}"
    )

    var (
        // SchemeBuilder initializes a scheme builder
      SchemeBuilder = runtime.NewSchemeBuilder(addKnowTypes)
	  localSchemeBuilder = &SchemeBuilder
        // AddToScheme is a global function that registers this API group & version to a scheme
      AddToScheme = SchemeBuilder.AddToScheme
    )

    var (
        // SchemeGroupPROJECT_VERSION is group version used to register these objects
      SchemeGroupVersion = schema.GroupVersion{Group:  GroupName, Version: Version}
    )

    // Resource takes an unqualified resource and returns a Group-qualified GroupResource.
    func Resource(resource string)schema.GroupResource{
      return SchemeGroupVersion.WithResource(resource).GroupResource()
    }

    // Kind takes an unqualified kind and returns back a Group qualified GroupKind
    func Kind(kind string)schema.GroupKind{
      return SchemeGroupVersion.WithKind(kind).GroupKind()
    }

    // addKnownTypes adds the set of types defined in this package to the supplied scheme.
    func addKnowTypes(scheme *runtime.Scheme)error{
      scheme.AddKnownTypes(SchemeGroupVersion,
        new(${RESOURCE_NAME}),
        new(${RESOURCE_NAME}List),)
      metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
      return nil
    }
EOF
    gofmt -w ${gendir}/register.go
}
function fn_gen_gofile_group_internal_register(){
    PROJECT_NAME=$1
    GROUP_NAME=$2
    GROUP_VERSION=$3
    RESOURCE_KIND=$4
    gendir="${OUTDIR}/pkg/apis/${GROUP_NAME}/"
    genfile="${gendir}/register.go"
    if [ -e "${genfile}" ]; then
        return 0
    fi
    mkdir -pv ${gendir}
    echo "Generate ${genfile}"
    cat >> ${genfile} << EOF
/*
Copyright `date "+%Y"` The ${PROJECT_NAME} Authors.
Licensed under the Apache License, PROJECT_VERSION 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

    package ${GROUP_NAME}

    import (

      "k8s.io/apimachinery/pkg/runtime"
      "k8s.io/apimachinery/pkg/runtime/schema"
    )
    const (
        GroupName = "${GROUP_NAME}"
    )

    var (
        // SchemeBuilder initializes a scheme builder
      SchemeBuilder = runtime.NewSchemeBuilder(addKnowTypes)
        // AddToScheme is a global function that registers this API group & version to a scheme
      AddToScheme = SchemeBuilder.AddToScheme
    )

    var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

    // Resource takes an unqualified resource and returns a Group-qualified GroupResource.
    func Resource(resource string)schema.GroupResource{
      return SchemeGroupVersion.WithResource(resource).GroupResource()
    }

    // Kind takes an unqualified kind and returns back a Group qualified GroupKind
    func Kind(kind string)schema.GroupKind{
      return SchemeGroupVersion.WithKind(kind).GroupKind()
    }

    // addKnownTypes adds the set of types defined in this package to the supplied scheme.
    func addKnowTypes(scheme *runtime.Scheme)error{
      scheme.AddKnownTypes(SchemeGroupVersion,
        new(${RESOURCE_NAME}),
        new(${RESOURCE_NAME}List),)
      return nil
    }
EOF
    gofmt -w ${genfile}
}
# install
function fn_gen_gofile_install_install(){
    PROJECT_NAME=$1
    GROUP_NAME=$2
    GROUP_VERSION=$3
    RESOURCE_KIND=$4
    gendir="${OUTDIR}/pkg/apis/install"
    genfile="${gendir}/install.go"
    echo "Generate ${genfile}"
      cat >> ${genfile} << EOF
/*
Copyright `date "+%Y"` The ${PROJECT_NAME} Authors.
Licensed under the Apache License, PROJECT_VERSION 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
      package install
      import (
      "$(fn_project_to_gomod ${PROJECT_NAME})/pkg/apis/${GROUP_NAME}/${GROUP_VERSION}"
      "k8s.io/apimachinery/pkg/runtime"
      utilruntime "k8s.io/apimachinery/pkg/util/runtime"
       )

    func Install(scheme *runtime.Scheme){
      utilruntime.Must(${GROUP_VERSION}.AddToScheme(scheme))
    }
EOF
    gofmt -w ${genfile} 
}


# generate hackscripts
mkdir -pv hack

function fn_gen_hack_script_docker() {
    PROJECT_NAME=$1
    GOVERSION=$(go env GOVERSION)
    mkdir -pv hack/docker
    cat >> hack/docker/codegen.dockerfile << EOF
FROM golang:${GOVERSION//go/}

ENV GO111MODULE=auto
ENV GOPROXY="https://goproxy.cn"

RUN go get k8s.io/code-generator; exit 0
WORKDIR /go/src/k8s.io/code-generator
RUN go get -d ./...

RUN mkdir -p /go/src/$(fn_project_to_gomod)
VOLUME /go/src/$(fn_project_to_gomod)

WORKDIR /go/src/$(fn_project_to_gomod)
EOF
}
# create kubernetes builder images



function fn_gen_hack_tools_gofile() {
    PROJECT_NAME=$(fn_strings_to_lower ${1})
    mkdir -pv hack/
    if [ -e "hack/tools.go" ]; then
        return 0
    fi
    cat >> hack/tools.go << EOF
// +build tools

package tools

import _ "k8s.io/code-generator"
EOF
}

function fn_gen_hack_boilerplate() {
    PROJECT_NAME=$(fn_strings_to_lower ${1})
    mkdir -pv hack/
    if [ -e "hack/boilerplate.go.txt" ]; then
        return 0
    fi
    cat >> hack/boilerplate.go.txt << EOF
/*
Copyright The Kubernetes Authors.

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
EOF
}

function fn_gen_hack_script_upgrade(){
    PROJECT_NAME=$1
    GROUP_NAME=$2
    GROUP_VERSION=$3
    mkdir -pv hack/scripts
    if [ -e "hack/codegen-update.sh" ]; then
        return 0
    fi
    cat >> hack/codegen-update.sh << EOF
#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=\$(dirname "\${BASH_SOURCE[0]}")/..
CODEGEN_PKG=\${CODEGEN_PKG:-\$(cd "\${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

source "\${CODEGEN_PKG}/kube_codegen.sh"

THIS_PKG="${PROJECT_NAME}"

kube::codegen::gen_helpers \\
    --boilerplate "\${SCRIPT_ROOT}/hack/boilerplate.go.txt" \\
    "\${SCRIPT_ROOT}/pkg/apis"


kube::codegen::gen_openapi \\
    --output-dir "\${SCRIPT_ROOT}/pkg/client/openapi" \\
    --output-pkg "\${THIS_PKG}/pkg/client/openapi" \\
    --report-filename "\${report_filename:-"/dev/null"}" \\
    \${update_report:+"\${update_report}"} \\
    --boilerplate "\${SCRIPT_ROOT}/hack/boilerplate.go.txt" \\
    "\${SCRIPT_ROOT}/pkg/apis"

kube::codegen::gen_client \\
    --with-watch \\
    --with-applyconfig \\
    --output-dir "\${SCRIPT_ROOT}/pkg/client" \\
    --output-pkg "\${THIS_PKG}/pkg/client" \\
    --boilerplate "\${SCRIPT_ROOT}/hack/boilerplate.go.txt" \\
    "\${SCRIPT_ROOT}/pkg/apis"
EOF
}





echo "Begin generate some necessary code file"
# 生成group doc文件
fn_gen_gofile_group_version_doc ${PROJECT_NAME} ${GROUP_NAME} ${GROUP_VERSION} ${RESOURCE_KIND}
fn_gen_gofile_group_internal_doc ${PROJECT_NAME} ${GROUP_NAME} ${GROUP_VERSION} ${RESOURCE_KIND}
# 生成group types文件
fn_gen_gofile_group_versioned_types  ${PROJECT_NAME} ${GROUP_NAME} ${GROUP_VERSION} ${RESOURCE_KIND}
fn_gen_gofile_group_internal_types  ${PROJECT_NAME} ${GROUP_NAME} ${GROUP_VERSION} ${RESOURCE_KIND}

# 生成 register文件
fn_gen_gofile_group_versioned_register ${PROJECT_NAME} ${GROUP_NAME} ${GROUP_VERSION} ${RESOURCE_KIND}
fn_gen_gofile_group_internal_register ${PROJECT_NAME} ${GROUP_NAME} ${GROUP_VERSION} ${RESOURCE_KIND}
# 生成配置文件
#
## 生成相关的更新脚本文件
fn_gen_hack_tools_gofile ${PROJECT_NAME}
fn_gen_hack_boilerplate
# 生成dockerfile文件
fn_gen_hack_script_docker ${PROJECT_NAME}
# 生成脚本文件
fn_gen_hack_script_upgrade ${PROJECT_NAME} ${GROUP_NAME} ${GROUP_VERSION}

# 开始自动生成相关的代码
# bash hack/scripts/codegen-update.sh
#
#
# go mod tidy && go mod vendor
