<template>
    <div class="page-content">
        <div class="portlet-body min-height150">
          <template>
            <div class="pv10" v-show="proShow">
              <el-row>
                <el-col :span="24">
                  <div class="setTitle">
                    <div class="f-r">
                      <el-button v-if="!projectInfo.status" @click="closeProject">结束项目</el-button>
                      <el-button type="primary" @click="doEdit(projectInfo.id)">编辑</el-button>
                    </div>
                    {{projectInfo.name}}<span class="setDescription">{{projectInfo.description}}</span>
                  </div>
                </el-col>
              </el-row>
              <el-row>
                <el-col class="setText">
                  管理员：{{projectInfo.owner}}
                </el-col>
                <el-col class="setText">项目成员：{{projectInfo.members || 0}}人</el-col>
              </el-row>
              <el-row>
                <el-col class="setText">
                  创建者：{{projectInfo.owner}}
                </el-col>
                <el-col class="setText">创建时间：{{projectInfo.create_at}}</el-col>
              </el-row>
              <el-row>
                <el-col class="setText">
                  开始时间：{{projectInfo.start_at}}
                </el-col>
                <el-col class="setText">结束时间：{{projectInfo.end_at}}</el-col>
              </el-row>
              <el-row>
                <el-col class="setText">
                  项目状态：{{projectInfo.status === 1 ? '运行中' : '已结束'}}
                </el-col>
              </el-row>
            </div>
          </template>
        </div>
        <div class="portlet-body">
          <template>
              <div class="clearfix pv10">
                <div class="f-r" v-if="isAdmin"><el-button type="primary" @click="doSubmitMember">加入项目</el-button></div>
                <el-form :model="form" ref="ruleForm" :rules="rules" class="form-number" v-if="isAdmin">
                  <el-row class="memberRow">
                    <el-col class="form-col">
                      <el-form-item prop="role">
                        <el-select v-model="form.role" placeholder="请选择角色" v-on:change="roleChange" filterable>
                          <el-option v-for="(item, index) in roleList" :key="index" :label="item.description" :value="item.id">
                          </el-option>
                        </el-select>
                      </el-form-item>
                    </el-col>
                    <el-col class="form-col">
                      <el-form-item prop="member">
                        <el-select v-model="form.member" placeholder="请选择成员" filterable>
                          <el-option v-for="(item, index) in memberList" :key="index" :label="item.name" :value="item.user">
                          </el-option>
                        </el-select>
                      </el-form-item>
                    </el-col>
                  </el-row>
                </el-form>
                <span class="font-title">成员列表</span>
                <!-- <el-button v-if="isAdmin" class="ml10" type="primary" @click="showAddUser">添加成员</el-button> -->
              </div>
              <el-table border :data="roleCheck">
                  <span slot="empty">
                    暂无
                  </span>
                  <el-table-column prop="user" label="用户名" min-width="15%"></el-table-column>
                  <el-table-column prop="realName" label="真实姓名" min-width="15%"></el-table-column>
                  <el-table-column prop="email" label="邮箱" min-width="15%"></el-table-column>
                  <el-table-column prop="role" label="角色" min-width="15%"></el-table-column>
                  <el-table-column prop="deleted" label="状态" min-width="15%">
                    <template slot-scope="scope">{{ scope.row.deleted ? '已删除' : '正常' }}</template>
                  </el-table-column>
                  <el-table-column label="操作" min-width="15%">
                    <template slot-scope="scope" v-if="scope.row.name!=='default'">
                      <el-button v-if="isAdmin" type="text" size="small" @click="delMemeber(scope.row.id)">移除</el-button>
                      </el-button>
                    </template>
                  </el-table-column>
              </el-table>
          </template>
        </div>
        <project-create ref="create" v-on:getlist="getProjectInfo"></project-create>
        <el-dialog custom-class="dialog-right width600" v-if="dialogVisible" title="添加成员" :close-on-click-modal="false"
          :visible.sync="dialogVisible">
          <template>
            <p class="font-gray pv20">提示：若添加多个成员，请输入空格或回车分割</p>
            <div class="add-user-panel" @click="inputFocus">
              <div class="mail-panel" v-for="(item,index) in mailList">
                {{item}}
                <div class="mail-close ml10" @click="delMail(index)"><i class="el-icon-close"></i></div>
              </div>
              <el-input ref="mailFocus" class="dialog-mail" v-model.trim="dialogMail" @blur="addMail" @keyup.enter.native="addMail" @keyup.space.native="addMail" placeholder="请输入邮箱"></el-input>
            </div>
          </template>
          <div slot="footer" class="dialog-footer">
            <el-button @click="dialogVisible = false">{{$t('bm.other.cancel')}}</el-button>
            <el-button type="primary" @click="dialogSubmit">确定</el-button>
          </div>
        </el-dialog>
    </div>
</template>

<style scoped>
  .min-height150 {
    min-height: 150px;
  }
  .pv10 {
    padding-top: 10px;
    padding-bottom: 10px;
  }
  .member-btn {
    width: 550px;
    text-align: right;
  }
  .f-r {
    float: right;
  }
  .pv20 {
    padding-top: 20px;
    padding-bottom: 20px;
  }
  .memberRow .el-select {
    width: 100%;
  }
  .containerMember {
    width: 550px;
    border: 1px solid #ccc;
    padding: 10px;
  }
  .el-tag {
    margin-right: 5px;
    margin-bottom: 3px;
  }
  .mb15 {
    margin-bottom: 15px;
  }
  .setTitle {
    font-size: 18px;
    color: #333;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: bold;
    line-height: 40px;
  }
  .setDescription {
    color: #606266;
    font-size: 14px;
    margin-left: 10px;
    font-weight: 400;
  }
  .setText {
    width: 400px;
    line-height: 20px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    margin-right: 10px;
    margin-top: 15px;
    color: #333;
    font-weight:400;
  }
  .form-number {
    width: 320px;
    float: right;
    margin-right: 5px;
  }
  .form-col {
    width: 150px;
    margin-right: 10px;
  }
  .font-title {
    font-size: 16px;
    color: #409EFF;
    display: inline-block;
    line-height: 50px;
  }
  .el-form-item {
    margin-bottom: 5px;
  }
  .proj-edit {
    font-size: 20px;
    color: #409EFF;
    vertical-align: middle;
    cursor: pointer;
  }
  .add-user-panel {
    border: 1px solid #ccc;
    padding: 10px 5px;
    height: 250px;
    overflow-y: auto;
  }
  .dialog-mail {
    width: 200px;
    border: 0;
  }
  .mail-panel {
    display: inline-block;
    height: 38px;
    line-height: 38px;
    padding: 0 8px;
    border-radius: 4px;
    background-color: #ecf5ff;
    border: 1px solid #d9ecff;
    color: #409eff;
    margin: 0 5px 10px 5px;
  }
  .mail-close {
    display: inline-block;
    height: 100%;
    padding-left: 8px;
    cursor: pointer;
    border-left: 1px solid #d9ecff;
  }
</style>
<style>
  .dialog-right {
    position: fixed;
    right: 0;
    height: 100%;
    margin: 0 !important;
    padding-bottom: 70px;
    border-radius: 0 !important;
    overflow-y: auto;
    z-index: 2;
  }

  .dialog-right .el-dialog__header {
    border-bottom: 1px solid #DCDFE6;
    padding: 20px;
  }

  .dialog-right .el-dialog__body {
    padding-top: 0;
  }

  .dialog-right .el-dialog__footer {
    position: fixed;
    bottom: 0;
    right: 0;
    padding: 10px;
    text-align: center;
    background-color: #fff;
    box-shadow: 0px -2px 4px 0px rgba(0, 0, 0, 0.12);
    z-index: 3;
  }
  .width600, .width600 .el-dialog__footer {
    width: 600px;
  }
  .dialog-mail .el-input__inner {
    border: 0;
    height: 38px;
  }
</style>
<script>
  import { mapGetters } from 'vuex';
  import { Message, MessageBox } from 'element-ui';
  import backend from '@/api/backend';
  import UtilsFn from '@/common/utils';
  import ProjectCreate from './dialog/ProjectCreate';

export default {
  data() {
    return {
      form: {
        role: '',
        member: '',
      },
      rules: {
        role: [
          { required: true, message: '请选择角色', trigger: 'blur' },
        ],
        member: [
          { required: true, message: '请选择成员', trigger: 'blur' },
        ],
      },
      proShow: false,
      roleList: [],
      memberList: [],
      roleCheck: [],
      member: [],
      stepsCheck: [],
      projectInfo: {},
      editInfo: {},
      isAdmin: false,
      dialogVisible: false,
      dialogMail: '',
      mailList: [],
    };
  },
  components: {
    ProjectCreate,
  },
  computed: {
    ...mapGetters({
      isSysAdmin: 'isAdmin',
      projectID: 'projectID',
    })
  },
  created() {
    //角色、成员下拉框初始化
    const groupName = 'system';
    if(this.isSysAdmin === 1) {
      backend.getGroupViewList(groupName, (data) => {
        if(data.roles) {
          this.roleList = data.roles;
        }
      });
      this.isAdmin = true;
    }
    this.getMemberList();
    this.getProjectInfo();
  },
  methods: {
    roleChange(val) {
      this.form.member = '';
      if(val) {
        let checkRole = '';
        this.roleList.map((i) => {
          if(i.id === val) {
            checkRole = i.role;
            return;
          }
        });
        backend.getProjectUser("system", checkRole, (data) => {
          this.memberList = data.users;
        });
      } else {
        this.memberList = [];
      }
    },
    getMemberList() {
      //获取成员列表
      backend.getProjectMember(this.projectID, (data) => {
        if(data) {
          this.roleCheck = data;
        }
      });
    },
    doSubmitMember() {
      this.$refs.ruleForm.validate((valid) => {
        if (valid) {
          const cl = {
            role_id: this.form.role,
            user: this.form.member
          };
          backend.updateProjectMember(this.projectID, cl, (data) => {
            Message.success('加入成功！');
            this.getMemberList();
          });
        }
      });
    },
    delMemeber(memberId) {
      MessageBox.confirm('是否确定删除当前成员', '提示', {type: 'warning'})
        .then(() => {
          backend.delProjectMember(this.projectID, memberId, (data) => {
            Message.success('删除成功！');
            this.getMemberList();
          });
        })
        .catch(() => {});
    },
    getProjectInfo() {
      backend.getProjectDetail(this.projectID, (data) => {
        if(data) {
          data.create_at = data.create_at ? UtilsFn.format(new Date(data.create_at), 'yyyy-MM-dd hh:mm:ss') : '';
          data.start_at = data.create_at ? UtilsFn.format(new Date(data.start_at), 'yyyy-MM-dd hh:mm:ss') : '';
          data.end_at = data.end_at ? UtilsFn.format(new Date(data.end_at), 'yyyy-MM-dd hh:mm:ss') : '';
          this.projectInfo = Object.assign({},data);
          this.proShow = true;
          this.editInfo = Object.assign({},data);
        }
      });
    },
    valChange(flag) {
      const val = this.projectInfo[flag];
      const change = this.editInfo[flag];
      this.$refs.proForm.validate((valid) => {
        if(valid) {
          // 这边预留调用接口，修改项目信息
          const params = {
            name: this.editInfo.name,
            description: this.editInfo.description,
            owner: this.editInfo.owner
          };
          backend.updateNewProject(this.projectID, params, (data) => {
            this.getProjectInfo();
          });
          //this.projectInfo[flag] = this.editInfo[flag];
        }
      });
    },
    closeProject() {
      MessageBox.confirm('确定关闭当前项目吗？', '提示', {type: 'warning'})
        .then(() => {
          // 这边预留关闭调用接口
          this.projectInfo.status = 1;
        })
        .catch(() => {});
    },
    doEdit(id) {
      backend.getProjectDetail(id, (data) => {
        if(data) {
          this.$refs.create.doCreate(true, data);
        }
      });
    },
    showAddUser() {
      this.mailList = [];
      this.dialogVisible = true;
    },
    inputFocus() {
      this.$refs.mailFocus.focus();
    },
    addMail() {
      if(this.dialogMail) {
        this.mailList.push(this.dialogMail);
        this.dialogMail = '';
      }
    },
    delMail(index) {
      this.mailList.splice(index, 1);
    },
    dialogSubmit() {
      console.info(this.mailList);
    },
  }
}
</script>
