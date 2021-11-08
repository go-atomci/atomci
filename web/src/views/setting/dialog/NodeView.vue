<style>
.createDialog .el-dialog__body .el-form-item .el-form-item__content {
  display: flex;
  flex-direction: column;
}
.nodeView .el-table__body-wrapper {
  max-height: 350px !important;
  overflow-y: auto !important;
}
</style>
<template>
  <el-dialog top='15vh' :title="title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog"  width='60%' :before-close="doClose">
    <div class="nodeView">
      <el-table border :data="colData">
        <span slot="empty">
          {{loading?$t('bm.add.dataLoading'):'暂无引用流程'}}
        </span>
        <el-table-column prop="name" label="流程名称" sortable min-width="15%" :show-overflow-tooltip=true />
        <el-table-column prop="enabled" label="状态" min-width="15%" :show-overflow-tooltip=true />
        <el-table-column prop="description" :label="$t('bm.serviceM.description')" min-width="15%" :show-overflow-tooltip=true />
      </el-table>
    </div>
    <div slot="footer" class="dialog-footer">
      <el-button @click="doClose">关闭</el-button>
    </div>
  </el-dialog>
</template>
<script>
import { mapGetters } from 'vuex';
import backend from '@/api/backend';
import createTemplate from '@/common/createTemplate';

export default {
  mixins: [createTemplate],
  data() {
    return {
      dialogFormVisible: false,
      title: '查看引用流程',
      colData: [],
    };
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
  },
  methods: {
    doView(stepId) {
      this.colData = [];
      backend.getPipeRow(stepId, (data) => {
        this.colData = data;
        if(this.colData) {
          this.colData.map((i) => {
            i.enabled = i.enabled === false ? '禁用' : '启用';
          });
        }
      });
      this.dialogFormVisible = true;
    },
    doClose() {
      this.dialogFormVisible = false;
    }
  },
};
</script>
