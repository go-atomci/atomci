<style>
.createDialog .el-dialog__body .el-form-item .el-form-item__content {
  display: flex;
  flex-direction: column;
}
</style>
<template>
  <el-dialog top='25vh' :title="title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog"  width='50%' :before-close="doCancelCreate">
    <el-form :model="form" ref="ruleForm" :rules="rules">
      <div v-show="!isFirst">
        <el-form-item label="上一节点信息">
          {{name}}
        </el-form-item>
        <el-form-item label="处理人" prop="user">
          <el-input v-model.trim="form.user" :disabled="true"></el-input>
        </el-form-item>
        <el-form-item label="处理备注" prop="description">
          <el-input v-model.trim="form.description" type="textarea" :rows="3" :disabled="true"></el-input>
        </el-form-item>
      </div>
      <el-form-item label="处理备注" prop="remark">
        <el-input v-model.trim="form.remark" auto-complete="off" type="textarea" :rows="3" maxlength="120" placeholder="请输入备注" ></el-input>
      </el-form-item>
      <el-form-item label="是否发送邮件" prop="enable">
        <el-checkbox label="是" v-model="form.enable"></el-checkbox>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="doSubmit(false)">驳回</el-button>
      <el-button type="primary" @click="doSubmit(true)" :loading="loading">通过</el-button>
    </div>
  </el-dialog>
</template>
<script>
import { mapGetters } from 'vuex';
import { Message } from 'element-ui';
import backend from '../../../api/backend';
import createTemplate from '../../../common/createTemplate';
import validate from '../../../common/validate';

export default {
  mixins: [createTemplate, validate],
  data() {
    return {
      // 是否属于编辑状态
      isFirst: true,
      name: '',
      dialogFormVisible: false,
      form: {
        user: '',
        description: '',
        remark: '',
        enable: false,
      },
      title: '',
      rules: {
        remark: [
          { required: true, message: '备注不能为空', trigger: 'blur' },
        ],
      },
      list: {},
      publishId: '',
      stageId: '',
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
    show(item) {
      this.publishId = item.id;
      this.stageId = item.stage_id;
      backend.getMark(this.projectID, item.id, item.stage_id, (data) => {
        this.list = data;
        const nev = data.previous_step;
        if(nev && JSON.stringify(nev) != '{}') {
          this.form.user = nev.creator;
          this.form.description = nev.message;
          this.name = nev.name;
          this.form.remark = '';
          this.form.enable = false;
          this.isFirst = false;
        } else {
          this.form = {
            user: '',
            description: '',
            remark: '',
            enable: false
          };
          this.name = '';
          this.isFirst = true;
        }
        this.title = item.step || '';
        this.dialogFormVisible = true;
      });
      
    },
    doSubmit(flag) {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          const successCallBack = () => {
            this.$emit('getprojectReleaseList', true);
            Message.success(this.$t('bm.add.optionSuc'));
            this.dialogFormVisible = false;
          };
          const cl = {
            message: this.form.remark,
            enable_notify: this.form.enable || false,
          };
          if(flag) {
            cl.status = 'success'
          } else {
            cl.status = 'failed';
          }
          backend.setMark(this.projectID, this.publishId, this.stageId, cl, (data) => {
            successCallBack();
          });
        }
      });
    },
  },
};
</script>
