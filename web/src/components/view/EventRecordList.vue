<style>
.normalcolor {
  color: #49ca01;
}
.warningcolor {
  color: #f1ad2d;
}
</style>
<template>
  <div>
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
      <el-table border :data="dataList" id='ingressTb' element-loading-target="ingressTb" :element-loading-text="$t('bm.add.loading')">
        <span slot="empty">
          {{loading?$t('bm.add.dataLoading'):noDataTxt}}
        </span>
        <el-table-column
          class="recordtdcolor"
          :label="$t('bm.add.eventRank')"
          sortable
          min-width="12%"
          :sort-method="sortRecord"
          :show-overflow-tooltip=true>
          <template slot-scope="scope">
            <span :class="{normalcolor:scope.row.event_level === 'Normal',warningcolor:scope.row.event_level === 'Warning'}">
              <!-- {{scope.row.event_level==='Normal'?'标准':$t('bm.add.warning')}} -->
              {{scope.row.event_level}}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="event_object" :label="$t('bm.add.eventObject')" min-width="12%" :show-overflow-tooltip=true />
        <el-table-column prop="event_type" :label="$t('bm.add.eventType')" min-width="12%" :show-overflow-tooltip=true />
        <el-table-column prop="event_message" :label="$t('bm.add.eventMessage')" min-width="44%" :show-overflow-tooltip=true />
        <el-table-column prop="event_time" :label="$t('bm.add.eventTime')" sortable min-width="20%" :show-overflow-tooltip=true />
      </el-table>
    </template>
    <page-nav ref="page" :list=filteredList></page-nav>
  </div>
</template>
<script>
import { mapGetters } from 'vuex';
import backend from '@/api/backend';
import PageNav from '@/components/utils/Page';
import ListSearch from '@/components/utils/ListSearch';
import Refresh from '@/components/utils/Refresh';
import listTemplate from '@/common/listTemplate';

export default {
  mixins: [listTemplate],
  props: ['appName', 'nodeId', 'appNamespace', 'clusterName'],
  data() {
    return {
      curList: [],
      searchList: [
        { key: 'event_level', txt: this.$t('bm.add.eventRank') },
        { key: 'event_object', txt: this.$t('bm.add.eventObject') },
        { key: 'event_type', txt: this.$t('bm.add.eventType') },
        { key: 'event_message', txt: this.$t('bm.add.eventMessage') },
        { key: 'event_time', txt: this.$t('bm.add.eventTime') },
      ],
      noDataTxt: this.$t('bm.add.noResult'),
      filterTxt: '',
    };
  },
  components: {
    PageNav,
    ListSearch,
    Refresh,
  },
  computed: {
    ...mapGetters({
      loading: 'getLoading',
    }),
  },
  methods: {
    getList(isRefresh) {
      if (isRefresh && this.$refs.length) {
        this.$refs.page.currentPage = 1;
      }
      if (isRefresh === 'clear') {
        this.$refs.clear.searchSelectChange();
        // this.currentPage = 1;
      }
      if (this.$props.clusterName === '') return;
      const postData = {
        clusterName: this.$props.clusterName,
        namespace: this.$props.appNamespace,
        appName: this.$props.appName,
      };
      if (this.$props.nodeId) {
        backend.getHostEventRecord(postData.clusterName, postData.appName, (data) => {
          this.curList = data;
        });
      } else {
        backend.getAppEventRecord(
          postData.clusterName,
          postData.namespace,
          postData.appName,
          (data) => {
            this.curList = data;
          }
        );
      }
    },
    // 排序“事件级别”
    sortRecord(a, b) {
      return a.event_level > b.event_level;
    },
  },
};
</script>
