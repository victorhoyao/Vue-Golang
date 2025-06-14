<template>
  <div class="container">
    <Breadcrumb :items="['menu.taskDistribution', 'menu.taskDistribution.management']" />
    <a-card class="general-card">
      <a-space direction="vertical" :size="16" style="width: 100%;">
        <a-typography-title :heading="5">
          {{ $t('menu.taskDistribution.management') }}
        </a-typography-title>
        
        <!-- Create Distribution Button -->
        <a-button type="primary" @click="handleCreateDistribution">
          Create Task Distribution
        </a-button>

        <!-- Empty State -->
        <a-empty v-if="tableData.length === 0" description="No task distributions found">
          <a-button type="primary" @click="handleCreateDistribution">
            Create Your First Distribution
          </a-button>
        </a-empty>

        <!-- Distributions Table -->
        <a-table 
          :columns="columns" 
          :data="tableData" 
          :pagination="paginationProps"
          @page-change="handlePageChange"
          @page-size-change="handlePageSizeChange"
        >
          <template #batchName="{ record }">
            {{ record.batchName || 'N/A' }}
          </template>

          <template #taskItem="{ record }">
            {{ record.taskItem?.taskName || record.taskItemName || 'N/A' }}
          </template>

          <template #totalTasks="{ record }">
            {{ record.totalTasks || 0 }}
          </template>

          <template #realMachineTasks="{ record }">
            {{ record.realMachineTasks || 0 }}
          </template>

          <template #protocolTasks="{ record }">
            {{ record.protocolTasks || 0 }}
          </template>

          <template #manualTasks="{ record }">
            {{ record.manualTasks || 0 }}
          </template>

          <template #status="{ record }">
            <a-tag 
              :color="getStatusColor(record.status || 'pending')"
            >
              {{ (record.status || 'pending').toUpperCase() }}
            </a-tag>
          </template>
          
          <template #progress="{ record }">
            <a-progress 
              :percent="record.totalTasks > 0 ? ((record.distributedTasks || 0) / record.totalTasks) * 100 : 0"
              :show-text="true"
              size="small"
            />
            <div style="font-size: 12px; color: #666; margin-top: 4px;">
              {{ record.distributedTasks || 0 }} / {{ record.totalTasks || 0 }} tasks
            </div>
          </template>

          <template #createdAt="{ record }">
            {{ record.createdAt ? new Date(record.createdAt).toLocaleString() : 'N/A' }}
          </template>

          <template #actions="{ record }">
            <a-space>
              <a-button 
                v-if="(record.status || 'pending') === 'pending'" 
                type="text" 
                size="small" 
                @click="handleActivate(record.id)"
              >
                Activate
              </a-button>
              
              <!-- Complete active distributions -->
              <a-button 
                v-if="(record.status || 'pending') === 'active'" 
                type="text" 
                size="small" 
                status="success"
                @click="handleComplete(record.id)"
              >
                Complete
              </a-button>
              
              <!-- Edit - available for pending and active -->
              <a-button 
                v-if="(record.status || 'pending') !== 'completed'"
                type="text" 
                size="small" 
                @click="handleEdit(record)"
              >
                Edit
              </a-button>
              
              <!-- Delete - available for pending and active -->
              <a-button 
                v-if="(record.status || 'pending') !== 'completed'"
                type="text" 
                size="small" 
                status="danger"
                @click="handleDelete(record.id)"
              >
                Delete
              </a-button>
            </a-space>
          </template>
        </a-table>
      </a-space>
    </a-card>

    <!-- Create/Edit Distribution Modal -->
    <a-modal
      :visible="modalVisible"
      :title="isEditMode ? 'Edit Task Distribution' : 'Create Task Distribution'"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
      :ok-loading="modalLoading"
    >
      <a-form :model="currentDistribution" layout="vertical">
        <a-form-item label="Batch Name" required>
          <a-input 
            v-model="currentDistribution.batchName" 
            placeholder="Enter batch name"
          />
        </a-form-item>

        <a-form-item label="Task Item" required>
          <a-select 
            v-model="currentDistribution.taskItemId" 
            placeholder="Select task item"
            :loading="taskItemsLoading"
          >
            <a-option 
              v-for="item in taskItems" 
              :key="item.id" 
              :value="item.id"
            >
              {{ item.taskName }} ({{ item.taskId }})
            </a-option>
          </a-select>
        </a-form-item>

        <a-form-item label="Total Tasks" required>
          <a-input-number 
            v-model="currentDistribution.totalTasks" 
            :min="1"
            placeholder="Enter total number of tasks"
            style="width: 100%"
          />
        </a-form-item>

        <a-divider>Task Distribution by User Type</a-divider>

        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item label="Real Machine Tasks">
              <a-input-number 
                v-model="currentDistribution.realMachineTasks" 
                :min="0"
                :max="currentDistribution.totalTasks"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="Protocol Tasks">
              <a-input-number 
                v-model="currentDistribution.protocolTasks" 
                :min="0"
                :max="currentDistribution.totalTasks"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="Manual Tasks">
              <a-input-number 
                v-model="currentDistribution.manualTasks" 
                :min="0"
                :max="currentDistribution.totalTasks"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <!-- Distribution Summary -->
        <a-alert 
          v-if="distributionSummary.total > 0"
          :type="distributionSummary.isValid ? 'success' : 'warning'"
          :message="distributionSummary.message"
          show-icon
        />
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, computed } from 'vue';
  import { Message, Modal } from '@arco-design/web-vue';
  import { useI18n } from 'vue-i18n';
  import { 
    getTaskDistributions, 
    createTaskDistribution, 
    updateTaskDistribution, 
    deleteTaskDistribution,
    activateTaskDistribution,
    type TaskDistribution,
    type TaskDistributionCreateData,
    type TaskDistributionUpdateData
  } from '../../api/taskDistribution';
  import { getTaskItems } from '../../api/taskItem';

  const { t } = useI18n();

  // Table data and pagination
  const tableData = ref<TaskDistribution[]>([]);
  const currentPage = ref(1);
  const pageSize = ref(10);
  const total = ref(0);

  // Modal state
  const modalVisible = ref(false);
  const modalLoading = ref(false);
  const isEditMode = ref(false);

  // Task items for dropdown
  const taskItems = ref<any[]>([]);
  const taskItemsLoading = ref(false);

  // Current distribution being edited/created
  const currentDistribution = ref<TaskDistributionCreateData | (TaskDistribution & TaskDistributionUpdateData)>({
    batchName: '',
    taskItemId: 0,
    totalTasks: 100,
    realMachineTasks: 0,
    protocolTasks: 0,
    manualTasks: 0,
  });

  // Table columns
  const columns = [
    { title: 'Batch Name', slotName: 'batchName', width: 150 },
    { title: 'Task Item', slotName: 'taskItem', width: 150 },
    { title: 'Total Tasks', slotName: 'totalTasks', width: 100 },
    { title: 'Real Machine', slotName: 'realMachineTasks', width: 100 },
    { title: 'Protocol', slotName: 'protocolTasks', width: 100 },
    { title: 'Manual', slotName: 'manualTasks', width: 100 },
    { title: 'Progress', slotName: 'progress', width: 150 },
    { title: 'Status', slotName: 'status', width: 100 },
    { title: 'Created At', slotName: 'createdAt', width: 150 },
    { title: 'Actions', slotName: 'actions', width: 150 },
  ];

  // Pagination props
  const paginationProps = computed(() => ({
    current: currentPage.value,
    pageSize: pageSize.value,
    total: total.value,
    showTotal: true,
    showPageSize: true,
  }));

  // Distribution summary
  const distributionSummary = computed(() => {
    const totalDistributed = 
      currentDistribution.value.realMachineTasks + 
      currentDistribution.value.protocolTasks + 
      currentDistribution.value.manualTasks;
    
    const remaining = currentDistribution.value.totalTasks - totalDistributed;
    const isValid = totalDistributed <= currentDistribution.value.totalTasks;
    
    let message = `Distributed: ${totalDistributed}, Remaining: ${remaining}`;
    if (!isValid) {
      message = `Error: Distributed tasks (${totalDistributed}) exceed total tasks (${currentDistribution.value.totalTasks})`;
    }
    
    return {
      total: totalDistributed,
      remaining,
      isValid,
      message,
    };
  });

  // Status color mapping
  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      pending: 'orange',
      active: 'green',
      completed: 'blue',
      cancelled: 'red',
    };
    return colors[status] || 'gray';
  };

  // Fetch distributions
  const fetchDistributions = async () => {
    try {
      const res = await getTaskDistributions(currentPage.value, pageSize.value);
      console.log('Full API Response:', res.data);
      
      // Handle different possible response structures
      let distributions = [];
      let totalCount = 0;
      
      if ((res.data as any).distributions) {
        distributions = (res.data as any).distributions;
        totalCount = (res.data as any).pagination?.total || 0;
      } else if (Array.isArray(res.data)) {
        distributions = res.data;
        totalCount = res.data.length;
      } else if ((res.data as any).data?.distributions) {
        distributions = (res.data as any).data.distributions;
        totalCount = (res.data as any).data.pagination?.total || 0;
      }
      
      console.log('Processed distributions:', distributions);
      console.log('Total count:', totalCount);
      
      // Transform data to ensure all required fields exist
      tableData.value = distributions.map((item: any) => ({
        id: item.id || 0,
        batchName: item.batchName || 'Unknown Batch',
        taskItemId: item.taskItemId || item.taskItemID || 0,
        taskItem: item.taskItem || null,
        taskItemName: item.taskItem?.taskName || item.taskItemName || 'Unknown Task',
        totalTasks: item.totalTasks || 0,
        realMachineTasks: item.realMachineTasks || 0,
        protocolTasks: item.protocolTasks || 0,
        manualTasks: item.manualTasks || 0,
        distributedTasks: item.distributedTasks || (item.realMachineTasks + item.protocolTasks + item.manualTasks) || 0,
        remainingTasks: item.remainingTasks || 0,
        status: item.status || 'pending',
        createdBy: item.createdBy || 0,
        createdAt: item.createdAt || new Date().toISOString(),
        updatedAt: item.updatedAt || new Date().toISOString(),
      })) as any;
      
      total.value = totalCount;
      
      console.log('Final table data:', tableData.value);
    } catch (error) {
      Message.error('Failed to fetch task distributions');
      console.error('Error fetching distributions:', error);
      // Set empty data on error
      tableData.value = [];
      total.value = 0;
    }
  };

  // Fetch task items for dropdown
  const fetchTaskItems = async () => {
    taskItemsLoading.value = true;
    try {
      const res = await getTaskItems();
      taskItems.value = (res.data as any).taskItems || [];
    } catch (error) {
      Message.error('Failed to fetch task items');
    } finally {
      taskItemsLoading.value = false;
    }
  };

  // Handle page change
  const handlePageChange = (page: number) => {
    currentPage.value = page;
    fetchDistributions();
  };

  // Handle page size change
  const handlePageSizeChange = (size: number) => {
    pageSize.value = size;
    currentPage.value = 1;
    fetchDistributions();
  };

  // Handle create distribution
  const handleCreateDistribution = () => {
    isEditMode.value = false;
    currentDistribution.value = {
      batchName: '',
      taskItemId: 0,
      totalTasks: 100,
      realMachineTasks: 0,
      protocolTasks: 0,
      manualTasks: 0,
    };
    modalVisible.value = true;
  };

  // Handle edit distribution
  const handleEdit = (distribution: TaskDistribution) => {
    isEditMode.value = true;
    currentDistribution.value = {
      ...distribution,
      taskItemId: distribution.taskItemId,
    };
    modalVisible.value = true;
  };

  // Handle activate distribution
  const handleActivate = async (id: number) => {
    try {
      await activateTaskDistribution(id);
      Message.success('Task distribution activated successfully');
      fetchDistributions();
    } catch (error) {
      Message.error('Failed to activate task distribution');
      console.error('Error activating distribution:', error);
    }
  };

  // Handle complete distribution
  const handleComplete = async (id: number) => {
    const confirmed = await new Promise((resolve) => {
      const modal = Modal.confirm({
        title: 'Complete Distribution',
        content: 'Are you sure you want to mark this distribution as completed? This action cannot be undone and the distribution cannot be edited or deleted afterwards.',
        okText: 'Complete',
        okButtonProps: { status: 'success' },
        cancelText: 'Cancel',
        onOk: () => {
          resolve(true);
          modal.close();
        },
        onCancel: () => {
          resolve(false);
          modal.close();
        }
      });
    });
    
    if (!confirmed) return;
    
    try {
      // Find the current distribution to get existing values
      const distribution = tableData.value.find(item => item.id === id);
      if (!distribution) return;
      
      await updateTaskDistribution(id, {
        batchName: distribution.batchName,
        totalTasks: distribution.totalTasks,
        realMachineTasks: distribution.realMachineTasks,
        protocolTasks: distribution.protocolTasks,
        manualTasks: distribution.manualTasks,
        status: 'completed'
      });
      Message.success('Task distribution completed successfully');
      fetchDistributions();
    } catch (error) {
      Message.error('Failed to complete task distribution');
      console.error('Error completing distribution:', error);
    }
  };

  // Handle delete distribution
  const handleDelete = async (id: number) => {
    // Find the distribution to check its status
    const distribution = tableData.value.find(item => item.id === id);
    const status = distribution?.status || 'pending';
    
    let confirmMessage = 'Are you sure you want to delete this task distribution?';
    let confirmTitle = 'Delete Distribution';
    
    if (status === 'active') {
      confirmMessage = 'This distribution is currently ACTIVE and may have tasks assigned to users. Are you sure you want to delete it? This action cannot be undone.';
      confirmTitle = 'Delete Active Distribution';
    }
    
    // Show confirmation dialog
    const confirmed = await new Promise((resolve) => {
      const modal = Modal.confirm({
        title: confirmTitle,
        content: confirmMessage,
        okText: 'Delete',
        okButtonProps: { status: 'danger' },
        cancelText: 'Cancel',
        onOk: () => {
          resolve(true);
          modal.close();
        },
        onCancel: () => {
          resolve(false);
          modal.close();
        }
      });
    });
    
    if (!confirmed) return;
    
    try {
      await deleteTaskDistribution(id);
      Message.success('Task distribution deleted successfully');
      fetchDistributions();
    } catch (error) {
      Message.error('Failed to delete task distribution');
      console.error('Error deleting distribution:', error);
    }
  };

  // Handle modal OK
  const handleModalOk = async () => {
    if (!distributionSummary.value.isValid) {
      Message.error(distributionSummary.value.message);
      return;
    }

    if (!currentDistribution.value.batchName || !currentDistribution.value.taskItemId) {
      Message.error('Please fill in all required fields');
      return;
    }

    // Check if at least one task is assigned
    const totalDistributed = 
      currentDistribution.value.realMachineTasks + 
      currentDistribution.value.protocolTasks + 
      currentDistribution.value.manualTasks;
    
    if (totalDistributed === 0) {
      Message.error('At least one user type must have tasks assigned');
      return;
    }

    modalLoading.value = true;
    try {
      if (isEditMode.value) {
        const updateData: TaskDistributionUpdateData = {
          batchName: currentDistribution.value.batchName,
          totalTasks: currentDistribution.value.totalTasks,
          realMachineTasks: currentDistribution.value.realMachineTasks,
          protocolTasks: currentDistribution.value.protocolTasks,
          manualTasks: currentDistribution.value.manualTasks,
        };
        await updateTaskDistribution((currentDistribution.value as TaskDistribution).id, updateData);
        Message.success('Task distribution updated successfully');
      } else {
        const createData: TaskDistributionCreateData = {
          batchName: currentDistribution.value.batchName,
          taskItemId: currentDistribution.value.taskItemId,
          totalTasks: currentDistribution.value.totalTasks,
          realMachineTasks: currentDistribution.value.realMachineTasks,
          protocolTasks: currentDistribution.value.protocolTasks,
          manualTasks: currentDistribution.value.manualTasks,
        };
        await createTaskDistribution(createData);
        Message.success('Task distribution created successfully');
      }
      
      modalVisible.value = false;
      fetchDistributions();
    } catch (error) {
      Message.error(`Failed to ${isEditMode.value ? 'update' : 'create'} task distribution`);
      console.error('Error saving distribution:', error);
    } finally {
      modalLoading.value = false;
    }
  };

  // Handle modal cancel
  const handleModalCancel = () => {
    modalVisible.value = false;
  };

  // Initialize data
  onMounted(() => {
    fetchDistributions();
    fetchTaskItems();
  });
</script>

<style scoped>
  .container {
    padding: 0 20px 20px 20px;
  }
</style> 