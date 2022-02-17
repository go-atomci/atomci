<template>
  <div class="page-content buManage">
    <div class="portlet-body">
      <div class="table-toolbar">
        <el-row>
          <el-col :span="10">
            <refresh v-on:getlist="getList"></refresh>
            <el-button :plain="true"
                      type="primary"
                      @click="$refs.create.doCreate(false)">
              <i class='icon-plus' /> {{ isType === 'new' ? $t('bm.add.newRole') : $t('bm.add.bindingRole') }}</el-button>
          </el-col>
          <el-col :span="6">
            &nbsp;
          </el-col>
          <el-col :span="8">
            <list-search :searchList="searchList"
                        v-on:changeFilterTxt="changeFilterTxt"></list-search>
          </el-col>
        </el-row>
      </div>
      <template>
        <el-table stripe :data="curList">
          <el-table-column prop="role"
                           :label="$t('bm.add.roleName')"
                           sortable
                           min-width="15%"
                           :show-overflow-tooltip=true />
          <el-table-column prop="description"
                           :label="$t('bm.serviceM.description')"
                           sortable
                           min-width="15%"
                           :show-overflow-tooltip=true />
          <el-table-column prop="create_at"
                           :label="$t('bm.serviceM.creationTime')"
                           sortable
                           min-width="15%"
                           :show-overflow-tooltip=true />
          <el-table-column :label="$t('bm.deployCenter.operation')"
                           min-width="10%"
                           v-if="isUser">
            <template slot-scope="scope">
              <el-button @click="$refs.commonDelete.doDeleteBody('deleteGroupUserRole', scope.row.role,group,$route.params.user)"
                         type="text"
                         size="small"
                         :title="$t('bm.depManage.remove')">
                {{$t('bm.depManage.remove')}}
              </el-button>
            </template>
          </el-table-column>
          <el-table-column :label="$t('bm.deployCenter.operation')"
                           min-width="15%"
                           v-else>
            <template slot-scope="scope">
              <el-button @click="goDetail(scope.row.role)"
                         type="text"
                         size="small"
                         :title="$t('bm.authorManage.manage')">
                管理
              </el-button>
              <el-button @click="$refs.create.doCreate(true, scope.row)"
                         type="text"
                         size="small"
                         :title="$t('bm.infrast.edit')">
                {{$t('bm.infrast.edit')}}
              </el-button>
              <el-button @click="$refs.commonDelete.doDelete('delGroupRole',group,scope.row.role)"
                         type="text"
                         size="small"
                         :title="$t('bm.other.delete')">
                {{$t('bm.other.delete')}}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </template>
      <page-nav ref="page"
                :list="filteredList"></page-nav>
      <common-delete ref="commonDelete"
                     v-on:getlist="getList"></common-delete>
      <role-create ref="create" v-on:getlist="getList" 
                    :isTypeTitle="isType"
                    :isUser="isUser"></role-create>
    </div>
  </div>
</template>
<script>
import { mapGetters } from 'vuex';
import backend from '@/api/backend';
import PageNav from '@/components/utils/Page';
import ListSearch from '@/components/utils/ListSearch';
import RoleCreate from '@/components/utils/user/RoleCreate';
import CommonDelete from '@/components/utils/Delete';
import Refresh from '@/components/utils/Refresh';
import listTemplate from '@/common/listTemplate';
import UtilsFn from '@/common/utils';

export default {
  mixins: [listTemplate],
  data() {
    return {
      activeName: 'user',
      curList: [],
      searchList: [
        {
          key: 'role',
          txt: this.$t('bm.add.roleName'),
        },
        {
          key: 'description',
          txt: this.$t('bm.serviceM.description'),
        },
      ],
      filterTxt: '',
      // TODO: default group use system, tmp
      group: 'system',
      isUser: false,
      isType: 'new',
    };
  },
  components: {
    PageNav,
    ListSearch,
    Refresh,
    RoleCreate,
    CommonDelete,
  },
  computed: {
    ...mapGetters({
      userInfo: 'getUserInfo',
    }),
  },
  created() {
    this.searchList = [
      {
        key: 'role',
        txt: this.$t('bm.add.roleName'),
      },
      {
        key: 'description',
        txt: this.$t('bm.serviceM.description'),
      },
      {
        key: 'create_at',
        txt: this.$t('bm.serviceM.creationTime'),
      },
    ];
    if (this.userInfo.groups) {
      this.group = this.userInfo.groups[0].group;
    }
  },
  mounted() {
    this.getList();
  },
  methods: {
    goDetail(roles) {
      this.$router.push({
        name: 'listPermission',
        params: {
          role: roles,
        },
      });
    },
    getList() {
        backend.getGroupRoleList((data) => {
          if (data) {
            this.curList = data.map(item => {
              item.create_at = UtilsFn.format(new Date(item.create_at), 'yyyy-MM-dd hh:mm');
              return item;
            });
          }
        });
      }
    },
};
</script>
