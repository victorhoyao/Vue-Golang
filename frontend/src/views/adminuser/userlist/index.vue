<template>
  <div class="container">
    <a-card class="general-card" :title="$t('menu.list.searchTable')">
      <a-divider style="margin-top: 0" />

      <!-- 优化按钮布局 -->
      <div class="toolbar-wrapper">
        <!-- 左侧操作区 -->
        <div class="left-toolbar">
          <!-- 新建用户按钮 -->
          <a-button type="primary" @click="NewUser">
            <template #icon>
              <icon-plus />
            </template>
            {{ $t('searchTable.operation.create') }}
          </a-button>
        </div>
        
        <!-- 右侧搜索区 -->
        <div class="right-toolbar">
          <a-input-search
            v-model="searchKeyword"
            placeholder="请输入用户名搜索"
            search-button
            allow-clear
            @search="handleSearch"
            style="width: 320px"
          >
            <template #button-icon>
              <icon-search />
            </template>
            <template #button-default>
              搜索
            </template>
          </a-input-search>
          <a-button class="reset-button" @click="resetSearch">
            <template #icon>
              <icon-refresh />
            </template>
            重置
          </a-button>
        </div>
      </div>

      <!-- 用户数据表格 -->
      <a-table row-key="id" :loading="loading" :columns="columns" :data="renderData" :bordered="{ cell: true }"
        :pagination="pagination" @page-change="onPageChange"> 
        <!-- 用户状态列自定义渲染 -->
        <template #status="{ record }">
          <a-tag v-if="record.status === 0" color="#00b42a">
            {{ `正常` }}</a-tag>
          <a-tag v-else color="#165dff">
            {{ `关注` }}</a-tag>
        </template>
        <!-- 操作列自定义渲染 -->
        <template #operations="{ record }">
          <!-- 编辑按钮 -->
          <a-button type="text" size="small" @click="userEdit(record)">
            {{ $t('searchTable.columns.operations.Edit') }}
          </a-button>
          <!-- 删除按钮 -->
          <a-button type="text" size="small" status="danger" @click="delUserBtn(record)">
            {{ $t('searchTable.columns.operations.Delete') }}
          </a-button>
        </template>
      </a-table>
    </a-card>
  </div>
  <!-- 新建用户模态框 -->
  <a-modal v-model:visible="visible1" :title="title1" @cancel="handleCancel" @before-ok="handleBeforeOk"
    unmount-on-close>
    <a-form :model="userFrom" ref="formRef">
      <a-form-item field="userName" label="用户名">
        <a-input v-model.trim="userFrom.userName" />
      </a-form-item>
      <a-form-item field="passWord" label="密码">
        <a-input v-model.trim="userFrom.passWord" />
      </a-form-item>
    </a-form>
  </a-modal>

  <!-- 编辑用户模态框 -->
  <a-modal v-model:visible="visible" :title="title" @cancel="handleCancel" @before-ok="handleBeforeOkEdit"
    unmount-on-close>
    <a-form :model="EditFrom" ref="formRef">
      <a-form-item field="userName" label="用户名">
        <a-input v-model.trim="EditFrom.userName" />
      </a-form-item>
      <a-form-item field="newPassword" label="密码">
        <a-input v-model.trim="EditFrom.newPassword" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';
import { reactive, ref } from 'vue';
import { Message } from '@arco-design/web-vue';

import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
import type { UserListRecord, Pagination } from '../../../api/userlist';
import { addUser, editPass, getUserList, delUser, findUser } from '../../../api/userlist';
import useLoading from '../../../hooks/loading';

// 加载状态控制
const { loading, setLoading } = useLoading();

// 表格列类型定义
type Column = TableColumnData & { checked?: true };

// 模态框显示控制
const visible = ref<boolean>(false);    // 编辑用户模态框
const visible1 = ref<boolean>(false);   // 新建用户模态框

// 模态框标题
const title = ref<string>('');
const title1 = ref<string>('');

// 新建用户表单数据
const userFrom = ref({
  userName: '',
  passWord: '',
});

const formRef = ref();
const { t } = useI18n();

// 表格列定义
const columns = ref<Column[]>([
  {
    title: t('searchTable.columns.id'),
    dataIndex: 'id',
    slotName: 'id',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: t('searchTable.columns.userName'),
    dataIndex: 'userName',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: t('searchTable.columns.status'),
    dataIndex: 'status',
    ellipsis: true,
    tooltip: true,
    slotName: 'status',
  },
  {
    title: t('searchTable.columns.userKey'),
    dataIndex: 'userKey',
    slotName: 'userKey',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: t('searchTable.columns.money'),
    dataIndex: 'money',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: t('searchTable.columns.tranType'),
    dataIndex: 'tranType',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: t('searchTable.columns.tranAccount'),
    dataIndex: 'tranAccount',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: t('searchTable.columns.tranName'),
    dataIndex: 'tranName',
    slotName: 'tranName',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: t('searchTable.columns.createTime'),
    dataIndex: 'createTime',
    slotName: 'createTime',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: t('searchTable.columns.operations'),
    dataIndex: 'operations',
    slotName: 'operations',
  },
]);

// 表格数据
const renderData = ref<UserListRecord[]>([]);

// 分页配置
const basePagination: Pagination = { pageNum: 1, pageSize: 10, current: 1 };
const pagination = reactive({
  ...basePagination,
});

// 搜索关键词
const searchKeyword = ref('');

/**
 * 获取用户列表数据
 * @param params 分页参数
 */
const fetchData = async (
  params: Pagination = { pageNum: 1, pageSize: 10 }
) => {
  setLoading(true);
  try {
    const { data } = await getUserList(params);
    console.log('用户列表数据', data);

    renderData.value = data.result;
    pagination.current = params.pageNum;
    pagination.total = data.total;
  } catch (err) {
    // 错误处理
  } finally {
    setLoading(false);
  }
};

/**
 * 搜索用户
 * 根据用户名关键词搜索用户
 */
const handleSearch = async () => {
  if (!searchKeyword.value.trim()) {
    Message.warning('请输入搜索关键词');
    return;
  }
  
  setLoading(true);
  try {
    const res = await findUser(searchKeyword.value.trim());
    if (res.code === 200 && res.data && res.data.userList) {
      renderData.value = res.data.userList;
      // 搜索时不考虑分页
      pagination.total = res.data.userList.length;
    } else {
      Message.error(res.msg || '搜索失败');
    }
  } catch (err) {
    Message.error('搜索出错');
  } finally {
    setLoading(false);
  }
};

/**
 * 重置搜索
 * 清空搜索关键词并重新加载数据
 */
const resetSearch = () => {
  searchKeyword.value = '';
  fetchData();
};

// 初始加载数据
fetchData();

/**
 * 打开新建用户模态框
 */
const NewUser = () => {
  visible1.value = true;
  title1.value = '新建用户';
};

/**
 * 新建用户提交处理
 * @param done 模态框关闭回调
 */
const handleBeforeOk = async (done: any) => {
  setLoading(true);
  const res = await addUser(userFrom.value);
  if (res.code === 200) {
    Message.success(res.msg || '操作成功');
  } else {
    Message.error(res.msg || '操作失败');
  }
  setLoading(false);

  // 重置表单
  userFrom.value = {
    userName: '',
    passWord: '',
  };
  // 刷新数据
  fetchData();
  done();
};

/**
 * 关闭模态框
 */
const handleCancel = () => {
  visible.value = false;
};

/**
 * 分页变化处理
 * @param current 当前页码
 */
const onPageChange = (current: number) => {
  pagination.pageNum = current;
  fetchData(pagination);
};

// 编辑用户表单数据
const EditFrom = ref({
  userName: '',
  newPassword: '',
});

/**
 * 打开编辑用户模态框
 * @param record 用户记录
 */
const userEdit = async (record: any) => {
  console.log(record.userName);
  visible.value = true;
  EditFrom.value = { ...record };
  title.value = '修改用户';
};

/**
 * 编辑用户提交处理
 * @param done 模态框关闭回调
 */
const handleBeforeOkEdit = async (done: any) => {
  console.log(EditFrom.value);

  const data = {
    userId: EditFrom.value.id,
    newPassword: EditFrom.value.newPassword,
  };
  console.log(data);

  const res = await editPass(data);
  if (res.code === 200) {
    Message.success(res.msg);
  } else {
    Message.error(res.msg);
  }
  setLoading(false);

  // 重置表单
  userFrom.value = {
    userName: '',
    passWord: '',
  };
  // 刷新数据
  fetchData();
  done();
};

/**
 * 删除用户
 * @param record 用户记录
 */
const delUserBtn = async (record: any) => {
  console.log(record);

  const params = {
    userId: record.id
  }
  const res = await delUser(params)
  if (res.code === 200) {
    Message.success(res.msg);
  } else {
    Message.error(res.msg);
  }
  // 刷新数据
  fetchData();
}
</script>
<style scoped lang="less">
.container {
  padding: 0 20px 20px 20px;
}

// 工具栏样式
.toolbar-wrapper {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  
  .left-toolbar {
    display: flex;
    align-items: center;
  }
  
  .right-toolbar {
    display: flex;
    align-items: center;
    
    .reset-button {
      margin-left: 8px;
    }
  }
}

:deep(.arco-table-th) {
  &:last-child {
    .arco-table-th-item-title {
      margin-left: 16px;
    }
  }
}

.action-icon {
  margin-left: 12px;
  cursor: pointer;
}

.active {
  color: #0960bd;
  background-color: #e3f4fc;
}

.setting {
  display: flex;
  align-items: center;
  width: 200px;

  .title {
    margin-left: 12px;
    cursor: pointer;
  }
}
</style>
