<style>
.createDialog .el-dialog__body .el-form-item .el-form-item__content {
  display: flex;
  flex-direction: column;
}
</style>
<template>
  <el-dialog top='25vh' :title="title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog"  width='50%' :before-close="doCancelCreate">
    <el-form :model="form" ref="ruleForm" :rules="rules">
      <el-form-item label="环境名称" prop="type">
        <el-select v-model="form.type" placeholder="请选择环境名称" filterable>
          <el-option v-for="(item, index) in envList" :key="index" :label="item.name" :value="index">
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="doCancelCreate">{{$t('bm.other.cancel')}}</el-button>
      <el-button type="primary" @click="doSubmit" :loading="loading">{{$t('bm.other.confirm')}}</el-button>
    </div>
  </el-dialog>
</template>
<script>
import { mapGetters } from 'vuex';
import backend from '@/api/backend';
import createTemplate from '@/common/createTemplate';
import validate from '@/common/validate';

export default {
  mixins: [createTemplate, validate],
  data() {
    return {
      name: '',
      // 是否属于编辑状态
      isEdit: false,
      dialogFormVisible: false,
      form: {type: ''},
      title: '添加环境',
      rules: {
        type: [
          { required: true, message: '请选择环境', trigger: 'blur' },
        ],
      },
      rowId: '',
      thisIndex: 0,
    };
  },
  props: {
    envList: {
      type: Array,
      default: []
    }
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
      projectID: 'projectID',
    }),
  },
  methods: {
    doCreate(index) {
      this.thisIndex = index + 1;
      this.form = {
        type: '',
      };
      this.dialogFormVisible = true;
    },
    doSubmit() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          const formCheck = this.envList[this.form.type];
          const obj = {'name': formCheck.name,'description': formCheck.description, 'stage_id': formCheck.id, 'steps': []};
          this.$emit('listAdd',this.thisIndex, obj);
        }
      });
    },
  },
};
</script>
