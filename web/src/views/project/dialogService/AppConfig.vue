<template>
  <div>
    <el-row class='app-panel'>
      <el-col :span="24">
        <codemirror v-model.trim="appSpec" :options="editorOption"></codemirror>
      </el-col>
    </el-row>
    <el-row class='app-panel'>
      <el-col :span="24">
        <el-button :plain="true" type="primary" @click="doUpdate" icon="el-icon-edit">{{$t('bm.authorManage.update')}}</el-button>
        <el-button :plain="true" type="primary" @click="doReset">
          <i class='icon-energy' />{{$t('bm.add.reset')}}</el-button>
      </el-col>
    </el-row>
  </div>
</template>
<script>
import { codemirror } from 'vue-codemirror';
import 'codemirror/lib/codemirror.css';
import 'codemirror/theme/rubyblue.css';
import backend from '../../../api/backend';

// let yaml = require('js-yaml')
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
        mode: 'x-yaml',
        lineWrapping: true,
        theme: 'rubyblue',
      },
      originData: null,
      appSpec: '',
    };
  },
  components: {
    codemirror,
  },
  methods: {
    render(data) {
      this.originData = data;
      this.appSpec = data;
    },
    // 配置更新
    doUpdate() {
      backend.updateServiceInspect(
        this.$route.params.clusterName,
        this.$route.params.namespace,
        this.$route.params.appName,
        { template: this.appSpec },
        (data) => {
          this.$notify({
            title: this.$t('bm.add.success'),
            message: this.$t('bm.add.appConfigSuc'),
            type: 'success',
          });
          this.$emit('appCallBack', data);
        }
      );
    },
    // 配置重置
    doReset() {
      this.appSpec = this.originData;
    },
  },
};
</script>
