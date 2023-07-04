<style>
.createDialog .el-dialog__body .el-form-item .el-form-item__content {
  display: flex;
  flex-direction: column;
}
</style>
<template>
  <el-dialog top='25vh' :title="title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog"  width='50%' :before-close="doCancelCreate">
    <el-form :model="form" ref="ruleForm" :rules="rules">
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" auto-complete="off" maxlength="64" placeholder="请输入名称" :disabled="systemReserved"></el-input>
      </el-form-item>
      <el-form-item label="镜像地址" prop="image">
        <el-input v-model.trim="form.image" auto-complete="off" maxlength="256" placeholder="请输入镜像地址"></el-input>
      </el-form-item>
      <el-form-item label="启动命令" prop="command">
        <el-input v-model="form.command" auto-complete="off" maxlength="128" placeholder="请输入启动命令"></el-input>
      </el-form-item>
      <el-form-item label="启动参数" prop="args">
        <el-input v-model="form.args" auto-complete="off" maxlength="128" placeholder="请输入启动参数"></el-input>
      </el-form-item>
      <el-form-item label="描述" prop="description">
        <el-input v-model="form.description" auto-complete="off" maxlength="256gukmf " placeholder="请输入描述" ></el-input>
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
  type: '',
  description: '',
  task_item: '',
};

export default {
  mixins: [createTemplate, validate],
  data() {
    return {
      name: '',
      systemReserved: false,
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
        image: [
          { required: true, message: '请输入镜像地址', trigger: 'blur' },
        ],
        command: [
          { required: false, message: '请输入启动命令', trigger: 'blur' },
        ],
        args: [
          { required: false, message: '请输入启动参数', trigger: 'blur' },
        ],
        description: [
          { required: false, message: '描述信息不能为空', trigger: 'blur' },
        ]
      },
      rowId: '',
    };
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
  },
  mounted() {
    backend.getStepComponent((data) => {
      if(data){
        this.typeList = data;
      }
    });
    
  },
  methods: {

    doCreate(flag, item) {
      this.isEdit = flag;
      if (flag) {
        this.title = '编辑';
        this.form = {
          name: item.name || '',
          image: item.image || '',
          command: item.command || '',
          args: item.args || '',
          description: item.description || '',
        };
        if (this.form.name == 'jnlp' || this.form.name == 'checkout' || this.form.name == 'kaniko') {
          this.systemReserved = true
        } else {
          this.systemReserved = false
        }
        this.rowId = item.id;
      } else {
        this.title = '新增';
        this.systemReserved = false;
        this.form = {
          name: '',
          image: '',
          command: '',
          args: '',
          description: '',
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
            image: this.form.image,
            command: this.form.command,
            args: this.form.args,
            description: this.form.description,
          };
          if (this.isEdit) {
            backend.editCompileEnv(this.rowId, cl, () => {
              successCallBack();
            });
          } else {
            backend.AddCompileEnv(cl, () => {
              successCallBack();
            });
          }
        }
      });
    },
  },
};
</script>
