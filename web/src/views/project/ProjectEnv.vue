<template>
    <div class="portlet-body">
      <template>
        <div class="table-toolbar">
          <el-row>
            <el-col :span="10">
              <refresh v-on:getlist="getList"></refresh>
              <el-button :plain="true" type="primary" @click="$refs.create.doCreate(false)">
                <i class='icon-plus' /> 新建</el-button>
            </el-col>
            <el-col :span="6">
              &nbsp;
            </el-col>
            <el-col :span="8">
              <list-search ref="userSh" :searchList="searchList" v-on:changeFilterTxt="changeFilterTxt"></list-search>
            </el-col>
          </el-row>
        </div>
        <template>
          <el-table :data="dataList">
            <span slot="empty">
              {{loading?$t('bm.add.dataLoading'):noDataTxt}}
            </span>
            <el-table-column prop="name" :label="$t('bm.operatingCenter.envName')" sortable min-width="10%" :show-overflow-tooltip=true />
            <el-table-column prop="clusterName" label="集群" sortable min-width="10%" :show-overflow-tooltip=true />
            <el-table-column prop="namespace" label="命名空间" sortable min-width="10%" :show-overflow-tooltip=true />
            <el-table-column prop="description" :label="$t('bm.serviceM.description')" min-width="20%" :show-overflow-tooltip=true />
            <el-table-column prop="creator" :label="$t('bm.serviceM.creator')" sortable min-width="10%" :show-overflow-tooltip=true />
            <el-table-column prop="create_at" :label="$t('bm.serviceM.creationTime')" sortable min-width="15%" :show-overflow-tooltip=true />
            <el-table-column :label="$t('bm.deployCenter.operation')" min-width="15%">
              <template slot-scope="scope">
                <el-button type="text" size="small"
                    @click="$refs.create.doCreate(true, scope.row)">{{$t('bm.infrast.edit')}}
                </el-button>
                <el-button @click="$refs.commonDelete.doDelete('delProjectEnv', scope.row.project_id, scope.row.id)"
                  type="text" size="small">{{$t('bm.other.delete')}}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </template>
        <page-nav ref="pages"
                  :list="curList"
                  v-on:getlist="getList"></page-nav>
      </template>
      <common-delete ref="commonDelete"
                  v-on:getlist="getList"></common-delete>
      <env-create ref="create"
                  v-on:getlist="getList"></env-create>
    </div>
</template>
  <script>
import { mapGetters } from 'vuex';
import backend from '@/api/backend';
import ListSearch from '@/components/utils/ListSearch';
import PageNav from '@/components/utils/PageList';
import CommonDelete from '@/components/utils/Delete';
import Refresh from '@/components/utils/Refresh';
import listTemplate from '@/common/listTemplate';
import Utils from '@/common/utils';

import EnvCreate from './components/EnvCreate';


export default {
  mixins: [listTemplate],
  data() {
    return {
      curList: [],
      searchList: [
        { key: 'name', txt: this.$t('bm.operatingCenter.envName') },
        { key: 'creator', txt: this.$t('bm.serviceM.creator') },
      ],
      filterTxt: '',
      IntegrateSettings: [],
      param: {},
      searchVal: '',
      searchType: '',
    };
  },
  components: {
    ListSearch,
    PageNav,
    Refresh,
    EnvCreate,
    CommonDelete,
  },
  mounted() {
    // this.getList();
  },
  activated() {
    this.getAllIntegrateSettings();
    this.getList(true);
  },
  computed: {
    ...mapGetters({
      loading: 'getLoading',
      projectID: 'projectID',
    }),
    dataList() {
      // 强制替换dataList替代listtemplate中的方法
      return this.curList;
    },
  },
  methods: {
    goEdit(user) {
      this.$refs.create.doCreate(true, user);
    },
    getAllIntegrateSettings() {
      backend.getAllIntegrateSettings((data) => {
        if (data) {
          this.IntegrateSettings = data
        }
      });
    },
    getList(isRefresh) {
      if (isRefresh) {
        this.$refs.pages.currentPage = 1;
      }
      if (isRefresh === 'clear') {
        this.$refs.userSh.searchSelectChange();
        this.$refs.pages.currentPage = 1;
      }
      this.curList = [];
      backend.getProjectEnvs(
        this.projectID,
        JSON.stringify({
          page_size: this.$refs.pages.pageSize,
          page_index: this.$refs.pages.currentPage,
          filter_key: this.searchType,
          filter_val: this.searchVal,
        }),
        data => {
          this.$refs.pages.total = data.total;
          this.curList = data.item;
          this.curList = data.item.map(i => {
            i.create_at = Utils.format(new Date(i.create_at), 'yyyy-MM-dd hh:mm');
            // For project env list display
            for (const element of this.IntegrateSettings) {
              if (i.cluster === element.id) {
                i.clusterName = element.name
                break
              }
            }
            return i;
          });
        }
      );
    },
    changeFilterTxt(val, type) {
      this.searchVal = val;
      this.searchType = type;
      this.getList(false);
    },
  },
};
</script>
