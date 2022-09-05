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
        <el-form ref="ruleForm" :model="form" :rules="rules" label-width="100px" >
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

          <el-form-item label="仓库路径" prop="full_name">
            <el-autocomplete
              v-model="form.full_name"
              style="width: 300px"
              placeholder="请选择/输入路径"
              :fetch-suggestions="querySearch"
              @blur="setScmAppName"
              @select="setSelectPath">
              <i
                class="el-icon-edit el-input__icon"
                slot="suffix">
              </i>
              <template slot-scope="{ item }">
                <div class="name">{{ item.full_name }}</div>
              </template>
            </el-autocomplete>
          </el-form-item>

          <el-form-item label="应用名" prop="name">
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
          <el-button type="success" @click="doTestConnection" :loading="loading">测试连接</el-button>
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
      loading: false,
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
        name: [{ required: true, message: '请输入应用名', trigger: 'blur' }],
        type: [{ required: true, message: '请选择应用类型', trigger: 'change' }],
        repo_id: [{ required: true, message: '请选择代码源', trigger: 'change' }],
        path: [{ required: true, message: '请选择仓库地址', trigger: 'change' }],
        full_name: [{ required: true, message: '请选择/输入仓库路径', trigger: 'blur' }],
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
    querySearch(queryString, cb) {
        var scmProjects = this.scmProjects;
        var results = queryString ? scmProjects.filter(this.createFilter(queryString)) : scmProjects;
        // 调用 callback 返回建议列表的数据
        cb(results);
      },
    createFilter(queryString) {
      return (scmProjects) => {
        return (scmProjects.full_name.toLowerCase().indexOf(queryString.toLowerCase()) === 0);
      };
    },
    setSelectPath(item){
      this.form.full_name=item.full_name;
      for (let i = 0; i < this.scmProjects.length; i++) {
          if (this.scmProjects[i].full_name == this.form.full_name) {
            this.form.name = this.scmProjects[i].name;
            this.form.path = this.scmProjects[i].path;
            break
          }
      }
    },
    setScmAppName() {
      this.setFullName();
      if (this.form.full_name &&  !this.form.name) {
        let names = this.form.full_name.split('/');
        this.form.name = names[names.length-1].split('.')[0]
      }
      this.form.path = "";
      for (let i = 0; i < this.scmProjects.length; i++) {
        if (this.scmProjects[i].full_name == this.form.full_name) {
          this.form.path = this.scmProjects[i].path;
          return
        }
      }
    },
    setFullName(){
      if (!this.form.full_name){
        return;
      }
      let fullName=this.form.full_name;
      // 优化路径，路径只需要[域名]和[.git]之间部分就行
      // 若是完整路径，如 http://reap.com/xx.git,则去除域名
      if(fullName.toLowerCase().indexOf('http://') > -1 || fullName.toLowerCase().indexOf('https://') > -1){
          let index = this.form.full_name.indexOf('/',8);
          fullName = this.form.full_name.substring(index + 1);
      }
      // 若路径包含[.git],也要去掉
      if(fullName.length > 3 && fullName.substring(fullName.length - 4).toLowerCase() == '.git'){
          this.form.full_name = fullName.substring(0,fullName.length - 4);
          return
      }
      this.form.full_name = fullName;
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
    doTestConnection() {
      this.setFormPath();
      const cl = {
          full_name: this.form.full_name,
          path: this.form.path,
          repo_id: this.form.repo_id,
        };
      if ( !cl.repo_id ){
        Message.error("请先选择一个代码源");
        return;
      }
      if ( !cl.full_name || !cl.path){
        Message.error("请先选择或输入仓库路径");
        return;
      }
      backend.verifyAppConnetion(cl, () => {
        Message.success("源码仓库地址连接成功");
      });
    },
    setFormPath(){
        // 完整的路径地址不存在的话，则组装一个出来
        if(this.form.path==""){
          for (let i = 0; i < this.integrateRepos.length; i++) {
            if (this.integrateRepos[i].id == this.form.repo_id) {
              let domain= this.integrateRepos[i].config.url;
              if(domain && domain.length>0 && domain.substring(domain.length-1)=='/'){
                this.form.path = domain + this.form.full_name + '.git';
              }else {
                this.form.path = domain + '/' + this.form.full_name + '.git';
              }
            }
          }
        }
    },
    addApp() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          this.setFormPath();
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