
<style>
.app-panel {
  margin-bottom: 20px;
}

.appDetail .CodeMirror.cm-s-rubyblue {
  min-height: calc(100vh - 442px);
  max-width: calc(100vw - 240px);
}

.appDetail .el-tabs__header {
  margin-bottom: 5px;
}

.appDetail .portlet-body.mt10 .footer {
  margin-bottom: 0;
}

.react-grid-layout {
  overflow: hidden;
  margin-bottom: 20px;
  margin-top: 10px;
  height: 240px;
}
.react-grid-item {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  border: 1px solid #ebeef5;
  width: 24%;
  height: 230px;
  margin-right: 1%;
  float: left;
  color: #303133;
}
.react-grid-item:last-child {
  margin-right: 0px;
}
.panel-container {
  text-align: center;
}
.panel-container .panel-title {
  font-weight: 500;
  line-height: 2;
  color: #494c50;
}
.panel-container .panel-content {
  position: relative;
  display: table;
  width: 100%;
  height: 200px;
  line-height: 1;
}
.panel-container .panel-content .panel-txt {
  display: table-cell;
  vertical-align: middle;
  text-align: center;
  position: relative;
  z-index: 1;
  font-size: 3em;
  font-weight: 700;
  color: #494c50;
}
.gau_pod_chart {
  width: 100%;
  height: 100%;
  margin: 0 auto;
}
</style>

<template>
  <div class="page-content appDetail">
    <el-dialog top='25vh' :title="$t('bm.add.extendApp')" :visible.sync="dialogSacleFormVisible" size='tiny' class="commonDialog">
      <el-form>
        <el-form-item :label="$t('bm.add.podInsNum')">
          <el-input-number v-model="extendReplicas" :min="0" :max="100"></el-input-number>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogSacleFormVisible = false">{{$t('bm.other.cancel')}}</el-button>
        <el-button type="primary" @click="doPostScacleService">{{$t('bm.other.confirm')}}</el-button>
      </div>
    </el-dialog>
    <el-dialog top='25vh' :title="$t('bm.add.resetDeploy')" :visible.sync="dialogRollingUpdateFormVisible" size='tiny' class="commonDialog">
      <el-form>
        <!-- <label>当前镜像: {{ detailInfo.image }}</label> -->
        <el-form-item :label="$t('bm.add.selectDeployMarVer')">
          <el-select v-model="currImageTags" :placeholder="$t('bm.add.selectNeedMarVer')" style="width:100%" multiple @change="tagsChange">
            <el-option v-for="(item, index) in imageTagList" :key="index" :label="item.desc" :value="item.name">
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogRollingUpdateFormVisible = false">{{$t('bm.other.cancel')}}</el-button>
        <el-button type="primary" @click="doPostRollingUpdateService">{{$t('bm.other.confirm')}}</el-button>
      </div>
    </el-dialog>
    <div class="portlet-body">
      <common-delete ref="commonDelete" routerName='projectService' :routerParams="{projectID:$route.params.projectID}" isRefresh='true'></common-delete>
      <div class="title-content">
        <el-row>
          <el-col :span="9" class="title-panel">
            <span class="base-title"> {{detailInfo.name}}</span>{{$t('bm.deployCenter.baseInfo')}}
          </el-col>
          <el-col :span="1">
            &nbsp;
          </el-col>
          <el-col :span="14">
            <refresh v-on:getlist="getDetail"></refresh>
            <template>
              <!--<el-button :plain="true" type="primary" @click="doFlow">{{$t('bm.add.flowContrl')}}</el-button>-->
              <el-button :plain="true" type="danger" @click="$refs.commonDelete.doDelete('removeService',detailInfo.cluster,detailInfo.namespace,detailInfo.name)" icon="el-icon-delete" :title="$t('bm.other.delete')">{{$t('bm.other.delete')}}</el-button>
              <el-button :plain="true" type="warning" @click="reStart(detailInfo.cluster, detailInfo.namespace, detailInfo.name)" icon="el-icon-warning" :title="$t('bm.add.restart')">{{$t('bm.add.restart')}}</el-button>
              <el-button :plain="true" type="primary" @click="doSacle">{{$t('bm.add.levelExtend')}}</el-button>
              <!-- TODO: disabled 重新部署 tmp -->
              <!-- <el-button :plain="true" type="primary" @click="doRollingUpdateService">{{$t('bm.add.resetDeploy')}}</el-button> -->
            </template>
          </el-col>
        </el-row>
        <template :element-loading-text="$t('bm.add.loading')">
          <el-row>
            <el-col :span="8">
              {{$t('bm.deployCenter.type')}}：{{detailInfo.kind}}
            </el-col>
            <el-col :span="8">
              {{$t('bm.add.insNum')}}：{{detailInfo.replicas}}
            </el-col>
            <!-- TODO: disable creator -->
            <!-- <el-col :span="8">
              {{$t('bm.serviceM.creator')}}：{{detailInfo.creator}}
            </el-col> -->
          </el-row>
          <el-row>
            <el-col :span="8">
              {{$t('bm.add.beCluster')}}：{{detailInfo.cluster}}
            </el-col>
            <el-col :span="8">
              {{$t('bm.serviceM.resourceSpace')}}：{{detailInfo.namespace}}
            </el-col>
            <el-col :span="8">
              {{$t('bm.serviceM.creationTime')}}：{{detailInfo.create_at}}
            </el-col>
          </el-row>
          <el-row>
              <el-tooltip class="item" effect="dark" :content="detailInfo.image" placement="top-start">
                <el-col :span="8">
                  {{$t('bm.serviceM.mirror')}}：{{detailInfo.image}}
                </el-col>
              </el-tooltip>

              <el-col :span="8">
                {{$t('bm.serviceM.creationTime')}}：{{detailInfo.create_at}}
              </el-col>
            <el-col :span="8">
              {{$t('bm.add.updateTime')}}：{{detailInfo.update_at}}
            </el-col>
          </el-row>
        </template>
      </div>
    </div>
    <div class="portlet-body mt10">
      <div class="tabs-content">
        <el-tabs v-model="activeName" @tab-click="handleClick">
          <el-tab-pane :label="$t('bm.serviceM.podInstance')" name="pod">
            <app-pod :list="detailInfo.pods" :cluster="detailInfo.cluster" :appName="detailInfo.name"></app-pod>
          </el-tab-pane>
          <!-- 应用配置　临时禁用掉 -->
          <!-- <el-tab-pane :label="$t('bm.add.appConfig')" name="config">
            <app-config ref='appConfig' v-on:appCallBack="appCallBack"></app-config>
          </el-tab-pane> -->
          <el-tab-pane :label="$t('bm.add.queryLog')" name="log">
            <app-log ref='appLog' :appName="getAppName"></app-log>
          </el-tab-pane>
          <el-tab-pane :label="$t('bm.add.serviceAddress')" name="server">
            <div v-for="service in detailInfo.services" :key="service.name">
                <server-list ref='ServerList' :list="service"  :appName="getAppName" :appNamespace='$route.params.namespace'></server-list>
            </div>
          </el-tab-pane>
          <el-tab-pane :label="$t('bm.add.eventRecord')" name="recordApp">
            <event-record-List ref='eventrecordList' :appName="getAppName" :appNamespace='$route.params.namespace' :clusterName="$route.params.clusterName"></event-record-List>
          </el-tab-pane>
        </el-tabs>
      </div>
      <flow-control ref="create" :grayList="detailInfo.version_list" v-on:getlist="getDetail"></flow-control>
    </div>
  </div>
</template>
<script>
import { mapGetters } from 'vuex';
import backend from '@/api/backend';
import Refresh from '@/components/utils/Refresh';
import CommonDelete from '@/components/utils/Delete';


import AppPod from '../dialogService/AppPod';
import AppLog from '../dialogService/AppLog';
import FlowControl from '../dialogService/FlowControl';

import EventRecordList from '@/components/view/EventRecordList';

import ServerList from '../dialogService/ServerList';

const formData = {
  cpu_quota: '',
  memory_quota: '',
};
export default {
  data() {
    return {
      isFirstConfig: false,
      // 是否是第一次点击日志查询
      isFirstLog: false,
      isFirstMonitor: false,
      maxValueUnit: '',
      timeDetail: '',
      form: JSON.parse(JSON.stringify(formData)),
      // 是否是第一次点击应用告警
      isFirstWarning: false,
      // 是否是第一次点击负载均衡
      isFirstIngress: false,
      // 是否是第一次点击服务地址
      isFirstService: false,
      // 是否是第一次点击事件记录
      isFirstRecord: false,
      intervalId: null,
      curTabName: 'pod',
      activeName: 'pod',
      // 详细信息
      detailInfo: {},
      clusterInfo: {},
      imageTagList: [],
      currImageTag: '',
      currImageTags: [],
      // 水平扩展的个数，之所以要重新声明一个变量是因为，修改时不会影响界面其它显示
      extendReplicas: 0,
      dialogSacleFormVisible: false,
      dialogRollingUpdateFormVisible: false,
      rangeInfo: null,
    };
  },
  components: {
    Refresh,
    AppPod,
    // AppConfig,
    AppLog,
    ServerList,
    EventRecordList,
    FlowControl,
    CommonDelete
  },
  activated() {
    this.getDetail();
  },
  computed: {
    // 获得路由中当中appName
    getAppName() {
      const curName = this.$route.params.appName;
      return curName;
    },
    ...mapGetters({
      loading: 'getLoading',
    }),
    appDeploymentName() {
      const { name, version } = this.detailInfo;
      return version ? `${name}-${version}` : name;
    },
  },
  methods: {
    handleClick(tab) {
      if (this.curTabName !== tab.name) {
        this.curTabName = tab.name;
        switch (true) {
          case tab.index === '2' && !this.isFirstService:
            // this.$refs.ServerList.getList();
            this.isFirstService = true;
            break;
          case tab.index === '1' && !this.isFirstLog:
            // this.$refs.appLog.getList();
            this.renderLog();
            this.isFirstLog = true;
            break;
          case tab.index === '3' && !this.isFirstRecord:
            this.$refs.eventrecordList.getList();
            this.isFirstRecord = true;
            break;
        }
      }
    },
    renderLog() {
      if (this.detailInfo.pods) {
        this.$refs.appLog.doSelectPodName(this.detailInfo);
        this.isFirstLog = true;
        // window.clearInterval(this.intervalId);
      } else {
        // const curInterId = window.setInterval(() => {
        // this.renderLog();
        // window.clearInterval(curInterId);
        // }, 500);
      }
    },
    // 水平扩展
    doPostScacleService() {
      const replicas = this.extendReplicas;
      backend.scaleService(
        this.$route.params.clusterName,
        this.$route.params.namespace,
        this.getAppName,
        replicas,
        () => {
          this.$notify({
            title: this.$t('bm.add.success'),
            message: `${this.$t('bm.add.podsInsNumEx')}: ${replicas}`,
            type: 'success',
          });
          this.getDetail();
          // 这句话一定要写在refresh后面，因为服务器端有延时
          this.detailInfo.replicas = this.extendReplicas;
          this.dialogSacleFormVisible = false;
        }
      );
    },
    // 弹出水平扩展对话框
    doSacle() {
      this.extendReplicas = this.detailInfo.replicas;
      this.dialogSacleFormVisible = true;
    },
    tagsChange(keys) {
      const key = keys[keys.length - 1];
      for (let i = this.currImageTags.length - 2; i > -1; i--) {
        if (this.currImageTags[i].split('^')[0] === key.split('^')[0]) {
          this.currImageTags.splice(i, 1);
        }
      }
    },
    // 确认是否重启对话框
    reStart(cluster, namespace, appname) {
      this.$confirm(this.$t('bm.add.isResetAllIns'), this.$t('bm.add.hint'), {
        confirmButtonText: this.$t('bm.other.confirm'),
        cancelButtonText: this.$t('bm.other.cancel'),
        type: 'warning',
      }).then(() => {
        backend.reStart(cluster, namespace, appname, () => {
          this.$message({
            type: 'success',
            message: this.$t('bm.add.resetSuc'),
          });
        });
      });
    },
    // 弹出滚动升级对话框，获得升级的镜像版本
    doRollingUpdateService() {
      this.$store.dispatch('setNeedLoading', false);
      this.dialogRollingUpdateFormVisible = true;
      const ps = [];
      if (this.detailInfo.pods === undefined) {
        return;
      }
      if (this.detailInfo.pods.length === 0) {
        return;
      }
      for (const a of this.detailInfo.pods[0].containers) {
        ps.push(
          new Promise((resolve) => {
            backend.getRepositoryTags(
              this.clusterInfo.registry,
              `${a.image.split('/')[1]}/${a.image.split('/')[2].split(':')[0]}`,
              (data) => {
                resolve(data);
              }
            );
          })
        );
      }
      Promise.all(ps).then((result) => {
        this.$store.dispatch('setNeedLoading', true);
        this.imageTagList = [];
        const arr = [];
        for (let i = 0; i < this.detailInfo.pods[0].containers.length; i++) {
          const cur = this.detailInfo.pods[0].containers[i];
          const image = cur.image.substr(0, cur.image.lastIndexOf(':'));
          for (const a of result[i]) {
            arr.push({
              desc: `${this.$t('bm.add.container')}：${cur.name}  --> ${this.$t('bm.serviceM.mirror')}：${image}:${a.tag}  --> 创建时间：${a.create_at}`,
              name: `${cur.name}^${image}:${a.tag}`,
            });
          }
        }
        this.imageTagList = arr;
      });
    },
    // 滚动升级
    doPostRollingUpdateService() {
      // let containerImage = this.detailInfo.image.split("/")[0] + "/" + this.detailInfo.image.split("/")[1] + "/" + this.detailInfo.image.split("/")[2].split(":")[0] + ":" + this.currImageTag
      if (this.currImageTags.length === 0) {
        this.$notify.closeAll();
        this.$notify({
          title: this.$t('bm.add.error'),
          message: this.$t('bm.add.mSelMarVer'),
          type: 'error',
        });
        return;
      }
      const postData = [];
      for (const a of this.currImageTags) {
        const arr = a.split('^');
        postData.push({ name: arr[0], image: arr[1] });
      }
      backend.rollingUpdateService(
        this.$route.params.clusterName,
        this.$route.params.namespace,
        this.getAppName,
        postData,
        () => {
          this.$notify({
            title: this.$t('bm.add.success'),
            message: this.$t('bm.add.appVerResetDeploy'),
            type: 'success',
          });
          this.dialogRollingUpdateFormVisible = false;
          this.getDetail();
        }
      );
    },
    // 获得应用数据返回时的回调函数
    appCallBack(data) {
      this.detailInfo = data;
      if (this.isFirstLog) {
        this.renderLog();
      }
    },
    getDetail() {
      if (this.detailInfo.deploy_status !== 'switching') {
        clearInterval(this.timeDetail);
      }
      backend.getServiceInspect(
        this.$route.params.clusterName,
        this.$route.params.namespace,
        this.getAppName,
        (data) => {
          this.appCallBack(data);
        }
      );
      // TODO: cluster inspect get disabled tmp.
      // backend.getClusterInspect(this.$route.params.clusterName, (data) => {
      //   this.clusterInfo = data;
      // });
    },
  },
};
</script>
