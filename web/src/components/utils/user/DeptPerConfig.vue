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
  <el-dialog top='25vh' :title="$t('bm.add.perPolicy')" :visible.sync="dialogFormVisible" :before-close="doCancelCreate">
    <el-form :model="form" ref="ruleForm">
      <el-form-item class='template-content'>
        <codemirror  v-model="form.text" :options="editorOption"></codemirror>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>
<script>
import { codemirror } from 'vue-codemirror';
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
        styleActiveLine: false,
        lineNumbers: true,
        line: true,
        mode: 'x-yaml',
        lineWrapping: true,
        theme: 'rubyblue',
        readOnly: true,
      },
      dialogFormVisible: false,
      form: {
        name: '',
        text: '',
      },
    };
  },

  methods: {
    getConfigMapDetail(data) {
      this.dialogFormVisible = true;
      this.form.text = data;
    },
    doCancelCreate() {
      this.dialogFormVisible = false;
      this.isEdit = false;
      this.isView = false;
    },
  },
};
</script>
