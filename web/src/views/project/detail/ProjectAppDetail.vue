<template>
  <div class="page-content">
    <div class="portlet-body mt10 pb-25 min-height150">
      <el-row>
        <el-col :span="24" class="app-title clearfix">
          <div class="f-r">
            <el-button type="primary" @click="appDetail(true)">编辑</el-button>
            <el-button type="danger" @click="$refs.commonDelete.doDelete('delProjectApp',$route.params.projectId,$route.params.appId)">删除</el-button>
          </div>
          {{detailInfo.name}}
        </el-col>
      </el-row>
      <el-row class="mt-15">
        <el-col class="w-400">开发语言：{{detailInfo.language}}</el-col>
        <el-col class="w-400">默认分支: {{detailInfo.branch_name}}</el-col>
      </el-row>
      <el-row class="mt-15">
        <el-col class="w-400">编译环境: {{detailInfo.compile_env != '' ? detailInfo.compile_env:'未配置'}}</el-col>
        <el-col class="w-400">构建目录: {{detailInfo.build_path}}</el-col>
      </el-row>
      <el-row class="mt-15">
        <el-col class="w-400">创建人：{{detailInfo.creator}}</el-col>
        <el-col class="w-400">创建时间：{{detailInfo.create_at}}</el-col>
      </el-row>
    </div>
    <div class="portlet-body mt10 min-height150">
        <el-tabs v-model="activeName" @tab-click="handleClick">
          <el-tab-pane label="分支信息" name="first">
            <div class="table-toolbar" style="float: left">
              <el-button type="primary" :plain="true" @click="synBranch()">
                <i class="el-icon-refresh"></i> {{$t('bm.deployCenter.synchBran')}}
              </el-button>
            </div>
            <template>
              <el-table :data="listCol" class="mt16">
                <el-table-column prop="branch_name" :label="$t('bm.deployCenter.branchName')" min-width="12%" :show-overflow-tooltip=true />
                <el-table-column prop="path" :label="$t('bm.deployCenter.path')" min-width="12%" :show-overflow-tooltip=true>
                  <template slot-scope="scope">
                    <a :href='scope.row.path' target="_blank">{{scope.row.path}}</a>
                  </template>
                </el-table-column>
                <el-table-column :label="$t('bm.deployCenter.operation')" min-width="10%">
                  <template slot-scope="scope">
                    <el-button type="text" size="small" @click="checkBranch(scope.row)">
                      {{$t('bm.deployCenter.setupCurrentBranch')}}
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </template>
            <page-nav ref="pages" :list="listCol" v-on:getlist="getList"></page-nav>
          </el-tab-pane>
        </el-tabs>
      <common-delete ref="commonDelete" v-on:getlist="backTo"></common-delete>
      <project-app-edit ref="projectEdit" v-on:getList="appDetail"></project-app-edit>
    </div>
  </div>
</template>
<style scoped>
  .app-title {
    font-size: 20px;
    line-height: 40px;
  }
  .f-r {
    float: right;
  }
  .text-r {
    text-align: right;
  }
  .text-null {
    text-align: center;
    font-size: 14px;
    color: #999;
    padding: 20px;
  }
  .w-400 {
    width: 400px;
  }
  .mt-15 {
    margin-top: 15px;
  }
  .pb-25 {
    padding-bottom: 25px;
  }
  .env-ul li {
    float: left;
    width: 33.3%;
    padding: 15px;
  }
  .env-ul li>div {
    border: 1px solid #ccc;
    box-shadow: 0px 2px 4px 0px rgba(64,158,255,0.12);
    border: 1px solid rgba(235,238,245,1);
    padding: 25px;
    line-height: 48px;
    color: #333;
    font-size: 14px;
    font-weight: bold;
  }
  .bg-normal {
    width: 48px;
    height: 48px;
    line-height: 48px;
    text-align: center;
    border-radius: 48px;
    display: inline-block;
    vertical-align: middle;
    margin-right: 10px;
    font-size: 22px;
    color: #fff;
    background-color: #409EFF;
  }
  .bg-normal-0 {
    background-color: #409EFF;
  }
  .bg-normal-1 {
    background-color: #FEA202;
  }
  .bg-normal-2 {
    background-color: #67C23A;
  }
  .min-height150 {
	min-height: 150px;
  }
</style>
<script>
import { Message, MessageBox } from 'element-ui';
import { mapGetters } from 'vuex';
import backend from '@/api/backend';
import ListSearch from '@/components/utils/ListSearch';
import Refresh from '@/components/utils/Refresh';
import listTemplate from '@/common/listTemplate';
import CommonDelete from '@/components/utils/Delete';
import Utils from '@/common/utils';
import ProjectAppEdit from '../dialog/ProjectAppEdit';
import PageNav from '@/components/utils/PageList';

export default {
  mixins: [listTemplate],
  data() {
    return {
      listCol: [],
      dependList: [],
      filterTxt: '',
      activeName: '',
      envList: [],
      icon: ['', 'app-test', 'app-prod'],
      detailInfo: {},
    };
  },
  components: {
    ListSearch,
    PageNav,
    Refresh,
    CommonDelete,
    ProjectAppEdit,
  },
  computed: {
    ...mapGetters({
      projectID: 'projectID',
    })
  },
  created() {
    this.activeName = this.$route.params.tabs != 1 ? 'second' : 'first';
    this.appDetail();
  },
  mounted() {
  },
  activated() {
    this.getList();
  },
  methods: {
    getList() {
      const params = {
        page_size: this.$refs.pages.pageSize,
        page_index: this.$refs.pages.currentPage,
      };
      backend.getProjectBranch(this.$route.params.projectId, this.$route.params.appId, params, (data) => {
        this.listCol = data.item;
        this.$refs.pages.total = data.total;
      });
    },
    handleClick(val) {
      if(val.name === 'first') {
        this.getList();
      }
    },
    // 切换新分支
    checkBranch(item) {
      MessageBox.confirm('确定切换为当前分支吗？', this.$t('bm.infrast.tips'), { type: 'warning' })
        .then(() => {
          const params = {
            app_id: parseInt(this.$route.params.appId),
            branch_name: item.branch_name
          };
          backend.changeBranch(
            this.$route.params.projectId,
            this.$route.params.appId,
            params,
            () => {
              this.getList();
            }
          );
        })
        .catch(() => {});
    },
    // 同步远程分支
    synBranch() {
      MessageBox.confirm('确定同步远程分支吗？', this.$t('bm.infrast.tips'), { type: 'warning' })
        .then(() => {
          backend.asyncBranch(this.$route.params.projectId, this.$route.params.appId, (data) => {
            Message.success(this.$t('bm.add.optionSuc'));
            this.getList();
          });
        })
        .catch(() => {});
    },
    appDetail(flag) {
      backend.getAppDetail(this.$route.params.projectId, this.$route.params.appId, (data) => {
        if(flag) {
          let history = [];
          this.listCol.map((i) => {
            history.push(i.branch_name);
          });
          const cl = Object.assign({"branch_history_list": history},data);
          this.$refs.projectEdit.doCreate(true, cl);
        } else {
          data.create_at = Utils.format(new Date(data.create_at), 'yyyy-MM-dd hh:mm:ss');
          this.detailInfo = Object.assign({},data);
        }
      });
    },
    backTo() {
      this.$router.push({
        name: 'projectApp', params: {projectId: this.$route.params.projectId}
      });
    },
  }
}
</script>
