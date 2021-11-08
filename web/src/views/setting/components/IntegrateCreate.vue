<style>
.createDialog .el-dialog__body .el-form-item .el-form-item__content {
  display: flex;
  flex-direction: column;
}
.form-item {
  margin-left: 22px;
  margin-right: 22px;
}

.el-drawer__body {
  overflow-y: scroll;
  margin-bottom: 20px;
}

</style>
<template>
  <el-drawer :title="title" :visible.sync="dialogFormVisible" :direction="direction" :before-close="handleClose">
    <div style="overflow-y: scroll">
    <el-form :model="form" ref="ruleForm" :rules="rules">
      <el-form-item label="名称" prop="name" class="form-item">
        <el-input v-model.trim="form.name" auto-complete="off" maxlength="60" placeholder="请输入名称" :disabled="isKubernetes"></el-input>
      </el-form-item>
      <el-form-item label="配置类型" prop="type" class="form-item">
        <el-select v-model="form.type" placeholder="请选择" filterable :disabled="isEdit">
          <el-option v-for="(item, index) in settingTypeList" :key="index" :label="item.name" :value="item.name">
          </el-option>
        </el-select>
      </el-form-item>
      <div v-if="form.type ==='kubernetes'">
        <el-form-item label="Kubernetes URL" prop="config.url" class="form-item">
          <el-input v-model.trim="form.config.url" auto-complete="off" placeholder="请输入Kubernetes地址"></el-input>
        </el-form-item>
        <el-form-item label="Kubernetes Config" prop="config.conf" class="form-item">
          <el-input type="textarea" :rows="8" v-model="form.config.conf" placeholder="请输入Kubernetes Config"></el-input>
        </el-form-item>
      </div>
      <div v-else-if="form.type ==='jenkins'">
        <el-form-item label="Jenkins URL" prop="config.url" class="form-item">
          <el-input v-model.trim="form.config.url" auto-complete="off" placeholder="请输入, 如http://jenkins.example.com"></el-input>
        </el-form-item>

        <el-form-item label="Jenkins 用户名" prop="config.user" class="form-item">
          <el-input v-model.trim="form.config.user" auto-complete="off" maxlength="60" placeholder="请输入Jeknins 用户名"></el-input>
        </el-form-item>

        <el-form-item label="Jenkins 用户Token" prop="config.token" class="form-item">
          <el-input v-model.trim="form.config.token" auto-complete="off" maxlength="120" placeholder="请输入Jeknins 用户Token"></el-input>
        </el-form-item>

        <el-form-item label="Jenkins WorkSpace" prop="config.workspace" class="form-item">
          <el-input v-model.trim="form.config.workspace" auto-complete="off" maxlength="120" placeholder="请输入Jeknins workspace"></el-input>
        </el-form-item>
      </div >
      <div v-else-if="form.type ==='harbor'">
        <el-form-item label="Harbor URL" prop="config.url" class="form-item">
          <el-input v-model.trim="form.config.url" auto-complete="off" placeholder="请输入 Harbor 地址"></el-input>
        </el-form-item>

        <el-form-item label="用户名" prop="config.user" class="form-item">
          <el-input v-model.trim="form.config.user" auto-complete="off" maxlength="60" placeholder="请输入 Harbor 用户名"></el-input>
        </el-form-item>

        <el-form-item label="用户密码" prop="config.password" class="form-item">
          <el-input v-model.trim="form.config.password" auto-complete="off" maxlength="120" placeholder="请输入 Harbor 密码"></el-input>
        </el-form-item>
      </div>
      <el-form-item label="描述" prop="description" class="form-item">
        <el-input v-model="form.description" auto-complete="off" maxlength="80" placeholder="请输入描述" ></el-input>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button type="success" @click="doTestConnection" :loading="loading">测试连接</el-button>
      <el-button type="primary" @click="doSubmit" :loading="loading">{{$t('bm.other.confirm')}}</el-button>
    </div>
    </div>
  </el-drawer>
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
  config: {},
  description: '',
};

export default {
  mixins: [createTemplate, validate],
  data() {
    return {
      name: '',
      groupRoleList: [],
      settingTypeList: [
        {"name": "kubernetes"},
        {"name": "jenkins"},
        {"name": "harbor"}
      ],
      direction: 'rtl',
      // 是否属于编辑状态
      isEdit: false,
      dialogFormVisible: false,
      form: JSON.parse(JSON.stringify(formData)),
      title: '新增',
      isKubernetes: false,
      rules: {
        name: [
          { required: true, message: '请输入名称', trigger: 'blur' },
        ],
        type: [
          { required: true, message: '请选择集成服务的类型', trigger: 'blur' },
        ],
        'config.url': [
          { required: true, message: '请输入url信息', trigger: 'blur' },
        ],
        'config.user': [
          { required: true, message: '请输入用户名', trigger: 'blur' },
        ],
        'config.password': [
          { required: true, message: '请输入密码', trigger: 'blur' },
        ],
        'config.conf': [
          { required: true, message: '请输入kubernetes conf', trigger: 'blur' },
        ],
        'config.token': [
          { required: true, message: '请输入token信息', trigger: 'blur' },
        ],
        description: [
          { required: false, message: '描述信息不能为空', trigger: 'blur' },
        ],
      },
      rowId: '',
    };
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
  },
  created() {
  },
  methods: {
    handleClose(done) {
      this.$confirm('确认关闭？')
        .then(_ => {
          done();
        })
        .catch(_ => {});
    },
    doCreate(flag, item) {
      this.isEdit = flag;
      this.isKubernetes = false
      if (flag) {
        this.title = '编辑配置';
        this.form = {
          name: item.name || '',
          type: item.type || '',
          config: item.config || {},
          description: item.description || '',
        };
        if (item.type == "kubernetes") {
          this.isKubernetes = true
        }
        this.rowId = item.id;
      } else {
        this.title = '新增配置';
        this.form = {
          name: '',
          type: '',
          config: {},
          description: '',
        };
        this.rowId = '';
      }
      this.dialogFormVisible = true
      this.isEdit = flag;
    },
    doTestConnection() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          const cl = {
            name: this.form.name,
            config: this.form.config,
            type: this.form.type,
            description: this.form.description,
          };
          backend.VerifyIntegrateSetting(cl, (data) => {
            Message.success(data);
          });
        }
      });
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
            config: this.form.config,
            type: this.form.type,
            description: this.form.description,
          };
          if (this.isEdit) {
            backend.editIntegrateSetting(this.rowId, cl, () => {
              successCallBack();
            });
          } else {
            backend.AddIntegrateSetting(cl, () => {
              successCallBack();
            });
          }
        }
      });
    },
  },
};
</script>
