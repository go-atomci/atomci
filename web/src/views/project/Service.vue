<style>
@import '../../style/layout.css';
@import '../../style/baseDetail.css';
</style>
<style>
.apt-item{
  width:220px;
  height:160px;
  border:1px solid #ddd;
  text-align:center;
  display:inline-block;
  margin:30px 30px;
  cursor: pointer;
  background:#fff;
  transition: all .25s ease-out 0s;
}
.apt-item h4{
  padding-top:40px;
  padding-bottom:15px;
  font-size:17px;
}
.apt-item:hover{
  background: -webkit-gradient(linear,0 100%,100% 0,from(#2a89ff),to(#09bbff));
  box-shadow: 3px 12px 12px rgba(48,48,77,.1);
  color:#fff;
}
</style>
<template>
  <div class="page-content">
    <div class="portlet-body">
      <cluster-tab v-on:changeCluster="changeCluster" ref="clusterTabs" :fwC="99" :clusterName="cluster">></cluster-tab>
      <div class="table-toolbar">
        <el-row>
          <el-col :span="16">
            <refresh v-on:getlist="getList('clear')"></refresh>
          </el-col>
          <el-col :span="8">
            <list-search ref="clear" :searchList="searchList" v-on:changeFilterTxt="changeFilterTxt"></list-search>
          </el-col>
        </el-row>
      </div>
      <template>
        <el-table border :data="dataList" :element-loading-text="$t('bm.add.loading')" class="mt16">
          <span slot="empty">
            {{loadings?$t('bm.add.dataLoading'):noDataTxt}}
          </span>
          <el-table-column prop="name" :label="$t('bm.serviceM.appName')" min-width="20%" :show-overflow-tooltip=true>
            <template slot-scope="scope">
              <span v-if="scope.row.status === 'NotReady'">
                <router-link :to="{name: 'projectServiceDetail', params: {clusterName: scope.row.cluster, namespace: scope.row.namespace, appName: scope.row.name } }">
                  <el-button type="text"><i class="el-icon-warning" style="color: red"></i> {{scope.row.name}}</el-button>
                </router-link>
              </span>
              <span v-else>
                <router-link :to="{name: 'projectServiceDetail', params: {clusterName: scope.row.cluster, namespace: scope.row.namespace, appName: scope.row.name } }">
                  <el-button type="text"><i class="el-icon-success" style="color: green"></i> {{scope.row.name}}</el-button>
                </router-link>
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="namespace" :label="$t('bm.serviceM.resourceSpace')" min-width="10%" :show-overflow-tooltip=true />
          <el-table-column prop="kind" :label="$t('bm.deployCenter.type')" min-width="10%" :show-overflow-tooltip=true />
          <el-table-column prop="pods" :label="$t('bm.serviceM.podInstance')" min-width="10%" :show-overflow-tooltip=true />
          <el-table-column prop="image" :label="$t('bm.serviceM.mirror')" min-width="20%" :show-overflow-tooltip=true />
          <!-- <el-table-column prop="creator" :label="$t('bm.serviceM.creator')" min-width="10%" :show-overflow-tooltip=true /> -->
          <el-table-column prop="update_at" :label="$t('bm.serviceM.updateTime')" sortable min-width="16%" :show-overflow-tooltip=true />
          <el-table-column :label="$t('bm.deployCenter.operation')" min-width="10%">
            <template slot-scope="scope">
              <el-button type="text" size="small"  @click="goDetail(scope.row.cluster,scope.row.namespace,scope.row.name)">{{$t('bm.other.view')}}</el-button>
              <el-button size="small" @click="$refs.commonDelete.doDelete('removeService',scope.row.cluster,scope.row.namespace,scope.row.name)" type="text">{{$t('bm.other.delete')}}</el-button>
            </template>
          </el-table-column>
        </el-table>
      </template>
      <page-nav ref="pages" :list="curList" v-on:getlist="getList"></page-nav>
      <common-delete ref="commonDelete" v-on:getlist="getList"></common-delete>
    </div>
  </div>
</template>
<script>
import { mapGetters } from 'vuex';
import backend from '@/api/backend';
import PageNav from '@/components/utils/PageList';
import ListSearch from '@/components/utils/ListSearch';
import CommonDelete from '@/components/utils/Delete';
import Refresh from '@/components/utils/Refresh';
import ClusterTab from '@/components/utils/ClusterTab';
import listTemplate from '@/common/listTemplate';

export default {
  mixins: [listTemplate],
  data() {
    return {
      cluster: this.$route.params.clusterName || '',
      curList: [],
      searchList: [
        { key: 'name', txt: this.$t('bm.serviceM.appName') },
        { key: 'namespace', txt: this.$t('bm.serviceM.resourceSpace') },
        { key: 'image', txt: this.$t('bm.serviceM.mirror') },
        { key: 'create_at', txt: this.$t('bm.serviceM.creationTime') },
        { key: 'kind', txt: this.$t('bm.deployCenter.type') },
      ],
      filterTxt: '',
      filter_key: '',
      filter_val: '',
    };
  },
  components: {
    ListSearch,
    PageNav,
    Refresh,
    CommonDelete,
    ClusterTab,
  },
  computed: {
    ...mapGetters({
      loadings: 'getLoading',
      projectID: 'projectID',
    }),
    dataList() {
      // 强制替换dataList替代listtemplate中的方法
      return this.curList;
    },
  },
  mounted() {
    this.getList();
  },
  watch: {
    // 如果是从应用编排提交过来则必须重新刷新列表数据
    $route(to, from) {
      if (this.$route.query.isRefresh && from.name === 'projectAppDetail') {
        this.getList(true);
      }
    },
  },
  methods: {
    goDeploy() {
      this.$router.push({
        name: 'deployMirror',
        params: {
          cluster: this.cluster,
        },
      });
    },
    changeCluster(cluster) {
      if (this.cluster !== cluster) {
        this.cluster = cluster;
        this.getList();
      }
    },
    changeFilterTxt(val, type) {
      this.filter_key = type;
      this.filter_val = val;
      this.getList(true);
    },
    goDetail(cluster, namespace, appName) {
      this.$router.push({
        name: 'projectServiceDetail',
        params: {
          clusterName: cluster,
          namespace,
          appName,
        },
      });
    },
    getList(isRefresh) {
      if (isRefresh) {
        this.$refs.pages.currentPage = 1;
      }
      if (isRefresh === 'clear') {
        this.$refs.clear.searchSelectChange();
        this.$refs.pages.currentPage = 1;
      }
      const params = {
        page_size: this.$refs.pages.pageSize || 10,
        page_index: this.$refs.pages.currentPage || 1,
        filter_key: this.filter_key,
        filter_val: this.filter_val,
      };
      if (this.cluster === '') return;
      backend.getProjectServiceList(this.projectID, this.cluster, params, (data) => {
        this.curList = data.item;
        this.$refs.pages.total = data.total;
      });
    },
  },
};
</script>
