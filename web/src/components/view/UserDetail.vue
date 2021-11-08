<template>
  <div>
     <div class="table-toolbar">
        <el-row>
          <el-col :span="10">
            <refresh v-on:getlist="getList"></refresh>
            <el-button :plain="true" type="primary" @click="$refs.create.doCreate(false)">
              <i class='icon-plus' /> {{ isType === 'new' ? $t('bm.add.newRole') : $t('bm.add.bindingRole') }}</el-button>
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
          <el-table-column prop="role" :label="$t('bm.add.roleName')" sortable min-width="15%" :show-overflow-tooltip="true" />
          <el-table-column prop="description" :label="$t('bm.serviceM.description')" sortable min-width="15%" :show-overflow-tooltip="true" />
          <el-table-column prop="create_at" :label="$t('bm.serviceM.creationTime')" sortable min-width="15%" :show-overflow-tooltip="true" />
          <el-table-column :label="$t('bm.deployCenter.operation')" min-width="10%" v-if="isUser">
            <template slot-scope="scope">
              <el-button @click="$refs.commonDelete.doDeleteBody('deleteGroupUserRole', scope.row.role,'system',$route.params.user)" type="text" size="small" :title="$t('bm.depManage.remove')">
                {{$t('bm.depManage.remove')}}
              </el-button>
            </template>
          </el-table-column>
          <el-table-column :label="$t('bm.deployCenter.operation')" min-width="10%" v-else>
            <template slot-scope="scope">
              <el-button @click="goDetail($route.params.dept, scope.row.role)" type="text" size="small" :title="$t('bm.authorManage.manage')">
                {{$t('bm.authorManage.manage')}}
              </el-button>
              <el-button @click="$refs.create.doCreate(true, scope.row)" type="text" size="small" :title="$t('bm.infrast.edit')">
                {{$t('bm.infrast.edit')}}
              </el-button>
              <el-button @click="$refs.commonDelete.doDelete('delGroupRole','system',scope.row.role)" type="text" size="small" :title="$t('bm.other.delete')">
                {{$t('bm.other.delete')}}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </template>
      <page-nav ref="page" :list="filteredList"></page-nav>
      <common-delete ref="commonDelete" v-on:getlist="getList"></common-delete>
      <role-create ref="create" v-on:getlist="getList" :isUser="isUser" :isTypeTile="isType"></role-create>
  </div>
</template>
<script>
import backend from '../../api/backend';
import PageNav from '../utils/Page';
import ListSearch from '../utils/ListSearch';
import RoleCreate from '../utils/user/RoleCreate';
import CommonDelete from '../utils/Delete';
import Refresh from '../utils/Refresh';
import listTemplate from '../../common/listTemplate';
import UtilsFn from '../../common/utils';

export default {
  mixins: [listTemplate],
  props: ['isUser', 'isType'],
  data() {
    return {
      curList: [],
      searchList: [
        { key: 'role', txt: this.$t('bm.add.roleName') },
        { key: 'description', txt: this.$t('bm.serviceM.description') },
        { key: 'policy', txt: this.$t('bm.add.perPolicy') },
      ],
      filterTxt: '',
    };
  },
  components: {
    PageNav,
    ListSearch,
    Refresh,
    RoleCreate,
    CommonDelete,
  },
  created() {
    if (!this.$props.isUser) {
      this.searchList = [
        { key: 'role', txt: this.$t('bm.add.roleName') },
        { key: 'description', txt: this.$t('bm.serviceM.description') },
        { key: 'create_at', txt: this.$t('bm.serviceM.creationTime') },
      ];
    }
  },
  methods: {
    goDetail(group, roles) {
      this.$router.push({
        name: 'managementRole',
        params: {
          dept: group,
          role: roles,
        },
      });
    },
    getList() {
      if (this.$props.isUser) {
        // TODO: group's name use system , tmp
        backend.getGroupUserRole("system", this.$route.params.user, (data) => {
          if (data) {
            this.curList = data.map((item) => {
              item.create_at = item.create_at ? UtilsFn.format(new Date(item.create_at), 'yyyy-MM-dd hh:mm') : '';
              return item;
            });
          }
        });
      } else {
        // TODO: group's name use system , tmp
        backend.getGroupRoleList("system", (data) => {
          if (data) {
            this.curList = data.map((item) => {
              if (item.policies) {
                const policies = item.policies.map((subItem) => {
                  return subItem.policy_name;
                });
                item.policy = policies.join(' ');
              } else {
                item.policy = '';
              }
              item.create_at = item.create_at ? UtilsFn.format(new Date(item.create_at), 'yyyy-MM-dd hh:mm') : '';
              return item;
            });
          }
        });
      }
    },
  },
};
</script>
