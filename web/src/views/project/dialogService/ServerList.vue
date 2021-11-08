<template>
  <div>
    <div id="top" v-if="list">
      <ul>
        <li>{{$t('bm.add.serviceType')}}:{{list.type}}</li>
        <li>{{$t('bm.serviceM.cluster')}}IP:{{list.cluster_ip}}</li>
      </ul>
      <span id="title1" v-if="list.address_list===null?false:true">{{$t('bm.add.serviceAdrChart')}}</span>
      <div class="table-toolbar"></div>
      <el-table
      border
      v-if="list.address_list===null?false:true"
      :data="list.address_list"
      style="width: 100%">
      <el-table-column
          prop="cluster_addr"
          :label="$t('bm.add.inSideAdr')"
          min-width="30%"
          >
      </el-table-column>
      <el-table-column
          prop="external_addr"
          :label="$t('bm.add.outSideAdr')"
          min-width="20%"
        >
      </el-table-column>
      <el-table-column
          prop="node_port_addr"
          :label="$t('bm.add.nodeportAdr')"
          min-width="20%">
      </el-table-column>
      <el-table-column
          prop="target_port"
          :label="$t('bm.serviceM.containerPort')"
          min-width="15%">
      </el-table-column>
      <el-table-column
          prop="protocol"
          :label="$t('bm.infrast.agreement')"
          min-width="15%">
      </el-table-column>
          <page-nav ref="page" :list=filteredList></page-nav>
          </el-table>
          <span id="title2" v-if="list.podsvc_addr_list===null?false:true">{{$t('bm.add.podAdrChart')}}</span>
          <el-table
          v-if="list.podsvc_addr_list===null?false:true"
          :data="list.podsvc_addr_list"
          style="width: 100%">
          <el-table-column
              prop="cluster_addr"
              :label="$t('bm.add.podClusterDomain')"
              min-width="50%">
          </el-table-column>
          <el-table-column
              prop="target_port"
              :label="$t('bm.serviceM.containerPort')"
              min-width="20%">
          </el-table-column>
          <el-table-column
              prop="protocol"
              :label="$t('bm.infrast.agreement')"
              min-width="30%">
          </el-table-column>
          </el-table>
        <page-nav ref="page" :list=filteredList></page-nav>
    </div>
  </div>
</template>
<script>
import listTemplate from '@/common/listTemplate';
import PageNav from '@/components/utils/Page';

export default {
  mixins: [listTemplate],
  props: ['list'],
  data() {
    return {
      searchList: [
        { key: 'cluster_addr', txt: this.$t('bm.add.serviceCluster') },
        { key: 'external_addr', txt: this.$t('bm.add.externalAdr') },
        { key: 'node_port_addr', txt: this.$t('bm.add.nodeportAdr') },
        { key: 'target_port', txt: this.$t('bm.serviceM.containerPort') },
        { key: 'protocol', txt: this.$t('bm.infrast.agreement') },
      ],
      filterTxt: '',
    };
  },
  components: {
    PageNav,
  },
};
</script>
