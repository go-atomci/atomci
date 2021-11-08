<template>
    <div class="page-content">
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
            <el-table stripe :data="dataList">
              <span slot="empty">
                {{loading?$t('bm.add.dataLoading'):noDataTxt}}
              </span>
              <el-table-column prop="name" :label="$t('bm.operatingCenter.stepName')" sortable min-width="15%" :show-overflow-tooltip=true />
              <el-table-column prop="type_display" :label="$t('bm.operatingCenter.stepType')" min-width="15%" :show-overflow-tooltip=true />
              <el-table-column prop="description" :label="$t('bm.serviceM.description')" min-width="15%" :show-overflow-tooltip=true />
              <el-table-column prop="creator" :label="$t('bm.serviceM.creator')" sortable min-width="10%" :show-overflow-tooltip=true />
              <el-table-column prop="create_at" :label="$t('bm.serviceM.creationTime')" sortable min-width="15%" :show-overflow-tooltip=true />
              <el-table-column :label="$t('bm.deployCenter.operation')" min-width="15%">
                <template slot-scope="scope" v-if="scope.row.name!=='default'">
                  <!-- <el-button type="text" size="small" @click="$refs.view.doView(scope.row.id)">查看引用流程</el-button> -->
                  <el-button type="text" size="small" @click="$refs.create.doCreate(true,scope.row)">{{$t('bm.infrast.edit')}}
                  </el-button>
                  <el-button @click="$refs.commonDelete.doDelete('delStep',scope.row.id)" type="text" size="small">{{$t('bm.other.delete')}}
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </template>
          <page-nav ref="pages" :list="curList" v-on:getlist="getList"></page-nav>
        </template>
        <common-delete ref="commonDelete" v-on:getlist="getList"></common-delete>
        <node-create ref="create" v-on:getlist="getList"></node-create>
        <node-view ref="view"></node-view>
      </div>
    </div>
  </template>
<script>
  import { mapGetters } from 'vuex';
  import backend from '@/api/backend';
  import ListSearch from '@/components/utils/ListSearch';
  import PageNav from '@/components/utils/PageList';
  import NodeCreate from './dialog/NodeCreate';
  import NodeView from './dialog/NodeView';
  import CommonDelete from '@/components/utils/Delete';
  import Refresh from '@/components/utils/Refresh';
  import listTemplate from '@/common/listTemplate';
  import Utils from '@/common/utils';

  export default {
    mixins: [listTemplate],
    data() {
      return {
        curList: [],
        searchList: [
          { key: 'name', txt: this.$t('bm.operatingCenter.stepName') },
          { key: 'type_display', txt: this.$t('bm.operatingCenter.stepType') },
        ],
        filterTxt: '',
        detailInfo: [],
        param: {},
        searchVal: '',
        searchType: '',
      };
    },
    components: {
      ListSearch,
      Refresh,
      NodeCreate,
      NodeView,
      CommonDelete,
      PageNav,
    },
    mounted() {
      this.getList();
    },
    activated() {
      this.getList(true);
    },
    computed: {
      ...mapGetters({
        loading: 'getLoading',
      }),
      dataList() {
        // 强制替换dataList替代listtemplate中的方法
        return this.curList;
      },
    },
    methods: {
      goEdit(user) {
        this.$refs.create.doCreate(true,user);
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
        backend.getStep(
          JSON.stringify({
            page_size: this.$refs.pages.pageSize,
            page_index: this.$refs.pages.currentPage,
            filter_key: this.searchType,
            filter_val: this.searchVal
        }),(data) => {
          this.$refs.pages.total = data.total;
          this.curList = data.item;
          this.curList = data.item.map((i) => {
            i.create_at = Utils.format(new Date(i.create_at), 'yyyy-MM-dd hh:mm');
            return i;
          });
        });
      },
      changeFilterTxt(val, type) {
        this.searchVal = val;
        this.searchType = type;
        this.getList(false);
      },
    },
  };
</script>
