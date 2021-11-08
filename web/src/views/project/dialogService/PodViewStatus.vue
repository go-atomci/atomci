<template>
  <el-dialog top='25vh' :title="$t('bm.add.viewStatus')" :close-on-click-modal="false" :visible.sync="dialogFormVisible" width='40%' :before-close="doCancelCreate">
    <codemirror v-model="spec" :options="editorOption"></codemirror>
  </el-dialog>
</template>
<script>
import { mapGetters } from 'vuex';
import { Message } from 'element-ui';
import backend from '../../../api/backend';
import 'codemirror/lib/codemirror.css';
import 'codemirror/theme/rubyblue.css';
import base64 from '../../../common/base64';

const yaml = require('js-yaml');
// require active-line.js
require('codemirror/addon/selection/active-line.js');

export default {
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
      dataList: [],
      spec: '',
      dialogFormVisible: false,
    };
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
  },
  methods: {
    doCreate(cluster, namespace, appname, podname) {
      // this.form = Object.assign({}, flag ? item : formData);
      this.dialogFormVisible = true;
      this.getList(cluster, namespace, appname, podname);
    },
    getList(cluster, namespace, appname, podname) {
      backend.getPodStatusViews(cluster, namespace, appname, podname, (data) => {
        this.spec = yaml.safeDump(data);
      });
    },
    doCancelCreate() {
      this.dialogFormVisible = false;
    },
  },
};
</script>
