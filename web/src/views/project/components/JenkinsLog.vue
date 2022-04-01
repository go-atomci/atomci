<template>
  <el-dialog
    ref="atomciDialog"
    class="atomciDialog"
    title="详细日志"
    :fullscreen="viewLog.fullscreen"
    :visible.sync="viewLog.visible"
    :append-to-body="true"
    :close-on-click-modal="false"
    :show-close="false"
    width="90%"
    center
  >
    <div slot="title" class="medium">
      <div class="selftitle">
        <span>详细日志</span>
        <div class="icons">
          <i :class="viewLog.fullscreen ? 'el-icon-bottom-left' : 'el-icon-full-screen'" @click="viewLog.fullscreen=!viewLog.fullscreen"></i>
          <i class="el-icon-close"  @click="viewLog.visible=false"></i>
        </div>
      </div>
    </div>
    <div class="dialogBody">
      <slot>
        <iframe :src="viewLog.url" frameborder="0" width="100%" scrolling="auto" class="viewlog" />
      </slot>
    </div>
  </el-dialog>
</template>

<script>

import backend from "@/api/backend";


export default {
  name: "JenkinsLog",
  data(){
    return {
      viewLog:{
        visible: false,
        fullscreen: false,
        url: null,
      }
    }
  },
  methods:{
    doCreate(jobName,runId,stageId){
      backend.getJenkinsServer((data) => {
        let jenkinsURL = data.jenkins;
        if (jenkinsURL.endsWith('/')) {
          jenkinsURL = jenkinsURL.slice(0, -1);
        }
        this.viewLog.url = `${jenkinsURL}/blue/organizations/jenkins/${jobName}/detail/${jobName}/${runId}/pipeline/`;
        this.viewLog.visible = true;
      }, stageId);
    },
  }
}
</script>

<style scoped>
  .viewlog {
    width: 100%;
    height: 100%;
    min-height: 600px;
  }
  .selftitle>.icons{
    display: flex;
    float: right;
    font-size:24px
  }
</style>
