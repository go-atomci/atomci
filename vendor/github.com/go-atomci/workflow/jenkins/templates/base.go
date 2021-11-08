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

// BaseXML ..
const BaseXML = `
<flow-definition plugin="workflow-job@2.32">
    <actions>
        <io.jenkins.blueocean.service.embedded.BlueOceanUrlAction_-DoNotShowPersistedBlueOceanUrlActions plugin="blueocean-rest-impl@1.16.0"/>
        <org.jenkinsci.plugins.pipeline.modeldefinition.actions.DeclarativeJobPropertyTrackerAction plugin="pipeline-model-definition@1.3.8">
            <jobProperties/>
            <triggers/>
            <parameters/>
            <options/>
        </org.jenkinsci.plugins.pipeline.modeldefinition.actions.DeclarativeJobPropertyTrackerAction>
        <org.jenkinsci.plugins.pipeline.modeldefinition.actions.DeclarativeJobAction plugin="pipeline-model-definition@1.3.8"/>
    </actions>
    <description>atomci jenkins pipeline</description>
    <keepDependencies>false</keepDependencies>
    <properties>
        <jenkins.model.BuildDiscarderProperty>
            <strategy class="hudson.tasks.LogRotator">
                <daysToKeep>-1</daysToKeep>
                <numToKeep>10</numToKeep>
                <artifactDaysToKeep>-1</artifactDaysToKeep>
                <artifactNumToKeep>-1</artifactNumToKeep>
            </strategy>
        </jenkins.model.BuildDiscarderProperty>
    </properties>
    <definition class="org.jenkinsci.plugins.workflow.cps.CpsFlowDefinition" plugin="workflow-cps@2.69">
        <script>
            {{ .Pipeline }}
        </script>
        <sandbox>false</sandbox>
    </definition>
    <triggers/>
    <disabled>false</disabled>
</flow-definition>
`

// DeployBaseXML ..
const DeployBaseXML = `
<flow-definition plugin="workflow-job@2.32">
    <actions>
        <io.jenkins.blueocean.service.embedded.BlueOceanUrlAction_-DoNotShowPersistedBlueOceanUrlActions plugin="blueocean-rest-impl@1.16.0"/>
        <org.jenkinsci.plugins.pipeline.modeldefinition.actions.DeclarativeJobPropertyTrackerAction plugin="pipeline-model-definition@1.3.8">
            <jobProperties/>
            <triggers/>
            <parameters/>
            <options/>
        </org.jenkinsci.plugins.pipeline.modeldefinition.actions.DeclarativeJobPropertyTrackerAction>
        <org.jenkinsci.plugins.pipeline.modeldefinition.actions.DeclarativeJobAction plugin="pipeline-model-definition@1.3.8"/>
    </actions>
    <description>deploy pipeline</description>
    <keepDependencies>false</keepDependencies>
    <properties>
        <jenkins.model.BuildDiscarderProperty>
            <strategy class="hudson.tasks.LogRotator">
                <daysToKeep>-1</daysToKeep>
                <numToKeep>60</numToKeep>
                <artifactDaysToKeep>-1</artifactDaysToKeep>
                <artifactNumToKeep>-1</artifactNumToKeep>
            </strategy>
        </jenkins.model.BuildDiscarderProperty>
    </properties>
    <definition class="org.jenkinsci.plugins.workflow.cps.CpsFlowDefinition" plugin="workflow-cps@2.69">
        <script>
            {{ .Pipeline }}
        </script>
        <sandbox>false</sandbox>
    </definition>
    <triggers/>
    <disabled>false</disabled>
</flow-definition>
`
