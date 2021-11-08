

<template>
  <div>
     <div class="table-toolbar">
        <el-row>
          <el-col :span="10">
            <refresh v-on:getlist="getList"></refresh>
            <el-button :plain="true" type="primary" @click="$refs.create.doCreate(false)">
              <i class='icon-plus' /> 添加约束</el-button>
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
          <el-table-column prop="key" label="约束" sortable min-width="15%" :show-overflow-tooltip=true />
          <el-table-column label="值" sortable min-width="15%" :show-overflow-tooltip=true>
            <template slot-scope="scope">
              {{ scope.row.value.join('、') }}
            </template>
          </el-table-column>
          <el-table-column :label="$t('bm.deployCenter.operation')" min-width="10%">
            <template slot-scope="scope">
              <el-button @click="$refs.create.doCreate(true, scope.row)" type="text" size="small" :title="$t('bm.infrast.edit')">
                {{$t('bm.infrast.edit')}}
              </el-button>
              <el-button @click="$refs.commonDelete.doDelete('deleteUserConstraints', 'system', $route.params.user, scope.row.key)" type="text" size="small" :title="$t('bm.other.delete')">
                {{$t('bm.other.delete')}}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </template>
      <page-nav ref="page" :list=filteredList></page-nav>
      <common-delete ref="commonDelete" v-on:getlist="getList"></common-delete>
      <constraints-create ref="create" v-on:getlist="getList" :isUser="isUser"></constraints-create>
      <dept-per-config ref="perConfig"></dept-per-config>
  </div>
</template>
<script>
import { mapGetters } from 'vuex';
import backend from '../../api/backend';
import PageNav from '../utils/Page';
import ListSearch from '../utils/ListSearch';
import ConstraintsCreate from '../utils/user/UserConstraints';
import CommonDelete from '../utils/Delete';
import Refresh from '../utils/Refresh';
import listTemplate from '../../common/listTemplate';
import DeptPerConfig from '../utils/user/DeptPerConfig';
import UtilsFn from '../../common/utils';

export default {
  mixins: [listTemplate],
  props: ['isUser'],
  data() {
    return {
      curList: [],
      searchList: [
        { key: 'key', txt: '约束' },
        { key: 'value', txt: '值' },
      ],
      filterTxt: '',
      resourceTypeList: [],
      resourceOpList: [],
      resourceConList: [],
    };
  },
  components: {
    PageNav,
    ListSearch,
    Refresh,
    ConstraintsCreate,
    CommonDelete,
    DeptPerConfig,
  },
  methods: {
    getList() {
      this.curList = [];
      const defaultGroup = "system"
      backend.getConstraintsList(defaultGroup, this.$route.params.user, (data) => {
        if (data) {
          Object.keys(data).forEach((item, index) => {
            this.curList.push({
              key: Object.keys(data)[index],
              value: Object.values(data)[index],
            });
          });
        }
      });
    },
  },
};
</script>
