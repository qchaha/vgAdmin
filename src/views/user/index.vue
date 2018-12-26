
<template>
  <div class="app-container">
    <el-row :gutter="20">
      <el-col :span="18" :offset="2">
        <el-input style="width: 200px;" class="filter-item"/>
        <el-button class="filter-item" type="primary" icon="el-icon-search" @click="$message({message:'敬请期待...',type:'success'})">搜索</el-button>
      </el-col>
      <el-col :span="2" :offset="0">
        <el-button class="filter-item" type="primary" icon="el-icon-plus" @click="newUserOrEditUser('')">新建</el-button>
      </el-col>
    </el-row>
    <el-row :gutter="22">
      <el-col :span="20" :offset="2">
        <br><br>
        <el-table
          v-loading="listLoading"
          :data="list"
          element-loading-text="Loading"
          border
          fit
          highlight-current-row>
          <el-table-column align="center" label="序号" min-width="50">
            <template slot-scope="scope">
              {{ scope.$index + 1 }}
            </template>
          </el-table-column>
          <el-table-column label="用户名" align="center" min-width="150">
            <template slot-scope="scope">
              <span>{{ scope.row.username }}</span>
            </template>
          </el-table-column>
          <el-table-column label="角色" align="center" min-width="150">
            <template slot-scope="scope">
              <span>{{ scope.row.role }}</span>
            </template>
          </el-table-column>
          <el-table-column class-name="status-col" label="账号状态" align="center" min-width="200">
            <template slot-scope="scope">
              <el-tag :type="scope.row.status | statusFilter">{{ userStatusTranslate(scope.row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column align="center" prop="created_at" label="邮箱" min-width="250">
            <template slot-scope="scope">
              <i class="el-icon-mail"/>
              <span>{{ scope.row.email }}</span>
            </template>
          </el-table-column>
          <el-table-column align="center" prop="created_at" label="最后修改时间" min-width="250">
            <template slot-scope="scope">
              <i class="el-icon-time"/>
              <span>{{ scope.row.last_modify_time }}</span>
            </template>
          </el-table-column>
          <el-table-column align="center" prop="created_at" label="操作" min-width="200">
            <template slot-scope="scope">
              <span>
                <el-button type="primary" @click="newUserOrEditUser(scope.row)">编辑</el-button>
                <el-button type="danger" @click="deleteRowData(scope.row)">删除</el-button>
              </span>
            </template>
          </el-table-column>
        </el-table>
      </el-col>
    </el-row>
    <el-dialog :visible.sync="dialogFormVisible" :before-close="cancel" width="30%">
      <el-form ref="formRowData" :model="formRowData" :rules="rules" label-width="120px">
        <el-form-item label="用户名" prop="username" required>
          <el-input v-model="formRowData.username" style="width: 300px; padding: 0 0 10px 20px"/>
        </el-form-item>
        <el-form-item label="密码" prop="password" required>
          <el-input :type="passwordType" v-model="formRowData.password" style="width: 300px; padding: 0 0 10px 20px"/>
          <span class="show-pwd" @click="showPwd">
            <svg-icon :icon-class="eyeType" />
          </span>
        </el-form-item>
        <el-form-item label="确认密码" prop="passwordRepeat" required>
          <el-input :type="passwordType" v-model="formRowData.passwordRepeat" style="width: 300px; padding: 0 0 10px 20px"/>
          <span class="show-pwd" @click="showPwd">
            <svg-icon :icon-class="eyeType" />
          </span>
        </el-form-item>
        <el-form-item label="角色" prop="role" required>
          <el-select v-model="formRowData.role" style="width: 300px; padding: 0 0 10px 20px">
            <el-option label="super" value="super"/>
            <el-option label="developer" value="developer"/>
          </el-select>
        </el-form-item>
        <el-form-item label="账号状态" prop="status" required>
          <el-select v-model="formRowData.status" style="width: 300px; padding: 0 0 10px 20px">
            <el-option label="Open" value="1"/>
            <el-option label="Locked" value="0"/>
          </el-select>
        </el-form-item>
        <el-form-item label="邮箱" prop="email" required>
          <el-input v-model="formRowData.email" style="width: 300px; padding: 0 0 10px 20px"/>
        </el-form-item>
        <el-row :gutter="20" style="padding-top: 10px">
          <el-col :span="6" :offset="4"><el-button type="primary" @click="submitForm()">确认</el-button></el-col>
          <el-col :span="6" :offset="6"><el-button type="warning" @click="cancel()">取消</el-button></el-col>
        </el-row>
      </el-form>
    </el-dialog>
  </div>
</template>

<script>
import { getList } from '@/api/user'
import { newUser } from '@/api/user'
import { deleteUser } from '@/api/user'
import { editUser } from '@/api/user'
export default {
  filters: {
    statusFilter(status) {
      const statusMap = {
        1: 'success',
        0: 'danger'
      }
      return statusMap[status]
    }
  },
  data() {
    var validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入密码'))
      } else {
        if (this.formRowData.passwordRepeat !== '') {
          this.$refs.formRowData.validateField('passwordRepeat')
        }
        callback()
      }
    }
    var validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'))
      } else if (value !== this.formRowData.password) {
        callback(new Error('两次输入的密码不一致'))
      } else {
        callback()
      }
    }
    return {
      list: null,
      listLoading: true,
      dialogFormVisible: false,
      passwordType: 'password',
      eyeType: 'eye',
      newOrEdit: '',
      formRowData: {
        username: '',
        password: '',
        passwordRepeat: '',
        role: '',
        status: '',
        email: '',
        origin_username: ''
      },
      rules: {
        username: [
          { message: '请输入用户名', trigger: 'blur' }
        ],
        password: [
          { validator: validatePass, trigger: 'blur' }
        ],
        passwordRepeat: [
          { validator: validatePass2, trigger: 'blur' }
        ],
        role: [
          { message: '请选择用户角色', trigger: 'blur' }
        ],
        status: [
          { message: '请选择账号状态', trigger: 'blur' }
        ],
        email: [
          { message: '请输入邮箱地址', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
        this.eyeType = 'eye-open'
      } else {
        this.passwordType = 'password'
        this.eyeType = 'eye'
      }
    },
    deleteRowData(row) {
      this.$confirm('确认要删除用户 "' + row.username + '" 吗？', '警告', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteUser(row.username).then(response => {
          if (response.code === 20000) {
            this.$message({
              message: '删除成功...',
              type: 'success'
            })
            this.fetchData()
          }
        })
      }).catch(() => {
      })
    },
    submitForm() {
      this.$refs.formRowData.validate(valid => {
        if (valid) {
          if (this.newOrEdit === 'new') {
            newUser(
              this.formRowData.username,
              this.formRowData.password,
              this.formRowData.role,
              this.formRowData.status,
              this.formRowData.email).then(response => {
              if (response.code === 20000) {
                this.$message({
                  message: '添加成功...',
                  type: 'success'
                })
                this.fetchData()
                this.cancel()
              }
            })
          } else {
            editUser(
              this.formRowData.username,
              this.formRowData.password,
              this.formRowData.role,
              this.formRowData.status,
              this.formRowData.email,
              this.formRowData.origin_username
            ).then(response => {
              if (response.code === 20000) {
                this.$message({
                  message: '修改成功...',
                  type: 'success'
                })
                this.fetchData()
                this.cancel()
              }
            })
          }
        } else {
          this.$message({
            message: '操作失败...',
            type: 'error'
          })
        }
      })
    },
    fetchData() {
      this.listLoading = true
      getList().then(response => {
        this.list = response.data.items
        this.listLoading = false
      })
    },
    cancel() {
      this.dialogFormVisible = false
      this.formRowData = Object.assign({}, {})
      setTimeout(() => { this.$refs.formRowData.clearValidate() }, 500)
    },
    confirm(row) {
      this.$message({
        message: '修改成功...',
        type: 'success'
      })
    },
    userStatusTranslate(status) {
      const us = {
        1: 'Open',
        0: 'Locked'
      }
      return us[status]
    },
    newUserOrEditUser(row) {
      if (row !== '') {
        this.dialogFormVisible = true
        this.newOrEdit = 'edit'
        this.formRowData = Object.assign({}, row)
        this.formRowData.origin_username = this.formRowData.username
      } else {
        this.dialogFormVisible = true
        this.newOrEdit = 'new'
      }
    }
  }
}
</script>

<style scoped>
.show-pwd {
    cursor: pointer;
    user-select: none;
  }
.cancel-btn {
  position: absolute;
  right: 15px;
  top: 16px;
}
</style>
