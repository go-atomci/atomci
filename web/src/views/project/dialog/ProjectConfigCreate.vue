<style>
@import "../../../style/dialog.css";
.createDialog.hostCreate .lbsList.el-col-20 .el-input__inner {
  vertical-align: -1px;
}
</style>
<template>
  <el-dialog top='25vh' :title="$t('bm.add.createConfigGroup')" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog hostCreate" width='50%' :before-close="doCancelCreate">
    <el-form :model="form" :rules="rules" ref="ruleForm">
      <el-form-item :label="$t('bm.serviceM.resourceSpace')" prop='namespace' v-if='!isView'>
        <el-select filterable v-model.trim="form.namespace" :placeholder="$t('bm.infrast.sResourceSpace')" :disabled="isEdit">
          <el-option v-for="(item, index) in namesapceList" :key="index" :label="item.name" :value="item.name">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="$t('bm.serviceM.configG')" prop='configmap'>
        <el-input v-model.trim="form.configmap" :placeholder="$t('bm.add.inputConfigGroup')" auto-complete="off" :disabled="isEdit"></el-input>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="doCancelCreate">{{$t('bm.other.cancel')}}</el-button>
      <el-button type="primary" @click="doSubmit" :loading="loading">{{$t('bm.other.confirm')}}</el-button>
    </div>
  </el-dialog>
</template>
<script>
import { Message } from 'element-ui';
import { mapGetters } from 'vuex';
import backend from '../../../api/backend';
import createTemplate from '../../../common/createTemplate';

const formData = {
  namespace: '',
  configmap: '',
};

export default {
  mixins: [createTemplate],
  props: ['user'],
  data() {
    return {
      // 是否属于编辑状态
      isEdit: false,
      dialogFormVisible: false,
      form: JSON.parse(JSON.stringify(formData)),
      // 是否属于查看状态
      isView: false,
      namesapceList: [],
      cluster: '',
      rules: {
        namespace: [
          { required: true, message: this.$t('bm.infrast.sResourceSpace'), trigger: 'blur' },
        ],
        configmap: [
          { required: true, message: this.$t('bm.add.inputConfigGroup'), trigger: 'blur' },
        ],
      },
    };
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
  },
  created() {
    this.$store.dispatch('setNeedLoading', false);
  },
  methods: {
    doCreate(flag, cluster, item) {
      this.cluster = cluster;
      this.dialogFormVisible = true;
      this.isEdit = flag;
      if (flag) {
        this.form.configmap = item.configmap;
      } else {
        this.form = Object.assign({}, formData);
      }
      backend.getNamespaceList(cluster, (data) => {
        this.namesapceList = data.map((it) => {
          return {
            name: it.name,
            desc: it.desc,
          };
        });
        this.$store.dispatch('setNeedLoading', true);
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
          backend.createConfigMap(this.cluster, this.form.namespace, JSON.stringify({
            metadata: {
              name: this.form.configmap,
            },
          }), (data) => {
            successCallBack();
          });
          return false;
        }
      });
    },
  },
};
</script>
