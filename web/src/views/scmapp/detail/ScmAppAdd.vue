<template>
  <div class="page-content">
    <div class="portlet-body projectMember">
      <template>
        <el-form ref="ruleForm" :model="form" :rules="rules">
          <el-form-item label="代码源" prop="integrate_repo_id">
            <el-select
              v-model="form.integrate_repo_id"
              clearable
              filterable
              placeholder="请选择代码源"
              style="width: 300px"
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
      <div class="member-btn pv10">
        <el-button type="default" @click="cancelApp">取消</el-button>
        <el-button type="primary" @click="addApp">添加</el-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.pv10 {
  padding-top: 10px;
  padding-bottom: 10px;
}
.pv30 {
  padding-top: 30px;
  padding-bottom: 30px;
}
.editBtn {
  float: right;
  padding-top: 3px;
  font-size: 12px;
  padding-right: 10px;
}
.member-btn {
  width: 650px;
  text-align: right;
}
.projectMember {
  padding: 30px;
}
.el-select {
  width: 550px;
}
.el-tabs .el-select {
  width: 310px;
}
.containerMember {
  width: 620px;
  padding: 20px;
  line-height: 40px;
}
.labelSize .el-input {
  width: 550px;
}
.el-tabs .labelSize .el-input {
  width: 500px;
}
.el-tabs .member-btn {
  width: 590px;
  text-align: right;
}
.sel-fl {
  width: 400px;
  float: left;
}
</style>
<style>
.projectMember .el-tabs__item {
  font-size: 18px;
}
.projectMember .el-form-item__error {
  left: 100px;
}
.projectMember .el-form-item__label {
  width: 100px;
}
.projectMember .el-tabs {
  border: 1px solid #ccc;
  padding: 10px;
  width: 650px;
}
.projectMember .el-tabs .el-form-item__label {
  width: 90px;
}
.el-form-item__content {
  float: left;
}
</style>
<script>
import { Message, MessageBox } from 'element-ui';
import backend from '@/api/backend';

export default {
  data() {
    return {
      integrateRepos: [],
      compileEnvs: [],
      proStep: {},
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
        type: 'app',
        language: 'Java',
        build_path: '/',
        dockerfile: 'Dockerfile',
      },
      getRepoLoading: true,
      rules: {
        type: [{ required: true, message: '请选择应用类型', trigger: 'change' }],
        integrate_repo_id: [{ required: true, message: '请选择代码源', trigger: 'change' }],
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
    handleClick(index) {},
    getIntegrateRepos() {
      // TODO: use 2 for tmp
      backend.getRepos((data) => {
        if (data) {
          this.integrateRepos = data;
        }
      });
    },
    getList(data) {
      data.map((item, index) => {
        if (item.base_url) {
          const cl = {
            addr: item.base_url,
            user: item.user,
          };
          // TODO: use 2 for tmp
          backend.getReposList(item.repo_id, 2, cl, (col) => {
            this.getTabs[index].proCol = col;
            this.getRepoLoading = false;
          });
        }
      });
    },
    firstStep(index) {
      this.getTabs[index].stepsNum = 2;
    },
    secondStep(index) {
      this.$refs['listRef' + index][0].validate((valid) => {
        if (valid) {
          if (this.getTabs[index].type == 'gitlab') {
            let url = new URL(this.getTabs[index].base_url);
            if (url.pathname !== '/') {
              let new_base_url = url.origin;
              MessageBox.confirm(
                '看起来您的地址不是仓库的主地址，需要主动为您修改为[' + new_base_url + ']吗？',
                this.$t('bm.infrast.tips'),
                {
                  confirmButtonText: '确定修改',
                  cancelButtonText: '保持原状',
                  type: 'warning',
                }
              )
                .then(() => {
                  this.getTabs[index].base_url = new_base_url;
                  this.getReposList(index);
                })
                .catch(() => {
                  this.getReposList(index);
                });
            } else {
              this.getReposList(index);
            }
          } else {
            this.getReposList(index);
          }
        }
      });
    },
    getReposList(index) {
      const cl = {
        base_url: this.getTabs[index].base_url,
        user: this.getTabs[index].user,
        token: this.getTabs[index].token,
      };
      // TODO: use 2 for tmp
      backend.getReposList(this.getTabs[index].repo_id, 2, cl, (data) => {
        if (data) {
          this.getTabs[index].proCol = data;
          this.getRepoLoading = false;
          this.getTabs[index].stepsNum = 3;
        }
      });
    },
    addApp() {
      const nums = parseInt(this.activeName);
      const arr = this.getTabs[nums];
      if (this.getTabs[nums].stepsNum < 3) {
        Message.error('请先同步代码源！');
        return;
      }
      this.$refs['ruleForm'].validate((valid) => {
        if (valid) {
          this.$refs['listRefs' + nums][0].validate((valid) => {
            if (valid) {
              const pathNum = arr.path;
              let cl = arr.proCol[pathNum];
              cl.language = this.form.language;
              cl.type = 'app';
              cl.build_path = this.form.build_path;
              cl.dockerfile = this.form.dockerfile || 'Dockerfile';
              cl.compile_env_id = this.form.compile_env_id || 0;
              if (this.form.name !== '') {
                cl.name = this.form.name;
              }
              // TODO: use 2 for tmp
              backend.addAppPro(2, cl, (data) => {
                Message.success('添加成功！');
                this.$router.push({ name: 'scmappIndex' });
              });
            }
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
          this.$router.push({
            name: 'projectApp',
            params: { projectID: this.$route.params.projectID },
          });
        })
        .catch(() => {});
    },
  },
};
</script>
