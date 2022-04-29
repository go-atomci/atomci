<style>
.arrange-body {
  padding-left: 20px;
  padding-right: 20px;
}

.arrange-row {
  padding-bottom: 10px;
}

.apparrange-layout {
  position: fixed;
  z-index: 9;
  height: 60px;
  bottom: 1px;
  width: 700px;
  background: #fff;
  box-shadow: 0 1px 5px 0 rgba(50, 50, 50, 0.5);
  line-height: 60px;
  overflow: hidden;
}
</style>

<template>
  <div>
    <el-drawer
      title="创建应用"
      :visible.sync="dialogFormVisible"
      class="arrangement"
      size="40%"
      :before-close="handleClose"
    >
      <div class="arrange-body">
      <template>
        <el-form ref="ruleForm" :model="form" :rules="rules">
          <el-form-item label="代码源" prop="repo_id">
            <el-select
              v-model="form.repo_id"
              clearable
              filterable
              placeholder="请选择代码源"
              style="width: 300px"
              @change="getReposList()"
            >
              <el-option
                v-for="(item, index) in integrateRepos"
                :key="index"
                :label="item.name"
                :value="item.id"
              ></el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="仓库地址" prop="path">
            <el-select
              v-model="form.path"
              filterable
              placeholder="请选择语言类型"
              style="width: 300px"
              @change="setScmAppName()"
            >
              <el-option
                v-for="(item, index) in scmProjects"
                :key="index"
                :label="item.full_name"
                :value="item.path"
              >
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="应用名" prop="build_path">
            <el-input
              v-model="form.name"
              placeholder="请输入应用名"
              style="width: 300px"
            ></el-input>
          </el-form-item>

          <el-form-item label="语言类型" prop="language">
            <el-select
              v-model="form.language"
              filterable
              placeholder="请选择语言类型"
              style="width: 300px"
            >
              <el-option
                v-for="(item, index) in languageList"
                :key="index"
                :label="item.description"
                :value="item.name"
              >
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="构建目录" prop="build_path">
            <el-input
              v-model="form.build_path"
              placeholder="请输入构建目录"
              style="width: 300px"
            ></el-input>
          </el-form-item>

          <el-form-item label="Dockerfile" prop="dockerfile">
            <el-input
              v-model="form.dockerfile"
              placeholder="请输入dockerfile,默认是Dockerfile"
              style="width: 300px"
            ></el-input>
          </el-form-item>

          <el-form-item label="编译环境" prop="compile_env_id">
            <el-select
              v-model="form.compile_env_id"
              clearable
              filterable
              placeholder="请选择编译环境"
              style="width: 300px"
            >
              <el-option
                v-for="(item, index) in compileEnvs"
                :key="index"
                :label="item.name"
                :value="item.id"
              ></el-option>
            </el-select>
          </el-form-item>
        </el-form>
      </template>
        <div class="apparrange-layout">
          <el-button
            type="primary"
            class="fb-ly-rbtn"
            icon="el-icon-edit"
            @click="addApp"
            >{{ '保存配置' }}</el-button
          >
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script>
import { Message, MessageBox } from 'element-ui';
import backend from '@/api/backend';

export default {
  data() {
    return {
      dialogFormVisible: false,
      integrateRepos: [],
      compileEnvs: [],
      scmProjects: [],
      languageList: [
        { description: 'Static', name: 'static' },
        { description: 'Java', name: 'Java' },
        { description: 'Node', name: 'Node' },
        { description: 'Go', name: 'go' },
        { description: 'Python', name: 'python' },
        { description: 'C#', name: 'C#' },
      ],
      form: {
        name: '',
        full_name: '',
        type: 'app',
        language: 'Java',
        build_path: '/',
        dockerfile: 'Dockerfile',
        repo_id: undefined,
      },
      getRepoLoading: true,
      rules: {
        type: [{ required: true, message: '请选择应用类型', trigger: 'change' }],
        repo_id: [{ required: true, message: '请选择代码源', trigger: 'change' }],
        path: [{ required: true, message: '请选择仓库地址', trigger: 'change' }],
        compile_env_id: [{ required: false, message: '请选择应用编译环境', trigger: 'change' }],
        language: [{ required: true, message: '请选择语言类型', trigger: 'change' }],
      },
      ruleProj: {
        proj: [{ required: true, message: '请选择代码库', trigger: 'blur' }],
      },
      activeName: '0',
      getTabs: [],
    };
  },
  components: {},
  created() {},
  mounted() {
    backend.getCompileEnvAll((data) => {
      if (data) {
        this.compileEnvs = data;
      }
    });
    this.getIntegrateRepos();
  },
  methods: {
    setScmAppName() {
      if (this.form.path == undefined) {
        return
      }
      for (let i = 0; i < this.scmProjects.length; i++) {
          if (this.scmProjects[i].path == this.form.path) {
            this.form.name = this.scmProjects[i].name
            this.form.full_name = this.scmProjects[i].full_name
            break
          }
      }
    },
    getIntegrateRepos() {
      backend.getIntegrateRepos((data) => {
        if (data) {
          this.integrateRepos = data;
        }
      });
    },
    getReposList() {
      backend.getReposList(this.form.repo_id, (data) => {
        if (data) {
          this.scmProjects = data;
          this.getRepoLoading = false;
        }
      });
    },
    doCreate() {
      this.form = {
        name: '',
        full_name: '',
        type: 'app',
        language: 'Java',
        build_path: '/',
        dockerfile: 'Dockerfile',
        repo_id: undefined,
      }
      this.dialogFormVisible = true
    },
    handleClose(done) {
      this.$confirm('确定返回我的应用?')
        .then((_) => {
          done();
        })
        .catch((_) => {});
    },
    addApp() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          const cl = this.form
          cl.type = 'app';
          cl.dockerfile = this.form.dockerfile || 'Dockerfile';
          cl.compile_env_id = this.form.compile_env_id || 0;
          if (this.form.name !== '') {
            cl.name = this.form.name;
          }
  
          backend.addScmAppPro(cl, (data) => {
            Message.success('添加成功！');
            this.$emit('getlist');
            this.dialogFormVisible = false
          });
        }
      });
    },
    editPath(index) {
      this.getTabs[index].stepsNum = 2;
    },
    cancelApp() {
      MessageBox.confirm('确定取消添加？', this.$t('bm.infrast.tips'), { type: 'warning' })
        .then(() => {
          this.$router.push({name: 'scmappIndex'});
        })
        .catch(() => {});
    },
  },
};
</script>