<template>
  <div>
    <el-row style="margin-bottom:5px;font-size:14px; margin-top: 10px;">
      <el-col class="label-100">
        {{$t('bm.serviceM.podInstance')}}
      </el-col>
      <el-col :span="8">
        <el-select v-model="podName" :placeholder="$t('bm.add.select')" @change="podChange" style="margin-right:10px;">
          <el-option v-for="(item, index) in podList" :key="index" :label="item.name" :value="item.name">
          </el-option>
        </el-select>
      </el-col>
      <el-col class="label-100">
        {{$t('bm.add.container')}}
      </el-col>
      <el-col :span="8">
        <el-select v-model="containerName" :placeholder="$t('bm.add.select')">
          <el-option v-for="(item, index) in containerList" :key="index" :label="item.name" :value="item.name">
          </el-option>
        </el-select>
      </el-col>
      <el-col :span="2">
        <el-button class="ml10" style="vertical-align:-1px;" :plain="true" type="primary" @click="doSearch" icon="search">{{$t('bm.serviceM.search')}}</el-button>
      </el-col>
    </el-row>
    <el-row>
      <el-col :span="24">
        <codemirror v-model="podContainerLog" :options="{ styleActiveLine: true, mode: 'text/plain',  tabSize: 10, autoCloseTags: false, theme: 'ambiance', lineNumbers: false }"></codemirror>
      </el-col>
    </el-row>
  </div>
</template>
<style scoped>
.label-100 {
  line-height: 35px;
  width: 100px;
  padding-left: 10px;
}
</style>
<script>
import { codemirror } from 'vue-codemirror';
import backend from '@/api/backend';

// import base style
import 'codemirror/lib/codemirror.css'

// theme css
import 'codemirror/theme/ambiance.css'

// require active-line.js
require('codemirror/addon/selection/active-line.js');

export default {
  props: ['appName'],
  data() {
    return {
      detailInfo: null,
      containerName: '',
      containerList: [],
      podName: '',
      podList: [],
      podContainerLog: '',
      logConfigUrl: '',
    };
  },
  components: {
    codemirror,
  },
  created() {
  },
  methods: {
    doSearch() {
      backend.getServiceLog(
        this.$route.params.clusterName,
        this.$route.params.namespace,
        this.$props.appName,
        this.podName,
        this.containerName,
        (data) => {
          this.podContainerLog = data;
        }
      );
    },
    // 资源空间select 下拉列表改变时
    podChange(key) {
      const cList = [];
      for (const a of this.detailInfo.pods) {
        if (a.name === key) {
          for (const b of a.containers) {
            cList.push({ name: b.name });
          }
          break;
        }
      }
      this.containerList = cList;
      if (cList.length) {
        this.containerName = cList[0].name;
      }
    },
    doSelectPodName(data) {
      const pList = [];
      for (const a of data.pods) {
        pList.push({ name: a.name });
      }
      this.podList = pList;
      this.detailInfo = data;
      if (pList.length) {
        this.podName = pList[0].name;
        this.podChange(this.podName);
      }
    },
  },
};
</script>
