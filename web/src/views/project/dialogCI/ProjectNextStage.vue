<template>
  <el-dialog z-index="1100" :close-on-click-modal="false"  style="width:50%;margin-left:25%" top='25vh' :visible.sync="dialogFormVisible" class="createDialog"  :before-close="doCancelCreate">
    <el-form>
      <el-form-item :label="$t('bm.add.nextStage')">
        <el-select v-model="stage" :placeholder="$t('bm.add.select')" @change="change" filterable style="padding-bottom:20px">
          <el-option v-for="(item, index) in tableData" :key="index" :label="item.name" :value="item.id">
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>
<script>
import { mapGetters } from 'vuex';
import { Message } from 'element-ui';
import backend from '@/api/backend';

export default {
  // inject: [getList],
  data() {
    return {
      idstage: '',
      stage: '',
      tableData: [],
      // 是否属于编辑状态
      isEdit: false,
      // 是否属于查看状态
      isView: false,
      dialogFormVisible: false,
      showandhide: true,
      id: '',
    };
  },
  computed: {
    ...mapGetters({
      projectIDgetter: 'projectID',
    }),
    projectID() {
        if (this.projectIDgetter === 0 || this.projectIDgetter === undefined) {
          this.$store.dispatch('project/setProjectID', this.$route.params.projectID);
          return this.$route.params.projectID
        } else {
          return this.projectIDgetter
        }
    },
  },
  methods: {
    change(val) {
      const body = {
        stage_id: this.stage,
      };
      const that = this;
      backend.gocontinue(this.projectID, this.id, this.idstage, 'next-stage', JSON.stringify(body), (data) => {
        Message.success(this.$t('bm.add.optionSuc'));
        that.$emit('getprojectReleaseList');
      });
      // this.getList()
      this.stage = '';
      this.dialogFormVisible = false;
    },
    show(id, idstage) {
      this.idstage = idstage;
      this.id = id;
      backend.goregression(this.projectID, id, this.idstage, 'next-stage', (data) => {
        this.tableData = data;
        this.dialogFormVisible = true;
      });
    },
    doCancelCreate() {
      this.dialogFormVisible = false;
      this.isEdit = false;
      this.isView = false;
      this.description = '';
    },
  },
};
</script>
