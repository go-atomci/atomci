<template>
  <div class="page-content">
    <div class="portlet-body">
      <template>
        <div class="table-toolbar">
          <el-row>
            <el-col :span="10">
              <refresh v-on:getlist="getList"></refresh>
              <el-button :plain="true"
                type="primary"
                @click="$refs.create.doCreate(false)">
                <i class='icon-plus' /> 新建</el-button>
            </el-col>
            <el-col :span="6">
              &nbsp;
            </el-col>
            <el-col :span="8">
              <list-search ref="userSh"
                :searchList="searchList"
                v-on:changeFilterTxt="changeFilterTxt">
              </list-search>
            </el-col>
          </el-row>
        </div>
        <template>
          <el-table :data="curList">
            <span slot="empty">
              {{loading?$t('bm.add.dataLoading'):noDataTxt}}
            </span>
            <el-table-column prop="name" :label="$t('bm.operatingCenter.stageName')" sortable  min-width="15%" :show-overflow-tooltip=true />
            <el-table-column prop="description" :label="$t('bm.serviceM.description')" min-width="20%" :show-overflow-tooltip=true />
            <el-table-column prop="creator" :label="$t('bm.serviceM.creator')" min-width="15%" :show-overflow-tooltip=true />
            <el-table-column prop="create_at" :label="$t('bm.serviceM.creationTime')" sortable min-width="15%" :show-overflow-tooltip=true />
            <el-table-column :label="$t('bm.deployCenter.operation')" min-width="15%">
              <template slot-scope="scope">
                <el-button type="text" size="small" @click="goSetting(scope.row.id)">配置流程</el-button>
                <el-button type="text" size="small" 
                  @click="$refs.create.doCreate(true, scope.row)">{{$t('bm.infrast.edit')}}
                </el-button>
                <el-button 
                  @click="$refs.commonDelete.doDelete('delProjectPipe', scope.row.project_id, scope.row.id)" 
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
      <pipe-create ref="create"
                   v-on:getlist="getList"></pipe-create>
    </div>
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

import PipeCreate from './components/PipeCreate';

export default {
  mixins: [listTemplate],
  data() {
    return {
      curList: [],
      searchList: [{ key: 'name', txt: this.$t('bm.operatingCenter.stageName') }],
      filterTxt: '',
      detailInfo: [],
      param: {},
      searchVal: '',
      searchType: '',
    };
  },
  components: {
    ListSearch,
    PageNav,
    Refresh,
    PipeCreate,
    CommonDelete,
  },
  mounted() {
    // this.getList();
  },
  activated() {
    this.getList(true);
  },
  computed: {
    ...mapGetters({
      loading: 'getLoading',
      projectIDgetter: 'projectID',
    }),
    projectID() {
      if (this.projectIDgetter === 0 || this.projectIDgetter === undefined) {
        this.$store.dispatch('project/setProjectID', this.$route.params.projectID);
        return this.$route.params.projectID
      } else {
        return this.projectIDgetter
      }
    },
  },
  methods: {
    goEdit(id) {
      this.$refs.create.doCreate(true, id);
    },
    goSetting(pipeId) {
      this.$router.push({
        name: 'pipelinesAdd',
        params: {
          pipeId,
        },
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
      backend.getProjectPipeline(this.projectID,
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
            i.enabledStatus = i.enabled === false ? '禁用' : '启用';
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
