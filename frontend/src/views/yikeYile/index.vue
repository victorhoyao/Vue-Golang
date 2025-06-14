<template>
  <div class="container">
    <Breadcrumb :items="[t('menu.yikeYile'), t('menu.yikeYile.management')]" />
    <a-card class="general-card">
      <a-space direction="vertical" :size="16" style="width: 100%;">
        <a-typography-title :heading="5">
          {{ $t('menu.yikeYile.management') }}
        </a-typography-title>
        <a-button type="primary" @click="handleAddSupplier">Add Supplier</a-button>
        <a-table :columns="columns" :data="tableData" :pagination="false">
          <template #columns>
            <a-table-column title="Supplier ID" data-index="supplierId"></a-table-column>
            <a-table-column title="Supplier Name" data-index="supplierName"></a-table-column>
            <a-table-column title="Contact" data-index="contact"></a-table-column>
            <a-table-column title="Status" data-index="status"></a-table-column>
            <a-table-column title="Tasks Assigned" data-index="tasksAssigned"></a-table-column>
            <a-table-column title="Actions">
              <template #cell="{ record }">
                <a-button type="text" @click="handleEdit(record)">Edit</a-button>
                <a-popconfirm content="Are you sure you want to delete this supplier?" @ok="handleDelete(record.id)">
                  <a-button type="text" status="danger">Delete</a-button>
                </a-popconfirm>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </a-space>
    </a-card>

    <!-- Add/Edit Supplier Modal -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEditMode ? 'Edit Supplier' : 'Add Supplier'"
      @ok="handleModalOk"
      @cancel="handleModalCancel"
    >
      <a-form :model="currentSupplier" layout="vertical">
        <a-form-item field="supplierId" label="Supplier ID" :rules="[{ required: true, message: 'Supplier ID is required' }]" v-if="!isEditMode">
          <a-input v-model="currentSupplier.supplierId" placeholder="Enter Supplier ID"></a-input>
        </a-form-item>
        <a-form-item field="supplierName" label="Supplier Name" :rules="[{ required: true, message: 'Supplier Name is required' }]">
          <a-input v-model="currentSupplier.supplierName" placeholder="Enter Supplier Name"></a-input>
        </a-form-item>
        <a-form-item field="contact" label="Contact">
          <a-input v-model="currentSupplier.contact" placeholder="Enter Contact Info"></a-input>
        </a-form-item>
        <a-form-item field="status" label="Status">
          <a-select v-model="currentSupplier.status">
            <a-option value="Active">Active</a-option>
            <a-option value="Inactive">Inactive</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="tasksAssigned" label="Tasks Assigned">
          <a-input-number v-model="currentSupplier.tasksAssigned" :min="0"></a-input-number>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import { useI18n } from 'vue-i18n';
  import { 
    Supplier, 
    SupplierCreateData, 
    SupplierUpdateData, 
    addSupplier, 
    getSuppliers, 
    updateSupplier, 
    deleteSupplier 
  } from '@/api/supplier';

  const { t } = useI18n();

  const columns = [
    {
      title: t('yikeYile.supplierId'),
      dataIndex: 'supplierId',
    },
    {
      title: t('yikeYile.supplierName'),
      dataIndex: 'supplierName',
    },
    {
      title: t('yikeYile.contact'),
      dataIndex: 'contact',
    },
    {
      title: t('yikeYile.status'),
      dataIndex: 'status',
    },
    {
      title: t('yikeYile.tasksAssigned'),
      dataIndex: 'tasksAssigned',
    },
    {
      title: t('yikeYile.actions'),
      slotName: 'actions',
    },
  ];

  const tableData = ref<Supplier[]>([]);
  const modalVisible = ref(false);
  const isEditMode = ref(false);

  const currentSupplier = ref<SupplierCreateData | Supplier>({ // Adjusted type
    supplierId: '',
    supplierName: '',
    contact: '',
    status: 'Active',
    tasksAssigned: 0,
  });

  const fetchSuppliers = async () => {
    try {
      const res = await getSuppliers();
      tableData.value = res.data.suppliers;
      Message.success('Suppliers fetched successfully');
    } catch (error) {
      Message.error('Failed to fetch suppliers');
      console.error('Error fetching suppliers:', error);
    }
  };

  const handleAddSupplier = () => {
    isEditMode.value = false;
    currentSupplier.value = {
      supplierId: '',
      supplierName: '',
      contact: '',
      status: 'Active',
      tasksAssigned: 0,
    };
    modalVisible.value = true;
  };

  const handleEdit = (record: Supplier) => {
    isEditMode.value = true;
    currentSupplier.value = { ...record };
    modalVisible.value = true;
  };

  const handleDelete = async (id: number | undefined) => { // Updated type for id
    if (id === undefined) {
        Message.error('Cannot delete: Supplier ID is missing.');
        return;
    }
    try {
      await deleteSupplier(id);
      Message.success('Supplier deleted successfully');
      fetchSuppliers(); // Refresh table
    } catch (error) {
      Message.error('Failed to delete supplier');
      console.error('Error deleting supplier:', error);
    }
  };

  const handleModalOk = async () => {
    // Basic validation
    if (!currentSupplier.value.supplierName || (!isEditMode.value && !currentSupplier.value.supplierId)) {
      Message.error('Please fill in all required fields.');
      return;
    }

    try {
      if (isEditMode.value) {
        // Update existing supplier
        if ('id' in currentSupplier.value && currentSupplier.value.id !== undefined) {
          await updateSupplier(currentSupplier.value.id, currentSupplier.value as SupplierUpdateData);
          Message.success('Supplier updated successfully');
        } else {
          Message.error('Supplier ID is missing for update.');
          return;
        }
      } else {
        // Add new supplier
        await addSupplier(currentSupplier.value as SupplierCreateData);
        Message.success('Supplier added successfully');
      }
      modalVisible.value = false;
      fetchSuppliers(); // Refresh table
    } catch (error) {
      Message.error(`Failed to ${isEditMode.value ? 'update' : 'add'} supplier`);
      console.error(`Error ${isEditMode.value ? 'updating' : 'adding'} supplier:`, error);
    }
  };

  const handleModalCancel = () => {
    modalVisible.value = false;
  };

  onMounted(() => {
    fetchSuppliers();
  });
</script>

<script lang="ts">
  export default {
    name: 'YikeYileManagement',
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