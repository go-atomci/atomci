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
      title="应用编排"
      :visible.sync="dialogFormVisible"
      class="arrangement"
      size="40%"
      :before-close="handleClose"
    >
      <div class="arrange-body">
        <template>
          <el-tabs v-model="editableTabsValue" @tab-click="handleClick">
            <el-tab-pane v-for="item in envList" :key="item.id" :label="item.name" :name="item.id">
              <div style="overflow: scroll-y; height: 600px">
                <el-form :model="form" ref="ruleForm">
                  <el-row class="arrange-row">
                    <el-col :span="3"> 编排文件 </el-col>
                    <el-col :span="18">
                      <el-input
                        type="textarea"
                        :rows="15"
                        v-model="form.config"
                        placeholder="请输入yaml"
                        @blur="parseYamlImages()"
                      ></el-input>
                    </el-col>
                  </el-row>
                  <el-row class="arrange-row">
                    <el-col :span="4"> 检测到的镜像 </el-col>
                  </el-row>
                  <el-row class="arrange-row">
                    <el-table :data="form.image_mapings" style="width: 100%">
                      <el-table-column prop="name" label="名称" :show-overflow-tooltip="true" ></el-table-column>
                      <el-table-column prop="image" label="镜像信息" :show-overflow-tooltip="true" width="120"></el-table-column>
                      <el-table-column prop="project_app_id" label="关联应用">
                        <template slot-scope="scope">
                          <el-select
                            v-model="scope.row.project_app_id"
                            placeholder="请选择关联应用"
                            clearable
                            filterable
                            @change="linkApp(scope.row)"
                          >
                            <el-option
                              v-for="item in appList"
                              :key="item.id"
                              :label="item.name"
                              :value="item.id"
                            ></el-option>
                          </el-select>
                        </template>
                      </el-table-column>
                      <el-table-column prop="image_tag_type" label="镜像版本规则">
                        <template slot-scope="scope">
                          <el-select
                            v-model="scope.row.image_tag_type"
                            placeholder="镜像版本规则"
                            clearable
                            filterable
                            @change="linkAppImageTagType(scope.row)"
                          >
                            <el-option
                              v-for="item in imageTagType"
                              :key="item.value"
                              :label="item.name"
                              :value="item.value"
                            ></el-option>
                          </el-select>
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-row>
                  <el-row class="arrange-row">
                    <el-col :span="3"> 复制到 </el-col>
                    <el-col :span="12">
                      <el-select
                        v-model="form.copy_to_env_ids"
                        filterable
                        placeholder="请选择要复制到的环境"
                        multiple="true"
                      >
                        <el-option
                          v-for="(item, index) in envList"
                          :key="index"
                          :label="item.name"
                          :value="item.id"
                        ></el-option>
                      </el-select>
                    </el-col>
                  </el-row>
                </el-form>
              </div>
            </el-tab-pane>
          </el-tabs>
        </template>
        <div class="apparrange-layout">
          <el-button
            type="primary"
            class="fb-ly-rbtn"
            icon="el-icon-edit"
            @click="setArrangement"
            >{{ '保存配置' }}</el-button
          >
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import { Message } from 'element-ui';
import backend from '@/api/backend';

export default {
  data() {
    return {
      editableTabsValue: 0,
      dialogFormVisible: false,
      form: {
        config: '',
        image_mapings: [],
        copy_to_env_ids: [],
      },
      currentImageMappings: [],
      dialogVisible: false,
      appID: 0,
      appList: [],
      imageTagType: [
        {
          name: '保持不变',
          value: 1,
        },
        {
          name: '系统默认规则',
          value: 2,
        },
        {
          name: 'latest',
          value: 3,
        },
      ],
      imagesData: [],
    };
  },
  components: {},
  computed: {
    ...mapGetters({
      projectID: 'projectID',
    }),
  },
  // watch: {
  //   $route(to, from) {
  //     this.lastUrl = from.name;
  //   },
  //   prd_mem_requests(val) {
  //     this.form.test.mem_limits = val;
  //     this.prd_mem_limits = val;
  //   },
  //   prd_mem_limits(val) {
  //     this.form.test.mem_requests = val;
  //     this.prd_mem_requests = val;
  //   },
  // },
  // destroyed() {
  //   formData = _.cloneDeep(FORMDATA);
  //   this.$set(this.form, 'form', formData);
  //   this.$nextTick(() => { });
  // },
  created() {},
  props: {
    envList: {
      type: Array,
      default: [],
    },
    appList: {
      type: Array,
      default: [],
    },
  },
  methods: {
    parseYamlImages() {
      backend.parseYamlImages(this.form, (data) => {
        this.form.image_mapings = this.currentImageMappings;
        this.form.image_mapings.forEach((currentItem) => {
          data.forEach((parseItem, index) => {
            if (
              currentItem != undefined &&
              currentItem.name == parseItem.name &&
              currentItem.image == parseItem.image
            ) {
              data.splice(index, 1);
            }
          });
        });
        this.form.image_mapings = this.form.image_mapings.concat(data);
        this.$forceUpdate();
      });
    },
    linkApp(row) {
      this.form.image_mapings.forEach((item) => {
        if (row.name == item.name && row.image == item.image) {
          item.project_app_id = row.project_app_id;
        }
      });
      console.log('change form image_mapping..');
      console.log(this.form.image_mapings);
    },
    linkAppImageTagType(row) {
      this.form.image_mapings.forEach((item) => {
        if (row.name == item.name && row.image == item.image) {
          item.image_tag_type = row.image_tag_type;
        }
      });
      console.log('change form image_mapping..');
      console.log(this.form.image_mapings);
    },
    doSetup(row) {
      this.appID = row.id;
      this.dialogFormVisible = true;
      if (this.envList.length > 0) {
        this.editableTabsValue = this.envList[0].id;
      } else {
        Message.error('请先通过"项目设置"--"项目环境" 新增自定义的环境信息，然后再配置 应用编排');
        return;
      }
      backend.getProjectArrange(this.projectID, this.appID, this.editableTabsValue, (data) => {
        if (data !== null) {
          this.form = data;
          if (data.image_mapings == undefined) {
            this.form.image_mapings = [];
            this.currentImageMappings = [];
          } else {
            this.currentImageMappings = data.image_mapings;
          }
        } else {
          resetForm()
        }
      });
    },
    handleClick() {
      this.form.copy_to_env_ids = [];
      this.form.config = '';
      backend.getProjectArrange(this.projectID, this.appID, this.editableTabsValue, (data) => {
        if (data !== null) {
          this.form = data;
          if (data.image_mapings == undefined) {
            this.form.image_mapings = [];
            this.currentImageMappings = [];
          } else {
            this.currentImageMappings = data.image_mapings;
          }
        } else {
          resetForm()
        }
      });
    },
    handleClose(done) {
      this.$confirm('确定返回应用代码?')
        .then((_) => {
          done();
        })
        .catch((_) => {});
    },
    // 保存应用编排
    setArrangement() {
      const submitData = { ...this.form };
      if (submitData.config == '') {
        Message.warning('请正确配置应用编排后，再点击 保存');
        return;
      }
      // reset project_app_id/ image_tag_type empty value to 0
      submitData.image_mapings.forEach(function (item, index) {
        if (item.project_app_id === '' || item.project_app_id === undefined) {
          item.project_app_id = 0;
        }
        if (item.image_tag_type === '' || item.image_tag_type === undefined) {
          item.image_tag_type = 1;
        }
      });
      backend.setProjectArrange(
        this.projectID,
        this.appID,
        this.editableTabsValue,
        JSON.stringify(submitData),
        () => {
          Message.success(this.$t('bm.add.setSuc'));
        }
      );
    },
    resetForm() {
      this.form = {
        config: '',
        image_mapings: [],
        copy_to_env_ids: [],
      }
    }
  },
};
</script>
