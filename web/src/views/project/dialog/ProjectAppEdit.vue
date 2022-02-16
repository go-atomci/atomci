<style>
.createDialog .el-dialog__body .el-form-item .el-form-item__content {
  display: flex;
  flex-direction: column;
}
</style>
<template>
  <el-dialog top='25vh' :title="title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog"  width='50%' :before-close="doCancelCreate">
    <el-form :model="form" ref="ruleForm" :rules="rules">
      <el-form-item :label="$t('bm.deployCenter.repositoryName')" prop="name">
        <el-input v-model="form.name" auto-complete="off" maxlength="20" placeholder="请输入仓库名称"></el-input>
      </el-form-item>
      <el-form-item label="语言类型" prop="language">
        <el-select v-model="form.language" placeholder="请选择语言类型" filterable>
          <el-option v-for="(item, index) in languageList" :key="index" :label="item.description" :value="item.name">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="$t('bm.deployCenter.path')" prop="path">
        <el-input v-model.trim="form.path" auto-complete="off" placeholder="请输入仓库路径"></el-input>
      </el-form-item>
      <el-form-item label="编译环境" prop="compile_env_id">
        <el-select v-model="form.compile_env_id" placeholder="请选择编译环境" clearable filterable>
          <el-option v-for="(item, index) in compileEnvs" :key="index" :label="item.name" :value="item.id">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="构建目录" prop="build_path">
        <el-input v-model="form.build_path" placeholder="请输入构建目录"></el-input>
      </el-form-item>
      <el-form-item label="Dockerfile" prop="dockerfile">
        <el-input v-model="form.dockerfile" placeholder="请输入Dockerfile,默认是根目录下的Dockerfile"></el-input>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="doCancelCreate">{{$t('bm.other.cancel')}}</el-button>
      <el-button type="primary" @click="doSubmit">{{$t('bm.other.confirm')}}</el-button>
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
  type: 'app',
  language: '',
  path: '',
  branch_name: '',
  branchList: []
};

export default {
  mixins: [createTemplate, validate],
  data() {
    return {
      name: '',
      groupRoleList: [],
      compileEnvs: [],
      languageList: [
        {'description': 'C','name': 'c'},
        {'description': 'Java','name': 'Java'},
        {'description': 'C#','name': 'C#'},
        {'description': 'go','name': 'go'},
        {'description': 'Node','name': 'Node'},
        {'description': '其他','name': 'other'},
      ],
      // 是否属于编辑状态
      isEdit: false,
      dialogFormVisible: false,
      form: JSON.parse(JSON.stringify(formData)),
      title: this.$t('bm.other.edit'),
      rules: {
        name: [
          { required: true, message: '请输入编号', trigger: 'blur' },
        ],
        type: [
          { required: false, message: '请选择模块类型', trigger: 'blur' },
        ],
        language: [
          { required: true, message: '请选择语言类型', trigger: 'blur' },
        ],
      },
      rowId: '',
      defaultActiveIndex: '/git',
    };
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
  },
  mounted() {
    backend.getCompileEnvAll((data) => {
      if(data){
        this.compileEnvs = data;
      }
    });
  },
  mounted() {

  },
  methods: {
    handleSelect(index) {
      
    },
    doCreate(flag, item) {
      this.isEdit = flag;
      this.form = {
        name: item.name || '',
        // TODO: comment app type tmp
        language: item.language || '',
        path: item.path || '',
        compile_env_id: item.compile_env_id || 0,
        branch_name: item.branch_name || '',
        build_path: item.build_path || '/',
        dockerfile: item.dockerfile || 'Dockerfile',
        branchList: item.branch_history_list || []
      };
      this.rowId = item.id;
      this.dialogFormVisible = true;
      this.isEdit = flag;
    },
    doSubmit() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          const successCallBack = () => {
            this.$emit('getList');
            Message.success(this.$t('bm.add.optionSuc'));
            this.dialogFormVisible = false;
          };
          const cl = {
            name: this.form.name,
            compile_env_id: this.form.compile_env_id || 0, 
            language: this.form.language,
            path: this.form.path,
            branch_name: this.form.branch_name,
            build_path: this.form.build_path,
            dockerfile:  this.form.dockerfile || 'Dockerfile',
          };
          backend.updateAppInfo(this.$route.params.projectID, this.rowId, cl, (data) => {
            successCallBack();
          });

        }
      });
    },
  },
};
</script>
