<style>
.createDialog .el-dialog__body .el-form-item .el-form-item__content {
  display: flex;
  flex-direction: column;
}
</style>
<template>
  <el-dialog top='25vh' :title="title" :close-on-click-modal="false" :visible.sync="dialogFormVisible" class="createDialog"  width='50%' :before-close="doCancelCreate">
    <el-form :model="form" ref="ruleForm" :rules="rules">
      <el-form-item label="节点名称" prop="name">
        <el-input v-model.trim="form.name" auto-complete="off" maxlength="10" placeholder="请输入节点名称"></el-input>
      </el-form-item>
      <el-form-item label="节点类型" prop="type">
        <el-select v-model="form.type" placeholder="请选择节点类型" filterable @change="stepTypeChange">
          <el-option v-for="(item, index) in typeList" :key="index" :label="item.name" :value="item.type">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="节点描述" prop="description">
        <el-input v-model.trim="form.description" auto-complete="off" maxlength="50" placeholder="请输入描述" ></el-input>
      </el-form-item>
      
      <el-form-item label="子任务" v-show="isShow">
        <div v-if="form.sub_task" class="tag">
          <el-tag style="float: left; margin-bottom: 10px;"
            :key="item.name"
            v-for="(item,index) in form.sub_task"
            closable
            @close="handleClose(index)" >
            {{item.name}}
          </el-tag>
        </div>
        <el-select v-model="form.task_item" placeholder="请输入子任务名称" filterable style="width: 60%">
          <el-option v-for="(item) in subTaskList" :key="item.type" :label="item.name" :value="item.type"></el-option> 
        </el-select>
        <div>
          <span style="float: left"><i class="el-icon-plus" @click="addTaskItem"></i></span>
        </div>
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
import backend from '../../../api/backend';
import createTemplate from '../../../common/createTemplate';
import validate from '../../../common/validate';

const formData = {
  name: '',
  type: '',
  description: '',
  task_item: '',
};

export default {
  mixins: [createTemplate, validate],
  data() {
    return {
      name: '',
      groupRoleList: [],
      typeList: [],
      // 是否属于编辑状态
      isEdit: false,
      dialogFormVisible: false,
      form: JSON.parse(JSON.stringify(formData)),
      title: '新增',
      rules: {
        name: [
          { required: true, message: '请输入名称', trigger: 'blur' },
        ],
        type: [
          { required: true, message: '请选择节点类型', trigger: 'blur' },
        ],
        description: [
          { required: true, message: '描述信息不能为空', trigger: 'blur' },
        ],
        task_item: [
           { required: false },
        ]
      },
      rowId: '',
      isShow: false,
      subTaskList: [
        { 'name': '检出代码', 'type': 'checkout'},
        { 'name': '编译', 'type': 'compile', 
          'params': [
            {
                 "language": "java",
                 "version": "jdk8",
                 "image": "harbor.com/library/jdk8:latest",
                 "compile_command": "mvn install ..."
            },
            {
                 "language": "node",
                 "version": "12",
                 "image": "harbor.com/library/node:latest",
                "compile_command": "yarn install .."
            }
          ]
        },
        { 'name': '制作镜像', 'type': 'build-image'},
        // { 'name': '自定义脚本', 'type': 'custom-script'}
      ]
    };
  },
  computed: {
    ...mapGetters({
      loading: 'getPopLoading',
    }),
  },
  mounted() {
    backend.getStepComponent((data) => {
      if(data){
        this.typeList = data;
      }
    });
    
  },
  methods: {
    stepTypeChange(item) {
      if(item) {
        if(item === 'build'){
          this.isShow = true;
        } else {
          this.isShow = false;
        }
      }
    },
    addTaskItem() {
      if (this.form.task_item.trim() === "") {
        return 
      }

      let subTaskName = this.form.task_item
      for (let i=0; i< this.subTaskList.length; i++) {
        if (this.subTaskList[i].type === this.form.task_item ) {
          subTaskName = this.subTaskList[i].name
          break
        }
      } 

      let item = {
        type: this.form.task_item,
        name: subTaskName
      }
      let sub_task = this.form.sub_task
      if (Object.prototype.toString.call(sub_task) === '[object Array]' && sub_task.length > 0 ) {
        item.index = sub_task[sub_task.length - 1].index + 1
        this.form.sub_task.push(item)
      } else {
        item.index = 1
        this.form.sub_task = [item]
      }
      this.form.task_item = ''
      console.log(this.form.sub_task)
    },
    handleClose(index) {
      // TODO: after delete sub task, need trigger component, then vue page can see the change
      console.log('will close index: ' + index)
      this.form.sub_task.splice(index, 1)
      let sub_task = this.form.sub_task
      sub_task.forEach((item, index) => {
         item.index = index + 1
      })
      console.log('current sub_task:')
      console.log(sub_task)
      this.form.sub_task = sub_task
    },
    doCreate(flag, item) {
      this.isEdit = flag;
      if (flag) {
        this.title = '编辑';
        this.form = {
          name: item.name || '',
          type: item.type || '',
          description: item.description || '',
          sub_task: item.sub_task || [],
        };
        this.stepTypeChange(item.type)
        this.rowId = item.id;
      } else {
        this.title = '新增';
        this.form = {
          name: '',
          type: '',
          description: '',
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
            type: this.form.type,
            description: this.form.description,
            sub_task: this.form.sub_task,
          };
          if (this.isEdit) {
            backend.editStep(this.rowId, cl, () => {
              successCallBack();
            });
          } else {
            backend.AddStep(cl, () => {
              successCallBack();
            });
          }
        }
      });
    },
  },
};
</script>
