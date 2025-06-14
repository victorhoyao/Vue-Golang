import axios from 'axios';
import { getToken, setToken } from './auth';

// Debug function to test login and get token
export const debugLogin = async () => {
  try {
    const response = await axios.post('/login', {
      userName: 'testadmin',
      passWord: 'test123',
    });

    if (response.data.code === 200 && response.data.data.token) {
      setToken(response.data.data.token);
      return true;
    }
    return false;
  } catch (error: any) {
    return false;
  }
};

// Debug function to test TaskItem creation
export const debugCreateTaskItem = async () => {
  const token = getToken();
  if (!token) {
    return false;
  }

  try {
    const testData = {
      taskId: `TEST_${Date.now()}`,
      taskName: 'Debug Test Task',
      supplier: 'Debug Supplier',
      userTypes: 'premium',
      priceReal: 10.99,
      priceProtocol: 12.99,
      priceManual: 15.99,
      selectionMode: 'Single-select',
      status: 'active',
    };

    const response = await axios.post('/TaskItem/add', testData);

    if (response.data.code === 200) {
      return response.data.data.taskItem;
    }
    return false;
  } catch (error: any) {
    return false;
  }
};

// Debug function to test TaskItem fetching
export const debugFetchTaskItems = async () => {
  const token = getToken();
  if (!token) {
    return false;
  }

  try {
    const response = await axios.get('/TaskItem/list');

    if (response.data.code === 200) {
      return response.data.data.taskItems;
    }
    return false;
  } catch (error: any) {
    return false;
  }
};

// Debug function to test TaskItem update
export const debugUpdateTaskItem = async (id: number) => {
  const token = getToken();
  if (!token) {
    return false;
  }

  try {
    const testData = {
      taskId: `UPDATED_${Date.now()}`,
      taskName: 'Updated Debug Test Task',
      supplier: 'Updated Debug Supplier',
      userTypes: 'premium,basic',
      priceReal: 20.99,
      priceProtocol: 22.99,
      priceManual: 25.99,
      selectionMode: 'Multi-select',
      status: 'active',
    };

    const response = await axios.put(`/TaskItem/update/${id}`, testData);

    if (response.data.code === 200) {
      return response.data.data.taskItem;
    }
    return false;
  } catch (error: any) {
    return false;
  }
};

// Complete debug test sequence
export const runCompleteDebugTest = async () => {
  // Step 1: Login
  const loginSuccess = await debugLogin();
  if (!loginSuccess) {
    return;
  }

  // Step 2: Fetch existing items
  await debugFetchTaskItems();

  // Step 3: Create new item
  const createdItem = await debugCreateTaskItem();
  if (!createdItem) {
    return;
  }

  // Step 4: Update the created item
  if (createdItem.id) {
    await debugUpdateTaskItem(createdItem.id);
  }

  // Step 5: Fetch items again to see changes
  await debugFetchTaskItems();
};

// Make functions available globally for browser console
if (typeof window !== 'undefined') {
  (window as any).debugLogin = debugLogin;
  (window as any).debugCreateTaskItem = debugCreateTaskItem;
  (window as any).debugFetchTaskItems = debugFetchTaskItems;
  (window as any).debugUpdateTaskItem = debugUpdateTaskItem;
  (window as any).runCompleteDebugTest = runCompleteDebugTest;
}
