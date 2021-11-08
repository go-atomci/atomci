<style>
.createDialog .el-dialog__body .el-form-item .el-form-item__content {
  display: flex;
  flex-direction: column;
}
</style>
<template>
  <el-dialog top='25vh' :title="title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog"  width='50%' :before-close="doCancelCreate">
    <el-form :model="form" ref="ruleForm" :rules="rules">
      <el-form-item label="环境名称" prop="name">
        <el-input v-model.trim="form.name" auto-complete="off" maxlength="60" placeholder="请输入环境名称"></el-input>
      </el-form-item>
      <el-form-item label="k8s集群" prop="cluster">
        <el-select v-model="form.cluster" placeholder="请选择k8s集群" filterable>
          <el-option v-for="(item, index) in clusterList" :key="index" :label="item.name" :value="item.id">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="命名空间" prop="namespace">
        <el-select v-model="form.type" placeholder="请选择命名空间" filterable>
          <el-option v-for="(item, index) in namespaceList" :key="index" :label="item" :value="item">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="ENV标识" prop="arrange_env">
        <el-input v-model.trim="form.arrange_env" auto-complete="off" maxlength="10" placeholder="请输入ENV标识"></el-input>
      </el-form-item>
      <el-form-item label="镜像仓库" prop="harbor">
        <el-select v-model="form.harbor" placeholder="请选择镜像仓库" filterable>
          <el-option v-for="(item, index) in harborList" :key="index" :label="item.name" :value="item.id">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="$t('bm.add.buildCluster')" prop="ci_server">
        <el-select v-model="form.ci_server" :placeholder="$t('bm.add.selectBuildCluster')" filterable>
          <el-option v-for="(item, index) in jenkinsList" :key="index" :label="item.name" :value="item.id">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="环境描述" prop="description">
        <el-input v-model.trim="form.description" auto-complete="off" maxlength="50" placeholder="请输入描述" ></el-input>
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
      clusterList: [],
      jenkinsList: [],
      harborList: [],
      namespaceList: [
        'default'
      ],
      // 是否属于编辑状态
      isEdit: false,
      dialogFormVisible: false,
      form: JSON.parse(JSON.stringify(formData)),
      title: '新增',
      rules: {
        name: [
          { required: true, message: '请输入名称', trigger: 'blur' },
        ],
        cluster: [
          { required: true, message: '请选择集群', trigger: 'blur' },
        ],
        description: [
          { required: false, message: '描述信息不能为空', trigger: 'blur' },
        ],
        arrange_env: [
          { required: true, message: '请输入ENV标识', trigger: 'blur' },
        ],
        harbor: [
          { required: true, message: '请输入镜像仓库', trigger: 'blur' },
        ],
        ci_server: [
          { required: true, message: '请输入构建主机', trigger: 'blur' },
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
  created() {
    backend.getAllIntegrateSettings((data) => {
      if (data) {
        data.forEach(element => {
          switch (element.type) {
          case 'kubernetes':
            this.clusterList.push(element)
            break;
          case 'jenkins':
            this.jenkinsList.push(element)
            break;
          case 'harbor':
            this.harborList.push(element)
            break;
          default:
            console.log("this type not support", element.type)
          }
        });
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
          cluster: item.cluster || undefined,
          ci_server: item.ci_server || undefined,
          namespace: item.namespace || 'default',
          harbor: item.harbor || undefined,
          description: item.description || '',
          arrange_env: item.arrange_env || '',
        };
        this.rowId = item.id;
      } else {
        this.title = '新增';
        let defaultJenkins, defaultHarbor
        if (this.jenkinsList.length > 0){
          defaultJenkins = this.jenkinsList[0].id
        }
        if (this.harborList.length > 0){
          defaultHarbor = this.harborList[0].id
        }
        this.form = {
          name: '',
          cluster: undefined,
          ci_server: defaultJenkins,
          harbor: defaultHarbor,
          namespace: 'default',
          description: '',
          arrange_env: '',
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
            cluster: this.form.cluster,
            namespace: this.form.namespace || 'default',
            ci_server: this.form.ci_server,
            harbor: this.form.harbor,
            description: this.form.description,
            arrange_env: this.form.arrange_env
          };
          if (this.isEdit) {
            backend.editProjectEnv(this.projectID, this.rowId, cl, () => {
              successCallBack();
            });
          } else {
            backend.AddProjectEnv(this.projectID, cl, () => {
              successCallBack();
            });
          }
        }
      });
    },
  },
};
</script>
