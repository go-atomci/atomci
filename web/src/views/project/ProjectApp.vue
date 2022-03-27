<template>
  <div class="page-content">
    <div class="portlet-body">
      <el-row>
          <el-col :span="2" class="search-name">应用名</el-col>
          <el-col :span="4">
            <el-input v-model="form.name" @keyup.enter.native="getList" placeholder="请输入应用名" filterable auto-complete="off"></el-input>
          </el-col>
        <!-- <el-col class="w320 mb10">
          <el-col :span="6" class="search-name">类型</el-col>
          <el-col :span="18">
            <el-select v-model="form.type" placeholder="请选择类型" filterable>
              <el-option v-for="(item, index) in typeList" :key="index" :label="item.description" :value="item.name">
              </el-option>
            </el-select>
          </el-col>
        </el-col> -->
          <el-col :span="2" class="search-name">开发语言</el-col>
          <el-col :span="4">
            <el-select v-model="form.language" placeholder="请选择开发语言" filterable @change="getList">
              <el-option v-for="(item, index) in languageList" :key="index" :label="item.description" :value="item.name">
              </el-option>
            </el-select>
          </el-col>
          <el-col :span="2" class="search-name">仓库地址</el-col>
          <el-col :span="4">
            <el-input v-model="form.address" placeholder="请输入仓库地址" @keyup.enter.native="getList" filterable auto-complete="off"></el-input>
          </el-col>
      </el-row>
      <el-row>
        <div class="mt10">
          <el-col :span="2" class="search-name">创建者</el-col>
          <el-col :span="4">
            <el-input v-model="form.creator" placeholder="请输入创建者" @keyup.enter.native="getList" filterable auto-complete="off"></el-input>
          </el-col>
        </div>
          <div class="mt10">
            <el-col :span="2" class="search-name">创建时间</el-col>
            <el-col :span="4" >
              <el-date-picker v-model="form.time" style="width:100%" :clearable="false" type="daterange" @change="getList"
              range-separator="~" value-format="yyyy-MM-dd" :picker-options="pickerOptions" start-placeholder="开始日期" end-placeholder="结束日期" />
            </el-col>
          </div>
          <div class="mt10">
            <el-col :span="4">
              <el-button type="primary" @click="getList">搜索</el-button>
              <el-button class="font-gray" type="text" @click="getList('clear')">重置</el-button>
            </el-col>
          </div>
      </el-row>
      <div class="table-toolbar">
        <el-row class="mt16">
          <el-col :span="16">
            <el-button :plain="false" type="primary" @click="$refs.appRegister.doCreate(true)">+{{$t('bm.deployCenter.addRepository')}}</el-button>
          </el-col>
        </el-row>
      </div>
      <template>
        <el-table border :data="listCol" class="mt16">
          <el-table-column prop="name" :label="$t('bm.deployCenter.repositoryName')" min-width="12%" :show-overflow-tooltip=true></el-table-column>
          <el-table-column prop="full_name" :label="$t('bm.deployCenter.repositoryFullName')" min-width="14%" :show-overflow-tooltip=true />
          <el-table-column prop="language" :label="$t('bm.deployCenter.buildLang')" sortable min-width="8%" :show-overflow-tooltip=true />
          <el-table-column prop="build_path" :label="$t('bm.deployCenter.buildPath')" min-width="6%" :show-overflow-tooltip=true />
          <el-table-column prop="path" :label="$t('bm.infrast.repositoryAdr')" min-width="25%" :show-overflow-tooltip=true>
            <template slot-scope="scope">
              <a class="font-blue" :href='scope.row.path' target="_blank">{{scope.row.path}}</a>
            </template>
          </el-table-column>
          <el-table-column :label="$t('bm.deployCenter.operation')" min-width="15%">
            <template slot-scope="scope">
              <el-button type="text" size="small" @click="$refs.appArrange.doSetup(scope.row)">
                {{$t('bm.deployCenter.appOrchestration')}}
              </el-button>
              <el-button type="text" size="small" @click="$refs.appRegister.doCreate(false, scope.row)">
                {{$t('bm.deployCenter.edit')}}
              </el-button> 
              <el-button type="text" size="small" @click="$refs.commonDelete.doDelete('delProjectApp',$route.params.projectID, scope.row.id)">
                {{$t('bm.deployCenter.delete')}}
              </el-button>           
              <el-button type="text" size="small" @click="appDetail(scope.row.scm_id)">
                {{$t('bm.deployCenter.details')}}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </template>
      <page-nav ref="pages" :list="listCol" v-on:getlist="getList"></page-nav>
      <common-delete ref="commonDelete" v-on:getlist="getList"></common-delete>
      <app-arrange ref="appArrange" :envList="envStageList" :appList="projectAppList"></app-arrange>
      <project-app-register ref="appRegister" :scmAppList="scmAppList" v-on:getlist="getList"></project-app-register>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
  import { MessageBox } from 'element-ui';
import backend from '@/api/backend';
import ListSearch from '@/components/utils/ListSearch';
import CommonDelete from '@/components/utils/Delete';
import PageNav from '@/components/utils/PageList';
import Refresh from '@/components/utils/Refresh';
import listTemplate from '@/common/listTemplate';
import AppArrange from './components/AppArrange';
import ProjectAppRegister from './components/ProjectAppRegister';

export default {
  mixins: [listTemplate],
  data() {
    return {
      listCol: [],
      envStageList: [],
      projectAppList: [],
      scmAppList: [],
      filterTxt: '',
      searchVal: '',
      searchType: '',
      pickerOptions: {
        disabledDate(time) {
          return time.getTime() > Date.now() - 8.64e6;
        },
      },
      form: {
        name: '',
        type: '',
        language: '',
        address: '',
        creator: '',
        time: ['', '']
      },
      typeList: [
        {'description': '应用','name': 'app'},
        {'description': '依赖','name': 'module'},
      ],
      languageList: [
        {'description': '全部','name': ''},
        {'description': 'Static','name': 'static'},
        {'description': 'Java','name': 'Java'},
        {'description': 'Node','name': 'Node'},
        {'description': 'Go','name': 'go'},
        {'description': 'C#','name': 'C#'},
      ],
    };
  },
  components: {
    ListSearch,
    PageNav,
    CommonDelete,
    Refresh,
    AppArrange,
    ProjectAppRegister,
  },
  computed: {
    ...mapGetters({
      loading: 'getLoading',
      projectID: 'projectID',
    }),
    dataList() {
      // 强制替换dataList替代listtemplate中的方法
      return this.listCol;
    },
  },
  mounted() {
    backend.getProjectEnvsAll(this.projectID, (data) => {
      this.envStageList = data
    }),
    backend.getAppAll(this.projectID, (data) => {
      this.projectAppList = data
    }),
    backend.getAllScmApps( (data) => {
      this.scmAppList = data
    }),
    this.getList(true);
  },
  methods: {
    getList(isRefresh) {
      if (isRefresh) {
        this.$refs.pages.currentPage = 1;
      }
      if (isRefresh === 'clear') {
        this.clearSearch();
      }
      this.listCol = [];
      const params = {
        page_size: this.$refs.pages.pageSize,
        page_index: this.$refs.pages.currentPage,
        name: this.form.name,
        creator: this.form.creator,
        language: this.form.language,
        type: this.form.type,
        path: this.form.address,
        createAtStart: this.form.time[0],
        createAtEnd: this.form.time[1]
      }
      backend.getApp(this.projectID, params, (data) => {
        this.$refs.pages.total = data.total;
        this.listCol = data.item
      });
    },
    handleSelectAll(val) {
      this.selectlist = val;
    },
    // 添加应用模块
    addApp(flag) {
      this.$router.push({
        name: 'addApp',
        params: {
          projectID: this.projectID
        }
      });
    },
    // 我的应用-应用详情
    appDetail(scmAppId) {
      MessageBox.confirm('确定要进入“我的应用 / 应用详情” 页面 ？', '提示', { type: 'warning' }).then(() => {
        this.$router.push({
          name: 'scmAppDetail',
          params: {
            appId: scmAppId,
          },
        });
      }).catch(() => { });
      
    },
    changeFilterTxt(val, type) {
      this.searchVal = val;
      this.searchType = type;
      this.getList(false);
    },
    clearSearch() {
      this.form = {
        name: '',
        type: '',
        language: '',
        address: '',
        creator: '',
        time: ['', '']
      };
    },
  }
}
</script>
