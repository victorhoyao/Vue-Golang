<template>
  <div class="container">
    <Breadcrumb :items="[t('menu.taskItem'), t('menu.taskItem.management')]" />
    <a-card class="general-card">
      <a-space direction="vertical" :size="16" style="width: 100%;">
        <a-typography-title :heading="5">
          {{ $t('menu.taskItem.management') }}
        </a-typography-title>
        <a-button type="primary" @click="handleAddTaskItem">Add Task Item</a-button>
        <a-table :columns="columns" :data="tableDataWithPrice" :pagination="false">
          <template #columns>
            <a-table-column title="Task ID" data-index="taskId"></a-table-column>
            <a-table-column title="Task Name" data-index="taskName"></a-table-column>
            <a-table-column title="Supplier" data-index="supplier"></a-table-column>
            <a-table-column title="User Types" data-index="userTypes"></a-table-column>
            <a-table-column title="Price (Real/Protocol/Manual)" data-index="price"></a-table-column>
            <a-table-column title="Selection Mode" data-index="selectionMode"></a-table-column>
            <a-table-column title="Status" data-index="status"></a-table-column>
            <a-table-column title="Actions">
              <template #cell="{ record }">
                <a-button type="text" @click="handleEdit(record)">Edit</a-button>
                <a-popconfirm content="Are you sure you want to delete this task item?" @ok="handleDelete(record.id)">
                  <a-button type="text" status="danger">Delete</a-button>
                </a-popconfirm>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </a-space>
    </a-card>

    <!-- Add/Edit Task Item Modal -->
    <a-modal
      :visible="modalVisible"
      :title="isEditMode ? 'Edit Task Item' : 'Add Task Item'"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
    >
      <a-form :model="currentTaskItem" layout="vertical">
        <a-form-item field="taskId" label="Task ID" :rules="[{ required: true, message: 'Task ID is required' }]" v-if="!isEditMode">
          <a-input v-model="currentTaskItem.taskId" placeholder="Enter Task ID"></a-input>
        </a-form-item>
        <a-form-item field="taskName" label="Task Name" :rules="[{ required: true, message: 'Task Name is required' }]">
          <a-input v-model="currentTaskItem.taskName" placeholder="Enter Task Name"></a-input>
        </a-form-item>
        <a-form-item field="supplier" label="Supplier">
          <a-input v-model="currentTaskItem.supplier" placeholder="Enter Supplier"></a-input>
        </a-form-item>
        <a-form-item field="userTypes" label="User Types">
          <a-input v-model="currentTaskItem.userTypes" placeholder="Enter User Types"></a-input>
        </a-form-item>
        <a-form-item field="priceReal" label="Price (Real)">
          <a-input-number v-model="currentTaskItem.priceReal" :min="0" :precision="2"></a-input-number>
        </a-form-item>
        <a-form-item field="priceProtocol" label="Price (Protocol)">
          <a-input-number v-model="currentTaskItem.priceProtocol" :min="0" :precision="2"></a-input-number>
        </a-form-item>
        <a-form-item field="priceManual" label="Price (Manual)">
          <a-input-number v-model="currentTaskItem.priceManual" :min="0" :precision="2"></a-input-number>
        </a-form-item>
        <a-form-item field="selectionMode" label="Selection Mode">
          <a-select v-model="currentTaskItem.selectionMode">
            <a-option value="Single-select">Single-select</a-option>
            <a-option value="Multi-select">Multi-select</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="status" label="Status">
          <a-select v-model="currentTaskItem.status">
            <a-option value="Active">Active</a-option>
            <a-option value="Inactive">Inactive</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted, computed } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import { useI18n } from 'vue-i18n';
  import { 
    getTaskItems, 
    addTaskItem, 
    updateTaskItem, 
    deleteTaskItem,
    type TaskItemCreateData,
    type TaskItemUpdateData
  } from '../../api/taskItem';
  import type { TaskItem, TaskItemFormData } from '../../types/taskItem';

  const { t } = useI18n();

  const columns = [
    {
      title: 'Task ID',
      dataIndex: 'taskId',
    },
    {
      title: 'Task Name',
      dataIndex: 'taskName',
    },
    {
      title: 'Supplier',
      dataIndex: 'supplier',
    },
    {
      title: 'User Types',
      dataIndex: 'userTypes',
    },
    {
      title: 'Price (Real/Protocol/Manual)',
      dataIndex: 'price',
    },
    {
      title: 'Selection Mode',
      dataIndex: 'selectionMode',
    },
    {
      title: 'Status',
      dataIndex: 'status',
    },
    {
      title: 'Actions',
      slotName: 'actions',
    },
  ];

  const tableData = ref<TaskItem[]>([]);
  const modalVisible = ref(false);
  const isEditMode = ref(false);

  const currentTaskItem = ref<TaskItemCreateData | TaskItem>({
    taskId: '',
    taskName: '',
    supplier: '',
    userTypes: '',
    priceReal: 0,
    priceProtocol: 0,
    priceManual: 0,
    selectionMode: 'Single-select',
    status: 'Active',
  });

  // Computed property to format price display
  const tableDataWithPrice = computed(() => {
    return tableData.value.map(item => ({
      ...item,
      price: `¥${item.priceReal} / ¥${item.priceProtocol} / ¥${item.priceManual}`
    }));
  });

  const fetchTaskItems = async () => {
    try {
      const res = await getTaskItems();
      tableData.value = res.data.taskItems;
      Message.success('Task items fetched successfully');
    } catch (error) {
      Message.error('Failed to fetch task items');
      console.error('Error fetching task items:', error);
    }
  };

  const handleAddTaskItem = () => {
    isEditMode.value = false;
    currentTaskItem.value = {
      taskId: '',
      taskName: '',
      supplier: '',
      userTypes: '',
      priceReal: 0,
      priceProtocol: 0,
      priceManual: 0,
      selectionMode: 'Single-select',
      status: 'Active',
    };
    modalVisible.value = true;
  };

  const handleEdit = (record: TaskItem) => {
    isEditMode.value = true;
    currentTaskItem.value = { ...record };
    modalVisible.value = true;
  };

  const handleDelete = async (id: number | undefined) => {
    if (id === undefined) {
      Message.error('Cannot delete: Task Item ID is missing.');
      return;
    }
    try {
      await deleteTaskItem(id);
      Message.success('Task item deleted successfully');
      fetchTaskItems(); // Refresh table
    } catch (error) {
      Message.error('Failed to delete task item');
      console.error('Error deleting task item:', error);
    }
  };

  const handleModalOk = async () => {
    // Basic validation
    if (!currentTaskItem.value.taskName || (!isEditMode.value && !currentTaskItem.value.taskId)) {
      Message.error('Please fill in all required fields.');
      return;
    }

    try {
      if (isEditMode.value) {
        // Update existing task item
        if ('id' in currentTaskItem.value && currentTaskItem.value.id !== undefined) {
          await updateTaskItem(currentTaskItem.value.id, currentTaskItem.value as TaskItemUpdateData);
          Message.success('Task item updated successfully');
        } else {
          Message.error('Task Item ID is missing for update.');
          return;
        }
      } else {
        // Add new task item
        await addTaskItem(currentTaskItem.value as TaskItemCreateData);
        Message.success('Task item added successfully');
      }
      modalVisible.value = false;
      fetchTaskItems(); // Refresh table
    } catch (error) {
      Message.error(`Failed to ${isEditMode.value ? 'update' : 'add'} task item`);
      console.error(`Error ${isEditMode.value ? 'updating' : 'adding'} task item:`, error);
    }
  };

  const handleModalCancel = () => {
    modalVisible.value = false;
  };

  onMounted(() => {
    fetchTaskItems();
  });
</script>

<script lang="ts">
  export default {
    name: 'TaskItemManagement',
  };
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px 20px;
  }

  .general-card {
    margin-top: 20px;
  }
</style> 