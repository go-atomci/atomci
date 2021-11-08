<template>
    <div class="portlet-body">
      <template>
        <div class="table-toolbar">
          <el-row>
            <el-col :span="10">
              <refresh v-on:getlist="getList('clear')"></refresh>
              <el-button :plain="true"
                         type="primary"
                         @click="$refs.create.doCreate(false)">
                <i class='icon-plus' /> {{$t('bm.authorManage.createUser')}}</el-button>
            </el-col>
            <el-col :span="6">
              &nbsp;
            </el-col>
            <el-col :span="8">
              <list-search ref="userSh"
                           :searchList="searchList"
                           v-on:changeFilterTxt="changeFilterTxt"></list-search>
            </el-col>
          </el-row>
        </div>
        <template>
          <el-table stripe :data="dataList">
            <span slot="empty">
              {{loading?$t('bm.add.dataLoading'):noDataTxt}}
            </span>
            <el-table-column prop="user"
                             label="帐号"
                             min-width="24%"
                             :show-overflow-tooltip=true />
            <el-table-column prop="name"
                             label="用户名"
                             sortable
                             min-width="20%"
                             :show-overflow-tooltip=true />
            <el-table-column prop="email"
                             :label="$t('bm.infrast.email')"
                             min-width="24%"
                             :show-overflow-tooltip=true />
            <el-table-column sortable prop="login_type" min-width="20%" label="登录方式">
            <template slot-scope="scope">
              <span v-if="scope.row.login_type === 2">LDAP</span>
              <span v-else>本地验证</span>
            </template>
          </el-table-column>
            <el-table-column prop="create_at"
                             :label="$t('bm.serviceM.creationTime')"
                             sortable
                             min-width="18%"
                             :show-overflow-tooltip=true />
            <el-table-column :label="$t('bm.deployCenter.operation')"
                             min-width="14%">
              <template slot-scope="scope"
                        v-if="scope.row.name!=='default'">
                <el-button type="text"
                           size="small"
                           @click="goUserAuthorized(scope.row.user)">{{$t('bm.authorManage.manage')}}
                </el-button>
                <el-button @click="$refs.commonDelete.doDelete('delUser',scope.row.user)"
                           type="text"
                           size="small">{{$t('bm.other.delete')}}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </template>
        <page-nav ref="page" :list="filteredList"></page-nav>
      </template>
      <common-delete ref="commonDelete" v-on:getlist="getList"></common-delete>
      <user-create ref="create" v-on:getlist="getList"></user-create>
    </div>
</template>
<script>
import { mapGetters } from 'vuex';
import backend from '@/api/backend';
import PageNav from '@/components/utils/Page';
import ListSearch from '@/components/utils/ListSearch';
import UserCreate from '@/components/view/UserCreate';
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
        { key: 'name', txt: this.$t('bm.authorManage.userName') },
        { key: 'user', txt: this.$t('bm.depManage.userAccount') },
        { key: 'create_at', txt: this.$t('bm.serviceM.creationTime') },
      ],
      searchListResource: [
        { key: 'resource', txt: this.$t('bm.authorManage.resourceType') },
        { key: 'description', txt: this.$t('bm.serviceM.description') },
        { key: 'create_at', txt: this.$t('bm.serviceM.creationTime') },
      ],
      filterTxt: '',
      detailInfo: [],
      activeName: 'user',
    };
  },
  components: {
    PageNav,
    ListSearch,
    Refresh,
    UserCreate,
    CommonDelete,
  },
  watch: {
    $route() {
      this.activeName = 'user';
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
    // changeFilterTxt(val) {
    // },
    handleClick(tab) {
      this.activeName = tab.name;
      this.getList();
      this.$refs.page.currentPage = 1;
      this.$refs.userSh.filterTxt = '';
      this.$refs.groupSh.filterTxt = '';
    },
    goUserAuthorized(user) {
      this.$router.push({
        name: 'managementUser',
        params: {
          user,
        },
      });
    },
    goDetail(resourceType) {
      this.$router.push({
        name: 'resourceType',
        params: {
          resourceType,
        },
      });
    },
    // changeFilterTxt: function(val) {
    // },
    getList(isRefresh) {
      // this.changeFilterTxt()
      if (isRefresh) {
        this.$refs.page.currentPage = 1;
      }
      if (isRefresh === 'clear') {
        this.$refs.userSh.searchSelectChange();
        this.currentPage = 1;
      }
      this.curList = [];
      if (this.activeName === 'user') {
        backend.getUserList(data => {
          this.curList = data.map(item => {
            item.create_at = UtilsFn.format(new Date(item.create_at), 'yyyy-MM-dd hh:mm');
            return item;
          });
        });
      } else if (this.activeName === 'resource') {
        // 系统资源查询
        backend.getResourceTypeList(data => {
          this.curList = data.map(item => {
            item.create_at = UtilsFn.format(new Date(item.create_at), 'yyyy-MM-dd hh:mm');
            return item;
          });
        });
      }
    },
  },
};
</script>
