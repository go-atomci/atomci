<style>
.template-content .el-form-item__content {
  font-size: 12px;
  line-height: 17px;
  width: 100%;
}

.isView.createDialog .el-dialog__body .template-content {
  padding-right: 0;
}
</style>
<template>
  <el-dialog top='25vh' :title="$t('bm.add.configFile')" :visible.sync="dialogFormVisible" class="createDialog" :class="{'isView': isView}" :before-close="doCancelCreate">
    <el-form :model="form" ref="ruleForm" :rules="rules">
      <el-form-item :label="$t('bm.add.fileName')" prop='name' v-if='!isView'>
        <el-input v-model.trim="form.name" auto-complete="off" :disabled="isEdit"></el-input>
      </el-form-item>
      <el-form-item :label="!isView?$t('bm.add.fileContent'):''" class='template-content'>
        <codemirror  v-model.trim="form.text" :options="editorOption"></codemirror>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer" v-show="!isView">
      <el-button @click="doCancelCreate">{{$t('bm.other.cancel')}}</el-button>
      <el-button type="primary" @click="doSubmit" :loading="loading" v-if='!isView'>{{$t('bm.other.confirm')}}</el-button>
    </div>
  </el-dialog>
</template>
<script>
import { codemirror } from 'vue-codemirror';
import { mapGetters } from 'vuex';
import { Message, MessageBox } from 'element-ui';
import backend from '../../../api/backend';
import 'codemirror/lib/codemirror.css';
import 'codemirror/theme/rubyblue.css';
require('codemirror/addon/selection/active-line.js');

export default {
  components: {
    codemirror,
  },
  data() {
    return {
      editorOption: {
        tabSize: 4,
        styleActiveLine: true,
        lineNumbers: true,
        line: true,
        mode: 'text/x-yaml',
        lineWrapping: true,
        theme: 'rubyblue',
      },
      rules: {
        name: [{ required: true, message: this.$t('bm.add.inputFileName'), trigger: 'blur' }],
      },
      // 是否属于编辑状态
      isEdit: false,
      // 是否属于查看状态
      isView: false,
      dialogFormVisible: false,
      ConfigMapDetail: {},
      flag: '',
      form: {
        name: '',
        text: '',
      },
    };
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
  },
  created() {
    // 获得模板内容
    this.getConfigMapDetail();
  },
  methods: {
    getConfigMapDetail() {
      backend.getConfigMapDetail(
        this.$route.params.clusterName,
        this.$route.params.namespace,
        this.$route.params.configmapName,
        (data) => {
          this.ConfigMapDetail = data;
        }
      );
    },
    doUpdate(flag, name, text) {
      this.flag = flag;
      switch (flag) {
        // 添加
        case 'add':
          this.dialogFormVisible = true;
          this.isEdit = false;
          this.isView = false;
          this.form.name = '';
          this.form.text = '';
          break;
        // 更新
        case 'update':
          this.dialogFormVisible = true;
          this.isEdit = true;
          this.isView = false;
          this.form.name = name;
          this.form.text = text;
          break;
        // 删除
        case 'delete':
          this.dialogFormVisible = false;
          this.form.name = name;
          this.form.text = text;
          this.isEdit = true;
          MessageBox.confirm(this.$t('bm.add.sureDelete'), this.$t('bm.infrast.tips'), { type: 'warning' })
            .then(() => {
              this.doSubmit();
            })
            .catch(() => {});
          break;
        // 删除
        case 'view':
          this.dialogFormVisible = true;
          this.isView = true;
          this.form.name = name;
          this.form.text = text;
          break;
      }
    },
    doCancelCreate() {
      this.dialogFormVisible = false;
      this.isEdit = false;
      this.isView = false;
    },
    doSubmit() {
      const body = this.ConfigMapDetail;
      if (this.flag === 'delete') {
        delete body.data[this.form.name];
      } else if (body.data) {
        body.data[this.form.name] = this.form.text;
      } else {
        const map = {};
        map[this.form.name] = this.form.text;
        body.data = map;
      }
      const successCallBack = () => {
        this.$emit('getlist');
        Message.success(this.$t('bm.add.optionSuc'));
        this.dialogFormVisible = false;
      };
      backend.updateConfigMap(
        this.$route.params.clusterName,
        this.$route.params.namespace,
        this.$route.params.configmapName,
        body,
        () => {
          successCallBack();
        }
      );
    },
  },
};
</script>
