<style>
.createDialog .el-dialog__body .el-form-item .el-form-item__content {
  display: flex;
  flex-direction: column;
}
</style>
<template>
  <el-dialog top='25vh' :title="title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog"  width='50%' :before-close="doCancelCreate">
    <el-form :model="form" ref="ruleForm" :rules="rules">
      <el-form-item label="应用名" prop="name">
        <el-select  allow-create filterable default-first-option v-model="form.namespace" placeholder="请选择应用名">
          <el-option v-for="(item, index) in scmAppList" :key="index" :label="item.full_name" :value="item.id">
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
import { Message } from 'element-ui';
import backend from '@/api/backend';
import createTemplate from '@/common/createTemplate';
import validate from '@/common/validate';

const formData = {
  name: '',
  description: '',
  arrange_env: '',
};

export default {
  mixins: [createTemplate, validate],
  data() {
    return {
      name: '',
      // 是否属于编辑状态
      isEdit: false,
      dialogFormVisible: false,
      form: JSON.parse(JSON.stringify(formData)),
      title: '新增',
      rules: {
        name: [
          { required: true, message: '请输入名称', trigger: 'blur' },
        ],
      },
      rowId: '',
    };
  },
  props: {
    scmAppList: {
      type: Array,
      default: []
    }
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
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
  created() {},
  methods: {
    doCreate(flag, item) {
      this.isEdit = flag;
      if (flag) {
        this.title = '编辑';
        this.form = {
          name: item.name || '',
          app_id: 0,
        };
        this.rowId = item.id;
      } else {
        this.title = '新增';
        this.form = {
          name: '',
          app_id: 0,
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
            app_id: this.form.app_id,
          };
          if (this.isEdit) {
            backend.updateProjectApp(this.projectID, this.rowId, cl, () => {
              successCallBack();
            });
          } else {
            backend.addProjectApp(this.projectID, cl, () => {
              successCallBack();
            });
          }
        }
      });
    },
  },
};
</script>
