<template>
  <div class="page-content">
    <div class="portlet-body projectMember">
      <template>
        <el-form :model="form" ref="ruleForm" :rules="rules">
          <!-- <el-form-item label="类型" prop="type" >
            <el-select v-model="form.type" placeholder="请选择应用类型" filterable style="width: 300px">
              <el-option v-for="(item, index) in typeList" :key="index" :label="item.description" :value="item.name">
              </el-option>
            </el-select>
          </el-form-item> -->
          <el-form-item label="应用名" prop="build_path">
            <el-input v-model="form.name" placeholder="请输入应用名" style="width: 300px"></el-input>
          </el-form-item>

          <el-form-item label="语言类型" prop="language">
            <el-select v-model="form.language" placeholder="请选择语言类型" filterable style="width: 300px">
              <el-option v-for="(item, index) in languageList" :key="index" :label="item.description" :value="item.name">
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="构建目录" prop="build_path">
            <el-input v-model="form.build_path" placeholder="请输入构建目录" style="width: 300px"></el-input>
          </el-form-item>

          <el-form-item label="Dockerfile" prop="dockerfile">
            <el-input v-model="form.dockerfile" placeholder="请输入dockerfile,默认是Dockerfile" style="width: 300px"></el-input>
          </el-form-item>

          <el-form-item label="编译环境" prop="compile_env_id">
            <el-select v-model="form.compile_env_id" placeholder="请选择编译环境" clearable filterable style="width: 300px">
              <el-option v-for="(item, index) in compileEnvs" :key="index" :label="item.name" :value="item.id">
            </el-option>
        </el-select>
          </el-form-item>

        </el-form>
      </template>
      <template>
        <el-tabs v-model="activeName" @tab-click="handleClick">
          <template v-if="getTabs">
            <template v-for="(item,index) in getTabs">
              <el-tab-pane :label="item.type"  :key="`${index}`"   :name="`${index}`" ref="tabPanel">
                <div class="pv30">
                  <div v-show="item.stepsNum === 1" class="containerMember">
                    {{item.type}}提供丰富的代码管理功能，是目前主流的企业内部代码仓库解决方案
                    <el-button type="primary" @click="firstStep(index)">同步代码源</el-button>
                  </div>
                  <div v-show="item.stepsNum === 2">
                    <el-form :model="item" :ref="'listRef'+index">
                      <div class="labelSize">
                        <el-form-item label="地址" prop="base_url" :rules="[{ required: true, message: '请输入地址', trigger: 'blur' }]">
                          <el-input v-model.trim="item.base_url" auto-complete="off" placeholder="请输入地址, 如：https://gitlab.com"></el-input>
                        </el-form-item>
                      </div>
                      <div class="labelSize">
                        <el-form-item label="用户名" prop="user">
                          <el-input v-model.trim="item.user" auto-complete="off" maxlength="128" placeholder="请输入用户名"></el-input>
                        </el-form-item>
                      </div>
                      <div class="labelSize">
                        <el-form-item label="Token" prop="token">
                          <el-input v-model.trim="item.token" auto-complete="off" maxlength="128" placeholder="请输入Token"></el-input>
                        </el-form-item>
                      </div>
                    </el-form>
                    <div class="member-btn">
                      <el-button type="primary" @click="secondStep(index)">同步代码源</el-button>
                    </div>
                  </div>
                  <div class="" v-show="item.stepsNum === 3">
                    <div class="editBtn">要改变当前绑定的代码库？<el-button type="text" size="small" @click="editPath(index)">解除绑定</el-button></div>
                    <div class="sel-fl">
                      <el-form :model="item" :ref="'listRefs'+index">
                        <el-form-item label="仓库名称" prop="path" :rules="[{ required: true, message: '请选择代码库', trigger: 'change' }]">
                          <el-select v-model="item.path" placeholder="请选择代码库" filterable :loading="getRepoLoading" :key="`sel${index}`">
                            <el-option v-for="(term, indexs) in item.proCol" :key="indexs" :label="term.full_name" :value="indexs">
                            </el-option>
                          </el-select>
                        </el-form-item>
                      </el-form>
                    </div>
                  </div>
                </div>
              </el-tab-pane>
            </template>
          </template>
        </el-tabs>
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
    width:400px;
    float:left;
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
        compileEnvs: [],
        proStep: {},
        typeList: [
          {'description': '应用','name': 'app'},
          {'description': '依赖','name': 'module'},
        ],
        languageList: [
          {'description': 'Static','name': 'static'},
          {'description': 'Java','name': 'Java'},
          {'description': 'Node','name': 'Node'},
          {'description': 'Go','name': 'go'},
          {'description': 'Python','name': 'python'},
          {'description': 'C#','name': 'C#'},
        ],
        form: {
          name: '',
          type: 'app',
          language: 'Java',
          build_path: '/',
          dockerfile: 'Dockerfile'
        },
        getRepoLoading: true,
        rules: {
          type: [
            { required: true, message: '请选择应用类型', trigger: 'change' },
          ],
          compile_env_id: [
            { required: false, message: '请选择应用编译环境', trigger: 'change' },
          ],
          language: [
            { required: true, message: '请选择语言类型', trigger: 'change' },
          ],
        },
        ruleProj: {
          proj: [
            { required: true, message: '请选择代码库', trigger: 'blur' },
          ],
        },
        activeName: '0',
        getTabs: [],
      };
    },
    components: {
    },
    created() {},
    activated() {
      backend.getCompileEnvAll((data) => {
      if(data){
        this.compileEnvs = data;
      }
      });
    },
    mounted() {
      this.getRepos();
    },
    methods: {
      handleClick(index) {},
      getRepos() {
        backend.getRepos(this.$route.params.projectId, (data) => {
          if(data) {
            data.map((item, index) => {
              if(item.user && item.base_url) {
                item.stepsNum = 3;
              } else if(item.user && !item.base_url) {
                item.stepsNum = 2;
              } else {
                item.stepsNum = 1;
              }
              item.proCol = [];
            });
            this.getTabs = data;
            this.activeName = '0';
            this.$nextTick(() => {
              this.getList(data);
            });
          }
        });
      },
      getList(data) {
        data.map((item, index) => {
          if(item.base_url) {
            const cl = {
              'addr': item.base_url,
              'user': item.user,
            };
            backend.getReposList(item.repo_id, this.$route.params.projectId, cl, (col) => {
              this.getTabs[index].proCol = col;
              this.getRepoLoading = false
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
            const cl = {
              "base_url": this.getTabs[index].base_url,
              "user": this.getTabs[index].user,
              "token": this.getTabs[index].token
            };
            backend.getReposList(this.getTabs[index].repo_id, this.$route.params.projectId, cl, (data) => {
              if(data) {
                this.getTabs[index].proCol = data;
                this.getRepoLoading = false
                this.getTabs[index].stepsNum = 3;
              }
            });
          }
        });
      },
      addApp() {
        const nums = parseInt(this.activeName);
        const arr = this.getTabs[nums];
        if(this.getTabs[nums].stepsNum < 3) {
          Message.error('请先同步代码源！');
          return;
        }
        this.$refs['ruleForm'].validate((valid) => {
          if(valid) {
            this.$refs['listRefs' + nums][0].validate((valid) => {
              if(valid) {
                const pathNum = arr.path;
                let cl = arr.proCol[pathNum];
                cl.language = this.form.language;
                cl.type = 'app';
                cl.build_path = this.form.build_path;
                cl.dockerfile = this.form.dockerfile || 'Dockerfile';
                cl.compile_env_id = this.form.compile_env_id || 0;
                if (this.form.name !== '') {
                  cl.name = this.form.name
                }
                backend.addAppPro(this.$route.params.projectId, cl, (data) => {
                  Message.success('添加成功！');
                  this.$router.push({
                    name: 'projectApp', params: {projectId: this.$route.params.projectId}
                  });
                });
              }
            });
          }
        })
      },
      editPath(index) {
        this.getTabs[index].stepsNum = 2;
      },
      cancelApp() {
        MessageBox.confirm('确定取消添加？', this.$t('bm.infrast.tips'), { type: 'warning' })
        .then(() => {
          this.$router.push({
            name: 'projectApp', params: {projectId: this.$route.params.projectId}
          });
        })
        .catch(() => {});
      }
    }
  }
</script>
