<style>
.createDialog .el-dialog__body .el-form-item .el-form-item__content {
  display: flex;
  flex-direction: column;
}
</style>
<template>
  <el-dialog top='25vh' :title="title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog"  width='50%' :before-close="doCancelCreate">
    <el-form :model="form" ref="ruleForm" :rules="rules">
      <el-form-item label="流程名称" prop="name">
        <el-input v-model.trim="form.name" auto-complete="off" maxlength="63" placeholder="请输入名称"></el-input>
      </el-form-item>
      <el-form-item label="流程描述" prop="description">
        <el-input v-model.trim="form.description" auto-complete="off" maxlength="254" placeholder="请输入描述" ></el-input>
      </el-form-item>
      <el-form-item label="默认流程" prop="is_default">
        <el-checkbox label="" v-model="form.is_default"></el-checkbox>
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
import { Message } from 'element-ui';
import backend from '@/api/backend';
import createTemplate from '@/common/createTemplate';
import validate from '@/common/validate';

const formData = {
  name: '',
  description: '',
  enabled: false,
  is_default: false
};

export default {
  mixins: [createTemplate, validate],
  data() {
    return {
      name: '',
      groupRoleList: [],
      typeList: [],
      // 是否属于编辑状态
      isEdit: false,
      dialogFormVisible: false,
      form: JSON.parse(JSON.stringify(formData)),
      title: '新增',
      rules: {
        name: [
          { required: true, message: '请输入名称', trigger: 'blur' },
        ],
        description: [
          { required: true, message: '描述信息不能为空', trigger: 'blur' },
        ],
      },
      rowId: '',
    };
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
      projectID: 'projectID',
    }),
  },
  mounted() {
    
  },
  methods: {
    doCreate(flag, item) {
      this.isEdit = flag;
      if (flag) {
        this.title = '编辑';
        this.form = {
          name: item.name || '',
          description: item.description || '',
          enabled: item.enabled || false,
          is_default: item.is_default || false
        };
        this.rowId = item.id;
      } else {
        this.title = '新增';
        this.form = {
          name: '',
          description: '',
          enabled: false,
          is_default: false
        };
        this.rowId = '';
      }
      this.dialogFormVisible = true;
      this.isEdit = flag;
    },
    doSubmit() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          const successCallBack = () => {
            this.$emit('getlist');
            Message.success(this.$t('bm.add.optionSuc'));
            this.dialogFormVisible = false;
          };
          const cl = {
            name: this.form.name,
            project_id: this.projectID,
            description: this.form.description,
            is_default: this.form.is_default,
          };
          if (this.isEdit) {
            backend.editProjectPipe(this.projectID, this.rowId, cl, () => {
              successCallBack();
            });
          } else {
            backend.addProjectPipe(this.projectID, cl, () => {
              successCallBack();
            });
          }
        }
      });
    },
  },
};
</script>
