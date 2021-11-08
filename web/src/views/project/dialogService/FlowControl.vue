<style>
.createDialog.hostCreate .lbsList.el-col-20 .el-input__inner {
  vertical-align: -1px;
}
</style>
<template>
  <el-dialog top='25vh' :title="$t('bm.add.flowContrl')" :visible.sync="dialogFormVisible" class="createDialog hostCreate" width='50%' :before-close="handleClose">
    <el-table
      border
      :data="tableList"
      style="width: 100%">
      <el-table-column
        prop="version"
        :label="$t('bm.add.serviceVer')"
        min-width="35%">
      </el-table-column>
      <el-table-column
        prop="stage"
        :label="$t('bm.deployCenter.status')"
        min-width="30%">
      </el-table-column>
    </el-table>
    <div slot="footer" class="dialog-footer" style='margin-top:30px'>
      <el-button @click="handleClose">{{$t('bm.other.cancel')}}</el-button>
      <el-button type="primary" @click="doSubmit" :loading="loading">{{$t('bm.other.confirm')}}</el-button>
    </div>
  </el-dialog>
</template>
<script>
import { mapGetters } from 'vuex';
import { Message } from 'element-ui';
import backend from '../../../api/backend';
import createTemplate from '../../../common/createTemplate';

// 全局唯一变量表示加载次数
// let loadedData = 0
export default {
  mixins: [createTemplate],
  props: ['grayList'],
  data() {
    return {
      dataList: [],
      // 是否属于编辑状态
      isEdit: false,
      clusterName: '',
      namespace: '',
      appName: '',
      dialogFormVisible: false,
    };
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
    tableList() {
      return this.$props.grayList;
    },
  },
  methods: {
    doCreate(flag, item, clusterName, namespace, appName) {
      this.dialogFormVisible = true;
      this.clusterName = clusterName;
      this.namespace = namespace;
      this.appName = appName;
      for (const a of this.tableList) {
        this.tableList[a].weight = item[a].weight;
      }
    },
    handleClose() {
      this.dialogFormVisible = false;
      this.$emit('getlist');
    },
    handleChange(value, index) {
      if (index === 0) {
        this.tableList[1].weight = 100 - value[0];
      } else {
        this.tableList[0].weight = 100 - value[0];
      }
    },
    doSubmit() {
      const successCallBack = () => {
        this.$emit('getlist');
        Message.success(this.$t('bm.add.optionSuc'));
        this.dialogFormVisible = false;
      };
      const param = [];
      for (const a of this.tableList) {
        param.push({
          version: a.version,
          weight: a.weight,
        });
      }
    },
  },
};
</script>
