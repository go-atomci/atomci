<style scoped>
  .el-menu--horizontal>.el-menu-item.is-active {
    color: #000 !important;
  }

  .el-menu--horizontal>.el-menu-item:hover {
    color: #000 !important;
  }

  .statusList {
    margin-bottom: 10px;
  }

  .statusList li {
    display: flex;
    height: 55px;
    border-bottom: 1px dashed #d9edf6;
  }

  .statusList li .pubOperate {
    width: 226px;
    align-self: center;
    text-align: right;
  }

  .statusList li .pubName {
    flex: 2;
  }

  .statusList li .pubStatus {
    flex: 3;
    padding-right: 10px;
    align-self: center;
  }
  .ulSteps li {
    float: left;
    position: relative;
    width: 33.3%;
    text-align: center;
    font-size: 12px;
    font-family: PingFangSC-Regular,PingFang SC;
    font-weight: 400;
    color: rgba(192,196,204,1);
  }
  .grayBorder {
    height: 20px;
    overflow: hidden;
    border-bottom: 1px solid #C0C4CC;
    margin-bottom: 10px;
  }
  .grayCircle {
    position: absolute;
    left: 50%;
    margin-left: -8px;
    top: 11px;
    width: 16px;
    height: 16px;
    line-height: 15px;
    border-radius: 16px;
    background: rgba(192,196,204,1);
    color: #fff;
    font-size: 12px;
  }
  .ulSteps li.nowSteps {
    color: #303133;
  }
  .nextSteps .grayCircle {
    background-color: #fff;
    border: 1px dashed #C0C4CC;
  }
  .nextSteps .grayBorder {
    border-bottom: 1px dashed #C0C4CC;
  }
  .circleSuccess, .circleRunning {
    display: inline-block;
    width: 8px;
    height: 8px;
    border-radius: 8px;
    overflow: hidden;
    background-color: #2DCA93;
  }
  .circleRunning {
    background-color: red;
  }
  .deploy-done, .deploy-run {
    display: inline-block;
    width: 16px;
    height: 16px;
    background: url(../../assets/deploy_done.png) left top no-repeat;
    background-size: cover;
  }
  .deploy-run {
    background-image: url(../../assets/deploy_run.png);
  }
</style>
<template>
    <div class="portlet-body">
      <div class="table-toolbar" style="text-align: left">
        <el-button :plain="false" type="primary" @click="$refs.pubDialog.doShow()">
          <i class='icon-plus' /> +{{$t('bm.add.createPipeline')}}</el-button>
        <refresh v-on:getlist="getList('clear')"></refresh>
      </div>
      <template>
        <el-tabs v-model="activeName" @tab-click="tabClick" class="mt10">
          <el-tab-pane v-for="(item, index) in stageList" :key="index" :label="item.name" :name="`${item.id}`"></el-tab-pane>
        </el-tabs>
      </template>
      <el-row class="mt10">
          <el-col :span="2" class="search-name">{{$t('bm.deployCenter.pipelineName')}}</el-col>
          <el-col :span="4">
            <el-input v-model="form.versionNo" :placeholder="$t('bm.deployCenter.pipelineSecarchName')" @keyup.enter.native="getList" filterable auto-complete="off"></el-input>
          </el-col>

          <el-col :span="2" class="search-name">{{$t('bm.deployCenter.pipelineDesc')}}</el-col>
          <el-col :span="4">
            <el-input v-model="form.name" placeholder="请输入流水线描述" @keyup.enter.native="getList" filterable auto-complete="off"></el-input>
          </el-col>

          <el-col :span="2" class="search-name">状态</el-col>
          <el-col :span="4">
            <el-select v-model="form.status" placeholder="请选择状态" clearable filterable @change="getList">
              <el-option v-for="(item, index) in statePublish" :key="index" :label="item" :value="index">
              </el-option>
            </el-select>
          </el-col>
        <el-col :span="4">
            <el-button type="primary" @click="getList">搜索</el-button>
            <el-button class="font-gray" type="text" @click="getList('clear')">重置</el-button>
        </el-col>
      </el-row>
      <!-- <el-row>
        <el-col class="w320 mt10 mr6">
          <el-col :span="6" class="search-name">当前步骤</el-col>
          <el-col :span="18">
            <el-input v-model="form.step" placeholder="请输入当前步骤" @keyup.enter.native="getList" filterable auto-complete="off"></el-input>
          </el-col>
        </el-col>
        <el-col class="w320 mt10 mr16">
          <el-col :span="6" class="search-name">当前环境</el-col>
          <el-col :span="18">
            <el-select v-model="form.stage" placeholder="请选择环境" filterable @change="getList">
              <el-option v-for="(item, index) in stageList" :key="index" :label="item.name" :value="item.name">
              </el-option>
            </el-select>
          </el-col>
        </el-col>

      </el-row> -->
      <template>
        <el-table stripe="true" :data="projectReleaseListData" class="mt16">
          <el-table-column sortable prop="version_no" min-width="10%" :label="$t('bm.deployCenter.pipelineName')">
            <template slot-scope="scope">
              <el-button type="text" @click="gotoDetail(scope.row.id)">{{scope.row.version_no}}</el-button>
            </template>
          </el-table-column>
          <el-table-column :show-overflow-tooltip="true" prop="name" min-width="12%" :label="$t('bm.deployCenter.pipelineDesc')" />
          <el-table-column prop="creator" min-width="10%" :label="$t('bm.deployCenter.creator')" />
          <el-table-column sortable :show-overflow-tooltip="true" prop="step" min-width="20%" label="当前步骤">
            <template slot-scope="scope">
              <ul class="ulSteps clearfix">
                <li v-if="scope.row.previous" class="graySteps">
                  <div class="grayBorder"></div>
                  <div class="grayCircle"><span class="deploy-done"></span></div>
                  {{scope.row.previous}}
                </li>
                <li v-if="scope.row.step" class="nowSteps">
                  <div class="grayBorder"></div>
                  <div class="grayCircle"><span class="deploy-run"></span></div>
                  {{scope.row.step}}
                </li>
                <li v-if="scope.row.next_step" class="nextSteps">
                  <div class="grayBorder"></div>
                  <div class="grayCircle"></div>
                  {{scope.row.next_step}}
                </li>
              </ul>
            </template>
          </el-table-column>
          <el-table-column sortable prop="status" min-width="9%" :label="$t('bm.deployCenter.status')">
            <template slot-scope="scope">
              <span v-if="scope.row.status === 1" class="circleSuccess"></span>
              <span v-if="scope.row.status === 2" class="circleRunning"></span>
              {{scope.row.statusName}}
            </template>
          </el-table-column>
          <el-table-column sortable :show-overflow-tooltip=true prop="update_at" min-width="10%"
          :label="$t('bm.add.updateTime')" />
          <el-table-column :label="$t('bm.deployCenter.operation')" min-width="10%">
            <template slot-scope="scope">
              <el-button type="text" @click="$refs.mark.show(scope.row)" v-if="scope.row.operations['manual']">人工卡点
              </el-button>
              <el-button type="text" @click="$refs.goToBuild.doShows(scope.row.id, scope.row.stage_id, scope.row.step)"
                v-if="scope.row.operations['build']">构建</el-button>
              <el-button type="text" @click="$refs.deploy.doShows(scope.row)" v-if="scope.row.operations['deploy']">部署
              </el-button>
              <el-button type="text" @click="termination(scope.row.id, scope.row.stage_id, scope.row.step_type)"
                v-if="scope.row.operations['terminate']">终止</el-button>
              <el-button type="text" @click="$refs.nextstage.show(scope.row.id, scope.row.stage_id)"
              v-if="scope.row.operations['next-stage']">{{$t('bm.add.enterNextStage')}}</el-button>
              <el-button type="text" @click="$refs.backto.show(scope.row.id, scope.row.stage_id)"
                v-if="scope.row.operations['back-to']">{{$t('bm.add.rollBack')}}</el-button>
            </template>
          </el-table-column>
        </el-table>
      </template>
      <page-nav ref="pages" :list="projectReleaseListData" v-on:getlist="getList" />
      <project-mark ref="mark" v-on:getprojectReleaseList="getList"></project-mark>

      <back-to ref="backto" v-on:getprojectReleaseList="getList"></back-to>
      <project-deploy ref="deploy" v-on:getprojectReleaseList="getList"></project-deploy>
      <project-release-create v-on:getprojectReleaseList="getList" :cpData="$props.listData"
        :listData="$props.listData" ref='pubDialog'></project-release-create>
      <common-delete ref="commonDelete" v-on:getlist="getList"></common-delete>
      <!-- <next-stage v-on:getprojectReleaseList="getList" ref="nextstage"></next-stage> -->
      <to-build v-on:getprojectReleaseList="getList" ref="goToBuild"></to-build>
    </div>
</template>
<script>
  import { Message, MessageBox } from 'element-ui';
  import { mapGetters } from 'vuex';
  import backend from '@/api/backend';
  import PageNav from '@/components/utils/PageList';
  import ListSearch from '@/components/utils/ListSearch';
  import Refresh from '@/components/utils/Refresh';
  import listTemplate from '@/common/listTemplate';
  import utils from '@/common/utils';
  import CommonDelete from '@/components/utils/Delete';

  import ProjectReleaseCreate from './dialogCI/ProjectReleaseCreate';// 创建流水线
  import NextStage from './dialogCI/ProjectNextStage'; // 下一个阶段
  import ProjectMark from './dialogCI/ProjectMark'; // 人工卡点
  import BackTo from './dialogCI/ProjectBackTo'; // 回退
  import ToBuild from './dialogCI/ProjectBuild'; //构建
  import ProjectDeploy from './dialogCI/ProjectDeploy'; // 部署


  export default {
    mixins: [listTemplate],
    // 发布类型
    props: ['pubType', 'listData', 'cpData', 'status', 'pubItem', 'hide'],
    data() {
      return {
        projectReleaseListData: [],
        hasStop: false,
        execFlag: '',
        del_service: false,
        hasError: false,
        // 当前正在发布的项目
        curPubRow: {},
        statePublish: [
          '失败', '成功', '进行中', '待执行', '结束', '已归档', '未知', '终止成功', '终止失败', '合并失败', '不支持'
        ],
        form: {
          'versionNo': '',
          'name': '',
          //'creator': '',
          //'step': '',
          'status': '',
          //'stage': ''
        },
        stageList: [],
        activeName: '',
      };
    },
    components: {
      PageNav,
      ProjectReleaseCreate,
      ListSearch,
      Refresh,
      CommonDelete,
      BackTo,
      NextStage,
      ToBuild,
      ProjectMark,
      ProjectDeploy,
    },
    activated() {
      this.getStageList();
    },
    computed: {
      ...mapGetters({
        loading: 'getLoading',
        projectID: 'projectID',
      }),
    },
    destroyed() {
    },
    methods: {
      tabClick(val) {
        console.info(val.name);
        this.activeName = val.name;
        this.getList();
      },
      getStageList() {
        backend.getProjectEnvsAll(this.projectID, (data) => {
          if(data && data.length > 0){
            this.stageList = data;
            this.activeName = `${data[0].id}`;
            this.getList('clear');
          }
        });
      },
      // 终止发布
      stopPub() {
        MessageBox.confirm(this.$t('bm.add.sureStopCodeModule'), this.$t('bm.infrast.tips'), { type: 'warning' }).then(() => {
          backend.execStopPub(this.curPubRow.id, this.$route.params.projectId, (data) => {
            this.hasStop = true;
            this.hasError = true;
            this.getList(true);
          });
        }).catch(() => { });
      },
      goPubDetail(type, id) {
        window.open(
          `//${window.location.host}/project/projectPubDetail/${this.$route.params.projectId}/${
          this.$props.pubItem.id
          }/${id}/${type}`
        );
      },
      // 终止发布
      termination(id, stargid, name) {
        MessageBox.confirm(this.$t('bm.add.sureStopPublishCodeModule'), this.$t('bm.add.hint'), { type: 'warning' }).then(() => {
          const params = {
            action_name: 'terminate'
          };
          const that = this;
          if (name === 'build') {
            backend.setBuildMerge(this.projectID, id, stargid, params, (data) => {
              Message.success(this.$t('bm.add.optionSuc'));
              that.getList(true);
            });
          } else if (name === 'deploy') {
            backend.setDeploy(this.projectID, id, stargid, params, (data) => {
              Message.success(this.$t('bm.add.optionSuc'));
              that.getList(true);
            });
          }
        }).catch(() => {});
      },
      // 查询项目发布单列表
      getList(isRefresh) {
        if (isRefresh) {
          this.$refs.pages.currentPage = 1;
        }
        if (isRefresh === 'clear') {
          this.clearSearch();
        }
        this.curPubRow = {};
        const params = Object.assign({
          page_size: this.$refs.pages.pageSize,
          page_index: this.$refs.pages.currentPage
        },this.form,{'stage': parseInt(this.activeName)});
        params.status = parseInt(params.status);
        backend.getProjectCI(this.projectID, JSON.stringify(params), (data) => {
          this.projectReleaseListData = data.item;
          if (this.projectReleaseListData) {
            this.projectReleaseListData.map((i) => {
              if (i.status >= 0) i.statusName = this.statePublish[i.status];
              i.update_at = utils.format(new Date(i.update_at), 'yyyy-MM-dd hh:mm');
            });
          }
          if(this.$refs.pages) {
            this.$refs.pages.total = data.total;
          };
        });
      },
      clearSearch() {
        this.form = {
          'versionNo': '',
          'name': '',
          'status': '',
        };
      },
      // app详情页面
      gotoDetail(id) {
        this.$router.push({
          name: 'projectCIDetail',
          params: {
            projectId: this.projectID,
            versionId: id
          }
        });
      },
    },
  };
</script>