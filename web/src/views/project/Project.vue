<template>
<div class="content-container">
    <div class="page-content buManage">
      <notice-bar v-if="releaseSetingMsg" :msg="releaseSetingMsg" />
      <div class="portlet-body mt0">
        <el-row style="display: flex">
          <el-col class="w300 mt10 mr6">
            <el-col :span="4" class="search-name">状态</el-col>
            <el-col :span="20">
              <el-select v-model="searchStatus" placeholder="请选择状态" filterable @change="getList">
                <el-option v-for="(item, index) in statusList" :key="index" :label="item.name" :value="item.id">
                </el-option>
              </el-select>
            </el-col>
          </el-col>
          <el-col class="w320 mt10 mr16">
            <el-col :span="6" class="search-name">创建时间</el-col>
            <el-col :span="18">
              <el-date-picker style="width: 100%;" v-model="searchTime" :clearable="false" type="daterange" @change="getList" 
                range-separator="~" value-format="yyyy-MM-dd" :picker-options="pickerOptions" start-placeholder="开始日期" end-placeholder="结束日期" />
            </el-col>
          </el-col>
          <el-col class="w320 mt10 mr16">
            <el-input v-model="searchName" @keyup.enter.native="getList" auto-complete="off" maxlength="20" placeholder="项目名称/创建者"></el-input>
          </el-col>
          <el-col class="w200 mt10">
            <el-button type="primary" @click="getList">搜索</el-button>
            <el-button type="text" class="font-gray" @click="getList('clear')">重置</el-button>
          </el-col>
        </el-row>
        
        <div class="table-toolbar">
          <el-row class="mt16">
            <el-col :span="16" style='text-align:left;'>
              <el-button type="primary" @click="$refs.create.doCreate(false)">
                <i class='icon-plus' /> +创建项目</el-button>
            </el-col>
          </el-row>
        </div>
        <template v-if="curList.length>0">
          <template v-for="(item, index) in curList">
            <div class="project-panel">
              <el-row>
                <el-col :span="6">
                  <div class="proj-icon-blue">
                    <div :class="item.status == '1' ? 'proj-icon-bg': 'proj-icon'"><i :class="item.status == '1' ? 'el-icon-folder-opened': 'el-icon-folder'"></i></div>
                    <div class="proj-title font-blue" v-on:click="gotoProjectDetail(item.id)" >
                        {{item.name}}
                    </div>
                    <div class="proj-text mt3">项目成员数：{{item.members || 0}}人</div>
                  </div>
                </el-col>
                <el-col :span="3">
                  <div class="proj-text">状态</div>
                  <div class="proj-text proj-text-gray mt3">
                    <span :class="item.status == '1'? 'proj-circle-blue': 'proj-circle'"></span>
                    {{item.status == 1?$t('bm.dashboard.running'):'已结束'}}</div>
                </el-col>
                <el-col :span="6">
                  <div class="proj-text">描述</div>
                  <div class="proj-text proj-text-gray mt3">
                    <el-tooltip class="item" v-if="item.description" effect="dark" :content="item.description" placement="top-start">
                      <div class="text-wrap">{{item.description}}</div>
                    </el-tooltip>
                  </div>
                </el-col>
                <el-col :span="3">
                  <div class="proj-text">创建者</div>
                  <div class="proj-text proj-text-gray mt3">{{item.owner}}</div>
                </el-col>
                <el-col :span="4">
                  <div class="proj-text">创建时间</div>
                  <div class="proj-text proj-text-gray mt3">{{item.create_at}}</div>
                </el-col>
                <el-col :span="2">
                  <div class="proj-text">操作</div>
                  <div class="proj-text proj-text-gray mt3">
                    <el-button v-if="item.status!=2" type="text" size="small" @click="doEdit(item.id)">{{$t('bm.infrast.edit')}}</el-button>
                    <el-button v-if="item.status!=2" @click="gotoclose(item.id)" type="text" size="small">结束</el-button>
                    <el-button v-if="item.status==2" @click="gotoRestart(item.id)" type="text" size="small">重启</el-button>
                    <el-button v-if="item.status==2" @click="gotodelete(item.id)" type="text" size="small">{{$t('bm.other.delete')}}</el-button>
                  </div>
                </el-col>
              </el-row>
            </div>
          </template>
        </template>
        <template v-else>
          <div class="project-panel proj-null">暂无数据</div>
        </template>
        <page-nav ref="pages" :list='curList' v-on:getlist="getList"></page-nav>
        <common-delete ref="commonDelete" v-on:getlist="getList"></common-delete>
        <common-close ref="commonClose" v-on:getlist="getList"></common-close>
        <project-create ref="create" v-on:getlist="getList"></project-create>
      </div>
    </div>
</div>
</template>
<style lang="scss"  scoped>
  .content-container {
    left: 0; 
    top: 30px;
    width: 100%;
    height: 100%;
    height: calc(100% - 80px);
  }
  .portlet-body {
    box-shadow: none;
    margin-top: 0;
  }
  .font-blue, .font-blue:active, .font-blue:visited {
    color: #333;
  }
  .font-blue:hover {
    color: #40A0FF;
  }
  .project-panel {
    margin: 15px 0;
    padding: 15px;
    background:rgba(255,255,255,1);
    box-shadow:0px 2px 4px 0px rgba(64,158,255,0.12);
    border-radius:4px;
    border:1px solid rgba(235,238,245,1);
  }
  .project-panel:hover {
    background-color: #F2F6FC;
  }
  .project-panel:hover .font-blue {
    color: #40A0FF;
  }
  .proj-icon-blue {
    padding-left: 50px;
    position: relative;
  }
  .proj-icon, .proj-icon-bg {
    position: absolute;
    left: 0;
    top: 0;
    width: 44px;
    height: 44px;
    line-height: 44px;
    border-radius: 44px;
    background:rgba(192,196,204,1);
    text-align: center;
  }
  .proj-icon-bg {
    background:rgba(64,158,255,1);
  }
  .proj-icon i, .proj-icon-bg i {
    color: #fff;
    font-size: 20px;
    vertical-align: middle;
  }
  .proj-title {
    height: 22px;
    font-size: 16px;
    font-family: PingFangSC-Medium,PingFang SC;
    font-weight: 500;
    color: rgba(51,51,51,1);
    line-height: 22px;
  }
  .proj-text, .text-wrap {
    height: 22px;
    font-size: 14px;
    font-family: PingFangSC-Regular,PingFang SC;
    font-weight: 400;
    color: rgba(144,147,153,1);
    line-height: 22px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .proj-text-gray, .text-wrap {
    color: rgba(51,51,51,1);
  }
  .proj-text .el-button--small, .el-button--small.is-round {
    padding: 0;
  }
  .proj-circle, .proj-circle-blue {
    display: inline-block;
    width: 8px;
    height: 8px;
    border-radius: 8px;
    margin-right: 5px;
    background: rgba(192,196,204,1);
  }
  .proj-circle-blue {
    background: rgba(24,144,255,1);
  }
  .proj-null {
    line-height: 30px;
    color: #ccc;
    text-align: center;
  }
</style>
<script>
import { mapGetters } from 'vuex';
import { Message, MessageBox } from 'element-ui';
import backend from '@/api/backend';
import PageNav from '@/components/utils/PageList';
import ListSearch from '@/components/utils/ListSearch';
import CommonDelete from '@/components/utils/Delete';
import CommonClose from '@/components/utils/Close';
import Refresh from '@/components/utils/Refresh';
import utils from '@/common/utils';
import NoticeBar from '@/components/utils/NoticeBar';

import ProjectCreate from './dialog/ProjectCreate';

export default {
  data() {
    return {
      instance: '',
      example: '',
      curList: [],
      searchList: [
        { key: 'name', txt: this.$t('bm.deployCenter.proName') },
        { key: 'description', txt: this.$t('bm.serviceM.description') },
        { key: 'owner', txt: this.$t('bm.deployCenter.owner') },
        { key: 'create_at', txt: this.$t('bm.serviceM.creationTime') },
      ],
      filterTxt: '',
      releaseSetingMsg: '',
      statusList: [
        {'name': '全部', 'id': ''},
        {'name': '运行中', 'id': 1},
        {'name': '已结束', 'id': 2}
      ],
      pickerOptions: {
        disabledDate(time) {
          return time.getTime() > Date.now() - 8.64e6;
        },
      },
      searchStatus: 1,
      searchTime: ['', ''],
      searchName: '',
      
    };
  },
  watch: {
    // 如果是从项目详情删除过来的话则必须重新刷新列表数据
    $route(to, from) {
      if (this.$route.query.isRefresh && from.name === 'projectDetail') {
        this.getList();
      }
    },
  },
  components: {
    PageNav,
    ListSearch,
    Refresh,
    ProjectCreate,
    CommonDelete,
    CommonClose,
    NoticeBar,
  },
  mounted() {
    this.getList();
  },
  computed: {
    ...mapGetters({
      loading: 'getLoading',
    })
  },
  methods: {
    gotoProjectDetail(projectID) {
      this.$store.dispatch('project/setProjectID', projectID);
      this.$router.push({
        name: 'projectSummary',
        params: { projectID: projectID }
      })
    },
    changeFilterTxt(val, type) {
      this.filter_key = type;
      this.filter_val = val;
      this.getList(true);
    },
    gotodelete(id) {
      MessageBox.confirm('确定删除当前项目？', this.$t('bm.infrast.tips'), { type: 'warning' }).then(() => {
        backend.delNewProject(id, (data) => {
          Message.success('删除成功！');
          this.getList();
        });
      }).catch(() => { });

    },
    gotoclose(id) {
      MessageBox.confirm('确定结束当前项目？', this.$t('bm.infrast.tips'), { type: 'warning' }).then(() => {
        backend.updateNewProject(id,{ status: 2},(data) => {
          Message.success('关闭成功！');
          this.getList();
        });
      }).catch(() => { });
    },
    gotoRestart(id) {
      MessageBox.confirm('确定重启当前项目？', this.$t('bm.infrast.tips'), { type: 'warning' }).then(() => {
        backend.updateNewProject(id,{ status: 1},(data) => {
          Message.success('重启成功！');
          this.getList();
        });
      }).catch(() => { });
    },
    doEdit(id) {
      backend.getProjectDetail(id, (data) => {
        if(data) {
          this.$refs.create.doCreate(true, data);
        }
      });
    },
    getList(isRefresh) {
      if (isRefresh) {
        this.$refs.pages.currentPage = 1;
      }
      if (isRefresh === 'clear') {
        this.clearSearch();
      }
      const params = {
        page_index: this.$refs.pages.currentPage,
        page_size: this.$refs.pages.pageSize,
        name: this.searchName,
        createAtStart: this.searchTime[0],
        createAtEnd: this.searchTime[1]
      };
      if(this.searchStatus === 1 || this.searchStatus ===2) params.status = this.searchStatus;
      if (this.cluster === '') return;
      backend.getProject(params, (data) => {
        if (data) {
          this.$refs.pages.total = data.total;
          this.curList = data.item.map((i) => {
            i.create_at = utils.format(new Date(i.create_at), 'yyyy-MM-dd hh:mm');
            return i;
          });
        } else {
          this.$refs.pages.total = 0;
        }
      }
      );
    },
    clearSearch() {
      this.searchStatus = 1;
      this.searchTime = ['', ''];
      this.searchName = '';
    }
  },
};
</script>
