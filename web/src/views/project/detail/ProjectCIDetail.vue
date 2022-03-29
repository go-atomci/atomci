<template>
    <div class="page-content">
      <div class="portlet-body mt10 min-height150">
        <template>
          <div class="pv10" v-show="proShow">
            <el-row>
              <el-col :span="24">
                <div class="setTitle">
                  <div class="f-r" v-if="projectInfo.operations">
                    <el-button @click="$refs.mark.show(projectInfo)" v-if="projectInfo.operations['manual']">人工卡点</el-button>
                    <el-button @click="$refs.merge.show(projectInfo)" v-if="projectInfo.operations['merge-branch']">合并分支</el-button>
                    <el-button @click="$refs.goToBuild.doShows(projectInfo.id, projectInfo.stage_id, projectInfo.step)" v-if="projectInfo.operations['build']">构建</el-button>
                    <el-button @click="$refs.deploy.doShows(projectInfo)" v-if="projectInfo.operations['deploy']">部署</el-button>
                    <el-button @click="termination(projectInfo.id, projectInfo.stage_id, projectInfo.step_type)" v-if="projectInfo.operations['terminate']">终止</el-button>
                    <el-button @click="$refs.nextstage.show(projectInfo.id, projectInfo.stage_id)" v-if="projectInfo.operations['next-stage']">{{$t('bm.add.enterNextStage')}}</el-button>
                    <el-button @click="$refs.backto.show(projectInfo.id, projectInfo.stage_id)" v-if="projectInfo.operations['back-to']">{{$t('bm.add.rollBack')}}</el-button>
                    <el-button v-show="statusCheck" @click="gotoclose(projectInfo.id)">{{$t('bm.deployCenter.achieved')}}</el-button>
                    <el-button @click="$refs.publishEdit.doShows(projectInfo)">{{$t('bm.other.edit')}}</el-button>
                    <el-button type="danger" @click="gotodelete(projectInfo.id)">{{$t('bm.other.delete')}}</el-button>
                    <el-button type="el-button el-button--success is-plain" @click="getVersionInfo">{{$t('bm.other.refresh')}}</el-button>
                  </div>
                  {{projectInfo.version_no}}<br /><span class="setDescription">{{projectInfo.name}}</span>
                </div>
              </el-col>
            </el-row>
            <el-row>
              <el-col class="setText">
                创建者：{{projectInfo.creator}}
              </el-col>
              <el-col class="setText">
                创建时间：{{projectInfo.create_at}}
              </el-col>
            </el-row>
            <el-row>
              <el-col class="setText">
                开始时间：{{projectInfo.start_at}}
              </el-col>
              <el-col class="setText">结束时间：{{projectInfo.end_at}}</el-col>
            </el-row>
            <el-row>
              <el-col class="setText">版本状态：{{projectInfo.statusName}}</el-col>
              <el-col class="setText">发布流程：{{projectInfo.pipeline_name}}</el-col>
            </el-row>
          </div>
        </template>
      </div>
      <div class="portlet-body mt10 pv20 min-height150" v-if="statusCheck">
        <div class="setTitle" style="text-align: -webkit-left">当前环境 (<span class="setDescription">{{projectInfo.stage_name}}</span>)</div>
        <div class="mt10 clearfix">
          <template v-if="projectInfo.steps">
            <template v-for="(item,index) in projectInfo.steps">
                <div v-if="index != 0" class="arrow-right">
                  <div class="arrow-line"></div>
                </div>
                <div class="env-panel">
                  <i class="cricle-bg" :class="iconClass[item.status]"></i>
                  {{item.name}}
                </div>
            </template>
          </template>
          <template v-else>
            <div class="text-c pv20">暂无任务</div>
          </template>
        </div>
      </div>
      <div class="portlet-body mt10 pv10">
        <el-tabs v-model="activeName">
          <el-tab-pane label="操作历史" name="history">
            <template>
              <el-table
                stripe
                :data="historyData"
                style="width: 100%">
                <el-table-column
                  prop="creator"
                  :label="$t('bm.add.operater')"
                  :show-overflow-tooltip=true
                  min-width="10%">
                </el-table-column>
                <el-table-column
                  prop="type"
                  :label="$t('bm.deployCenter.type')"
                  :show-overflow-tooltip=true
                  min-width="10%">
                </el-table-column>
                <el-table-column
                  prop="update_at"
                  label="操作时间"
                  :show-overflow-tooltip=true
                  min-width="12%">
                </el-table-column>
                <el-table-column
                  prop="stage"
                  :label="$t('bm.add.stage')"
                  :show-overflow-tooltip=true
                  min-width="10%">
                </el-table-column>
                <el-table-column prop="step"
                  label="步骤"
                  :show-overflow-tooltip=true
                  min-width="10%">
                </el-table-column>
                <el-table-column
                  prop="status"
                  :label="$t('bm.deployCenter.status')"
                  :show-overflow-tooltip=true
                  min-width="10%">
                </el-table-column>
                <el-table-column
                  prop="message"
                  :label="$t('bm.add.info')"
                  :show-overflow-tooltip=true
                  min-width="10%">
                  <template slot-scope="scope">
                    {{scope.row.message}}
                    <el-button v-if="scope.row.run_id" @click="viewFile(scope.row)" type="text" size="small">详细日志</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </template>
            <page-nav ref="pages" :list="historyData" v-on:getlist="getList"></page-nav>
          </el-tab-pane>
          <el-tab-pane label="应用列表" name="appList">
            <p style="float: left" v-if="statusCheck"><el-button type="primary" @click="$refs.versionAdd.doShow($route.params.projectID, $route.params.versionId, projectInfo.version_no)">添加应用</el-button></p>
            <el-table class="mt10"  :data="projectInfo.apps">
              <el-table-column prop="name" :label="$t('bm.deployCenter.name')" sortable min-width="12%" :show-overflow-tooltip="true" />
              <el-table-column prop="type" :label="$t('bm.deployCenter.type')" sortable min-width="12%" :show-overflow-tooltip="true" />
              <el-table-column prop="language" :label="$t('bm.deployCenter.language')" sortable min-width="10%" :show-overflow-tooltip="true" />
              <el-table-column prop="branch_name" :label="$t('bm.deployCenter.releaseBran')" min-width="11%" :show-overflow-tooltip="true" />
              <el-table-column prop="compile_command" :label="$t('bm.deployCenter.buildCompile')" min-width="11%" :show-overflow-tooltip="true" />
              <el-table-column :label="$t('bm.deployCenter.operation')" min-width="10%" v-if="statusCheck">
                <template slot-scope="scope">
                  <el-button type="text" @click="removeApp(scope.row.id)">移除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </div>
      <project-mark ref="mark" v-on:getprojectReleaseList="getVersionInfo"></project-mark>
      <back-to ref="backto" v-on:getprojectReleaseList="getVersionInfo"></back-to>
      <project-deploy ref="deploy" v-on:getprojectReleaseList="getVersionInfo"></project-deploy>
      <next-stage v-on:getprojectReleaseList="getVersionInfo" ref="nextstage"></next-stage>
      <to-build v-on:getprojectReleaseList="getVersionInfo" ref="goToBuild"></to-build>
      <version-add v-on:getlist="getVersionInfo" ref="versionAdd"></version-add>
      <publish-edit v-on:getPublishBaseInfo="getVersionBaseInfo" ref="publishEdit"></publish-edit>
      <el-dialog
      v-dialogDrag
      ref="atomciDialog"
      class="atomciDialog"
      title="详细日志"
      :fullscreen="viewLog.fullscreen"
      :visible.sync="viewLog.visiable"
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
            <i class="el-icon-close"  @click="viewLog.visiable=false"></i>
          </div>
        </div>
      </div>
      <div class="dialogBody">
        <slot>
          <iframe :src="viewLog.url" frameborder="0" width="100%" scrolling="auto" class="viewlog" />
        </slot>
      </div>
    </el-dialog>
    </div>
</template>

<style scoped>
  .min-height150 {
    min-height: 150px;
  }
  .pv10 {
    padding-top: 10px;
    padding-bottom: 10px;
  }
  .pl16 {
    padding-left: 16px;
  }
  .f-r {
    float: right;
  }
  .pv20 {
    padding: 20px;
  }
  .mb15 {
    margin-bottom: 15px;
  }
  .setTitle {
    font-size: 18px;
    color: #333;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: bold;
    line-height: 25px;
  }
  .setDescription {
    color: #606266;
    font-size: 14px;
    font-weight: 400;
  }
  .fontSmall {
    font-size: 14px;
    font-family: PingFangSC-Regular,PingFang SC;
    font-weight: 400;
    line-height: 20px;
    color: #909399;
  }
  .setText {
    width: 400px;
    line-height: 20px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    margin-right: 10px;
    margin-top: 15px;
    color: #333;
    font-weight:400;
  }
  .font-title {
    font-size: 16px;
    color: #409EFF;
    display: inline-block;
    line-height: 50px;
  }
  .el-form-item {
    margin-bottom: 5px;
  }
  .env-panel {
    float: left;
    width: 180px;
    height: 65px;
    background: rgba(255,255,255,1);
    box-shadow: 0px 2px 4px 0px rgba(0,0,0,0.12);
    border-radius: 4px;
    border: 1px solid rgba(242,246,252,1);
    padding: 16px;
    margin-top: 15px;
    font-size: 14px;
    font-family: PingFangSC-Regular,PingFang SC;
    color: #333;
  }
  .arrow-right {
    float: left;
    position: relative;
    width: 50px;
    height: 65px;
    margin-top: 15px;
    padding: 0 5px;
  }
  .arrow-right:after {
    position: absolute;
    display: inline-block;
    top: 26px;
    right: 5px;
    width: 0;
    height: 0px;
    content: '';
    border-style: solid;
    border-width: 5px;
    border-color: #DCDFE6 transparent transparent #DCDFE6;
    transform: rotate(135deg);
  }
  .arrow-line {
    height: 1px;
    overflow: hidden;
    background-color: #DCDFE6;
    margin-top: 30px;
  }
  .cricle-bg {
    display: inline-block;
    width: 32px;
    height: 32px;
    border-radius: 32px;
    vertical-align: middle;
    margin-right: 10px;
    background: left top no-repeat;
    background-size: cover;
  }
  .env-error {
    background-image: url(../../../assets/env_error.png);
  }
  .env-done {
    background-image: url(../../../assets/env_done.png);
  }
  .env-next {
    background-image: url(../../../assets/env_next.png);
  }
  .env-run {
    background-image: url(../../../assets/env_run.png);
  }
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
<script>
  import { Message, MessageBox } from 'element-ui';
  import backend from '@/api/backend';
  import PageNav from '@/components/utils/PageList';
  import Utils from '@/common/utils';

  import NextStage from '../dialogCI/ProjectNextStage'; // 下一个阶段
  import ProjectMark from '../dialogCI/ProjectMark'; // 人工卡点
  import BackTo from '../dialogCI/ProjectBackTo'; // 回退
  import ToBuild from '../dialogCI/ProjectBuild'; //构建
  import ProjectDeploy from '../dialogCI/ProjectDeploy'; // 部署
  import versionAdd from '../dialogCI/ProjectVersionAdd'; //添加应用
  import PublishEdit from '../dialogCI/PublishEdit'; // 编辑版本



export default {
  data() {
    return {
      proShow: false,
      activeName: 'history',
      projectInfo: {},
      historyData: [],
      statePublish: [
        '失败', '成功', '进行中', '待执行', '结束', '已归档', '未知', '终止成功', '终止失败', '合并失败', '不支持'
      ],
      iconClass: ['env-error', 'env-done', 'env-run', 'env-next', 'env-done', 'env-done', 'env-next', 'env-done', 'env-error', 'env-error', 'env-error'],
      statusCheck: false,
      viewLog: {
        visiable: false,
        fullscreen: false,
        url: null,
      },
    };
  },
  components: {
    BackTo,
    NextStage,
    ToBuild,
    ProjectMark,
    ProjectDeploy,
    PageNav,
    versionAdd,
    PublishEdit,
  },
  mounted() {
    this.getVersionInfo();
  },
  methods: {
    getVersionInfo() {
      backend.getListdetail(this.$route.params.projectID, this.$route.params.versionId, (data) => {
        if(data) {
          data.create_at = Utils.format(new Date(data.create_at), 'yyyy-MM-dd hh:mm:ss');
          data.start_at = Utils.format(new Date(data.start_at), 'yyyy-MM-dd hh:mm:ss');
          data.end_at = data.end_at ? Utils.format(new Date(data.end_at), 'yyyy-MM-dd hh:mm:ss') : '';
          let statusName = '';
          if(data.status >= 0) statusName = this.statePublish[data.status];
          this.projectInfo = Object.assign({}, data, {'statusName': statusName});
          this.proShow = true;
          this.getList();
          this.statusCheck = data.status === 5 ? false : true;
        }
      });
    },
    getVersionBaseInfo() {
      backend.getListdetail(this.$route.params.projectID, this.$route.params.versionId, (data) => {
        if(data) {
         this.projectInfo.name = data.name
         this.projectInfo.version_no = data.version_no
        }
      });
    },
    // 终止发布
    termination(id, stargid, name) {
      MessageBox.confirm(this.$t('bm.add.sureStopPublishCodeModule'), this.$t('bm.add.hint'), { type: 'warning' }).then(() => {
        const params = {
          action_name: 'terminate'
        };
        const that = this;
        if(name === 'build') {
          backend.setBuildMerge(this.$route.params.projectID, id, stargid, params, (data) => {
            Message.success(this.$t('bm.add.optionSuc'));
            that.getVersionInfo();
          });
        } else if(name === 'deploy') {
          backend.setDeploy(this.$route.params.projectID, id, stargid, params, (data)=> {
            Message.success(this.$t('bm.add.optionSuc'));
            that.getVersionInfo();
          });
        }
      }).catch(() => {
      });
    },
    // 关闭-归档
    gotoclose(id) {
      MessageBox.confirm(this.$t('bm.add.isAgreeAchieved'), this.$t('bm.infrast.tips'), { type: 'warning' }).then(() => {
        const that = this;
        backend.goclose(this.$route.params.projectID, id, (data) => {
          Message.success(this.$t('bm.add.optionSuc'));
          that.getVersionInfo();
        });
      }).catch(() => {
      });
    },
    // 删除
    gotodelete(id) {
      MessageBox.confirm(this.$t('bm.add.isAgreeDelete'), this.$t('bm.infrast.tips'), { type: 'warning' }).then(() => {
        const that = this;
        backend.getDeletionPublish(this.$route.params.projectID, id, (data) => {
          Message.success(this.$t('bm.add.optionSuc'));
          that.$router.push({
            name: 'projectCI',
            params: {
              projectID: this.$route.params.projectID
            }
          });
        });
      }).catch(() => {});
    },
    getList() {
      const params = {
        page_size: this.$refs.pages.pageSize,
        page_index: this.$refs.pages.currentPage,
      };
      backend.getOperationLog(this.$route.params.projectID, this.$route.params.versionId, params, (data) => {
        this.historyData = data.item;
        this.historyData.map((i) => {
          if(i.status >= 0) i.status = this.statePublish[i.status];
          i.update_at = Utils.format(new Date(i.update_at), 'yyyy-MM-dd hh:mm');
          i.step = i.step === 'back-to' ? '回退' : i.step;
          i.step = i.step === 'next-stage' ? '流转至下一阶段' : i.step;
          i.type = i.type ? i.type : '手动触发';
          if(i.creator === 'system') i.type="系统触发";
        });
        this.$refs.pages.total = data.total;
      });
    },
    viewFile(item) {
      backend.getJenkinsServer((data) => {
        let jenkinsURL = data.jenkins;
        if (jenkinsURL.endsWith('/')) {
          jenkinsURL = jenkinsURL.slice(0, -1);
        }
        this.viewLog.url = `${jenkinsURL}/blue/organizations/jenkins/${item.job_name}/detail/${item.job_name}/${item.run_id}/pipeline/`;
        this.viewLog.visiable = true;
      }, item.stage_id);
    },
    // 移除应用
    removeApp(id) {
      MessageBox.confirm(this.$t('bm.add.sureDelete'), this.$t('bm.infrast.tips'), { type: 'warning' })
      .then(() => {
        backend.removeApp(this.$route.params.projectID, this.$route.params.versionId, id, (data) => {
          Message.success(this.$t('bm.add.optionSuc'));
          this.getVersionInfo();
        });
      })
      .catch(() => {});
    },
  }
}
</script>