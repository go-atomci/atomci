<style>
.clusterTab {
  margin-top: -10px;
}
.clusterTab .el-tabs__header {
  margin-bottom: 10px;
}
</style>

<template>
  <el-tabs class="clusterTab" v-model="cluster" @tab-click="tabClick">
    <el-tab-pane v-for="(item, index) in clusterList || []" :key="index" :label="item.name" :name="item.name">
    </el-tab-pane>
  </el-tabs>
</template>
<script>
import { mapGetters } from 'vuex';
import backend from '../../api/backend';
export default {
  props: ['clusterName', 'fwC'],
  data() {
    return {
      selectedCluster: '',
      cluster: '',
      clusterGive: this.$props.fwC ? this.$props.fwC : 0,
    };
  },
  watch: {
    cluster(newVal, oldVal) {
	  if(!newVal || !oldVal) return;
      if (newVal === '0') {
        if (this.curClusterList && this.curClusterList.length) {
          this.cluster = this.curClusterList[0].name;
        } else {
          this.cluster = this.clusterName;
        }
        return;
      }
      if (newVal !== oldVal) {
        this.$emit('changeCluster', this.cluster);
      }
    },
    clusterList(newVal) {
      this.cluster = newVal[0].name;
    },
  },
  computed: {
    ...mapGetters({
      clusterList: 'getClusterList',
    }),
  },
  created() {
    // 初始化的时候数据如果有，就直接获取，如果没有就通过watch监听获取
    this.getClusterList();
  },
  methods: {
    tabClick(tab) {
      if (this.selectedCluster !== tab.name) {
        this.selectedCluster = tab.name;
        this.cluster = tab.name;
      }
    },
    getClusterList() {
      // 首次页面加载，需要获取集群信息；如果再次加载，数据已缓存的情况，则使用缓存数据
      if (!this.clusterList || !this.clusterList.length) {
        backend.getClusterList((data) => {
          this.$store.dispatch('setClusterList', data);
          if(data) this.cluster = data[0].name;
        });
      } else {
        this.cluster = this.clusterList[0].name;
        this.$emit('changeCluster', this.cluster);
      }
    },
  },
};
</script>
