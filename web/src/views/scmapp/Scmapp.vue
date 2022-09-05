<template>
<div class="content-container">
    <div class="page-content buManage">
      <notice-bar v-if="releaseSetingMsg" :msg="releaseSetingMsg" />
      <div class="portlet-body mt0">
        <el-row style="display: flex">
          <el-col class="w320 mt10 mr16">
            <el-input v-model="searchName" @keyup.enter.native="getList" auto-complete="off" maxlength="20" placeholder="应用名称/创建者"></el-input>
          </el-col>
          <el-col class="w320 mt10 mr16">
            <el-col :span="6" class="search-name">创建时间</el-col>
            <el-col :span="18">
              <el-date-picker style="width: 100%;" v-model="searchTime" :clearable="false" type="daterange" @change="getList" 
                range-separator="~" value-format="yyyy-MM-dd" :picker-options="pickerOptions" start-placeholder="开始日期" end-placeholder="结束日期" />
            </el-col>
          </el-col>
          <el-col class="w200 mt10">
            <el-button type="primary" @click="getList">搜索</el-button>
            <el-button type="text" class="font-gray" @click="getList('clear')">重置</el-button>
          </el-col>
        </el-row>
        
        <div class="table-toolbar">
          <el-row class="mt16">
            <el-col :span="16" style='text-align:left;'>
              <el-button :plain="false" type="primary" @click="$refs.create.doCreate()"> +创建应用</el-button>
            </el-col>
          </el-row>
        </div>
        <template>
        <el-table :data="curList" class="mt16">
          <el-table-column prop="name" :label="$t('bm.deployCenter.repositoryName')" min-width="12%" :show-overflow-tooltip=true>
            <template slot-scope="scope">
                <el-button @click="appDetail(scope.row.id, 1)" type="text">{{scope.row.name}}</el-button>
            </template>
          </el-table-column>
          <el-table-column prop="tags" label="应用标签" sortable　min-width="12%" :show-overflow-tooltip=true />
          <el-table-column prop="full_name" :label="$t('bm.deployCenter.repositoryFullName')" min-width="10%" :show-overflow-tooltip=true />
          <el-table-column prop="language" :label="$t('bm.deployCenter.buildLang')" sortable min-width="8%" :show-overflow-tooltip=true />
          <el-table-column prop="build_path" :label="$t('bm.deployCenter.buildPath')" min-width="8%" :show-overflow-tooltip=true />
          <el-table-column prop="path" :label="$t('bm.infrast.repositoryAdr')" min-width="25%" :show-overflow-tooltip=true>
            <template slot-scope="scope">
              <a class="font-blue" :href='scope.row.path' target="_blank">{{scope.row.path}}</a>
            </template>
          </el-table-column>
          <el-table-column :label="$t('bm.deployCenter.operation')" min-width="15%">
            <template slot-scope="scope">
                <el-button type="text" size="small" @click="appDetail(scope.row.id)">详情</el-button>
            </template>
          </el-table-column>
        </el-table>
      </template>
        <page-nav ref="pages" :list='curList' v-on:getlist="getList"></page-nav>
        <app-add ref="create" v-on:getlist="getList"></app-add>
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
import backend from '@/api/backend';
import PageNav from '@/components/utils/PageList';
import ListSearch from '@/components/utils/ListSearch';
import CommonDelete from '@/components/utils/Delete';
import CommonClose from '@/components/utils/Close';
import Refresh from '@/components/utils/Refresh';
import utils from '@/common/utils';
import NoticeBar from '@/components/utils/NoticeBar';
import AppAdd from './components/AppAdd';

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
      pickerOptions: {
        disabledDate(time) {
          return time.getTime() > Date.now() - 8.64e6;
        },
      },
      searchTime: ['', ''],
      searchName: '',
      
    };
  },
  watch: {},
  components: {
    PageNav,
    ListSearch,
    Refresh,
    CommonDelete,
    CommonClose,
    NoticeBar,
    AppAdd,
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
    changeFilterTxt(val, type) {
      this.filter_key = type;
      this.filter_val = val;
      this.getList(true);
    },
    // 添加应用模块
    addApp(flag) {
      this.$router.push({
        name: 'addScmApp',
      });
    },
    // 应用详情
    appDetail(scmAppId) {
      this.$router.push({
        name: 'scmAppDetail',
        params: {
          appId: scmAppId,
        },
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
      backend.getScmApps(params, (data) => {
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
