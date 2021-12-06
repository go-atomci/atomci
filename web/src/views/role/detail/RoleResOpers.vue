<template>
  <div class="page-content buManage">
    <div class="portlet-body">
    <div class="table-toolbar">
        <el-row v-show="innerShow">
          <el-col :span="10">
            <refresh v-on:getlist="getList"></refresh>
            <el-button :plain="true" type="primary" @click="$refs.create.doCreate(false)">
              <i class='icon-plus' /> {{$t('bm.add.addResOperation')}}</el-button>
          </el-col>
          <el-col :span="6">
            &nbsp;
          </el-col>
          <el-col :span="8">
            <list-search :searchList="searchList" v-on:changeFilterTxt="changeFilterTxt"></list-search>
          </el-col>
        </el-row>
      </div>
      <template>
        <el-table border :data="dataList">
          <el-table-column prop="resource_operation" :label="$t('bm.authorManage.resourceOper')" sortable min-width="15%" :show-overflow-tooltip=true />
          <el-table-column prop="resource_type" :label="$t('bm.authorManage.resourceType')" sortable min-width="15%" :show-overflow-tooltip=true />
          <el-table-column prop="description" :label="$t('bm.serviceM.description')" sortable min-width="15%" :show-overflow-tooltip=true />
          <el-table-column prop="create_at" :label="$t('bm.serviceM.creationTime')" sortable min-width="15%" :show-overflow-tooltip=true />
          <el-table-column :label="$t('bm.deployCenter.operation')" min-width="10%">
            <template slot-scope="scope">
              <el-button @click="$refs.commonDelete.doDeleteBody('deleteRoleOperation', $route.params.role, scope.row.id)" type="text" size="small" :title="$t('bm.depManage.remove')">
                {{$t('bm.depManage.remove')}}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </template>
      <page-nav ref="page" :list="filteredList"></page-nav>
      <common-delete ref="commonDelete" v-on:getlist="getList"></common-delete>
      <role-permission-create ref="create" v-on:getlist="getList"></role-permission-create>
    </div>
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
import RolePermissionCreate from '../dialog/RolePermissionCreate';

export default {
  mixins: [listTemplate],
  data() {
    return {
      curList: [],
      searchList: [
        { key: 'policy_name', txt: this.$t('bm.add.policyName') },
        { key: 'description', txt: this.$t('bm.serviceM.description') },
        { key: 'create_at', txt: this.$t('bm.serviceM.creationTime') },
      ],
      filterTxt: '',
      resourceTypeList: [],
      resourceOpList: [],
      resourceConList: [],
      innerShow: true,
      group: '',
    };
  },
  components: {
    PageNav,
    ListSearch,
    Refresh,
    RolePermissionCreate,
    CommonDelete,
  },
  computed: {
    ...mapGetters({
      loading: 'getLoading',
    }),
  },
  created() {
    this.getList();
  },
  methods: {
    getList() {
      if(!this.$route.params.role) return;
      backend.getRoleOperations(this.$route.params.role, (data) => {
        if (data) {
          this.curList = data.map((item) => {
            item.create_at = UtilsFn.format(new Date(item.create_at), 'yyyy-MM-dd hh:mm');
            return item;
          });
        }
      });
    },
  },
  watch: {
    '$route': 'getList',
  },
};
</script>
