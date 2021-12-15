<style scoped="scoped">
.projectPubDetail .iframe {
  width: 100%;
  height: calc(100vh - 30px);
  transform: translate(0px, -130px);
  margin-bottom: -146px;
}
.projectPubDetail .portlet-body {
  overflow: hidden;
}

/* 下面的这一段先暂时做隐藏吧 */
.projectPubDetailContainer .page-header,
.projectPubDetailContainer .main-left,
.projectPubDetailContainer .page-footer,
.projectPubDetailContainer .page-bar {
  display: none;
}
.projectPubDetailContainer .main-right {
  margin-left: 0 !important;
}
.projectPubDetailContainer .page-main .main-right .view {
  height: calc(100vh - 20px);
  padding-bottom: 0;
}
.projectPubDetailContainer .portlet-body {
  height: calc(100vh - 30px);
  padding: 0;
  margin: 0;
}
.projectPubDetailContainer .iframe {
  height: calc(100vh + 100px);
}
.projectPubDetail .platformMenu {
  display: none;
}
</style>

<template>
  <div class="page-content projectPubDetail">
    <div class="portlet-body" v-loading.body="loading" :element-loading-text="$t('bm.add.loading')">
      <iframe class="iframe" :src='curSrc'>
      </iframe>
    </div>
  </div>
</template>
<script>
import backend from '@/api/backend';

export default {
  data() {
    return {
      loading: true,
      curSrc: '',
    };
  },
  components: {},
  mounted() {
    const container = document.getElementById('app');
    container.className = `${container.className} projectPubDetailContainer`;
    setTimeout(() => {
      this.loading = false;
    }, 1000);
  },
  created() {
    const projectId = this.$route.params.projectId;
    const jobName = this.$route.params.jobName;
    const runId = this.$route.params.runId;
    const stageId = this.$route.params.stageId;
   
    backend.getJenkinsServer((data) => {
      let jenkinsURL = data.jenkins
      if (jenkinsURL.endsWith('/')) {
        jenkinsURL = jenkinsURL.slice(0,-1);
      }
      this.curSrc = `${jenkinsURL}/blue/organizations/jenkins/${jobName}/detail/${jobName}/${runId}/pipeline/`;
    },stageId);
  },
};
</script>
