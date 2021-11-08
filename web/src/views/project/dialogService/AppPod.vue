<template>
  <div>
    <div class="table-toolbar">
      <el-row>
        <el-col :span="10">
          &nbsp;
        </el-col>
        <el-col :span="6">
          &nbsp;
        </el-col>
        <el-col :span="8">
          <list-search :searchList="searchList" v-on:changeFilterTxt="changeFilterTxt"></list-search>
        </el-col>
      </el-row>
    </div>
    <el-table border :data="dataList" :element-loading-text="$t('bm.add.loading')">
      <el-table-column prop="name" :label="$t('bm.deployCenter.name')" min-width="21%" :show-overflow-tooltip=true />
      <el-table-column prop="version" :label="$t('bm.add.version')" min-width="15%" :show-overflow-tooltip=true />
      <el-table-column prop="pod_ip" :label="$t('bm.add.insIP')" min-width="15%" :show-overflow-tooltip=true />
      <el-table-column prop="node_ip" :label="$t('bm.add.hostIP')" min-width="15%" :show-overflow-tooltip=true />
      <el-table-column :label="$t('bm.deployCenter.status')" min-width="10%">
        <template slot-scope="scope">
          <span v-if="scope.row.status === 'NotReady'" style="color: red">
            {{ scope.row.status }}
            <el-tooltip class="item" effect="light" placement="top">
              <div slot="content">
                <div style="color: red">{{scope.row.message}}</div>
              </div>
              <i style="margin-left:5px;color:gray;cursor:pointer" class="el-icon-info"></i>
            </el-tooltip>
          </span>
          <span v-else>{{ scope.row.status }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="restart_count" :label="$t('bm.add.resetNum')" min-width="12%" />
      <el-table-column prop="start_time" :label="$t('bm.add.startTime')" min-width="15%" :show-overflow-tooltip=true />
      <el-table-column :label="$t('bm.deployCenter.operation')" min-width="12%">
        <template slot-scope="scope">
          <el-button type="text" size="small"  @click="showContainerExec(scope.row.namespace, scope.row.name, scope.row.containers[0].name)" :disabled="scope.row.pod_ip==''">
            {{$t('bm.add.conConsole')}}
          </el-button>
          <!-- <el-button type="text" size="small" @click="$refs.view.doCreate($props.cluster, scope.row.namespace, $props.appName, scope.row.name)" :disabled="scope.row.pod_ip==''">
            {{$t('bm.add.viewStatus')}}
          </el-button> -->
        </template>
      </el-table-column>
    </el-table>

    <el-dialog :visible.sync="execModal.visible"
        width="70%"
        :close-on-click-modal="false"
        title="Pod 终端"
        class="log-dialog">
      <span slot="title"
            class="modal-title">
        <span class="unimportant">Pod 名称:</span>
        {{execModal.podName}}
        <span class="unimportant">容器:</span>
        {{execModal.containerName}}
        <!-- <i class="el-icon-full-screen screen"
           @click="fullScreen(execModal.podName +'-debug')"></i> -->
      </span>
      <xterm-debug :id="execModal.podName +'-debug'"
        :clusterName="cluster"
        :namespace="namespace"
        serviceName="serviceName"
        :containerName="execModal.containerName"
        :podName="execModal.podName"
        :visible="execModal.visible"
        ref="debug"></xterm-debug>
    </el-dialog>

    <page-nav ref="page" :list=filteredList></page-nav>
    <pod-status ref="view"></pod-status>
  </div>
</template>
<script>
import { mapGetters } from 'vuex';
import PageNav from '@/components/utils/Page';
import ListSearch from '@/components/utils/ListSearch';
import Refresh from '@/components/utils/Refresh';
import listTemplate from '@/common/listTemplate';
import XtermDebug from '@/components/utils/XtermDebug'
// import PodStatus from './PodViewStatus';

export default {
  mixins: [listTemplate],
  props: ['list', 'cluster', 'appName'],
  data() {
    return {
      searchList: [
        { key: 'name', txt: this.$t('bm.deployCenter.name') },
        { key: 'namespace', txt: this.$t('bm.serviceM.resourceSpace') },
        { key: 'node_ip', txt: this.$t('bm.add.insIP') },
        { key: 'pod_ip', txt: this.$t('bm.add.PodIP') },
        { key: 'image', txt: this.$t('bm.serviceM.mirror') },
        { key: 'status', txt: this.$t('bm.deployCenter.status') },
        { key: 'start_time', txt: this.$t('bm.add.startTime') },
      ],
      filterTxt: '',
      cluster: '',
      namespace: '',
      execModal: {
        visible: false,
        podName: null,
        containerName: null
      },
    };
  },
  computed: {
  },
  components: {
    PageNav,
    ListSearch,
    Refresh,
    XtermDebug,
    // PodStatus,
  },
  computed: {
    ...mapGetters({
      loading: 'getLoading',
    }),
    filteredList() {
      const filterVal = this.filterTxt;
      const list = this.$props.list || [];
      return filterVal.length > 0
        ? list.filter((element) => {
          return element[this.searchSelect].indexOf(filterVal) > -1;
        })
        : list;
    },
  },
  methods: {
    showContainerExec (namespace, pod_name, container_name) {
      this.cluster = this.$route.params.clusterName
      this.namespace = namespace
      this.execModal.visible = true
      this.execModal.podName = pod_name
      this.execModal.containerName = container_name
    }
  },
};
</script>
