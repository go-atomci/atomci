<style scoped>
  .pv10 {
    padding-top: 10px;
    padding-bottom: 10px;
    height: 40px;
    position: relative;
  }
  .radioStyle {
    position: absolute;
    left: 100px;
    top: 12px;
  }
  .radioStyles {
    position: absolute;
    left: 0px;
    top: 12px;
  }
  .el-form-item {
    margin: 20px 0;
  }
</style>
<template>
  <el-dialog top='15vh' v-if="dialogFormVisible" :close-on-click-modal="false" width='65%' :title="setname" :visible.sync="dialogFormVisible" class="createDialog">
    <el-table style="margin-top:2%" border :data="detailData" @select-all="handleSelectAll" @select='handleSelect'>
      <span slot="empty">
        {{noDataTxt}}
      </span>
      <el-table-column type="selection" min-width="7%" :show-overflow-tooltip=true></el-table-column>
      <el-table-column prop="name" :label="$t('bm.deployCenter.name')" sortable min-width="12%" :show-overflow-tooltip=true />
      <el-table-column prop="language" :label="$t('bm.deployCenter.language')" sortable min-width="10%" :show-overflow-tooltip=true />
      <el-table-column prop="branch_name" :label="$t('bm.deployCenter.releaseBran')" min-width="20%" :show-overflow-tooltip=true>
        <template slot-scope="scope">
          <el-select v-model.trim="scope.row.branch_name" filterable :placeholder="$t('bm.add.selectSubmitBra')">
            <el-option v-for="(item, index) in scope.row.branch_history_list" :key="index" :label="item" :value="item">
            </el-option>
          </el-select>
        </template>
      </el-table-column>
      <el-table-column :label="$t('bm.add.version')"  sortable min-width="10%" :show-overflow-tooltip=true>
        <template slot-scope="scope">
          <el-input v-model="scope.row.image_version"></el-input>
        </template>
      </el-table-column>
    </el-table>
    <div slot="footer" class="dialog-footer">
      <el-button @click="handleClose" style="margin-top:20px">{{$t('bm.other.cancel')}}</el-button>
      <el-button type="primary" @click="doSubmit()">{{$t('bm.other.confirm')}}</el-button>
    </div>
  </el-dialog>
</template>
<script>
  import { mapGetters } from 'vuex';
  import { Message } from 'element-ui';
  import backend from '../../../api/backend';
  import Refresh from '../../../components/utils/Refresh';

  export default {
    props: ['listData', 'pubType', 'cpData'],
    data() {
      return {
        setname: '添加应用',
        detailData: [],
        dialogFormVisible: false,
        selectList: [],
        version: [],
        tableList: [],
        dataList: [],
        projectId: '',
        versionId: '',
        noDataTxt: '当前暂无可添加应用，请先去代码仓库中添加后重试！',
      };
    },
    components: {
      Refresh,
    },
    created() {

    },
    methods: {
      handleSelectAll(val) {
        this.selectList = val;
      },
      handleSelect(val) {
        this.selectList = val;
      },
      doSubmit() {
        const apps = [];
        if(this.selectList.length === 0) {
          Message.error('请先选择应用！');
          return;
        }
        for (const a of this.selectList) {
          if(!a.image_version) {
            Message.error('请输入1-16位的版本号');
            return;
          }
          const at = {
            "app_id": a.id,
            "branch_name": a.branch_name,
            "image_version": a.image_version
          };
          apps.push(at);
        }
        const params = {
          'apps': apps
        };
        const that = this;
        backend.addVersionApp(this.projectId, this.versionId, params, (data) => {
          Message.success(this.$t('bm.add.optionSuc'));
          that.$emit('getlist');
          this.handleClose();
        }, () => {});
      },
      doShow(projectId, versionId, version) {
        this.projectId = projectId;
        this.versionId = versionId;
        this.dialogFormVisible = true;
        backend.versionApp(projectId, versionId, (data) => {
          if(data) {
            this.detailData = data.map((i) => {
              i.image_version = version;
              return i;
            });
          }
        });
      },
      handleClose() {
        this.dialogFormVisible = false;
      },
    },
  };
</script>
