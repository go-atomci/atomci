<template>
    <div class="portlet-body">
      <template>
          <div class="table-toolbar">
            <el-row>
              <el-col :span="10">
                <refresh v-on:getlist="getList('clear')"></refresh>
              </el-col>
              <el-col :span="6">
                &nbsp;
              </el-col>
              <el-col :span="8">
                <list-search ref="clear" :searchList="searchList" v-on:changeFilterTxt="changeFilterTxt"></list-search>
              </el-col>
            </el-row>
          </div>
          <template>
            <el-table :data="dataList">
              <el-table-column prop="user" :label="$t('bm.operCenter.account')" sortable min-width="25%" :show-overflow-tooltip="true" />
              <el-table-column prop="operation" :label="$t('bm.deployCenter.operation')" min-width="25%" :show-overflow-tooltip="true" />
              <el-table-column prop="operation_object" :label="$t('bm.operCenter.actionObject')" min-width="33%" :show-overflow-tooltip="true" />
              <el-table-column prop="method" :label="$t('bm.operCenter.requestMethod')" sortable min-width="25%" :show-overflow-tooltip="true" />
              <el-table-column prop="operation_status" :label="$t('bm.operCenter.responseCode')" min-width="14%" :show-overflow-tooltip="true" />
              <el-table-column prop="create_at" :label="$t('bm.operCenter.operTime')" sortable min-width="14%" :show-overflow-tooltip="true" />
            </el-table>
          </template>
        <page-nav ref="page" :list="filteredList"></page-nav>
      </template>
      <common-delete ref="commonDelete" v-on:getlist="getList"></common-delete>
    </div>
</template>
<script>
import { mapGetters } from 'vuex';
import backend from '@/api/backend';
import PageNav from '@/components/utils/Page';
import ListSearch from '@/components/utils/ListSearch';
import CommonDelete from '@/components/utils/Delete';
import Refresh from '@/components/utils/Refresh';
import listTemplate from '@/common/listTemplate';
import UtilsFn from '@/common/utils';

export default {
  mixins: [listTemplate],
  data() {
    return {
      curList: [],
      searchList: [
        { key: 'user', txt: this.$t('bm.operCenter.account') },
        { key: 'operation', txt: this.$t('bm.deployCenter.operation') },
        { key: 'operation_object', txt: this.$t('bm.operCenter.actionObject') },
        { key: 'method', txt: this.$t('bm.operCenter.requestMethod') },
        { key: 'operation_status', txt: this.$t('bm.operCenter.responseCode') },
        { key: 'create_at', txt: this.$t('bm.operCenter.operTime') },
      ],
      filterTxt: '',
    };
  },
  components: {
    PageNav,
    ListSearch,
    Refresh,
    CommonDelete,
  },
  watch: {
    $route() {
      this.getList();
    },
  },
  created() {
    this.getList();
  },
  computed: {
    ...mapGetters({
      loading: 'getLoading',
    }),
  },
  methods: {
    getList(isRefresh) {
      if (isRefresh) {
        this.$refs.page.currentPage = 1;
      }
      if (isRefresh === 'clear') {
        this.$refs.clear.searchSelectChange();
        this.currentPage = 1;
      }
      backend.getAudit((data) => {
        this.curList = data.map((item) => {
          item.create_at = UtilsFn.format(new Date(item.create_at), 'yyyy-MM-dd hh:mm');
          return item;
        });
      });
    },
  },
};
</script>
